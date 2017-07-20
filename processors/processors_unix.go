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

// Package processors gathers information about the physical processors on a
// system by parsing the information from /procs/cpuinfo. This package gathers
// basic information about each physical processor, cpu, on the system, with
// one entry per processor. 
//
// The CPUMHz field shouldn't be relied on; the CPU data of the first CPU
// processed for each processor is used. This value may be different than that
// of other cores on the processor and may also be higher or lower than the
// processor's base frequency because of dynamic frequency scaling and
// frequency boosts, like turbo. For more detailed information about each cpu
// core, use joefriday/cpuinfo, which returns an entry per core.
package processors

import (
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
)

const procFile = "/proc/cpuinfo"

// Processors holds information about a system's processors
type Processors struct {
	Timestamp int64  `json:"timestamp"`
	// The number of physical processors.
	Count     int16  `json:"count"`
	Socket     []Processor `json:"processor"`
}

// Processor holds the /proc/cpuinfo for a single physical cpu.
type Processor struct {
	PhysicalID int16    `json:"physical_id"`
	VendorID   string   `json:"vendor_id"`
	CPUFamily  string   `json:"cpu_family"`
	Model      string   `json:"model"`
	ModelName  string   `json:"model_name"`
	Stepping   string   `json:"stepping"`
	Microcode  string   `json:"microcode"`
	CPUMHz     float32  `json:"cpu_mhz"`
	CacheSize  string   `json:"cache_size"`
	CPUCores   int16    `json:"cpu_cores"`
	BogoMIPS   float32  `json:"bogomips"`
	Flags      []string `json:"flags"`
}

// Profiler is used to get the processor information by processing the
// /proc/cpuinfo file.
type Profiler struct {
	*joe.Proc
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.NewProc(procFile)
	if err != nil {
		return nil, err
	}
	return &Profiler{Proc: proc}, nil
}

// Get returns the processor information.
func (prof *Profiler) Get() (procs *Processors, err error) {
	var (
		i, pos, nameLen int
		priorID                 int16
		n                       uint64
		v                       byte
		proc                     Processor
		first                   = true // set to false after first proc
		add                     bool
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	procs = &Processors{Timestamp: time.Now().UTC().UnixNano()}
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
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
		v = prof.Val[0]
		if v == 'c' {
			v = prof.Val[1]
			if v == 'p' {
				v = prof.Val[4]
				if v == 'c' { // cpu cores
					n, err = helpers.ParseUint(prof.Val[nameLen:])
					if err != nil {
						return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
					}
					proc.CPUCores = int16(n)
					continue
				}
				if v == 'f' { // cpu family
					proc.CPUFamily = string(prof.Val[nameLen:])
					continue
				}
				if v == 'M' { // cpu MHz
					f, err := strconv.ParseFloat(string(prof.Val[nameLen:]), 32)
					if err != nil {
						return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
					}
					proc.CPUMHz = float32(f)
				}
				continue
			}
			if prof.Val[5] == ' ' { // cache size
				proc.CacheSize = string(prof.Val[nameLen:])
			}
			continue
		}
		if v == 'f' {
			if prof.Val[1] == 'l' { // flags
				proc.Flags = strings.Split(string(prof.Val[nameLen:]), " ")
			}
			continue
		}
		if v == 'm' {
			if prof.Val[1] == 'o' {
				if nameLen == 5 { // model
					proc.Model = string(prof.Val[nameLen:])
					continue
				}
				proc.ModelName = string(prof.Val[nameLen:])
				continue
			}
			if prof.Val[1] == 'i' {
				proc.Microcode = string(prof.Val[nameLen:])
			}
			continue
		}
		if v == 'p' {
			v = prof.Val[1]
			if v == 'h' { // physical id
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				proc.PhysicalID = int16(n)
				if first || proc.PhysicalID != priorID {
					add = true
				}
				priorID = proc.PhysicalID
				continue
			}
			// processor starts information about a logical processor; if there was a previously
			// processed processor, only add it if it is a different physical processor.
			if v == 'r' { // processor
				if add {
					procs.Socket = append(procs.Socket, proc)
					add = false
				}
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
			}
			continue
		}
		if v == 's' && prof.Val[1] == 't' { // stepping
			proc.Stepping = string(prof.Val[nameLen:])
			continue
		}
		if v == 'v' { // vendor_id
			proc.VendorID = string(prof.Val[nameLen:])
		}
		// also check 2nd name pos for o as some output also have a bugs line.
		if v == 'b' && prof.Val[1] == 'o' { // bogomips
			f, err := strconv.ParseFloat(string(prof.Val[nameLen:]), 32)
			if err != nil {
				return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
			}
			proc.BogoMIPS = float32(f)
			continue
		}
		
	}
	// append the current processor informatin
	if add {
		procs.Socket = append(procs.Socket, proc)
	}
	procs.Count = int16(len(procs.Socket))
	return procs, nil
}

var std *Profiler
var stdMu sync.Mutex

// Get returns the information about the processors using the package's global
// Profiler.
func Get() (procs *Processors, err error) {
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
