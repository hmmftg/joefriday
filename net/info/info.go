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

// Package net gets and processes /proc/net/dev, returning the infoFlat in the
// appropriate format.
package info

import (
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/net/structs"
)

const procFile = "/proc/net/dev"

type Profiler struct {
	*joe.Proc
}

func New() (prof *Profiler, err error) {
	proc, err := joe.New(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// GetInfo returns some of the results of /proc/net/dev
func (prof *Profiler) Get() (*structs.Info, error) {
	var l, i, pos, fieldNum, fieldVal int
	var v byte
	err := prof.Reset()
	if err != nil {
		return nil, err
	}
	// there's always at least 2 interfaces (I think)
	inf := &structs.Info{Timestamp: time.Now().UTC().UnixNano(), Interfaces: make([]structs.Interface, 0, 2)}
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		l++
		if l < 3 {
			continue
		}
		var iInfo structs.Interface

		// first grab the interface name (everything up to the ':')
		for i, v = range prof.Line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			prof.Val = append(prof.Val, v)
		}
		iInfo.Name = string(prof.Val[:])
		fieldNum = 0
		// process the rest of the line
		for {
			fieldNum++
			prof.Val = prof.Val[:0]
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
					pos += i
					break
				}
				prof.Val = append(prof.Val, v)
			}
			// any conversion error results in 0
			fieldVal, err = strconv.Atoi(string(prof.Val[:]))
			if err != nil {
				return nil, fmt.Errorf("%s: %s", iInfo.Name, err)
			}
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
		inf.Interfaces = append(inf.Interfaces, iInfo)
	}
	return inf, nil
}

var std *Profiler
var stdMu sync.Mutex

func Get() (inf *structs.Info, err error) {
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

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			info, err := prof.Get()
			if err != nil {
				errs <- err
				continue
			}
			out <- info
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
