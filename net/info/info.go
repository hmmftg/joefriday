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

// Package stat handles processing of network interface information:
// /proc/net/dev using JSON.
package info

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/net/structs"
)

// The proc file used by the Profiler.
const ProcFile = "/proc/net/dev"

// Profiler is used to process the network interface information using the
// /proc/net/dev file.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.New(ProcFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// Get returns the current network interface information.
func (prof *Profiler) Get() (*structs.Info, error) {
	var (
		l, i, pos, fieldNum int
		n                   uint64
		v                   byte
		iInfo               structs.Interface
	)
	err := prof.Reset()
	if err != nil {
		return nil, err
	}
	// there's always at least 2 interfaces (I think)
	inf := &structs.Info{Timestamp: time.Now().UTC().UnixNano(), Interfaces: make([]structs.Interface, 0, 2)}
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		l++
		if l < 3 {
			continue
		}
		prof.Val = prof.Val[:0]
		// first grab the interface name (everything up to the ':')
		for i, v = range prof.Line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			// skip spaces
			if v == 0x20 {
				continue
			}
			prof.Val = append(prof.Val, v)
		}
		iInfo.Name = string(prof.Val[:])
		fieldNum = 0
		// process the rest of the line
		for {
			fieldNum++
			// skip all spaces
			for i, v = range prof.Line[pos:] {
				if v != 0x20 {
					pos += i
					break
				}
			}

			// grab the numbers
			for i, v = range prof.Line[pos:] {
				if v == 0x20 || v == '\n' {
					break
				}
			}
			// any conversion error results in 0
			n, err = helpers.ParseUint(prof.Line[pos : pos+i])
			pos += i
			if err != nil {
				return nil, fmt.Errorf("%s: %s", iInfo.Name, err)
			}
			if fieldNum < 9 {
				if fieldNum < 5 {
					if fieldNum < 3 {
						if fieldNum == 1 {
							iInfo.RBytes = int64(n)
							continue
						}
						iInfo.RPackets = int64(n) // must be 2
						continue
					}
					if fieldNum == 3 {
						iInfo.RErrs = int64(n)
						continue
					}
					iInfo.RDrop = int64(n) // must be 4
					continue
				}
				if fieldNum < 7 {
					if fieldNum == 5 {
						iInfo.RFIFO = int64(n)
						continue
					}
					iInfo.RFrame = int64(n) // must be 6
					continue
				}
				if fieldNum == 7 {
					iInfo.RCompressed = int64(n)
					continue
				}
				iInfo.RMulticast = int64(n) // must be 8
				continue
			}
			if fieldNum < 13 {
				if fieldNum < 11 {
					if fieldNum == 9 {
						iInfo.TBytes = int64(n)
						continue
					}
					iInfo.TPackets = int64(n) // must be 10
					continue
				}
				if fieldNum == 11 {
					iInfo.TErrs = int64(n)
					continue
				}
				iInfo.TDrop = int64(n) // must be 12
				continue
			}
			if fieldNum < 15 {
				if fieldNum == 13 {
					iInfo.TFIFO = int64(n)
					continue
				}
				iInfo.TColls = int64(n) // must be 14
				continue
			}
			if fieldNum == 15 {
				iInfo.TCarrier = int64(n)
				continue
			}
			iInfo.TCompressed = int64(n)
			break
		}
		inf.Interfaces = append(inf.Interfaces, iInfo)
	}
	return inf, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current network interface information using the package's
// global Profiler.
func Get() (inf *structs.Info, err error) {
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

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *structs.Info
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *structs.Info), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	// ticker
	for {
		select {
		case <-t.Done:
			return
		case <-t.Ticker.C:
			s, err := t.Profiler.Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- s
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
