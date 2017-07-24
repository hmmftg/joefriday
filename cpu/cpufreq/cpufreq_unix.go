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

// Package cpufreq provides the current CPU frequency, in MHz, as reported by
// /proc/cpuinfo.
package cpufreq

import (
	"io"
	"strconv"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
)

const procFile = "/proc/cpuinfo"

// Frequency holds information about the frequency of a system's cpus, in MHz.
// The reported values are the current speeds as reported by /proc/cpuinfo.
type Frequency struct {
	Timestamp int64
	CPU       []CPU `json:"cpu"`
}

// CPU holds the clock info for a single processor.
type CPU struct {
	Processor       int16    `json:"processor"`
	CPUMHz          float32  `json:"cpu_mhz"`
	PhysicalID      int16    `json:"physical_id"`
	CoreID          int16    `json:"core_id"`
	APICID          int16    `json:"apicid"`
}

// Profiler is used to process the frequency information.
type Profiler struct {
	joe.Procer
	*joe.Buffer
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.NewProc(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Procer: proc, Buffer: joe.NewBuffer()}, nil
}

// Reset resources; after reset the profiler is ready to be used again.
func (prof *Profiler) Reset() error {
	prof.Buffer.Reset()
	return prof.Procer.Reset()
}

// Get returns the current Frequency.
func (prof *Profiler) Get() (f *Frequency, err error) {
	var (
		cpuCnt, i, pos, nameLen int
		n                       uint64
		v                       byte
		cpu                     CPU
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	f = &Frequency{Timestamp: time.Now().UTC().UnixNano()}
	for {
		prof.Line, err = prof.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, &joe.ReadError{Err: err}
		}
		prof.Val = prof.Val[:0]
		// First grab the attribute name; everything up to the ':'.  The key may have
		// spaces and has trailing spaces; that gets trimmed.
		for i, v = range prof.Line {
			if v == 0x3A {
				prof.Val = prof.Line[:i]
				pos = i + 1
				break
			}
			//prof.Val = append(prof.Val, v)
		}
		prof.Val = joe.TrimTrailingSpaces(prof.Val[:])
		nameLen = len(prof.Val)
		// if there's no name; skip.
		if nameLen == 0 {
			continue
		}
		// if there's anything left, the value is everything else; trim spaces
		if pos+1 < len(prof.Line) {
			prof.Val = append(prof.Val, joe.TrimTrailingSpaces(prof.Line[pos+1:])...)
		}
		if prof.Val[0] == 'a' {
			if prof.Val[1] == 'p' { // apicid
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.APICID = int16(n)
			}
			continue
		}
		if prof.Val[0] == 'c' {
			if prof.Val[1] == 'p' {
				if prof.Val[4] == 'M' { // cpu MHz
					f, err := strconv.ParseFloat(string(prof.Val[nameLen:]), 32)
					if err != nil {
						return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
					}
					cpu.CPUMHz = float32(f)
				}
				continue
			}
			if prof.Val[1] == 'o' { // core id
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.CoreID = int16(n)
			}
			continue
		}
		if prof.Val[0] == 'p' {
			if prof.Val[1] == 'h' { // physical id
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu.PhysicalID = int16(n)
				continue
			}
			// processor starts information about a processor.
			if prof.Val[1] == 'r' { // processor
				if cpuCnt > 0 {
					f.CPU = append(f.CPU, cpu)
				}
				cpuCnt++
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				cpu = CPU{Processor: int16(n)}
			}
		}
		continue
	}
	// append the current processor informatin
	f.CPU = append(f.CPU, cpu)
	return f, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns Frequency using the package's global Profiler.
func Get() (f *Frequency, err error) {
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

// Ticker delivers the CPU Frequencies at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *Frequency
	*Profiler
}

// NewTicker returns a new Ticker containing a Data channel that delivers the
// data at intervals and an error channel that delivers any errors encountered.
// Stop the ticker to signal the ticker to stop running. Stopping the ticker
// does not close the Data channel; call Close to close both the ticker and the
// data channel.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	p, err := NewProfiler()
	if err != nil {
		return nil, err
	}
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *Frequency), Profiler: p}
	go t.Run()
	return &t, nil
}

// Run runs the ticker.
func (t *Ticker) Run() {
	// ticker
	for {
		select {
		case <-t.Done:
			return
		case <-t.C:
			s, err := t.Get()
			if err != nil {
				t.Errs <- err
				continue
			}
			t.Data <- s
		}
	}
}

// Close closes the ticker resources.
func (t *Ticker) Close() {
	t.Ticker.Close()
	close(t.Data)
}
