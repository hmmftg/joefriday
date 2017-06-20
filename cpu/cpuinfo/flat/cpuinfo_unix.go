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

// Package cpuinfo (flat) handles Flatbuffer based processing of CPU info.
// Instead of returning a Go struct, it returns Flatbuffer serialized bytes.
// A function to deserialize the Flatbuffer serialized bytes into a 
// cpuinfo.CPUs struct is provided. After the first use, the flatbuffer builder
// is reused.
//
// Note: the package name is cpuinfo and not the final element of the import
// path (flat). 
package cpuinfo

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	inf "github.com/mohae/joefriday/cpu/cpuinfo"
	"github.com/mohae/joefriday/cpu/cpuinfo/flat/flat"
)

// Profiler is used to process the cpuinfo as Flatbuffers serialized bytes.
type Profiler struct {
	*inf.Profiler
	*fb.Builder
}

// Initializes and returns a cpuinfo profiler that utilizes FlatBuffers.
func NewProfiler() (p *Profiler, err error) {
	prof, err := inf.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current cpuinfo as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	cpus, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(cpus), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the current cpuinfo as Flatbuffer serialized bytes using the
// package global profiler.
func Get() (p []byte, err error) {
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

// Serialize serializes cpuinfo using Flatbuffers.
func (p *Profiler) Serialize(cpus *inf.CPUs) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	uoffs := make([]fb.UOffsetT, len(cpus.CPU))
	for i, cpu := range cpus.CPU {
		uoffs[i] = p.SerializeCPU(&cpu)
	}
	flat.CPUsStartCPUVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cpusV := p.Builder.EndVector(len(uoffs))
	flat.CPUsStart(p.Builder)
	flat.CPUsAddTimestamp(p.Builder, cpus.Timestamp)
	flat.CPUsAddCPU(p.Builder, cpusV)
	p.Builder.Finish(flat.CPUsEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// Serialize serializes a CPU's info using flatbuffers and returns the
// resulting UOffsetT.
func (p *Profiler) SerializeCPU(cpu *inf.CPU) fb.UOffsetT {
	vendorID := p.Builder.CreateString(cpu.VendorID)
	cpuFamily := p.Builder.CreateString(cpu.CPUFamily)
	model := p.Builder.CreateString(cpu.Model)
	modelName := p.Builder.CreateString(cpu.ModelName)
	stepping := p.Builder.CreateString(cpu.Stepping)
	microcode := p.Builder.CreateString(cpu.Microcode)
	cacheSize := p.Builder.CreateString(cpu.CacheSize)
	fpu := p.Builder.CreateString(cpu.FPU)
	fpuException := p.Builder.CreateString(cpu.FPUException)
	cpuIDLevel := p.Builder.CreateString(cpu.CPUIDLevel)
	wp := p.Builder.CreateString(cpu.WP)
	clFlushSize := p.Builder.CreateString(cpu.CLFlushSize)
	cacheAlignment := p.Builder.CreateString(cpu.CacheAlignment)
	addressSize := p.Builder.CreateString(cpu.AddressSizes)
	powerManagement := p.Builder.CreateString(cpu.PowerManagement)
	uoffs := make([]fb.UOffsetT, len(cpu.Flags))
	for i, flag := range cpu.Flags {
		uoffs[i] = p.Builder.CreateString(flag)
	}
	flat.CPUsStartCPUVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	flags := p.Builder.EndVector(len(uoffs))
	flat.CPUStart(p.Builder)
	flat.CPUAddProcessor(p.Builder, cpu.Processor)
	flat.CPUAddVendorID(p.Builder, vendorID)
	flat.CPUAddCPUFamily(p.Builder, cpuFamily)
	flat.CPUAddModel(p.Builder, model)
	flat.CPUAddModelName(p.Builder, modelName)
	flat.CPUAddStepping(p.Builder, stepping)
	flat.CPUAddMicrocode(p.Builder, microcode)
	flat.CPUAddCPUMHz(p.Builder, cpu.CPUMHz)
	flat.CPUAddCacheSize(p.Builder, cacheSize)
	flat.CPUAddPhysicalID(p.Builder, cpu.PhysicalID)
	flat.CPUAddSiblings(p.Builder, cpu.Siblings)
	flat.CPUAddCoreID(p.Builder, cpu.CoreID)
	flat.CPUAddCPUCores(p.Builder, cpu.CPUCores)
	flat.CPUAddApicID(p.Builder, cpu.ApicID)
	flat.CPUAddInitialApicID(p.Builder, cpu.InitialApicID)
	flat.CPUAddFPU(p.Builder, fpu)
	flat.CPUAddFPUException(p.Builder, fpuException)
	flat.CPUAddCPUIDLevel(p.Builder, cpuIDLevel)
	flat.CPUAddWP(p.Builder, wp)
	flat.CPUAddBogoMIPS(p.Builder, cpu.BogoMIPS)
	flat.CPUAddCLFlushSize(p.Builder, clFlushSize)
	flat.CPUAddCacheAlignment(p.Builder, cacheAlignment)
	flat.CPUAddAddressSizes(p.Builder, addressSize)
	flat.CPUAddPowerManagement(p.Builder, powerManagement)
	flat.CPUAddFlags(p.Builder, flags)
	return flat.CPUEnd(p.Builder)
}

// Serialize Facts using the package global profiler.
func Serialize(cpus *inf.CPUs) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(cpus), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as cpuinfo.CPUs.
func Deserialize(p []byte) *inf.CPUs {
	flatCPUs := flat.GetRootAsCPUs(p, 0)
	l := flatCPUs.CPULength()
	cpus := &inf.CPUs{}
	flatCPU := &flat.CPU{}
	cpu := inf.CPU{}
	cpus.Timestamp = flatCPUs.Timestamp()
	for i := 0; i < l; i++ {
		if !flatCPUs.CPU(flatCPU, i) {
			continue
		}
		cpu.Processor = flatCPU.Processor()
		cpu.VendorID = string(flatCPU.VendorID())
		cpu.CPUFamily = string(flatCPU.CPUFamily())
		cpu.Model = string(flatCPU.Model())
		cpu.ModelName = string(flatCPU.ModelName())
		cpu.Stepping = string(flatCPU.Stepping())
		cpu.Microcode = string(flatCPU.Microcode())
		cpu.CPUMHz = flatCPU.CPUMHz()
		cpu.CacheSize = string(flatCPU.CacheSize())
		cpu.PhysicalID = flatCPU.PhysicalID()
		cpu.Siblings = flatCPU.Siblings()
		cpu.CoreID = flatCPU.CoreID()
		cpu.CPUCores = flatCPU.CPUCores()
		cpu.ApicID = flatCPU.ApicID()
		cpu.InitialApicID = flatCPU.InitialApicID()
		cpu.FPU = string(flatCPU.FPU())
		cpu.FPUException = string(flatCPU.FPUException())
		cpu.CPUIDLevel = string(flatCPU.CPUIDLevel())
		cpu.WP = string(flatCPU.WP())
		cpu.Flags = make([]string, flatCPU.FlagsLength())
		for i := 0; i < len(cpu.Flags); i++ {
			cpu.Flags[i] = string(flatCPU.Flags(i))
		}
		cpu.BogoMIPS = flatCPU.BogoMIPS()
		cpu.CLFlushSize = string(flatCPU.CLFlushSize())
		cpu.CacheAlignment = string(flatCPU.CacheAlignment())
		cpu.AddressSizes = string(flatCPU.AddressSizes())
		cpu.PowerManagement = string(flatCPU.PowerManagement())
		cpus.CPU = append(cpus.CPU, cpu)
	}
	return cpus
}
