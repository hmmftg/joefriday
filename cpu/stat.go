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

package cpu

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	joe "github.com/mohae/joefriday"
)

const procStat = "/proc/stat"

var CLK_TCK int32    // the ticks per clock cycle
var tckMu sync.Mutex //protects CLK_TCK

// Set CLK_TCK.
func ClkTck() error {
	tckMu.Lock()
	defer tckMu.Unlock()
	var out bytes.Buffer
	cmd := exec.Command("getconf", "CLK_TCK")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return joe.Error{Type: "cpu", Op: "get conf CLK_TCK", Err: err}
	}
	b, err := out.ReadBytes('\n')
	if err != nil {
		return joe.Error{Type: "cpu", Op: "read conf CLK_TCK output", Err: err}
	}
	v, err := strconv.Atoi(string(b[:len(b)-1]))
	if err != nil {
		return joe.Error{Type: "cpu", Op: "processing conf CLK_TCK output", Err: err}
	}
	atomic.StoreInt32(&CLK_TCK, int32(v))
	return nil
}

type StatProfiler struct {
	joe.Proc
	ClkTck int16
}

func NewStatProfiler() (prof *StatProfiler, err error) {
	// if it hasn't been set, set it.
	if atomic.LoadInt32(&CLK_TCK) == 0 {
		err = ClkTck()
		if err != nil {
			return nil, err
		}
	}
	f, err := os.Open(procStat)
	if err != nil {
		return nil, err
	}
	return &StatProfiler{Proc: joe.Proc{File: f, Buf: bufio.NewReader(f)}, ClkTck: int16(atomic.LoadInt32(&CLK_TCK))}, nil
}

// It is expected that the caller as the lock.
func (prof *StatProfiler) reset() error {
	_, err := prof.File.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}
	prof.Buf.Reset(prof.File)
	return nil
}

func (prof *StatProfiler) Get() (stats *Stats, err error) {
	prof.Reset()
	prof.Lock()
	defer prof.Unlock()
	if CLK_TCK == 0 {
		err := ClkTck()
		if err != nil {
			return stats, err
		}
	}
	var (
		name                     string
		i, j, pos, val, fieldNum int
		v                        byte
		stop                     bool
	)

	stats = &Stats{Timestamp: time.Now().UTC().UnixNano(), ClkTck: prof.ClkTck, CPU: make([]Stat, 0, 2)}

	// read each line until eof
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, joe.Error{Type: "cpu stat", Op: "reading /proc/stat output", Err: err}
		}
		// Get everything up to the first space, this is the key.  Not all keys are processed.
		for i, v = range prof.Line {
			if v == 0x20 {
				name = string(prof.Line[:i])
				pos = i + 1
				break
			}
		}
		// skip the intr line
		if name == "intr" {
			continue
		}
		if name[:3] == "cpu" {
			j = 0
			// skip over any remaining spaces
			for i, v = range prof.Line[pos:] {
				if v != 0x20 {
					break
				}
				j++
			}
			stat := Stat{ID: name}
			fieldNum = 0
			pos, j = j+pos, j+pos
			// space is the field separator
			for i, v = range prof.Line[pos:] {
				if v == '\n' {
					stop = true
				}
				if v == 0x20 || stop {
					fieldNum++
					val, err = strconv.Atoi(string(prof.Line[j : pos+i]))
					if err != nil {
						return stats, joe.Error{Type: "cpu stat", Op: "convert cpu data", Err: err}
					}
					j = pos + i + 1
					if fieldNum < 4 {
						if fieldNum == 1 {
							stat.User = int64(val)
						} else if fieldNum == 2 {
							stat.Nice = int64(val)
						} else if fieldNum == 3 {
							stat.System = int64(val)
						}
					} else if fieldNum < 7 {
						if fieldNum == 4 {
							stat.Idle = int64(val)
						} else if fieldNum == 5 {
							stat.IOWait = int64(val)
						} else if fieldNum == 6 {
							stat.IRQ = int64(val)
						}
					} else if fieldNum < 10 {
						if fieldNum == 7 {
							stat.SoftIRQ = int64(val)
						} else if fieldNum == 8 {
							stat.Steal = int64(val)
						} else if fieldNum == 9 {
							stat.Quest = int64(val)
						}
					} else if fieldNum == 10 {
						stat.QuestNice = int64(val)
					}
				}
			}
			stats.CPU = append(stats.CPU, stat)
			stop = false
			continue
		}
		if name == "ctxt" {
			// rest of the line is the data
			val, err = strconv.Atoi(string(prof.Line[pos : len(prof.Line)-1]))
			if err != nil {
				return stats, joe.Error{Type: "cpu stat", Op: "convert ctxt data", Err: err}
			}
			stats.Ctxt = int64(val)
			continue
		}
		if name == "btime" {
			// rest of the line is the data
			val, err = strconv.Atoi(string(prof.Line[pos : len(prof.Line)-1]))
			if err != nil {
				return stats, joe.Error{Type: "cpu stat", Op: "convert btime data", Err: err}
			}
			stats.BTime = int64(val)
			continue
		}
		if name == "processes" {
			// rest of the line is the data
			val, err = strconv.Atoi(string(prof.Line[pos : len(prof.Line)-1]))
			if err != nil {
				return stats, joe.Error{Type: "cpu stat", Op: "convert processes data", Err: err}
			}
			stats.Processes = int64(val)
			continue
		}
	}
	return stats, nil
}

var stdStat *StatProfiler
var stdStatMu sync.Mutex

// GetStats gets the output of /proc/stat.
func GetStat() (stat *Stats, err error) {
	stdStatMu.Lock()
	defer stdStatMu.Unlock()
	if stdStat == nil {
		stdStat, err = NewStatProfiler()
		if err != nil {
			return nil, err
		}
	}
	return stdStat.Get()
}

// GetUtilization returns the cpu utilization.  Utilization calculations
// requires two pieces of data.  This func gets a snapshot of /proc/stat,
// sleeps for a second, takes another snapshot and calcualtes the utilization
// from the two snapshots.  If ongoing utilitzation information is desired,
// the UtilizationTicker should be used; it's better suited for ongoing
// utilization information being; using less cpu cycles and generating less
// garbage.
// TODO: should this be changed so that this calculates utilization since
// last time the stats were obtained.  If there aren't pre-existing stats
// it would get current utilization (which may be a separate method (or
// should be?))
func (prof *StatProfiler) GetUtilization() (*Utilization, error) {
	stat1, err := prof.Get()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	stat2, err := prof.Get()
	if err != nil {
		return nil, err
	}
	return calculateUtilization(stat1, stat2), nil
}

func GetUtilization() (*Utilization, error) {
	stdStatMu.Lock()
	defer stdStatMu.Unlock()
	if stdStat == nil {
		var err error
		stdStat, err = NewStatProfiler()
		if err != nil {
			return nil, err
		}
	}
	return stdStat.GetUtilization()
}

// UtilizationTicker processes CPU utilization information on a ticker.  The
// generated utilization data is sent to the outCh.  Any errors encountered
// are sent to the errCh.  Processing ends when either a done signal is
// received or the done channel is closed.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *StatProfiler) UtilizationTicker(interval time.Duration, outCh chan *Utilization, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(outCh)
	// predeclare some vars
	var (
		i, j, pos, val, fieldNum int
		v                        byte
		name                     string
		stop                     bool
		prior                    Stats
	)
	// first get stats as the baseline
	cur, err := prof.Get()
	if err != nil {
		errs <- err
	}
	// Lock to prevent deadlock
	prof.Lock()
	// ticker
tick:
	for {
		prof.Unlock()
		select {
		case <-done:
			return
		case <-ticker.C:
			err = prof.Reset()
			prof.Lock()
			if err != nil {
				errs <- joe.Error{Type: "cpu", Op: "utilization ticker: seek /proc/stat", Err: err}
				continue tick
			}
			prior.Ctxt = cur.Ctxt
			prior.BTime = cur.BTime
			prior.Processes = cur.Processes
			if len(prior.CPU) != len(cur.CPU) {
				prior.CPU = make([]Stat, len(cur.CPU))
			}
			copy(prior.CPU, cur.CPU)
			cur.Timestamp = time.Now().UTC().UnixNano()
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
						name = string(prof.Line[:i])
						pos = i + 1
						break
					}
				}
				// skip the intr line
				if name == "intr" {
					continue
				}
				if name[:3] == "cpu" {
					j = 0
					// skip over any remaining spaces
					for i, v = range prof.Line[pos:] {
						if v != 0x20 {
							break
						}
						j++
					}
					stat := Stat{ID: name}
					fieldNum = 0
					pos, j = j+pos, j+pos
					// space is the field separator
					for i, v = range prof.Line[pos:] {
						if v == '\n' {
							stop = true
						}
						if v == 0x20 || stop {
							fieldNum++
							val, err = strconv.Atoi(string(prof.Line[j : pos+i]))
							if err != nil {
								errs <- joe.Error{Type: "cpu stat", Op: "convert cpu data", Err: err}
								continue
							}
							j = pos + i + 1
							if fieldNum < 4 {
								if fieldNum == 1 {
									stat.User = int64(val)
								} else if fieldNum == 2 {
									stat.Nice = int64(val)
								} else if fieldNum == 3 {
									stat.System = int64(val)
								}
							} else if fieldNum < 7 {
								if fieldNum == 4 {
									stat.Idle = int64(val)
								} else if fieldNum == 5 {
									stat.IOWait = int64(val)
								} else if fieldNum == 6 {
									stat.IRQ = int64(val)
								}
							} else if fieldNum < 10 {
								if fieldNum == 7 {
									stat.SoftIRQ = int64(val)
								} else if fieldNum == 8 {
									stat.Steal = int64(val)
								} else if fieldNum == 9 {
									stat.Quest = int64(val)
								}
							} else if fieldNum == 10 {
								stat.QuestNice = int64(val)
							}
						}
					}
					cur.CPU = append(cur.CPU, stat)
					stop = false
					continue
				}
				if name == "ctxt" {
					// rest of the line is the data
					val, err = strconv.Atoi(string(prof.Line[pos : len(prof.Line)-1]))
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert ctxt data", Err: err}
					}
					cur.Ctxt = int64(val)
					continue
				}
				if name == "btime" {
					// rest of the line is the data
					val, err = strconv.Atoi(string(prof.Line[pos : len(prof.Line)-1]))
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert btime data", Err: err}
					}
					cur.BTime = int64(val)
					continue
				}
				if name == "processes" {
					// rest of the line is the data
					val, err = strconv.Atoi(string(prof.Line[pos : len(prof.Line)-1]))
					if err != nil {
						errs <- joe.Error{Type: "cpu stat", Op: "convert processes data", Err: err}
					}
					cur.Processes = int64(val)
					continue
				}
			}
			outCh <- calculateUtilization(&prior, cur)
		}
	}
}

func UtilizationTicker(interval time.Duration, out chan *Utilization, done chan struct{}, errs chan error) {
	p, err := NewStatProfiler()
	if err != nil {
		errs <- err
		return
	}
	p.UtilizationTicker(interval, out, done, errs)
}

// Stats holds the /proc/stat information
type Stats struct {
	ClkTck    int16  `json:"clk_tck"`
	Timestamp int64  `json:"timestamp"`
	Ctxt      int64  `json:"ctxt"`
	BTime     int64  `json:"btime"`
	Processes int64  `json:"processes"`
	CPU       []Stat `json:"cpu"`
}

// Stat is for capturing the CPU information of /proc/stat.
type Stat struct {
	ID        string `json:"ID"`
	User      int64  `json:"user"`
	Nice      int64  `json:"nice"`
	System    int64  `json:"system"`
	Idle      int64  `json:"idle"`
	IOWait    int64  `json:"io_wait"`
	IRQ       int64  `json:"irq"`
	SoftIRQ   int64  `json:"soft_irq"`
	Steal     int64  `json:"steal"`
	Quest     int64  `json:"quest"`
	QuestNice int64  `json:"quest_nice"`
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

// Util holds utilization information for a CPU.
type Util struct {
	ID        string  `json:"id"`
	Usage     float32 `json:"total"`
	User      float32 `json:"user"`
	Nice      float32 `json:"nice"`
	System    float32 `json:"system"`
	Idle      float32 `json:"idle"`
	IOWait    float32 `json:"io_wait"`
	IRQ       float32 `json:"irq"`
	SoftIRQ   float32 `json:"soft_irq"`
	Steal     float32 `json:"steal"`
	Quest     float32 `json:"quest"`
	QuestNice float32 `json:"quest_nice"`
}

// usage = ()(Δuser + Δnice + Δsystem)/(Δuser+Δnice+Δsystem+Δidle)) * CLK_TCK
func calculateUtilization(s1, s2 *Stats) *Utilization {
	u := &Utilization{
		Timestamp:  s2.Timestamp,
		BTimeDelta: int32(s2.Timestamp/1000000000 - s2.BTime),
		CtxtDelta:  s2.Ctxt - s1.Ctxt,
		Processes:  int32(s2.Processes),
		CPU:        make([]Util, len(s2.CPU)),
	}
	var dUser, dNice, dSys, dIdle, tot float32
	// Rest of the calculations are per core
	for i := 0; i < len(s2.CPU); i++ {
		v := Util{ID: s2.CPU[i].ID}
		dUser = float32(s2.CPU[i].User - s1.CPU[i].User)
		dNice = float32(s2.CPU[i].Nice - s1.CPU[i].Nice)
		dSys = float32(s2.CPU[i].System - s1.CPU[i].System)
		dIdle = float32(s2.CPU[i].Idle - s1.CPU[i].Idle)
		tot = dUser + dNice + dSys + dIdle
		v.Usage = (dUser + dNice + dSys) / tot * float32(s2.ClkTck)
		v.User = dUser / tot * float32(s2.ClkTck)
		v.Nice = dNice / tot * float32(s2.ClkTck)
		v.System = dSys / tot * float32(s2.ClkTck)
		v.Idle = dIdle / tot * float32(s2.ClkTck)
		v.IOWait = float32(s2.CPU[i].IOWait-s1.CPU[i].IOWait) / tot * float32(s2.ClkTck)
		u.CPU[i] = v
	}
	return u
}
