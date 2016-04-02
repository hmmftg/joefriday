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

// Package usage gets and processes /proc/net/dev usage stats; the difference,
// in bytes, between two snapshots of /proc/net/dev.
package usage

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/net/info"
	"github.com/mohae/joefriday/net/structs"
)

type Profiler struct {
	*info.Profiler
	prior *structs.Info
}

func New() (prof *Profiler, err error) {
	p, err := info.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: &structs.Info{}}, nil
}

// Get returns the network usage.  Usage calculations requires two pieces of
// data.  This func gets a snapshot of /proc/net/dev, sleeps for a/ second,
// and takes another snapshot and calcualtes the usage from the two snapshots.
// If ongoing usage information is desired, Ticker should be called; it's
// better suited for ongoing usage information: using less cpu cycles and
// generating less garbage.
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

func (prof *Profiler) Ticker(interval time.Duration, out chan *structs.Info, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var i, l, pos, fieldNum, fieldVal int
	var v byte
	var iInfo structs.Interface
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

				// skip leading spaces
				for i, v = range prof.Line {
					if v != 0x20 {
						pos = i
						break
					}
				}
				// first grab the interface name (everything up to the ':')
				for i, v = range prof.Line[pos:] {
					if v == 0x3A {
						iInfo.Name = string(prof.Line[pos : pos+i])
						pos += i + 1
						break
					}
				}
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
							prof.Val = prof.Line[pos : pos+i]
							pos += i + 1
							break
						}
					}
					// any conversion error results in 0
					fieldVal, err = strconv.Atoi(string(prof.Val[:]))
					if err != nil {
						errs <- fmt.Errorf("/proc/net/dev ticker: %s: %s", iInfo.Name, err)
						continue
					}
					prof.Val = prof.Val[:0]
					if fieldNum == 1 {
						iInfo.RBytes = int64(fieldVal)
						continue
					}
					if fieldNum == 2 {
						iInfo.RPackets = int64(fieldVal)
						continue
					}
					if fieldNum == 3 {
						iInfo.RErrs = int64(fieldVal)
						continue
					}
					if fieldNum == 4 {
						iInfo.RDrop = int64(fieldVal)
						continue
					}
					if fieldNum == 5 {
						iInfo.RFIFO = int64(fieldVal)
						continue
					}
					if fieldNum == 6 {
						iInfo.RFrame = int64(fieldVal)
						continue
					}
					if fieldNum == 7 {
						iInfo.RCompressed = int64(fieldVal)
						continue
					}
					if fieldNum == 8 {
						iInfo.RMulticast = int64(fieldVal)
						continue
					}
					if fieldNum == 9 {
						iInfo.TBytes = int64(fieldVal)
						continue
					}
					if fieldNum == 10 {
						iInfo.TPackets = int64(fieldVal)
						continue
					}
					if fieldNum == 11 {
						iInfo.TErrs = int64(fieldVal)
						continue
					}
					if fieldNum == 12 {
						iInfo.TDrop = int64(fieldVal)
						continue
					}
					if fieldNum == 13 {
						iInfo.TFIFO = int64(fieldVal)
						continue
					}
					if fieldNum == 14 {
						iInfo.TColls = int64(fieldVal)
						continue
					}
					if fieldNum == 15 {
						iInfo.TCarrier = int64(fieldVal)
						continue
					}
					if fieldNum == 16 {
						iInfo.TCompressed = int64(fieldVal)
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

func Ticker(interval time.Duration, out chan *structs.Info, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}

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
