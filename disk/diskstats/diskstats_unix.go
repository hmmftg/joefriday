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

// Package diskstats handles processing of information about block devices:
// /proc/diskstats.
package diskstats

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/disk/structs"
)

const procFile = "/proc/diskstats"

// Profiler is used to process the /proc/diskstats file.
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

// Get returns information about current disk activity.
func (prof *Profiler) Get() (stats *structs.DiskStats, err error) {
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	var (
		i, priorPos, pos, line, fieldNum int
		n                                uint64
		v                                byte
		dev                              structs.Device
	)

	stats = &structs.DiskStats{Timestamp: time.Now().UTC().UnixNano(), Devices: make([]structs.Device, 0, 2)}

	// read each line until eof
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, &joe.ReadError{Err: err}
		}
		line++
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
					return stats, &joe.ParseError{Info: fmt.Sprintf("line %d: field %d", line, fieldNum), Err: err}
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

// Get returns the current block device stats using the package's global
// Profiler.
func Get() (stat *structs.DiskStats, err error) {
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

// Ticker delivers the system's memory information at intervals.
type Ticker struct {
	*joe.Ticker
	Data chan *structs.DiskStats
	*Profiler
}

// NewTicker returns a new Ticker contianing a Data channel that delivers
// the data at intervals and an error channel that delivers any errors
// encountered.  Stop the ticker to signal the ticker to stop running; it
// does not close the Data channel.  Close the ticker to close all ticker
// channels.
func NewTicker(d time.Duration) (joe.Tocker, error) {
	p, err := NewProfiler()
	if err != nil {
		return nil, err
	}
	t := Ticker{Ticker: joe.NewTicker(d), Data: make(chan *structs.DiskStats), Profiler: p}
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
