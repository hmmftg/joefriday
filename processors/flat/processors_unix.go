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
// on the system, with one entry per processor. Instead of returning a Go
/// struct, Flatbuffer serialized bytes are returned. A function to deserialize
// the Flatbuffer serialized bytes into a processors.Processors struct is
// provided.
//
// CPUMHz currently provides the current speed of the first core encountered
// for each physical processor. Modern x86/x86-64 cores have the ability to
// shift their speed so this is just a point in time data point for that core;
// there may be other cores on the processor that are at higher and lower
// speeds at the time the data is read. This field is more useful for other
// architectures. For x86/x86-64 cores, the MHzMin and MHzMax fields provide
// information about the range of speeds that are possible for the cores.
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
	structs.ProcessorsAddCPUs(p.Builder, int64(procs.CPUs))
	structs.ProcessorsAddSockets(p.Builder, procs.Sockets)
	structs.ProcessorsAddCoresPerSocket(p.Builder, procs.CoresPerSocket)
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
	uoffs := make([]fb.UOffsetT, len(proc.CacheIDs))
	// serialize cache info in order; the flatbuffer table will have the info
	// in order so a separate cache ID list isn't necessary for flatbuffers.
	for i, id := range proc.CacheIDs {
		// If the ID doesn't exist, the 0 value will be used
		inf := proc.Cache[id]
		uoffs[i] = p.SerializeCache(id, inf)
	}
	structs.ProcessorStartCacheVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cache := p.Builder.EndVector(len(uoffs))

	uoffs = make([]fb.UOffsetT, len(proc.Flags))
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
	structs.ProcessorAddMHzMin(p.Builder, proc.MHzMin)
	structs.ProcessorAddMHzMax(p.Builder, proc.MHzMax)
	structs.ProcessorAddCPUCores(p.Builder, proc.CPUCores)
	structs.ProcessorAddThreadsPerCore(p.Builder, proc.ThreadsPerCore)
	structs.ProcessorAddBogoMIPS(p.Builder, proc.BogoMIPS)
	structs.ProcessorAddCacheSize(p.Builder, cacheSize)
	structs.ProcessorAddCache(p.Builder, cache)
	structs.ProcessorAddFlags(p.Builder, flags)
	return structs.ProcessorEnd(p.Builder)
}

// SerializeCache serializes a cache entry using flatbuffers and returns the
// resulting UOffsetT.
func (p *Profiler) SerializeCache(id, inf string) fb.UOffsetT {
	cID := p.Builder.CreateString(id)
	cInf := p.Builder.CreateString(inf)
	structs.CacheInfStart(p.Builder)
	structs.CacheInfAddID(p.Builder, cID)
	structs.CacheInfAddSize(p.Builder, cInf)
	return structs.CacheInfEnd(p.Builder)
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
	flatCache := &structs.CacheInf{}
	proc := processors.Processor{}
	procs.Timestamp = flatP.Timestamp()
	procs.CPUs = int(flatP.CPUs())
	procs.Sockets = flatP.Sockets()
	procs.CoresPerSocket = flatP.CoresPerSocket()
	procs.Socket = make([]processors.Processor, flatP.Sockets())
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
		proc.MHzMin = flatProc.MHzMin()
		proc.MHzMax = flatProc.MHzMax()
		proc.CPUCores = flatProc.CPUCores()
		proc.ThreadsPerCore = flatProc.ThreadsPerCore()
		proc.BogoMIPS = flatProc.BogoMIPS()
		proc.CacheSize = string(flatProc.CacheSize())
		proc.CacheIDs = make([]string, 0, flatProc.CacheLength())
		proc.Cache = make(map[string]string, flatProc.CacheLength())
		for j := 0; j < flatProc.CacheLength(); j++ {
			if !flatProc.Cache(flatCache, j) {
				continue
			}
			proc.CacheIDs = append(proc.CacheIDs, string(flatCache.ID()))
			proc.Cache[proc.CacheIDs[j]] = string(flatCache.Size())
		}
		proc.Flags = make([]string, flatProc.FlagsLength())
		for i := 0; i < len(proc.Flags); i++ {
			proc.Flags[i] = string(flatProc.Flags(i))
		}
		procs.Socket[i] = proc
	}
	return procs
}
