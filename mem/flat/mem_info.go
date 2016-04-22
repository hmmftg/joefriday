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

// Package flat handles Flatbuffer based processing of /proc/meminfo.
// Instead of returning a Go struct, it returns Flatbuffer serialized bytes.
// A function to deserialize the Flatbuffer serialized bytes into a
// mem.Info struct is provided.  After the first use, the flatbuffer
// builder is reused.
package flat

import (
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/mem"
)

// Profiler is used to process the /proc/meminfo file using Flatbuffers.
type Profiler struct {
	*mem.Profiler
	*fb.Builder
}

// Initializes and returns a mem info profiler that utilizes FlatBuffers.
func NewProfiler() (prof *Profiler, err error) {
	p, err := mem.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current meminfo as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	inf, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current meminfo as Flatbuffer serialized bytes using the
// package's global Profiler.
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

// Serialize mem.Info using Flatbuffers.
func (prof *Profiler) Serialize(inf *mem.Info) []byte {
	// ensure the Builder is in a usable state.
	std.Builder.Reset()
	InfoStart(prof.Builder)
	InfoAddTimestamp(prof.Builder, int64(inf.Timestamp))
	InfoAddMemTotal(prof.Builder, int64(inf.MemTotal))
	InfoAddMemFree(prof.Builder, int64(inf.MemFree))
	InfoAddMemAvailable(prof.Builder, int64(inf.MemAvailable))
	InfoAddBuffers(prof.Builder, int64(inf.Buffers))
	InfoAddCached(prof.Builder, int64(inf.Cached))
	InfoAddSwapCached(prof.Builder, int64(inf.SwapCached))
	InfoAddActive(prof.Builder, int64(inf.Active))
	InfoAddInactive(prof.Builder, int64(inf.Inactive))
	InfoAddSwapTotal(prof.Builder, int64(inf.SwapTotal))
	InfoAddSwapFree(prof.Builder, int64(inf.SwapFree))
	prof.Builder.Finish(InfoEnd(prof.Builder))
	return prof.Builder.Bytes[prof.Builder.Head():]
}

// Serialize mem.Info using Flatbuffers with the package global Profiler.
func Serialize(inf *mem.Info) (p []byte, err error) {
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
func Deserialize(p []byte) *mem.Info {
	infoFlat := GetRootAsInfo(p, 0)
	info := &mem.Info{}
	info.Timestamp = infoFlat.Timestamp()
	info.MemTotal = infoFlat.MemTotal()
	info.MemFree = infoFlat.MemFree()
	info.MemAvailable = infoFlat.MemAvailable()
	info.Buffers = infoFlat.Buffers()
	info.Cached = infoFlat.Cached()
	info.SwapCached = infoFlat.SwapCached()
	info.Active = infoFlat.Active()
	info.Inactive = infoFlat.Inactive()
	info.SwapTotal = infoFlat.SwapTotal()
	info.SwapFree = infoFlat.SwapFree()
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
		i, pos, line, nameLen int
		v                     byte
		n                     uint64
		err                   error
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
			InfoStart(t.Builder)
			InfoAddTimestamp(t.Builder, time.Now().UTC().UnixNano())
			for line = 0; line < 16; line++ {
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
				if line > 7 && line < 14 {
					continue
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
				if v == 'M' {
					v = t.Val[3]
					if v == 'T' {
						InfoAddMemTotal(t.Builder, int64(n))
					} else if v == 'F' {
						InfoAddMemFree(t.Builder, int64(n))
					} else {
						InfoAddMemAvailable(t.Builder, int64(n))
					}
				} else if v == 'S' {
					v = t.Val[4]
					if v == 'C' {
						InfoAddSwapCached(t.Builder, int64(n))
					} else if v == 'T' {
						InfoAddSwapTotal(t.Builder, int64(n))
					} else if v == 'F' {
						InfoAddSwapFree(t.Builder, int64(n))
					}
				} else if v == 'B' {
					InfoAddBuffers(t.Builder, int64(n))
				} else if v == 'I' {
					InfoAddInactive(t.Builder, int64(n))
				} else if v == 'C' {
					InfoAddMemAvailable(t.Builder, int64(n))
				} else if v == 'A' {
					InfoAddInactive(t.Builder, int64(n))
				}
			}
			t.Builder.Finish(InfoEnd(t.Builder))
			t.Data <- t.Profiler.Builder.Bytes[t.Builder.Head():]
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
