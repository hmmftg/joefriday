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

// Package usage calculates network usage.  Usage is calculated by taking the
// difference in two /proc/net/dev snapshots and reflect bytes received and
// transmitted since the prior snapshot.
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

// Profiler is used to process the network usage..
type Profiler struct {
	*info.Profiler
	prior *structs.Info
}

// Returns an initialized Profiler; ready to use.
func New() (prof *Profiler, err error) {
	p, err := info.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: &structs.Info{}}, nil
}

// Get returns the current network usage.
// TODO: should this be changed so that this calculates usage since the last
// time the network info was obtained.  If there aren't pre-existing info
// it would get current usage (which may be a separate method (or should be?))
func (prof *Profiler) Get() (u *structs.Info, err error) {
	prof.prior, err = prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	infCur, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.CalculateUsage(infCur), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current network usage using the package's global Profiler..
func Get() (u *structs.Info, err error) {
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

// Ticker calculates network usage on a ticker.  The generated data is sent
// to the out channel.  Any errors encountered are sent to the errs channel.
//  Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *Profiler) Ticker(interval time.Duration, out chan *structs.Info, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		i, l, pos, fieldNum int
		n                   uint64
		v                   byte
		iInfo               structs.Interface
	)
	// first get Info as the baseline
	cur, err := prof.Profiler.Get()
	if err != nil {
		errs <- err
		return
	}
	// ticker
tick:
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			prof.prior.Timestamp = cur.Timestamp
			if len(prof.prior.Interfaces) != len(cur.Interfaces) {
				prof.prior.Interfaces = make([]structs.Interface, len(cur.Interfaces))
			}
			copy(prof.prior.Interfaces, cur.Interfaces)
			cur.Timestamp = time.Now().UTC().UnixNano()
			err = prof.Reset()
			if err != nil {
				errs <- joe.Error{Type: "net", Op: "usage ticker", Err: err}
				continue tick
			}
			cur.Interfaces = cur.Interfaces[:0]
			// read each line until eof
			for {
				prof.Val = prof.Val[:0]
				prof.Line, err = prof.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errs <- fmt.Errorf("/proc/mem/dev: read output bytes: %s", err)
					break
				}
				l++
				if l < 3 {
					continue
				}
				// first grab the interface name (everything up to the ':')
				for i, v = range prof.Line {
					if v == 0x3A {
						pos = i + 1
						break
					}
					// skip spaces
					if v != 0x20 {
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
					n, err = helpers.ParseUint(prof.Line[pos : pos+1])
					pos += i
					if err != nil {
						errs <- fmt.Errorf("/proc/net/dev ticker: %s: %s", iInfo.Name, err)
						continue
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
					if fieldNum == 16 {
						iInfo.TCompressed = int64(n)
						break
					}
				}
				cur.Interfaces = append(cur.Interfaces, iInfo)
			}
			out <- prof.CalculateUsage(cur)
			l = 0
		}
	}
}

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan *structs.Info, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// CalculateUsage returns the difference between the current /proc/net/dev
// data and the prior one.
func (prof *Profiler) CalculateUsage(cur *structs.Info) *structs.Info {
	u := &structs.Info{Timestamp: cur.Timestamp, Interfaces: make([]structs.Interface, len(cur.Interfaces))}
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
