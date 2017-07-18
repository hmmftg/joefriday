// Copyright 2016 The JoeFriday authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package joefriday gets information about a system: platform, kernel,
// memory information, cpu information, cpu stats, cpu utilization, network
// information, and network usage.
//
// Ticker versions of non-static information are available to enable
// monitoring.
//
// The data can be returned as Go structs, Flatbuffer serialized bytes, or
// JSON serialized bytes.  For convenience, there are deserialization
// functions for all structs that are serialized.
package joefriday

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mohae/randchars"
)

type ResetError struct {
	Err error
}

func (e *ResetError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Err.Error()
}

type ParseError struct {
	Info string
	Err  error
}

func (e *ParseError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Info
	if e.Info != "" {
		s += ": "
	}
	s += e.Err.Error()
	return s
}

type ReadError struct {
	Info string
	Err  error
}

func (e *ReadError) Error() string {
	if e == nil {
		return "<nil>"
	}
	s := e.Info
	if e.Info != "" {
		s += ": "
	}
	s += e.Err.Error()
	return s
}

// IsReadError returns a boolean indicating whether the error is a result of
// a read problem.
func IsReadError(e error) bool {
	if _, ok := e.(*ReadError); ok {
		return true
	}
	return false
}

// IsResetError returns a boolean indicating whether the error is a result of
// a problem resetting the file buffer.
func IsResetError(e error) bool {
	if _, ok := e.(*ResetError); ok {
		return true
	}
	return false
}

// IsParseError r eturns a boolean indicating whether the error is a result of
// encountering a problem while trying to parse the file data.
func IsParseError(e error) bool {
	if _, ok := e.(*ParseError); ok {
		return true
	}
	return false
}

// Procer processes things.
type Procer interface {
	ReadSlice(byte) ([]byte, error)
	Reset() error
}
// A Proc holds everything related to a proc file and some processing vars.
type Proc struct {
	*os.File
	Buf  *bufio.Reader
}

// Creats a Proc using the file handle.
func New(fname string) (*Proc, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	return &Proc{File: f, Buf: bufio.NewReader(f)}, nil
}

// ReadSlice is a wrapper for bufio.Reader.ReadSlice.
func (p *Proc) ReadSlice(delim byte) (line []byte, err error) {
	return p.Buf.ReadSlice(delim)
}

// Reset reset's the profiler's resources.
func (p *Proc) Reset() error {
	_, err := p.File.Seek(0, os.SEEK_SET)
	if err != nil {
		return &ResetError{err}
	}
	p.Buf.Reset(p.File)
	return nil
}

type Tocker interface {
	Close() // Close the Tocker's resources
	Run()   // Run some code on an interval.
	Stop()  // Stop the Tocker.
}

type Ticker struct {
	*time.Ticker
	Done chan struct{} // done channel
	Errs chan error    // error channel
}

func NewTicker(d time.Duration) *Ticker {
	return &Ticker{Ticker: time.NewTicker(d), Errs: make(chan error), Done: make(chan struct{})}
}

// Stop sends a signal to the done channel; stopping the Ticker.  The Ticker
// can be restarted with Run.
func (t *Ticker) Stop() {
	t.Done <- struct{}{}
}

// Close stops the ticker and closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Stop()
	close(t.Done)
	close(t.Errs)
}

// Column returns a right justified string of width w.
// TODO: replace with text/tabwriter
func Column(w int, s string) string {
	pad := w - len(s)
	padding := make([]byte, pad)
	for i := 0; i < pad; i++ {
		padding[i] = 0x20
	}
	return fmt.Sprintf("%s%s", string(padding), s)
}

// Int64Column takes an int64 and returns a right justified string of width w.
func Int64Column(w int, v int64) string {
	s := strconv.FormatInt(v, 10)
	return Column(w, s)
}

// TrimTrailingSpaces removes the trailing spaces from a slice and returns
// it.  Only 0x20, tabs, NL are considered space characters.
func TrimTrailingSpaces(p []byte) []byte {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] != 0x20 && p[i] != '\n' && p[i] != '\t' {
			return p[:i+1]
		}
	}
	// it was all spaces
	return p[:0]
}

// TrimLeadingpaces removes the leading spaces from a slice and returns it.
// Only 0x20 and tabs are considered space characters.
func TrimLeadingSpaces(p []byte) []byte {
	for i := 0; i < len(p); i++ {
		if p[i] != 0x20 && p[i] != '\t' {
			return p[i:]
		}
	}
	// it was all spaces
	return p[:0]
}

// TempFileProc is used to do Proc processing off of a temp file. Prefer using
// the Proc type instead.
type TempFileProc struct {
	*Proc
	// The directory holding the temp file.
	Dir string
	// The name of the file.
	Name string
}

// NewTempFileProc creates a temporary file with data as its contents and
// returns a TempFileProc that uses the temporary file. The file will be saved
// in a randomly generated tempdir that starts with prefix. If prefix is empty
// the os.TempDir will be used as the save directory. Name is the name of the
// temporary file that will be created. If name is empty, the name will be 12
// randomly selected characters without an extension. The data will be used for
// the file and the file will be created with 0777 perms. If an error occurs,
// proc will be nil.
func NewTempFileProc(prefix, name string, data []byte) (proc *TempFileProc, err error) {
	var t TempFileProc
	if prefix == "" {
		t.Dir = os.TempDir()
	} else {
		t.Dir, err = ioutil.TempDir("", prefix)
		if err != nil {
			return nil, err
		}
	}
	if name == "" {
		name = string(randchars.AlphaNum(12))
	}
	t.Name = name

	err = ioutil.WriteFile(t.FullPath(), data, 0777)
	if err != nil {
		return nil, err
	}

	t.Proc, err = New(t.FullPath())
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// Returns the full path of the temp file for this TempFileProc.
func (p *TempFileProc) FullPath() string {
	return filepath.Join(p.Dir, p.Name)
}

// ReadSlice is a wrapper for bufio.Reader.ReadSlice.
func (p *TempFileProc) ReadSlice(delim byte) (line []byte, err error) {
	return p.Buf.ReadSlice(delim)
}

// Reset reset's the profiler's resources.
func (p *TempFileProc) Reset() error {
	_, err := p.File.Seek(0, os.SEEK_SET)
	if err != nil {
		return &ResetError{err}
	}
	p.Buf.Reset(p.File)
	return nil
}

// Remove removes the temp dir and temp file.
func (p *TempFileProc) Remove() error {
	// only remove the directory if it is a subdir of the default temp dir.
	if p.Dir != os.TempDir() {
		os.RemoveAll(p.Dir)
	}
	// otherwise just remove the file
	return os.RemoveAll(p.FullPath())
}
