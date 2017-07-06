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

// Package cpuutil handles processing of CPU utilization information. This
// information is calculated using the differences of two CPU stats snapshots,
// /proc/stat and represented as a percentage.
package cpuutil

import (
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	stats "github.com/mohae/joefriday/cpu/cpustats"
)

// Utilization holds information about cpu utilization.
type Utilization struct {
	Timestamp int64 `json:"timestamp"`
	// the time since the prior snapshot; the window that the utilization covers.
	TimeDelta int64 `json:"time_delta"`
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
	Usage  float32 `json:"usage"`
	User   float32 `json:"user"`
	Nice   float32 `json:"nice"`
	System float32 `json:"system"`
	Idle   float32 `json:"idle"`
	IOWait float32 `json:"io_wait"`
}

// Profiler is used to process the /proc/stats file and calculate Utilization
// information.
type Profiler struct {
	*stats.Profiler
	prior stats.Stats
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	p, err := stats.NewProfiler()
	if err != nil {
		return nil, err
	}
	s, err := p.Get()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: p, prior: *s}, nil
}

// Get returns the cpu utilization.  Utilization calculations requires two
// pieces of data.  This func gets a snapshot of /proc/stat, sleeps for a
// second, takes another snapshot and calcualtes the utilization from the
// two snapshots.  If ongoing utilitzation information is desired, the
// Ticker should be used; it's better suited for ongoing utilization
// information.
func (prof *Profiler) Get() (u *Utilization, err error) {
	stat, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	u = prof.calculateUtilization(stat)
	prof.prior = *stat
	return u, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current cpu utilization using the package's global
// Profiler.  The Profiler instantiated lazily; if it doesn't already exist,
// the first Utilization received may be inaccurate.
func Get() (*Utilization, error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		var err error
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Get()
}

// utilizaton =
//   ()(Δuser + Δnice + Δsystem)/(Δuser+Δnice+Δsystem+Δidle)) * CLK_TCK
func (prof *Profiler) calculateUtilization(cur *stats.Stats) *Utilization {
	u := &Utilization{
		Timestamp:  cur.Timestamp,
		TimeDelta:  cur.Timestamp - prof.prior.Timestamp,
		BTimeDelta: int32(cur.Timestamp/1000000000 - cur.BTime),
		CtxtDelta:  cur.Ctxt - prof.prior.Ctxt,
		Processes:  int32(cur.Processes),
		CPU:        make([]Util, len(cur.CPU)),
	}
	var dUser, dNice, dSys, dIdle, tot float32
	// Rest of the calculations are per core
	for i := 0; i < len(cur.CPU); i++ {
		v := Util{ID: cur.CPU[i].ID}
		dUser = float32(cur.CPU[i].User - prof.prior.CPU[i].User)
		dNice = float32(cur.CPU[i].Nice - prof.prior.CPU[i].Nice)
		dSys = float32(cur.CPU[i].System - prof.prior.CPU[i].System)
		dIdle = float32(cur.CPU[i].Idle - prof.prior.CPU[i].Idle)
		tot = dUser + dNice + dSys + dIdle
		v.Usage = (dUser + dNice + dSys) / tot * float32(cur.ClkTck)
		v.User = dUser / tot * float32(cur.ClkTck)
		v.Nice = dNice / tot * float32(cur.ClkTck)
		v.System = dSys / tot * float32(cur.ClkTck)
		v.Idle = dIdle / tot * float32(cur.ClkTck)
		v.IOWait = float32(cur.CPU[i].IOWait-prof.prior.CPU[i].IOWait) / tot * float32(cur.ClkTck)
		u.CPU[i] = v
	}
	return u
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *Utilization
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *Utilization), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	// predeclare some vars
	var (
		i, j, pos, fieldNum int
		n                   uint64
		v                   byte
		stop                bool
		err                 error
		cur                 stats.Stats
		cpu                stats.CPU
	)
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
			cur.CPU = cur.CPU[:0]
			// read each line until eof
			for {
				t.Line, err = t.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Errs <- &joe.ReadError{Err: err}
					break
				}
				// Get everything up to the first space, this is the key.  Not all keys are processed.
				for i, v = range t.Line {
					if v == 0x20 {
						t.Val = t.Line[:i]
						pos = i + 1
						break
					}
				}
				// skip the intr line
				if t.Val[0] == 'i' {
					continue
				}
				if t.Val[0] == 'c' {
					if t.Val[1] == 'p' { // process CPU
						cpu.ID = string(t.Val[:])
						j = 0
						// skip over any remaining spaces
						for i, v = range t.Line[pos:] {
							if v != 0x20 {
								break
							}
							j++
						}
						fieldNum = 0
						pos, j = j+pos, j+pos
						// space is the field separator
						for i, v = range t.Line[pos:] {
							if v == '\n' {
								stop = true
							}
							if v == 0x20 || stop {
								fieldNum++
								n, err = helpers.ParseUint(t.Line[j : pos+i])
								if err != nil {
									t.Errs <- &joe.ParseError{Info: string(t.Val[:]), Err: err}
									continue
								}
								j = pos + i + 1
								if fieldNum < 6 {
									if fieldNum < 4 {
										if fieldNum == 1 {
											cpu.User = int64(n)
											continue
										}
										if fieldNum == 2 {
											cpu.Nice = int64(n)
											continue
										}
										cpu.System = int64(n) // 3
										continue
									}
									if fieldNum == 4 {
										cpu.Idle = int64(n)
										continue
									}
									cpu.IOWait = int64(n) // 5
									continue
								}
								if fieldNum < 8 {
									if fieldNum == 6 {
										cpu.IRQ = int64(n)
										continue
									}
									cpu.SoftIRQ = int64(n) // 7
									continue
								}
								if fieldNum == 8 {
									cpu.Steal = int64(n)
									continue
								}
								if fieldNum == 9 {
									cpu.Quest = int64(n)
									continue
								}
								cpu.QuestNice = int64(n) // 10
							}
						}
						cur.CPU = append(cur.CPU, cpu)
						stop = false
						continue
					}
					// Otherwise it's ctxt info; rest of the line is the data.
					n, err = helpers.ParseUint(t.Line[pos : len(t.Line)-1])
					if err != nil {
						t.Errs <- &joe.ParseError{Info: string(t.Val[:]), Err: err}
						continue
					}
					cur.Ctxt = int64(n)
					continue
				}
				if t.Val[0] == 'b' {
					// rest of the line is the data
					n, err = helpers.ParseUint(t.Line[pos : len(t.Line)-1])
					if err != nil {
						t.Errs <- &joe.ParseError{Info: string(t.Val[:]), Err: err}
						continue
					}
					cur.BTime = int64(n)
					continue
				}
				if t.Val[0] == 'p' && t.Val[4] == 'e' { // processes info
					// rest of the line is the data
					n, err = helpers.ParseUint(t.Line[pos : len(t.Line)-1])
					if err != nil {
						t.Errs <- &joe.ParseError{Info: string(t.Val[:]), Err: err}
						continue
					}
					cur.Processes = int64(n)
					continue
				}
			}
			t.Data <- t.Profiler.calculateUtilization(&cur)
			t.Profiler.prior.Ctxt = cur.Ctxt
			t.Profiler.prior.BTime = cur.BTime
			t.Profiler.prior.Processes = cur.Processes
			if len(t.Profiler.prior.CPU) != len(cur.CPU) {
				t.Profiler.prior.CPU = make([]stats.CPU, len(cur.CPU))
			}
			copy(t.Profiler.prior.CPU, cur.CPU)
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
