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

// Package LoadAvg processes loadavg information from the /proc/loadavg file.
package loadavg

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
)

const procFile = "/proc/loadavg"

// LoadAvg holds loadavg information
type LoadAvg struct {
	LastMinute       float32
	LastFive         float32
	LastTen          float32
	RunningProcesses int32
	TotalProcesses   int32
	PID              int32
}

// Profiler processes the loadavg information.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.New(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// Get populates LoadAvg with /proc/loadavg information.
func (prof *Profiler) Get() (l LoadAvg, err error) {
	err = prof.Reset()
	if err != nil {
		return l, err
	}
	var (
		i, priorPos, pos, line, fieldNum int
		n                                uint64
		f                                float64
		v                                byte
	)

	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return l, &joe.ReadError{Err: err}
		}
		line++
		for {
			// space delimits the values
			for i, v = range prof.Line[pos:] {
				if v == 0x20 {
					priorPos, pos = pos, pos+i+1
					break
				}
			}
			fieldNum++
			if fieldNum <= 3 {
				f, err = strconv.ParseFloat(string(prof.Line[priorPos:pos-1]), 64)
				if err != nil {
					return l, &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
				}
				if fieldNum == 1 {
					l.LastMinute = float32(f)
					continue
				}
				if fieldNum == 2 {
					l.LastFive = float32(f)
					continue
				}
				l.LastTen = float32(f)
				continue
			}
			if fieldNum == 4 {
				// get the process information: separated by /
				for i, v = range prof.Line[priorPos:pos] {
					if v == '/' {
						break
					}
				}
				n, err = helpers.ParseUint(prof.Line[priorPos : priorPos+i])
				if err != nil {
					return l, &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
				}
				l.RunningProcesses = int32(n)
				n, err = helpers.ParseUint(prof.Line[priorPos+i+1 : pos-1])
				if err != nil {
					return l, &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
				}
				l.TotalProcesses = int32(n)
				continue
			}
			n, err = helpers.ParseUint(prof.Line[pos : len(prof.Line)-1])
			if err != nil {
				return l, err
			}
			l.PID = int32(n)
			break
		}
	}
	return l, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get gets the loadavg information using the package's global Profiler, which
// is lazily instantiated.
func Get() (l LoadAvg, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return l, err
		}
	}
	return std.Get()
}

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan LoadAvg
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
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan LoadAvg), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	var (
		i, priorPos, pos, line, fieldNum int
		n                                uint64
		f                                float64
		v                                byte
		l                                LoadAvg
		err                              error
	)
	// ticker
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			//			l.Timestamp = time.Now().UTC().UnixNano()
			err = t.Reset()
			line = 0
		tick:
			for {
				t.Line, err = t.Buf.ReadSlice('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Errs <- &joe.ReadError{Err: err}
					break
				}
				line++
				for {
					// space delimits the values
					for i, v = range t.Line[pos:] {
						if v == 0x20 {
							priorPos, pos = pos, pos+i+1
							break
						}
					}
					fieldNum++
					if fieldNum <= 3 {
						f, err = strconv.ParseFloat(string(t.Line[priorPos:pos-1]), 64)
						if err != nil {
							t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
							break tick
						}
						if fieldNum == 1 {
							l.LastMinute = float32(f)
							continue
						}
						if fieldNum == 2 {
							l.LastFive = float32(f)
							continue
						}
						l.LastTen = float32(f)
						continue
					}
					if fieldNum == 4 {
						// get the process information: separated by /
						for i, v = range t.Line[priorPos:pos] {
							if v == '/' {
								break
							}
						}
						n, err = helpers.ParseUint(t.Line[priorPos : priorPos+i])
						if err != nil {
							t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
							break tick
						}
						l.RunningProcesses = int32(n)
						n, err = helpers.ParseUint(t.Line[priorPos+i+1 : pos-1])
						if err != nil {
							t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
							break tick
						}
						l.TotalProcesses = int32(n)
						continue
					}
					n, err = helpers.ParseUint(t.Line[pos : len(t.Line)-1])
					if err != nil {
						t.Errs <- &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
						break tick
					}
					l.PID = int32(n)
					break
				}
			}
			t.Data <- l
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
