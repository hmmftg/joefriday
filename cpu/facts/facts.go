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

// Package facts handles processong of CPU facts: data gathered from
// /proc/cpuinfo.
package facts

import (
	//Flat "github.com/google/flatbuffers/go"

	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SermoDigital/helpers"
	joe "github.com/mohae/joefriday"
)

const procFile = "/proc/cpuinfo"

// Profiler is used to process the /proc/cpuinfo file.
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

// Reset resets the profiler resources so that it is ready to use.  The caller
// must hold the lock.
func (prof *Profiler) Reset() {
	prof.Val = prof.Val[:0]
	prof.NoLockReset()
}

// Get returns the current cpuinfo (Facts).
func (prof *Profiler) Get() (facts *Facts, err error) {
	var (
		cpuCnt, i, pos int
		n              uint64
		v              byte
		name, value    string
		cpu            Fact
	)
	prof.Lock()
	defer prof.Unlock()
	prof.Reset()
	facts = &Facts{Timestamp: time.Now().UTC().UnixNano()}
	for {
		prof.Line, err = prof.Buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		// First grab the attribute name; everything up to the ':'.  The key may have
		// spaces and has trailing spaces; that gets trimmed.
		for i, v = range prof.Line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			prof.Val = append(prof.Val, v)
		}
		name = strings.TrimSpace(string(prof.Val[:]))
		prof.Val = prof.Val[:0]
		// if there's anything left, the value is everything else; trim spaces
		if pos < len(prof.Line) {
			value = strings.TrimSpace(string(prof.Line[pos:]))
		}
		// check to see if this is flat.Facts for a different processor
		if name == "processor" {
			if cpuCnt > 0 {
				facts.CPU = append(facts.CPU, cpu)
			}
			cpuCnt++
			n, err = helpers.ParseUint([]byte(value))
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "fact: processor", Err: err}
			}
			cpu = Fact{Processor: int16(n)}
			continue
		}
		if name == "vendor_id" {
			cpu.VendorID = value
			continue
		}
		if name == "cpu family" {
			cpu.CPUFamily = value
			continue
		}
		if name == "model" {
			cpu.Model = value
			continue
		}
		if name == "model name" {
			cpu.ModelName = value
			continue
		}
		if name == "stepping" {
			cpu.Stepping = value
			continue
		}
		if name == "microcode" {
			cpu.Microcode = value
			continue
		}
		if name == "cpu MHz" {
			f, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: cpu MHz", Err: err}
			}
			cpu.CPUMHz = float32(f)
			continue
		}
		if name == "cache size" {
			cpu.CacheSize = value
			continue
		}
		if name == "physical id" {
			n, err = helpers.ParseUint([]byte(value))
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: physical id", Err: err}
			}
			cpu.PhysicalID = int16(n)
			continue
		}
		if name == "siblings" {
			n, err = helpers.ParseUint([]byte(value))
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: siblings", Err: err}
			}
			cpu.Siblings = int16(n)
			continue
		}
		if name == "core id" {
			n, err = helpers.ParseUint([]byte(value))
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: core id", Err: err}
			}
			cpu.CoreID = int16(n)
			continue
		}
		if name == "cpu cores" {
			n, err = helpers.ParseUint([]byte(value))
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: cpu cores", Err: err}
			}
			cpu.CPUCores = int16(n)
			continue
		}
		if name == "apicid" {
			n, err = helpers.ParseUint([]byte(value))
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: apicid", Err: err}
			}
			cpu.ApicID = int16(n)
			continue
		}
		if name == "initial apicid" {
			n, err = helpers.ParseUint([]byte(value))
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: initial apicid", Err: err}
			}
			cpu.InitialApicID = int16(n)
			continue
		}
		if name == "fpu" {
			cpu.FPU = value
			continue
		}
		if name == "fpu_exception" {
			cpu.FPUException = value
			continue
		}
		if name == "cpuid level" {
			cpu.CPUIDLevel = value
			continue
		}
		if name == "WP" {
			cpu.WP = value
			continue
		}
		if name == "flags" {
			cpu.Flags = value
			continue
		}
		if name == "bogomips" {
			f, err := strconv.ParseFloat(value, 32)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: bogomips", Err: err}
			}
			cpu.BogoMIPS = float32(f)
			continue
		}
		if name == "clflush size" {
			cpu.CLFlushSize = value
			continue
		}
		if name == "cache_alignment" {
			cpu.CacheAlignment = value
			continue
		}
		if name == "address sizes" {
			cpu.AddressSizes = value
			continue
		}
		if name == "power management" {
			cpu.PowerManagement = value
		}
	}
	facts.CPU = append(facts.CPU, cpu)
	return facts, nil
}

// TODO: is it even worth it to have this as a global?  Should Get just
// instantiate a local version and use that?
var stdProfiler *Profiler
var stdProfilerMu sync.Mutex

// Get returns the current cpuinfo (Facts) using the package's global
// Profiler.
func Get() (facts *Facts, err error) {
	stdProfilerMu.Lock()
	defer stdProfilerMu.Unlock()
	if stdProfiler == nil {
		stdProfiler, err = New()
		if err != nil {
			return nil, err
		}
	}
	return stdProfiler.Get()
}

// Facts are a collection of facts, cpuinfo, about the system's cpus.
type Facts struct {
	Timestamp int64
	CPU       []Fact `json:"cpu"`
}

// Fact holds the /proc/cpuinfo for a single processor.
type Fact struct {
	Processor       int16   `json:"processor"`
	VendorID        string  `json:"vendor_id"`
	CPUFamily       string  `json:"cpu_family"`
	Model           string  `json:"model"`
	ModelName       string  `json:"model_name"`
	Stepping        string  `json:"stepping"`
	Microcode       string  `json:"microcode"`
	CPUMHz          float32 `json:"cpu_mhz"`
	CacheSize       string  `json:"cache_size"`
	PhysicalID      int16   `json:"physical_id"`
	Siblings        int16   `json:"siblings"`
	CoreID          int16   `json:"core_id"`
	CPUCores        int16   `json:"cpu_cores"`
	ApicID          int16   `json:"apicid"`
	InitialApicID   int16   `json:"initial_apicid"`
	FPU             string  `json:"fpu"`
	FPUException    string  `json:"fpu_exception"`
	CPUIDLevel      string  `json:"cpuid_level"`
	WP              string  `json:"wp"`
	Flags           string  `json:"flags"` // should this be a []string?
	BogoMIPS        float32 `json:"bogomips"`
	CLFlushSize     string  `json:"clflush_size"`
	CacheAlignment  string  `json:"cache_alignment"`
	AddressSizes    string  `json:"address_sizes"`
	PowerManagement string  `json:"power_management"`
}
