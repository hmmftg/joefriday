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

// Package netusage calculates network devices usage. Usage is calculated by
// taking the difference of two /proc/net/dev snapshots; the elapsed time
// between the two snapshots is stored in the TimeDelta field.
package netusage

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/net/netdev"
	"github.com/mohae/joefriday/net/structs"
)

// Profiler is used to process the network devices usage.
type Profiler struct {
	*netdev.Profiler
	prior structs.DevInfo
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := netdev.NewProfiler()
	if err != nil {
		return nil, err
	}
	prior, err := p.Get()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: *prior}, nil
}

// Get returns the current network devices usage: the delta between the current
// snapshot and the prior one.
func (prof *Profiler) Get() (u *structs.DevUsage, err error) {
	infCur, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	u = prof.CalculateUsage(infCur)
	prof.prior = *infCur
	return u, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current network devices usage using the package's global
// Profiler.
func Get() (u *structs.DevUsage, err error) {
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

// CalculateUsage calculates the network devices usage: the difference between
// the current /proc/net/dev snapshot and the prior one.
func (prof *Profiler) CalculateUsage(cur *structs.DevInfo) *structs.DevUsage {
	u := &structs.DevUsage{
		Timestamp:  cur.Timestamp,
		TimeDelta:  cur.Timestamp - prof.prior.Timestamp,
		Devices: make([]structs.Device, len(cur.Devices)),
	}
	for i := 0; i < len(cur.Devices); i++ {
		u.Devices[i].Name = cur.Devices[i].Name
		u.Devices[i].RBytes = cur.Devices[i].RBytes - prof.prior.Devices[i].RBytes
		u.Devices[i].RPackets = cur.Devices[i].RPackets - prof.prior.Devices[i].RPackets
		u.Devices[i].RErrs = cur.Devices[i].RErrs - prof.prior.Devices[i].RErrs
		u.Devices[i].RDrop = cur.Devices[i].RDrop - prof.prior.Devices[i].RDrop
		u.Devices[i].RFIFO = cur.Devices[i].RFIFO - prof.prior.Devices[i].RFIFO
		u.Devices[i].RFrame = cur.Devices[i].RFrame - prof.prior.Devices[i].RFrame
		u.Devices[i].RCompressed = cur.Devices[i].RCompressed - prof.prior.Devices[i].RCompressed
		u.Devices[i].RMulticast = cur.Devices[i].RMulticast - prof.prior.Devices[i].RMulticast
		u.Devices[i].TBytes = cur.Devices[i].TBytes - prof.prior.Devices[i].TBytes
		u.Devices[i].TPackets = cur.Devices[i].TPackets - prof.prior.Devices[i].TPackets
		u.Devices[i].TErrs = cur.Devices[i].TErrs - prof.prior.Devices[i].TErrs
		u.Devices[i].TDrop = cur.Devices[i].TDrop - prof.prior.Devices[i].TDrop
		u.Devices[i].TFIFO = cur.Devices[i].TFIFO - prof.prior.Devices[i].TFIFO
		u.Devices[i].TColls = cur.Devices[i].TColls - prof.prior.Devices[i].TColls
		u.Devices[i].TCarrier = cur.Devices[i].TCarrier - prof.prior.Devices[i].TCarrier
		u.Devices[i].TCompressed = cur.Devices[i].TCompressed - prof.prior.Devices[i].TCompressed
	}
	return u
}

// Ticker delivers the system's network devices usage at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *structs.DevUsage
	*Profiler
}

// NewTicker returns a new Ticker containing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	p, err := NewProfiler()
	if err != nil {
		return nil, err
	}
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *structs.DevUsage), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	var (
		i, pos, line, fieldNum int
		n                      uint64
		v                      byte
		err                    error
		cur                    structs.DevInfo
		dev                 structs.Device
	)
	// ticker
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			cur.Timestamp = time.Now().UTC().UnixNano()
			err = t.Reset()
			if err != nil {
				t.Errs <- err
				break
			}
			line = 0
			cur.Devices = cur.Devices[:0]
			// read each line until eof
			for {
				t.Val = t.Val[:0]
				t.Line, err = t.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Errs <- &joe.ReadError{Err: err}
					break
				}
				line++
				if line < 3 {
					continue
				}
				// first grab the interface name (everything up to the ':')
				for i, v = range t.Line {
					if v == 0x3A {
						pos = i + 1
						break
					}
					// skip spaces
					if v != 0x20 {
						continue
					}
					t.Val = append(t.Val, v)
				}
				dev.Name = string(t.Val[:])
				fieldNum = 0
				// process the rest of the line
				for {
					fieldNum++
					// skip all spaces
					for i, v = range t.Line[pos:] {
						if v != 0x20 {
							pos += i
							break
						}
					}
					// grab the numbers
					for i, v = range t.Line[pos:] {
						if v == 0x20 || v == '\n' {
							break
						}
					}
					// any conversion error results in 0
					n, err = helpers.ParseUint(t.Line[pos : pos+1])
					pos += i
					if err != nil {
						t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
						continue
					}
					if fieldNum < 9 {
						if fieldNum < 5 {
							if fieldNum < 3 {
								if fieldNum == 1 {
									dev.RBytes = int64(n)
									continue
								}
								dev.RPackets = int64(n) // must be 2
								continue
							}
							if fieldNum == 3 {
								dev.RErrs = int64(n)
								continue
							}
							dev.RDrop = int64(n) // must be 4
							continue
						}
						if fieldNum < 7 {
							if fieldNum == 5 {
								dev.RFIFO = int64(n)
								continue
							}
							dev.RFrame = int64(n) // must be 6
							continue
						}
						if fieldNum == 7 {
							dev.RCompressed = int64(n)
							continue
						}
						dev.RMulticast = int64(n) // must be 8
						continue
					}
					if fieldNum < 13 {
						if fieldNum < 11 {
							if fieldNum == 9 {
								dev.TBytes = int64(n)
								continue
							}
							dev.TPackets = int64(n) // must be 10
							continue
						}
						if fieldNum == 11 {
							dev.TErrs = int64(n)
							continue
						}
						dev.TDrop = int64(n) // must be 12
						continue
					}
					if fieldNum < 15 {
						if fieldNum == 13 {
							dev.TFIFO = int64(n)
							continue
						}
						dev.TColls = int64(n) // must be 14
						continue
					}
					if fieldNum == 15 {
						dev.TCarrier = int64(n)
						continue
					}
					if fieldNum == 16 {
						dev.TCompressed = int64(n)
						break
					}
				}
				cur.Devices = append(cur.Devices, dev)
			}
			t.Data <- t.CalculateUsage(&cur)
			t.prior.Timestamp = cur.Timestamp
			if len(t.prior.Devices) != len(cur.Devices) {
				t.prior.Devices = make([]structs.Device, len(cur.Devices))
			}
			copy(t.prior.Devices, cur.Devices)
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
