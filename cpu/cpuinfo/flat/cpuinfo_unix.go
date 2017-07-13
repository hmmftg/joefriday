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

// Package cpuinfo (flat) handles Flatbuffer based processing of /proc/cpuinfo.
// Instead of returning a Go struct, it returns Flatbuffer serialized bytes. A
// function to deserialize the Flatbuffer serialized bytes into a
// cpuinfo.CPUInfo struct is provided.
//
// Note: the package name is cpuinfo and not the final element of the import
// path (flat). 
package cpuinfo

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	info "github.com/mohae/joefriday/cpu/cpuinfo"
	"github.com/mohae/joefriday/cpu/cpuinfo/flat/structs"
)

// Profiler is used to process the /proc/cpuinfo file as Flatbuffers serialized
// bytes.
type Profiler struct {
	*info.Profiler
	*fb.Builder
}

// Initializes and returns a cpuinfo profiler.
func NewProfiler() (p *Profiler, err error) {
	prof, err := info.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current cpuinfo as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	inf, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(inf), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the current cpuinfo as Flatbuffer serialized bytes using the
// package's global profiler.
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
func (p *Profiler) Serialize(inf *info.CPUInfo) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	uoffs := make([]fb.UOffsetT, len(inf.CPU))
	for i, cpu := range inf.CPU {
		uoffs[i] = p.SerializeCPU(&cpu)
	}
	structs.CPUInfoStartCPUVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cpusV := p.Builder.EndVector(len(uoffs))
	structs.CPUInfoStart(p.Builder)
	structs.CPUInfoAddTimestamp(p.Builder, inf.Timestamp)
	structs.CPUInfoAddCPU(p.Builder, cpusV)
	p.Builder.Finish(structs.CPUInfoEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// Serialize serializes a CPU using flatbuffers and returns the resulting
// UOffsetT.
func (p *Profiler) SerializeCPU(cpu *info.CPU) fb.UOffsetT {
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
	structs.CPUStartFlagsVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	flags := p.Builder.EndVector(len(uoffs))
	structs.CPUStart(p.Builder)
	structs.CPUAddProcessor(p.Builder, cpu.Processor)
	structs.CPUAddVendorID(p.Builder, vendorID)
	structs.CPUAddCPUFamily(p.Builder, cpuFamily)
	structs.CPUAddModel(p.Builder, model)
	structs.CPUAddModelName(p.Builder, modelName)
	structs.CPUAddStepping(p.Builder, stepping)
	structs.CPUAddMicrocode(p.Builder, microcode)
	structs.CPUAddCPUMHz(p.Builder, cpu.CPUMHz)
	structs.CPUAddCacheSize(p.Builder, cacheSize)
	structs.CPUAddPhysicalID(p.Builder, cpu.PhysicalID)
	structs.CPUAddSiblings(p.Builder, cpu.Siblings)
	structs.CPUAddCoreID(p.Builder, cpu.CoreID)
	structs.CPUAddCPUCores(p.Builder, cpu.CPUCores)
	structs.CPUAddApicID(p.Builder, cpu.ApicID)
	structs.CPUAddInitialApicID(p.Builder, cpu.InitialApicID)
	structs.CPUAddFPU(p.Builder, fpu)
	structs.CPUAddFPUException(p.Builder, fpuException)
	structs.CPUAddCPUIDLevel(p.Builder, cpuIDLevel)
	structs.CPUAddWP(p.Builder, wp)
	structs.CPUAddBogoMIPS(p.Builder, cpu.BogoMIPS)
	structs.CPUAddCLFlushSize(p.Builder, clFlushSize)
	structs.CPUAddCacheAlignment(p.Builder, cacheAlignment)
	structs.CPUAddAddressSizes(p.Builder, addressSize)
	structs.CPUAddPowerManagement(p.Builder, powerManagement)
	structs.CPUAddFlags(p.Builder, flags)
	return structs.CPUEnd(p.Builder)
}

// Serialize cpuinfo.CPUInfo using the package global profiler.
func Serialize(inf *info.CPUInfo) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(inf), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as cpuinfo.CPUInfo.
func Deserialize(p []byte) *info.CPUInfo {
	fInf := structs.GetRootAsCPUInfo(p, 0)
	l := fInf.CPULength()
	inf := &info.CPUInfo{}
	fCPU := &structs.CPU{}
	cpu := info.CPU{}
	inf.Timestamp = fInf.Timestamp()
	for i := 0; i < l; i++ {
		if !fInf.CPU(fCPU, i) {
			continue
		}
		cpu.Processor = fCPU.Processor()
		cpu.VendorID = string(fCPU.VendorID())
		cpu.CPUFamily = string(fCPU.CPUFamily())
		cpu.Model = string(fCPU.Model())
		cpu.ModelName = string(fCPU.ModelName())
		cpu.Stepping = string(fCPU.Stepping())
		cpu.Microcode = string(fCPU.Microcode())
		cpu.CPUMHz = fCPU.CPUMHz()
		cpu.CacheSize = string(fCPU.CacheSize())
		cpu.PhysicalID = fCPU.PhysicalID()
		cpu.Siblings = fCPU.Siblings()
		cpu.CoreID = fCPU.CoreID()
		cpu.CPUCores = fCPU.CPUCores()
		cpu.ApicID = fCPU.ApicID()
		cpu.InitialApicID = fCPU.InitialApicID()
		cpu.FPU = string(fCPU.FPU())
		cpu.FPUException = string(fCPU.FPUException())
		cpu.CPUIDLevel = string(fCPU.CPUIDLevel())
		cpu.WP = string(fCPU.WP())
		cpu.Flags = make([]string, fCPU.FlagsLength())
		for i := 0; i < len(cpu.Flags); i++ {
			cpu.Flags[i] = string(fCPU.Flags(i))
		}
		cpu.BogoMIPS = fCPU.BogoMIPS()
		cpu.CLFlushSize = string(fCPU.CLFlushSize())
		cpu.CacheAlignment = string(fCPU.CacheAlignment())
		cpu.AddressSizes = string(fCPU.AddressSizes())
		cpu.PowerManagement = string(fCPU.PowerManagement())
		inf.CPU = append(inf.CPU, cpu)
	}
	return inf
}
