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
// system by parsing the information from /procs/cpuinfo and sysfs. This
// package gathers basic information about sockets, physical processors, etc.
// on the system, with one entry per processor. 
//
// CPUMHz currently provides the current speed of the first core encountered
// for each physical processor. Modern x86/x86-64 cores have the ability to
// shift their speed so this is just a point in time data point for that core;
// there may be other cores on the processor that are at higher and lower
// speeds at the time the data is read. This field is more useful for other
// architectures. For x86/x86-64 cores, the MHzMin and MHzMax fields provide
// information about the range of speeds that are possible for the cores.
package processors

import (
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/cpu/cpux"
)

const (
	procFile = "/proc/cpuinfo"
)

// Processors holds information about a system's processors
type Processors struct {
	Timestamp int64 `json:"timestamp"`
	// The number of sockets.
	Sockets        int32 `json:"sockets"`
	CoresPerSocket int16 `json:"cores_per_socket"`
	// Information about each processor in each socket.
	Socket []Processor `json:"socket"`
	CPUs   int         `json:"cpus"` // number of cpus on the system
}

// Processor holds the /proc/cpuinfo for a single physical cpu.
type Processor struct {
	PhysicalID     int32             `json:"physical_id"`
	VendorID       string            `json:"vendor_id"`
	CPUFamily      string            `json:"cpu_family"`
	Model          string            `json:"model"`
	ModelName      string            `json:"model_name"`
	Stepping       string            `json:"stepping"`
	Microcode      string            `json:"microcode"`
	CPUMHz         float32           `json:"cpu_mhz"`
	MHzMin         float32           `json:"mhz_min"`
	MHzMax         float32           `json:"mhz_max"`
	Cache          map[string]string `json:"cache"`
	CacheSize      string            `json:"cache_size"`
	CacheIDs       []string          `json:"cache_ids"`
	CPUCores       int32             `json:"cpu_cores"`
	ThreadsPerCore int8              `json:"threads_per_core"`
	BogoMIPS       float32           `json:"bogomips"`
	Flags          []string          `json:"flags"`
}

// Profiler is used to get the processor information by processing the
// /proc/cpuinfo file.
type Profiler struct {
	joe.Procer
	*joe.Buffer
	*cpux.Profiler // This is created with the profiler for testability.
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (prof *Profiler, err error) {
	proc, err := joe.NewProc(procFile)
	if err != nil {
		return nil, err
	}
	p, err := cpux.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Procer: proc, Buffer: joe.NewBuffer(), Profiler: p}, nil
}

// Reset resources: after reset, the profiler is ready to be used again.
func (prof *Profiler) Reset() error {
	prof.Buffer.Reset()
	return prof.Procer.Reset()
}

// Get returns the processor information.
func (prof *Profiler) Get() (procs *Processors, err error) {
	procs, err = prof.getCPUInfo()
	if err != nil {
		return nil, err
	}
	// process the system cpu info
	err = prof.getSysDevicesCPUInfo(procs)
	if err != nil {
		return nil, err
	}
	// get the core count and calculate cores per socket
	var cores int32
	for i := range procs.Socket {
		cores += procs.Socket[i].CPUCores
	}
	procs.CoresPerSocket = int16(cores / procs.Sockets)
	return procs, nil
}

func (prof *Profiler) getCPUInfo() (procs *Processors, err error) {
	var (
		i, pos, nameLen, cpuCnt int
		siblings                int32
		ids                     []int32
		n                       uint64
		v                       byte
		proc                    Processor
		add                     bool
	)
	err = prof.Reset()
	if err != nil {
		return nil, err
	}
	procs = &Processors{Timestamp: time.Now().UTC().UnixNano()}
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
					proc.CPUCores = int32(n)
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
			if v == 'a' && prof.Val[5] == ' ' {
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
				var exists bool
				for _, v := range ids {
					if v == int32(n) {
						exists = true
						break
					}
				}
				if !exists {
					add = true
					ids = append(ids, int32(n))
				}
				proc.PhysicalID = int32(n)
				continue
			}
			// processor starts information about a logical processor; if there was a previously
			// processed processor, only add it if it is a different physical processor.
			if v == 'r' { // processor
				cpuCnt++ // increment counter
				if add {
					proc.ThreadsPerCore = int8(siblings / proc.CPUCores)
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
		if v == 's' {
			if prof.Val[1] == 'i' { // siblings
				n, err = helpers.ParseUint(prof.Val[nameLen:])
				if err != nil {
					return nil, &joe.ParseError{Info: string(prof.Val[:nameLen]), Err: err}
				}
				siblings = int32(n)
				continue
			}
			if prof.Val[1] == 't' { // stepping
				proc.Stepping = string(prof.Val[nameLen:])
				continue
			}
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
	// append the current processor information
	if add {
		proc.ThreadsPerCore = int8(siblings / proc.CPUCores)
		procs.Socket = append(procs.Socket, proc)
	}
	procs.CPUs = cpuCnt
	return procs, nil
}

func (prof *Profiler) getSysDevicesCPUInfo(procs *Processors) error {
	// get the cpux profiler
	cpus, err := prof.Profiler.Get()
	if err != nil {
		return err
	}
	var ids []int32 // holds the encountered PhysicalIDs; each physicalID should only be processed once.
	// go through the results and use the first match per physical id
	for i := range cpus.CPU {
		var exists bool
		for _, v := range ids {
			if v == cpus.CPU[i].PhysicalPackageID {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		id := cpus.CPU[i].PhysicalPackageID
		ids = append(ids, id)
		//find the matching entry
		for j := range procs.Socket {
			if procs.Socket[j].PhysicalID != id {
				continue
			}
			procs.Socket[j].MHzMin = cpus.CPU[i].MHzMin
			procs.Socket[j].MHzMax = cpus.CPU[i].MHzMax
			procs.Socket[j].Cache = make(map[string]string, len(cpus.CPU[i].Cache))
			procs.Socket[j].CacheIDs = make([]string, len(cpus.CPU[i].CacheIDs))
			for k, id := range cpus.CPU[i].CacheIDs {
				procs.Socket[j].CacheIDs[k] = id
				procs.Socket[j].Cache[id] = cpus.CPU[i].Cache[id]
			}
		}
	}
	procs.Sockets = int32(len(ids))
	return nil
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
