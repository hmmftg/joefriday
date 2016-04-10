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

// Package utilization handles processing of CPU utilization information.
// This information is calculated using the differences of two CPU stats
// snapshots and represented as a percentage.
package utilization

import (
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/cpu/stats"
)

// Profiler is used to process the /proc/stats file and calculate Utilization
// information.
type Profiler struct {
	*stats.Profiler
	prior *stats.Stats
}

// Returns an initialized Profiler; ready to use.
func New() (prof *Profiler, err error) {
	p, err := stats.New()
	if err != nil {
		return nil, err
	}
	s, err := p.Get()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: s}, nil
}

// Get returns the cpu utilization.  Utilization calculations requires two
// pieces of data.  This func gets a snapshot of /proc/stat, sleeps for a
// second, takes another snapshot and calcualtes the utilization from the
// two snapshots.  If ongoing utilitzation information is desired, the
// Ticker should be used; it's better suited for ongoing utilization
// information.
//
// TODO: should this be changed so that this calculates utilization since
// last time the stats were obtained.  If there aren't pre-existing stats
// it would get current utilization (which may be a separate method (or
// should be?)).  Also: rethink locking.
func (prof *Profiler) Get() (u *Utilization, err error) {
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	prof.prior, err = prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	stat2, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.calculateUtilization(stat2), nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current cpu utilization using the package's global Profiler.
func Get() (*Utilization, error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		var err error
		std, err = New()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// Ticker processes CPU utilization information on a ticker.  The generated
// utilization data is sent to the out channel.  Any errors encountered are
// sent to the errs.  Processing ends when a done signal is received.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *Profiler) Ticker(interval time.Duration, out chan *Utilization, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)
	// predeclare some vars
	var (
		i, j, pos, fieldNum int
		n                   uint64
		v                   byte
		stop                bool
	)
	// first get stats as the baseline
	cur, err := prof.Profiler.Get()
	if err != nil {
		errs <- err
	}
	// ticker
tick:
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			prof.prior.Ctxt = cur.Ctxt
			prof.prior.BTime = cur.BTime
			prof.prior.Processes = cur.Processes
			if len(prof.prior.CPU) != len(cur.CPU) {
				prof.prior.CPU = make([]stats.Stat, len(cur.CPU))
			}
			copy(prof.prior.CPU, cur.CPU)
			cur.Timestamp = time.Now().UTC().UnixNano()
			err = prof.Reset()
			if err != nil {
				errs <- joe.Error{Type: "cpu", Op: "utilization ticker: seek /proc/stat", Err: err}
				continue tick
			}
			cur.CPU = cur.CPU[:0]
			// read each line until eof
			for {
				prof.Line, err = prof.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					errs <- joe.Error{Type: "cpu stat", Op: "reading /proc/stat output", Err: err}
					break
				}
				// Get everything up to the first space, this is the key.  Not all keys are processed.
				for i, v = range prof.Line {
					if v == 0x20 {
						prof.Val = prof.Line[:i]
						pos = i + 1
						break
					}
				}
				// skip the intr line
				if prof.Val[0] == 'i' {
					continue
				}
				if prof.Val[0] == 'c' {
					if prof.Val[1] == 'p' { // process CPU
						stat := stats.Stat{ID: string(prof.Val[:])}
						j = 0
						// skip over any remaining spaces
						for i, v = range prof.Line[pos:] {
							if v != 0x20 {
								break
							}
							j++
						}
						fieldNum = 0
						pos, j = j+pos, j+pos
						// space is the field separator
						for i, v = range prof.Line[pos:] {
							if v == '\n' {
								stop = true
							}
							if v == 0x20 || stop {
								fieldNum++
								n, err = helpers.ParseUint(prof.Line[j : pos+i])
								if err != nil {
									errs <- joe.Error{Type: "cpu stat", Op: "convert cpu data", Err: err}
									continue
								}
								j = pos + i + 1
								if fieldNum < 6 {
									if fieldNum < 4 {
										if fieldNum == 1 {
											stat.User = int64(n)
											continue
										}
										if fieldNum == 2 {
											stat.Nice = int64(n)
											continue
										}
										stat.System = int64(n) // 3
										continue
									}
									if fieldNum == 4 {
										stat.Idle = int64(n)
										continue
									}
									stat.IOWait = int64(n) // 5
									continue
								}
								if fieldNum < 8 {
									if fieldNum == 6 {
										stat.IRQ = int64(n)
										continue
									}
									stat.SoftIRQ = int64(n) // 7
									continue
								}
								if fieldNum == 8 {
									stat.Steal = int64(n)
									continue
								}
								if fieldNum == 9 {
									stat.Quest = int64(n)
									continue
								}
								stat.QuestNice = int64(n) // 10
							}
						}
						cur.CPU = append(cur.CPU, stat)
						stop = false
						continue
					}
					// Otherwise it's ctxt info; rest of the line is the data.
					n, err = helpers.ParseUint(prof.Line[pos : len(prof.Line)-1])
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert ctxt data", Err: err}
						continue
					}
					cur.Ctxt = int64(n)
					continue
				}
				if prof.Val[0] == 'b' {
					// rest of the line is the data
					n, err = helpers.ParseUint(prof.Line[pos : len(prof.Line)-1])
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert btime data", Err: err}
						continue
					}
					cur.BTime = int64(n)
					continue
				}
				if prof.Val[0] == 'p' && prof.Val[4] == 'e' { // processes info
					// rest of the line is the data
					n, err = helpers.ParseUint(prof.Line[pos : len(prof.Line)-1])
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert processes data", Err: err}
						continue
					}
					cur.Processes = int64(n)
					continue
				}
			}
			out <- prof.calculateUtilization(cur)
		}
	}
}

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan *Utilization, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}

// utilizaton =
//   ()(Δuser + Δnice + Δsystem)/(Δuser+Δnice+Δsystem+Δidle)) * CLK_TCK
func (prof *Profiler) calculateUtilization(s2 *stats.Stats) *Utilization {
	u := &Utilization{
		Timestamp:  s2.Timestamp,
		BTimeDelta: int32(s2.Timestamp/1000000000 - s2.BTime),
		CtxtDelta:  s2.Ctxt - prof.prior.Ctxt,
		Processes:  int32(s2.Processes),
		CPU:        make([]Util, len(s2.CPU)),
	}
	var dUser, dNice, dSys, dIdle, tot float32
	// Rest of the calculations are per core
	for i := 0; i < len(s2.CPU); i++ {
		v := Util{ID: s2.CPU[i].ID}
		dUser = float32(s2.CPU[i].User - prof.prior.CPU[i].User)
		dNice = float32(s2.CPU[i].Nice - prof.prior.CPU[i].Nice)
		dSys = float32(s2.CPU[i].System - prof.prior.CPU[i].System)
		dIdle = float32(s2.CPU[i].Idle - prof.prior.CPU[i].Idle)
		tot = dUser + dNice + dSys + dIdle
		v.Usage = (dUser + dNice + dSys) / tot * float32(s2.ClkTck)
		v.User = dUser / tot * float32(s2.ClkTck)
		v.Nice = dNice / tot * float32(s2.ClkTck)
		v.System = dSys / tot * float32(s2.ClkTck)
		v.Idle = dIdle / tot * float32(s2.ClkTck)
		v.IOWait = float32(s2.CPU[i].IOWait-prof.prior.CPU[i].IOWait) / tot * float32(s2.ClkTck)
		u.CPU[i] = v
	}
	return u
}

// Utilization holds information about cpu utilization.
type Utilization struct {
	Timestamp int64 `json:"timestamp"`
	// time since last reboot, in seconds
	BTimeDelta int32 `json:"btime_delta"`
	// context switches since last snapshot
	CtxtDelta int64 `json:"ctxt_delta"`
	// current number of Processes
	Processes int32 `json:"processes"`
	// cpu specific utilization information
	CPU []Util `json:"cpu"`
}

// Util holds utilization information, as percentages, for a CPU.
type Util struct {
	ID     string  `json:"id"`
	Usage  float32 `json:"total"`
	User   float32 `json:"user"`
	Nice   float32 `json:"nice"`
	System float32 `json:"system"`
	Idle   float32 `json:"idle"`
	IOWait float32 `json:"io_wait"`
}
