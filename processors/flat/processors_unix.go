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
	"github.com/mohae/joefriday/processors"
	"github.com/mohae/joefriday/processors/flat/structs"
)

// Profiler is used to get the processor information, as Flatbuffers serialized
// bytes, by processing the /proc/cpuinfo file.
type Profiler struct {
	*processors.Profiler
	*fb.Builder
}

// Returns an initialized Profiler; ready to use.
func NewProfiler() (p *Profiler, err error) {
	prof, err := processors.NewProfiler()
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
func (p *Profiler) Serialize(procs *processors.Processors) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	uoffs := make([]fb.UOffsetT, len(procs.Socket))
	for i, proc := range procs.Socket {
		uoffs[i] = p.SerializeProcessor(&proc)
	}
	structs.ProcessorsStartSocketVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	procsV := p.Builder.EndVector(len(uoffs))
	structs.ProcessorsStart(p.Builder)
	structs.ProcessorsAddTimestamp(p.Builder, procs.Timestamp)
	structs.ProcessorsAddCount(p.Builder, procs.Count)
	structs.ProcessorsAddSocket(p.Builder, procsV)
	p.Builder.Finish(structs.ProcessorsEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

func (p *Profiler) SerializeProcessor(proc *processors.Processor) fb.UOffsetT {
	vendorID := p.Builder.CreateString(proc.VendorID)
	cpuFamily := p.Builder.CreateString(proc.CPUFamily)
	model := p.Builder.CreateString(proc.Model)
	modelName := p.Builder.CreateString(proc.ModelName)
	stepping := p.Builder.CreateString(proc.Stepping)
	microcode := p.Builder.CreateString(proc.Microcode)
	cacheSize := p.Builder.CreateString(proc.CacheSize)
	uoffs := make([]fb.UOffsetT, len(proc.Flags))
	for i, flag := range proc.Flags {
		uoffs[i] = p.Builder.CreateString(flag)
	}
	structs.ProcessorStartFlagsVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	flags := p.Builder.EndVector(len(uoffs))
	structs.ProcessorStart(p.Builder)
	structs.ProcessorAddPhysicalID(p.Builder, proc.PhysicalID)
	structs.ProcessorAddVendorID(p.Builder, vendorID)
	structs.ProcessorAddCPUFamily(p.Builder, cpuFamily)
	structs.ProcessorAddModel(p.Builder, model)
	structs.ProcessorAddModelName(p.Builder, modelName)
	structs.ProcessorAddStepping(p.Builder, stepping)
	structs.ProcessorAddMicrocode(p.Builder, microcode)
	structs.ProcessorAddCPUMHz(p.Builder, proc.CPUMHz)
	structs.ProcessorAddCacheSize(p.Builder, cacheSize)
	structs.ProcessorAddCPUCores(p.Builder, proc.CPUCores)
	structs.ProcessorAddBogoMIPS(p.Builder, proc.BogoMIPS)
	structs.ProcessorAddFlags(p.Builder, flags)
	return structs.ProcessorEnd(p.Builder)
}

// Serialize processors information.
func Serialize(proc *processors.Processors) (p []byte, err error) {
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
func Deserialize(p []byte) *processors.Processors {
	flatP := structs.GetRootAsProcessors(p, 0)
	procs := &processors.Processors{}
	flatProc := &structs.Processor{}
	proc := processors.Processor{}
	procs.Timestamp = flatP.Timestamp()
	procs.Socket = make([]processors.Processor, flatP.Count())
	for i := 0; i < len(procs.Socket); i++ {
		if !flatP.Socket(flatProc, i) {
			continue
		}
		proc.PhysicalID = flatProc.PhysicalID()
		proc.VendorID = string(flatProc.VendorID())
		proc.CPUFamily = string(flatProc.CPUFamily())
		proc.Model = string(flatProc.Model())
		proc.ModelName = string(flatProc.ModelName())
		proc.Stepping = string(flatProc.Stepping())
		proc.Microcode = string(flatProc.Microcode())
		proc.CPUMHz = flatProc.CPUMHz()
		proc.CacheSize = string(flatProc.CacheSize())
		proc.CPUCores = flatProc.CPUCores()
		proc.BogoMIPS = flatProc.BogoMIPS()
		proc.Flags = make([]string, flatProc.FlagsLength())
		for i := 0; i < len(proc.Flags); i++ {
			proc.Flags[i] = string(flatProc.Flags(i))
		}
		procs.Socket[i] = proc
	}
	return procs
}
