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

// Package usage calculates disk usage.  Usage is calculated by taking the
// difference in two /proc/diskstats snapshots and reflect the difference
// between the two snapshots.
package usage

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/disk/stats"
	"github.com/mohae/joefriday/disk/structs"
)

// Profiler is used to process the disk usage..
type Profiler struct {
	*stats.Profiler
	prior *structs.Stats
}

// Returns an initialized Profiler; ready to use.
func New() (prof *Profiler, err error) {
	p, err := stats.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: &structs.Stats{}}, nil
}

// Get returns the current disk usage.
// TODO: should this be changed so that this calculates usage since the last
// time the disk stats were obtained.  If there aren't pre-existing stats
// it would get current usage (which may be a separate method (or should be?))
func (prof *Profiler) Get() (u *structs.Usage, err error) {
	prof.prior, err = prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	st, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.CalculateUsage(st), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current disk usage using the package's global Profiler..
func Get() (u *structs.Usage, err error) {
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

// Ticker calculates disk usage on a ticker.  The generated data is sent to
// the out channel.  Any errors encountered are sent to the errs channel.
// Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *Profiler) Ticker(interval time.Duration, out chan *structs.Usage, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		i, priorPos, pos, fieldNum int
		n                          uint64
		v                          byte
		dev                        structs.Device
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
			if len(prof.prior.Devices) != len(cur.Devices) {
				prof.prior.Devices = make([]structs.Device, len(cur.Devices))
			}
			copy(prof.prior.Devices, cur.Devices)
			cur.Timestamp = time.Now().UTC().UnixNano()
			err = prof.Reset()
			if err != nil {
				errs <- joe.Error{Type: "disk", Op: "usage ticker", Err: err}
				continue tick
			}
			cur.Devices = cur.Devices[:0]
			// read each line until eof
			for {
				prof.Val = prof.Val[:0]
				prof.Line, err = prof.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errs <- fmt.Errorf("/proc/diskstats: read output bytes: %s", err)
					break
				}
				pos = 0
				fieldNum = 0
				// process the fields in the line
				for {
					// ignore spaces on the first two fields
					if fieldNum < 2 {
						for i, v = range prof.Line[pos:] {
							if v != 0x20 {
								break
							}
						}
						pos += i
					}
					fieldNum++
					for i, v = range prof.Line[pos:] {
						if v == 0x20 || v == '\n' {
							break
						}
					}
					if fieldNum != 3 {
						n, err = helpers.ParseUint(prof.Line[pos : pos+i])
						if err != nil {
							errs <- joe.Error{Type: "cpu stat", Op: "convert cpu data", Err: err}
							continue
						}
					}
					priorPos, pos = pos, pos+i+1
					if fieldNum < 8 {
						if fieldNum < 4 {
							if fieldNum < 2 {
								if fieldNum == 1 {
									dev.Major = uint32(n)
									continue
								}
								dev.Minor = uint32(n)
								continue
							}
							dev.Name = string(prof.Line[priorPos:pos])
							continue
						}
						if fieldNum < 6 {
							if fieldNum == 4 {
								dev.ReadsCompleted = n
								continue
							}
							dev.ReadsMerged = n
							continue
						}
						if fieldNum == 6 {
							dev.ReadSectors = n
							continue
						}
						dev.ReadingTime = n
						continue
					}
					if fieldNum < 12 {
						if fieldNum < 10 {
							if fieldNum == 8 {
								dev.WritesCompleted = n
								continue
							}
							dev.WritesMerged = n
							continue
						}
						if fieldNum == 10 {
							dev.WrittenSectors = n
							continue
						}
						dev.WritingTime = n
						continue
					}
					if fieldNum == 12 {
						dev.IOInProgress = int32(n)
						continue
					}
					if fieldNum == 13 {
						dev.IOTime = n
						continue
					}
					dev.WeightedIOTime = n
					break
				}
				cur.Devices = append(cur.Devices, dev)
			}
			out <- prof.CalculateUsage(cur)
		}
	}
}

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan *structs.Usage, done chan struct{}, errs chan error) {
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
func (prof *Profiler) CalculateUsage(cur *structs.Stats) *structs.Usage {
	u := &structs.Usage{Timestamp: cur.Timestamp, Devices: make([]structs.Device, len(cur.Devices))}
	u.TimeDelta = cur.Timestamp - prof.prior.Timestamp
	for i := 0; i < len(cur.Devices); i++ {
		u.Devices[i].Major = cur.Devices[i].Major
		u.Devices[i].Minor = cur.Devices[i].Minor
		u.Devices[i].Name = cur.Devices[i].Name
		u.Devices[i].ReadsCompleted = cur.Devices[i].ReadsCompleted - prof.prior.Devices[i].ReadsCompleted
		u.Devices[i].ReadsMerged = cur.Devices[i].ReadsMerged - prof.prior.Devices[i].ReadsMerged
		u.Devices[i].ReadSectors = cur.Devices[i].ReadSectors - prof.prior.Devices[i].ReadSectors
		u.Devices[i].ReadingTime = cur.Devices[i].ReadingTime - prof.prior.Devices[i].ReadingTime
		u.Devices[i].WritesCompleted = cur.Devices[i].WritesCompleted - prof.prior.Devices[i].WritesCompleted
		u.Devices[i].WritesMerged = cur.Devices[i].WritesMerged - prof.prior.Devices[i].WritesMerged
		u.Devices[i].WrittenSectors = cur.Devices[i].WrittenSectors - prof.prior.Devices[i].WrittenSectors
		u.Devices[i].WritingTime = cur.Devices[i].WritingTime - prof.prior.Devices[i].WritingTime
		u.Devices[i].IOInProgress = cur.Devices[i].IOInProgress - prof.prior.Devices[i].IOInProgress
		u.Devices[i].IOTime = cur.Devices[i].IOTime - prof.prior.Devices[i].IOTime
		u.Devices[i].WeightedIOTime = cur.Devices[i].WeightedIOTime - prof.prior.Devices[i].WeightedIOTime
	}
	return u
}
