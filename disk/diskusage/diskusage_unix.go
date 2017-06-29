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

// Package diskusage calculates block device usage. Usage is calculated by
// taking the difference of two /proc/diskstats snapshots.
package diskusage

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	stats "github.com/mohae/joefriday/disk/diskstats"
	"github.com/mohae/joefriday/disk/structs"
)

// Profiler is used to process the disk usage..
type Profiler struct {
	*stats.Profiler
	prior *structs.DiskStats
}

// Returns an initialized Profiler; ready to use.  The prior stats is set to
// the current stats snapshot.
func NewProfiler() (prof *Profiler, err error) {
	p, err := stats.NewProfiler()
	if err != nil {
		return nil, err
	}
	s, err := p.Get()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: s}, nil
}

// Get returns the current disk usage.  Usage is calculated as the difference
// between the prior stats snapshot and the current one.
func (prof *Profiler) Get() (u *structs.DiskUsage, err error) {
	st, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	u = prof.CalculateUsage(st)
	prof.prior = st
	return u, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current disk usage using the package's global Profiler.
// The profiler is lazily instantiated.  This means that probability of
// the first utilization snapshot returning inaccurate information is high
// due to the lack of time elapsing between the initial and current
// snapshot for utilization calculation.
func Get() (u *structs.DiskUsage, err error) {
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

// CalculateUsage returns the difference between the current /proc/net/dev
// data and the prior one.
func (prof *Profiler) CalculateUsage(cur *structs.DiskStats) *structs.DiskUsage {
	u := &structs.DiskUsage{Timestamp: cur.Timestamp, Device: make([]structs.Device, len(cur.Device))}
	u.TimeDelta = cur.Timestamp - prof.prior.Timestamp
	for i := 0; i < len(cur.Device); i++ {
		u.Device[i].Major = cur.Device[i].Major
		u.Device[i].Minor = cur.Device[i].Minor
		u.Device[i].Name = cur.Device[i].Name
		u.Device[i].ReadsCompleted = cur.Device[i].ReadsCompleted - prof.prior.Device[i].ReadsCompleted
		u.Device[i].ReadsMerged = cur.Device[i].ReadsMerged - prof.prior.Device[i].ReadsMerged
		u.Device[i].ReadSectors = cur.Device[i].ReadSectors - prof.prior.Device[i].ReadSectors
		u.Device[i].ReadingTime = cur.Device[i].ReadingTime - prof.prior.Device[i].ReadingTime
		u.Device[i].WritesCompleted = cur.Device[i].WritesCompleted - prof.prior.Device[i].WritesCompleted
		u.Device[i].WritesMerged = cur.Device[i].WritesMerged - prof.prior.Device[i].WritesMerged
		u.Device[i].WrittenSectors = cur.Device[i].WrittenSectors - prof.prior.Device[i].WrittenSectors
		u.Device[i].WritingTime = cur.Device[i].WritingTime - prof.prior.Device[i].WritingTime
		u.Device[i].IOInProgress = cur.Device[i].IOInProgress - prof.prior.Device[i].IOInProgress
		u.Device[i].IOTime = cur.Device[i].IOTime - prof.prior.Device[i].IOTime
		u.Device[i].WeightedIOTime = cur.Device[i].WeightedIOTime - prof.prior.Device[i].WeightedIOTime
	}
	return u
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *structs.DiskUsage
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *structs.DiskUsage), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	var (
		i, priorPos, pos, line, fieldNum int
		n                                uint64
		v                                byte
		err                              error
		dev                              structs.Device
		cur                              structs.DiskStats
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
			cur.Device = cur.Device[:0]
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
				pos = 0
				fieldNum = 0
				// process the fields in the line
				for {
					// ignore spaces on the first two fields
					if fieldNum < 2 {
						for i, v = range t.Line[pos:] {
							if v != 0x20 {
								break
							}
						}
						pos += i
					}
					fieldNum++
					for i, v = range t.Line[pos:] {
						if v == 0x20 || v == '\n' {
							break
						}
					}
					if fieldNum != 3 {
						n, err = helpers.ParseUint(t.Line[pos : pos+i])
						if err != nil {
							t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
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
							dev.Name = string(t.Line[priorPos:pos])
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
				cur.Device = append(cur.Device, dev)
			}
			t.Data <- t.CalculateUsage(&cur)
			// set prior info
			t.prior.Timestamp = cur.Timestamp
			if len(t.prior.Device) != len(cur.Device) {
				t.prior.Device = make([]structs.Device, len(cur.Device))
			}
			copy(t.prior.Device, cur.Device)
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
