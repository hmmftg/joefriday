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

// Package mem gets and processes mem info: information for the /proc/meminfo
// file.
package mem

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
)

const procFile = "/proc/meminfo"

// Profiler is used to process the /proc/meminfo file.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func New() (prof *Profiler, err error) {
	proc, err := joe.New(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// Get returns the current meminfo.
func (prof *Profiler) Get() (inf *Info, err error) {
	var (
		i, pos, nameLen int
		v               byte
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	inf = &Info{}
	for l := 0; l < 16; l++ {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return inf, fmt.Errorf("error reading output bytes: %s", err)
		}
		if l > 8 && l < 14 {
			continue
		}
		// first grab the key name (everything up to the ':')
		for i, v = range prof.Line {
			if v == ':' {
				pos = i + 1
				break
			}
			prof.Val = append(prof.Val, v)
		}
		nameLen = len(prof.Val)

		// skip all spaces
		for i, v = range prof.Line[pos:] {
			if v != ' ' {
				pos += i
				break
			}
		}

		// grab the numbers
		for _, v = range prof.Line[pos:] {
			if v == ' ' || v == '\n' {
				break
			}
			prof.Val = append(prof.Val, v)
		}
		// any conversion error results in 0
		n, err := helpers.ParseUint(prof.Val[nameLen:])
		if err != nil {
			return inf, fmt.Errorf("%s: %s", prof.Val[:nameLen], err)
		}

		v = prof.Val[0]

		// Reduce evaluations.
		if v == 'M' {
			v = prof.Val[3]
			if v == 'T' {
				inf.MemTotal = int64(n)
			} else if v == 'F' {
				inf.MemFree = int64(n)
			} else {
				inf.MemAvailable = int64(n)
			}
		} else if v == 'S' {
			v = prof.Val[4]
			if v == 'C' {
				inf.SwapCached = int64(n)
			} else if v == 'T' {
				inf.SwapTotal = int64(n)
			} else if v == 'F' {
				inf.SwapFree = int64(n)
			}
		} else if v == 'B' {
			inf.Buffers = int64(n)
		} else if v == 'I' {
			inf.Inactive = int64(n)
		} else if v == 'C' {
			inf.Cached = int64(n)
		} else if v == 'A' {
			inf.Active = int64(n)
		}
		prof.Val = prof.Val[:0]
	}
	inf.Timestamp = time.Now().UTC().UnixNano()
	return inf, nil
}

// TODO: is it even worth it to have this as a global?  Should GetInfo()
// just instantiate a local version and use that?  InfoTicker does...
var std *Profiler
var stdMu sync.Mutex //protects standard to preven data race on checking/instantiation

// Get returns the current meminfo using the package's global Profiler.
func Get() (inf *Info, err error) {
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
func (prof *Profiler) Ticker(interval time.Duration, out chan Info, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		l, i, pos, nameLen int
		v                  byte
		n                  uint64
		err                error
		inf                Info
	)
	// ticker
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err = prof.Reset()
			if err != nil {
				errs <- joe.Error{Type: "mem", Op: "seek byte 0: /proc/meminfo", Err: err}
				continue
			}
			prof.Line, err = prof.Buf.ReadSlice('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				errs <- fmt.Errorf("error reading output bytes: %s", err)
				continue
			}
			if l > 8 && l < 14 {
				continue
			}
			// first grab the key name (everything up to the ':')
			for i, v = range prof.Line {
				if v == ':' {
					pos = i + 1
					break
				}
				prof.Val = append(prof.Val, v)
			}
			nameLen = len(prof.Val)

			// skip all spaces
			for i, v = range prof.Line[pos:] {
				if v != ' ' {
					pos += i
					break
				}
			}

			// grab the numbers
			for _, v = range prof.Line[pos:] {
				if v == ' ' || v == '\n' {
					break
				}
				prof.Val = append(prof.Val, v)
			}
			// any conversion error results in 0
			n, err = helpers.ParseUint(prof.Val[nameLen:])
			if err != nil {
				errs <- fmt.Errorf("%s: %s", prof.Val[:nameLen], err)
			}
			v = prof.Val[0]

			// Reduce evaluations.
			if v == 'M' {
				v = prof.Val[3]
				if v == 'T' {
					inf.MemTotal = int64(n)
				} else if v == 'F' {
					inf.MemFree = int64(n)
				} else {
					inf.MemAvailable = int64(n)
				}
			} else if v == 'S' {
				v = prof.Val[4]
				if v == 'C' {
					inf.SwapCached = int64(n)
				} else if v == 'T' {
					inf.SwapTotal = int64(n)
				} else if v == 'F' {
					inf.SwapFree = int64(n)
				}
			} else if v == 'B' {
				inf.Buffers = int64(n)
			} else if v == 'I' {
				inf.Inactive = int64(n)
			} else if v == 'C' {
				inf.Cached = int64(n)
			} else if v == 'A' {
				inf.Active = int64(n)
			}
			prof.Val = prof.Val[:0]
		}
		inf.Timestamp = time.Now().UTC().UnixNano()
		out <- inf
	}
}

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan Info, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// Info holds the mem info information.
type Info struct {
	Timestamp    int64 `json:"timestamp"`
	MemTotal     int64 `json:"mem_total"`
	MemFree      int64 `json:"mem_free"`
	MemAvailable int64 `json:"mem_available"`
	Buffers      int64 `json:"buffers"`
	Cached       int64 `json:"cached"`
	SwapCached   int64 `json:"swap_cached"`
	Active       int64 `json:"active"`
	Inactive     int64 `json:"inactive"`
	SwapTotal    int64 `json:"swap_total"`
	SwapFree     int64 `json:"swap_free"`
}

func (i *Info) String() string {
	return fmt.Sprintf("Timestamp: %v\nMemTotal:\t%d\tMemFree:\t%d\tMemAvailable:\t%d\tActive:\t%d\tInactive:\t%d\nCached:\t\t%d\tBuffers\t:%d\nSwapTotal:\t%d\tSwapCached:\t%d\tSwapFree:\t%d\n", time.Unix(0, i.Timestamp).UTC(), i.MemTotal, i.MemFree, i.MemAvailable, i.Active, i.Inactive, i.Cached, i.Buffers, i.SwapTotal, i.SwapCached, i.SwapFree)
}
