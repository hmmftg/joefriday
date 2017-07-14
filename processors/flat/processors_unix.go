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
// one entry per processor. Instead of returning a Go struct, it returns
// Flatbuffer serialized bytes. A function to deserialize the Flatbuffer
// serialized bytes into a processors.Processors struct is provided.
//
// The CPUMHz field shouldn't be relied on; the CPU data of the first CPU
// processed for each processor is used. This value may be different than that
// of other cores on the processor and may also be higher or lower than the
// processor's base frequency because of dynamic frequency scaling and
// frequency boosts, like turbo. For more detailed information about each cpu
// core, use joefriday/cpuinfo, which returns an entry per core.
//
// Note: the package name is processors and not the final element of the import
// path (flat). 
package processors

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	procs "github.com/mohae/joefriday/processors"
	"github.com/mohae/joefriday/processors/flat/structs"
)

// Profiler is used to get the processor information, as Flatbuffers serialized
// bytes, by processing the /proc/cpuinfo file.
type Profiler struct {
	*procs.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (p *Profiler, err error) {
	prof, err := procs.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the processor information as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	proc, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(proc), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the current processor info as Flatbuffer serialized bytes using
// the package's global Profiler.
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

// Serialize serializes Processors using Flatbuffers.
func (p *Profiler) Serialize(proc *procs.Processors) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	uoffs := make([]fb.UOffsetT, len(proc.Socket))
	for i, cpu := range proc.Socket {
		uoffs[i] = p.SerializeCPU(&cpu)
	}
	structs.ProcessorsStartSocketVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cpus := p.Builder.EndVector(len(uoffs))
	structs.ProcessorsStart(p.Builder)
	structs.ProcessorsAddTimestamp(p.Builder, proc.Timestamp)
	structs.ProcessorsAddCount(p.Builder, proc.Count)
	structs.ProcessorsAddSocket(p.Builder, cpus)
	p.Builder.Finish(structs.ProcessorsEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

func (p *Profiler) SerializeCPU(c *procs.CPU) fb.UOffsetT {
	vendorID := p.Builder.CreateString(c.VendorID)
	cpuFamily := p.Builder.CreateString(c.CPUFamily)
	model := p.Builder.CreateString(c.Model)
	modelName := p.Builder.CreateString(c.ModelName)
	stepping := p.Builder.CreateString(c.Stepping)
	microcode := p.Builder.CreateString(c.Microcode)
	cacheSize := p.Builder.CreateString(c.CacheSize)
	uoffs := make([]fb.UOffsetT, len(c.Flags))
	for i, flag := range c.Flags {
		uoffs[i] = p.Builder.CreateString(flag)
	}
	structs.CPUStartFlagsVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	flags := p.Builder.EndVector(len(uoffs))
	structs.CPUStart(p.Builder)
	structs.CPUAddPhysicalID(p.Builder, c.PhysicalID)
	structs.CPUAddVendorID(p.Builder, vendorID)
	structs.CPUAddCPUFamily(p.Builder, cpuFamily)
	structs.CPUAddModel(p.Builder, model)
	structs.CPUAddModelName(p.Builder, modelName)
	structs.CPUAddStepping(p.Builder, stepping)
	structs.CPUAddMicrocode(p.Builder, microcode)
	structs.CPUAddCPUMHz(p.Builder, c.CPUMHz)
	structs.CPUAddCacheSize(p.Builder, cacheSize)
	structs.CPUAddCPUCores(p.Builder, c.CPUCores)
	structs.CPUAddBogoMIPS(p.Builder, c.BogoMIPS)
	structs.CPUAddFlags(p.Builder, flags)
	return structs.CPUEnd(p.Builder)
}

// Serialize processors information.
func Serialize(proc *procs.Processors) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(proc), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserializes them
// as processors.Processors.
func Deserialize(p []byte) *procs.Processors {
	flatP := structs.GetRootAsProcessors(p, 0)
	proc := &procs.Processors{}
	flatC := &structs.CPU{}
	cpu := procs.CPU{}
	proc.Timestamp = flatP.Timestamp()
	proc.Socket = make([]procs.CPU, flatP.Count())
	for i := 0; i < len(proc.Socket); i++ {
		if !flatP.Socket(flatC, i) {
			continue
		}
		cpu.PhysicalID = flatC.PhysicalID()
		cpu.VendorID = string(flatC.VendorID())
		cpu.CPUFamily = string(flatC.CPUFamily())
		cpu.Model = string(flatC.Model())
		cpu.ModelName = string(flatC.ModelName())
		cpu.Stepping = string(flatC.Stepping())
		cpu.Microcode = string(flatC.Microcode())
		cpu.CPUMHz = flatC.CPUMHz()
		cpu.CacheSize = string(flatC.CacheSize())
		cpu.CPUCores = flatC.CPUCores()
		cpu.BogoMIPS = flatC.BogoMIPS()
		cpu.Flags = make([]string, flatC.FlagsLength())
		for i := 0; i < len(cpu.Flags); i++ {
			cpu.Flags[i] = string(flatC.Flags(i))
		}
		proc.Socket[i] = cpu
	}
	return proc
}
