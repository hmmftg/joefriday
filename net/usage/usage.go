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

// Package usage calculates network interface usage.  Usage is calculated by
// taking the difference of two /proc/net/dev snapshots and reflect bytes
// received and transmitted since the prior snapshot.
package usage

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/net/info"
	"github.com/mohae/joefriday/net/structs"
)

// Profiler is used to process the network interface usage..
type Profiler struct {
	*info.Profiler
	prior structs.Info
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := info.NewProfiler()
	if err != nil {
		return nil, err
	}
	prior, err := p.Get()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: *prior}, nil
}

// Get returns the current network interface usage.
// TODO: should this be changed so that this calculates usage since the last
// time the network info was obtained.  If there aren't pre-existing info
// it would get current usage (which may be a separate method (or should be?))
func (prof *Profiler) Get() (u *structs.Usage, err error) {
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

// Get returns the current network interface usage using the package's global
// Profiler.
func Get() (u *structs.Usage, err error) {
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

// CalculateUsage calculates the network interface usage: the ference between
// the current /proc/net/dev data and the prior one.
func (prof *Profiler) CalculateUsage(cur *structs.Info) *structs.Usage {
	u := &structs.Usage{
		Timestamp:  cur.Timestamp,
		TimeDelta:  cur.Timestamp - prof.prior.Timestamp,
		Interfaces: make([]structs.Interface, len(cur.Interfaces)),
	}
	for i := 0; i < len(cur.Interfaces); i++ {
		u.Interfaces[i].Name = cur.Interfaces[i].Name
		u.Interfaces[i].RBytes = cur.Interfaces[i].RBytes - prof.prior.Interfaces[i].RBytes
		u.Interfaces[i].RPackets = cur.Interfaces[i].RPackets - prof.prior.Interfaces[i].RPackets
		u.Interfaces[i].RErrs = cur.Interfaces[i].RErrs - prof.prior.Interfaces[i].RErrs
		u.Interfaces[i].RDrop = cur.Interfaces[i].RDrop - prof.prior.Interfaces[i].RDrop
		u.Interfaces[i].RFIFO = cur.Interfaces[i].RFIFO - prof.prior.Interfaces[i].RFIFO
		u.Interfaces[i].RFrame = cur.Interfaces[i].RFrame - prof.prior.Interfaces[i].RFrame
		u.Interfaces[i].RCompressed = cur.Interfaces[i].RCompressed - prof.prior.Interfaces[i].RCompressed
		u.Interfaces[i].RMulticast = cur.Interfaces[i].RMulticast - prof.prior.Interfaces[i].RMulticast
		u.Interfaces[i].TBytes = cur.Interfaces[i].TBytes - prof.prior.Interfaces[i].TBytes
		u.Interfaces[i].TPackets = cur.Interfaces[i].TPackets - prof.prior.Interfaces[i].TPackets
		u.Interfaces[i].TErrs = cur.Interfaces[i].TErrs - prof.prior.Interfaces[i].TErrs
		u.Interfaces[i].TDrop = cur.Interfaces[i].TDrop - prof.prior.Interfaces[i].TDrop
		u.Interfaces[i].TFIFO = cur.Interfaces[i].TFIFO - prof.prior.Interfaces[i].TFIFO
		u.Interfaces[i].TColls = cur.Interfaces[i].TColls - prof.prior.Interfaces[i].TColls
		u.Interfaces[i].TCarrier = cur.Interfaces[i].TCarrier - prof.prior.Interfaces[i].TCarrier
		u.Interfaces[i].TCompressed = cur.Interfaces[i].TCompressed - prof.prior.Interfaces[i].TCompressed
	}
	return u
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *structs.Usage
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *structs.Usage), Profiler: p}
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
		cur                    structs.Info
		iUsage                 structs.Interface
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
			cur.Interfaces = cur.Interfaces[:0]
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
				iUsage.Name = string(t.Val[:])
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
									iUsage.RBytes = int64(n)
									continue
								}
								iUsage.RPackets = int64(n) // must be 2
								continue
							}
							if fieldNum == 3 {
								iUsage.RErrs = int64(n)
								continue
							}
							iUsage.RDrop = int64(n) // must be 4
							continue
						}
						if fieldNum < 7 {
							if fieldNum == 5 {
								iUsage.RFIFO = int64(n)
								continue
							}
							iUsage.RFrame = int64(n) // must be 6
							continue
						}
						if fieldNum == 7 {
							iUsage.RCompressed = int64(n)
							continue
						}
						iUsage.RMulticast = int64(n) // must be 8
						continue
					}
					if fieldNum < 13 {
						if fieldNum < 11 {
							if fieldNum == 9 {
								iUsage.TBytes = int64(n)
								continue
							}
							iUsage.TPackets = int64(n) // must be 10
							continue
						}
						if fieldNum == 11 {
							iUsage.TErrs = int64(n)
							continue
						}
						iUsage.TDrop = int64(n) // must be 12
						continue
					}
					if fieldNum < 15 {
						if fieldNum == 13 {
							iUsage.TFIFO = int64(n)
							continue
						}
						iUsage.TColls = int64(n) // must be 14
						continue
					}
					if fieldNum == 15 {
						iUsage.TCarrier = int64(n)
						continue
					}
					if fieldNum == 16 {
						iUsage.TCompressed = int64(n)
						break
					}
				}
				cur.Interfaces = append(cur.Interfaces, iUsage)
			}
			t.Data <- t.CalculateUsage(&cur)
			t.prior.Timestamp = cur.Timestamp
			if len(t.prior.Interfaces) != len(cur.Interfaces) {
				t.prior.Interfaces = make([]structs.Interface, len(cur.Interfaces))
			}
			copy(t.prior.Interfaces, cur.Interfaces)
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
