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
// facts.Facts struct is provided.  After the first use, the flatbuffer
// builder is reused.
package flat

import (
	"fmt"
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
	Prof    *mem.Profiler
	Builder *fb.Builder
}

// Initializes and returns a mem info profiler that utilizes FlatBuffers.
func New() (prof *Profiler, err error) {
	p, err := mem.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Prof: p, Builder: fb.NewBuilder(0)}, nil
}

func (prof *Profiler) reset() {
	prof.Builder.Reset()
}

// Get returns the current meminfo as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	prof.reset()
	inf, err := prof.Prof.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(inf), nil
}

// TODO: is it even worth it to have this as a global?  Should GetInfo()
// just instantiate a local version and use that?  InfoTicker does...
var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current meminfo as Flatbuffer serialized bytes using the
// package's global Profiler.
func Get() (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Ticker processes meminfo information on a ticker.  The generated data is
// sent to the out channel.  Any errors encountered are sent to the errs
// channel.  Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
func (prof *Profiler) Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		l, i, pos, nameLen int
		v                  byte
		n                  uint64
		err                error
	)
	// ticker
Tick:
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			prof.reset()
			err = prof.Prof.Reset()
			if err != nil {
				errs <- joe.Error{Type: "mem", Op: "seek byte 0: /proc/meminfo", Err: err}
				continue
			}
			InfoStart(prof.Builder)
			InfoAddTimestamp(prof.Builder, time.Now().UTC().UnixNano())
			for l = 0; l < 16; l++ {
				prof.Prof.Line, err = prof.Prof.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					// An error results in sending error message and stop processing of this tick.
					errs <- joe.Error{Type: "mem", Op: "read output bytes", Err: err}
					continue Tick
				}
				if l > 7 && l < 14 {
					continue
				}
				// first grab the key name (everything up to the ':')
				for i, v = range prof.Prof.Line {
					if v == 0x3A {
						prof.Prof.Val = prof.Prof.Line[:i]
						break
					}
				}
				nameLen = len(prof.Prof.Val)
				// skip all spaces
				for i, v = range prof.Prof.Line[pos:] {
					if v != 0x20 {
						pos += i
						break
					}
				}

				// grab the numbers
				for _, v = range prof.Prof.Line[pos:] {
					if v == 0x20 || v == '\n' {
						break
					}
					prof.Prof.Val = append(prof.Prof.Val, v)
				}
				// any conversion error results in 0
				n, err = helpers.ParseUint(prof.Prof.Val[nameLen:])
				if err != nil {
					errs <- joe.Error{Type: "mem", Op: fmt.Sprintf("convert %s", prof.Prof.Val[:nameLen]), Err: err}
					continue
				}
				v = prof.Prof.Val[0]
				if v == 'M' {
					v = prof.Prof.Val[3]
					if v == 'T' {
						InfoAddMemTotal(prof.Builder, int64(n))
					} else if v == 'F' {
						InfoAddMemFree(prof.Builder, int64(n))
					} else {
						InfoAddMemAvailable(prof.Builder, int64(n))
					}
				} else if v == 'S' {
					v = prof.Prof.Val[4]
					if v == 'C' {
						InfoAddSwapCached(prof.Builder, int64(n))
					} else if v == 'T' {
						InfoAddSwapTotal(prof.Builder, int64(n))
					} else if v == 'F' {
						InfoAddSwapFree(prof.Builder, int64(n))
					}
				} else if v == 'B' {
					InfoAddBuffers(prof.Builder, int64(n))
				} else if v == 'I' {
					InfoAddInactive(prof.Builder, int64(n))
				} else if v == 'C' {
					InfoAddMemAvailable(prof.Builder, int64(n))
				} else if v == 'A' {
					InfoAddInactive(prof.Builder, int64(n))
				}
			}
			prof.Builder.Finish(InfoEnd(prof.Builder))
			inf := prof.Builder.Bytes[prof.Builder.Head():]
			out <- inf
		}
	}
}

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan []byte, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// Serialize mem.Info using Flatbuffers.
func (prof *Profiler) Serialize(inf *mem.Info) []byte {
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
