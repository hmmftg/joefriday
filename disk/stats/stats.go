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

// Package stat handles processing of the /proc/stats file: information about
// kernel activity.
package stats

import (
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/disk/structs"
)

const procFile = "/proc/diskstats"

// Profiler is used to process the /proc/stats file.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func New() (prof *Profiler, err error) {
	proc, err := joe.New(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// Get returns information about current kernel activity.
func (prof *Profiler) Get() (stats *structs.Stats, err error) {
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	var (
		i, priorPos, pos, fieldNum int
		n                          uint64
		v                          byte
		dev                        structs.Device
	)

	stats = &structs.Stats{Timestamp: time.Now().UTC().UnixNano(), Devices: make([]structs.Device, 0, 2)}

	// read each line until eof
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, joe.Error{Type: "cpu stat", Op: "reading /proc/stat output", Err: err}
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
					return stats, joe.Error{Type: "cpu stat", Op: "convert cpu data", Err: err}
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
		stats.Devices = append(stats.Devices, dev)
	}
	return stats, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the current kernal activity information using the package's
// global Profiler.
func Get() (stat *structs.Stats, err error) {
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

// Ticker processes CPU utilization information on a ticker.  The generated
// utilization data is sent to the outCh.  Any errors encountered are sent
// to the errCh.  Processing ends when either a done signal is received or
// the done channel is closed.
//
// It is the callers responsibility to close the done and errs channels.
//
// TODO: better handle errors, e.g. restore cur from prior so that there
// isn't the possibility of temporarily having bad data, just a missed
// collection interval.
func (prof *Profiler) Ticker(interval time.Duration, out chan *structs.Stats, done chan struct{}, errs chan error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(out)

	// read each line until eof
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			s, err := prof.Get()
			if err != nil {
				errs <- err
				continue
			}
			out <- s
		}
	}
}

// Ticker gathers information on a ticker using the specified interval.
// This uses a local Profiler as using the global doesn't make sense for
// an ongoing ticker.
func Ticker(interval time.Duration, out chan *structs.Stats, done chan struct{}, errs chan error) {
	prof, err := New()
	if err != nil {
		errs <- err
		close(out)
		return
	}
	prof.Ticker(interval, out, done, errs)
}
