// Copyright 2016 Joel Scoble and The JoeFriday authors.
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

// Package flat handles Flatbuffer based processing of some of /proc/meminfo's
// data.  Instead of returning a Go struct, it returns Flatbuffer serialized
// bytes.  A function to deserialize the Flatbuffer serialized bytes into a
// MemInfo struct is provided.  After the first use, the flatbuffer builder is
// reused.
package flat

import (
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/mem/basic"
)

// Profiler is used to process the /proc/meminfo file, extracting basic info,
// using Flatbuffers.
type Profiler struct {
	*basic.Profiler
	*fb.Builder
}

// Initializes and returns a basic meminfo profiler that utilizes FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := basic.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current basic meminfo as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current basic meminfo as Flatbuffer serialized bytes using
// the package's global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Serialize MemInfo using Flatbuffers.
func (prof *Profiler) Serialize(inf *basic.MemInfo) []byte {
	// ensure the Builder is in a usable state.
	prof.Builder.Reset()
	MemInfoStart(prof.Builder)
	MemInfoAddTimestamp(prof.Builder, inf.Timestamp)
	MemInfoAddActive(prof.Builder, inf.Active)
	MemInfoAddInactive(prof.Builder, inf.Inactive)
	MemInfoAddMapped(prof.Builder, inf.Mapped)
	MemInfoAddMemAvailable(prof.Builder, inf.MemAvailable)
	MemInfoAddMemFree(prof.Builder, inf.MemFree)
	MemInfoAddMemTotal(prof.Builder, inf.MemTotal)
	MemInfoAddSwapCached(prof.Builder, inf.SwapCached)
	MemInfoAddSwapFree(prof.Builder, inf.SwapFree)
	MemInfoAddSwapTotal(prof.Builder, inf.SwapTotal)
	prof.Builder.Finish(MemInfoEnd(prof.Builder))
	p := prof.Builder.Bytes[prof.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(p))
	copy(tmp, p)
	return tmp
}

// Serialize MemInfo using Flatbuffers with the package global Profiler.
func Serialize(inf *basic.MemInfo) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(inf), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as mem.Info.
func Deserialize(p []byte) *basic.MemInfo {
	infoFlat := GetRootAsMemInfo(p, 0)
	info := &basic.MemInfo{}
	info.Timestamp = infoFlat.Timestamp()
	info.Active = infoFlat.Active()
	info.Inactive = infoFlat.Inactive()
	info.Mapped = infoFlat.Mapped()
	info.MemAvailable = infoFlat.MemAvailable()
	info.MemFree = infoFlat.MemFree()
	info.MemTotal = infoFlat.MemTotal()
	info.SwapCached = infoFlat.SwapCached()
	info.SwapFree = infoFlat.SwapFree()
	info.SwapTotal = infoFlat.SwapTotal()
	return info
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan []byte
	*Profiler
}

// NewTicker returns a new Ticker continaing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	p, err := NewProfiler()
	if err != nil {
		return nil, err
	}
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan []byte), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	// predeclare some vars
	var (
		i, pos, nameLen int
		v               byte
		n               uint64
		err             error
	)
	// ticker
Tick:
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			t.Builder.Reset()
			err = t.Profiler.Profiler.Reset()
			if err != nil {
				t.Errs <- err
				continue
			}
			MemInfoStart(t.Builder)
			MemInfoAddTimestamp(t.Builder, time.Now().UTC().UnixNano())
			for {
				t.Val = t.Val[:0]
				t.Line, err = t.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					// An error results in sending error message and stop processing of this tick.
					t.Errs <- &joe.ReadError{Err: err}
					continue Tick
				}
				// first grab the key name (everything up to the ':')
				for i, v = range t.Line {
					if v == 0x3A {
						t.Val = t.Line[:i]
						pos = i + 1 // skip the :
						break
					}
				}
				nameLen = len(t.Val)
				// skip all spaces
				for i, v = range t.Line[pos:] {
					if v != 0x20 {
						pos += i
						break
					}
				}

				// grab the numbers
				for _, v = range t.Line[pos:] {
					if v == 0x20 || v == '\n' {
						break
					}
					t.Val = append(t.Val, v)
				}
				// any conversion error results in 0
				n, err = helpers.ParseUint(t.Val[nameLen:])
				if err != nil {
					t.Errs <- &joe.ParseError{Info: string(t.Val[:nameLen]), Err: err}
					continue
				}
				v = t.Val[0]
				if v == 'A' {
					if t.Val[5] == 'e' && nameLen == 6 {
						MemInfoAddActive(t.Builder, n)
					}
					continue
				}
				if v == 'I' {
					if nameLen == 8 {
						MemInfoAddInactive(t.Builder, n)
					}
					continue
				}
				if v == 'M' {
					v = t.Val[3]
					if nameLen < 8 {
						if v == 'p' {
							MemInfoAddMapped(t.Builder, n)
							continue
						}
						if v == 'F' {
							MemInfoAddMemFree(t.Builder, n)
						}
						continue
					}
					if v == 'A' {
						MemInfoAddMemAvailable(t.Builder, n)
						continue
					}
					MemInfoAddMemTotal(t.Builder, n)
					continue
				}
				if v == 'S' {
					v = t.Val[1]
					if v == 'w' {
						if t.Val[4] == 'C' {
							MemInfoAddSwapCached(t.Builder, n)
							continue
						}
						if t.Val[4] == 'F' {
							MemInfoAddSwapFree(t.Builder, n)
							continue
						}
						MemInfoAddSwapTotal(t.Builder, n)
					}
				}
			}
			t.Builder.Finish(MemInfoEnd(t.Builder))
			t.Data <- t.Profiler.Builder.Bytes[t.Builder.Head():]
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
