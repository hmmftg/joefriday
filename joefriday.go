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
	"os"
	"strconv"
	"time"
)

// TODO: make current/better implementation
type Error struct {
	Type string
	Op   string
	Err  error
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %q: %s", e.Type, e.Op, e.Err)
}

// A Proc holds everything related to a proc file and some processing vars.
type Proc struct {
	*os.File
	Buf  *bufio.Reader
	Line []byte // current line
	Val  []byte
}

// Creats a Proc using the file handle.
func New(fname string) (*Proc, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	return &Proc{File: f, Buf: bufio.NewReader(f), Val: make([]byte, 0, 32)}, nil
}

// Reset reset's the profiler's resources.
func (p *Proc) Reset() error {
	_, err := p.File.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}
	p.Buf.Reset(p.File)
	p.Val = p.Val[:0]
	return nil
}

type Tocker interface {
	Run()  // runs some code on an interval
	Stop() // stops execution of the code
}

type Ticker struct {
	*time.Ticker
	Done chan struct{} // done channel
	Errs chan error    // error channel
}

func NewTicker(d time.Duration) *Ticker {
	return &Ticker{Ticker: time.NewTicker(d), Errs: make(chan error), Done: make(chan struct{})}
}

// Close stops the ticker and closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Stop()
	close(t.Done)
	close(t.Errs)
}

// ProfileSerializer is implemented by any profiler that has a Get method
// that returns the data as serialized bytes.
type ProfileSerializer interface {
	Get() ([]byte, error)
}

// ProfileSerializerTicker is implemented by any profiler that defines a
// Ticker method that uses a ticker to get the current data.  It is used
// for ongoing monitoring of something.  The current data is sent to the
// out channel as serialized bytes.
type ProfileSerializerTicker interface {
	Ticker(tick time.Duration, out chan []byte, done chan struct{}, errs chan error)
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
