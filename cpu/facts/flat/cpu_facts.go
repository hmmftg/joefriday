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

// Package flat handles Flatbuffer based processing of CPU facts.  Instead of
// returning a Go struct, it returns Flatbuffer serialized bytes.  A function
// to deserialize the Flatbuffer serialized bytes into a facts.Facts struct
// is provided.  After the first use, the flatbuffer builder is reused.
package flat

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/cpu/facts"
)

// profiler is used to process the cpuinfo (facts) as Flatbuffers serialized
// bytes.
type Profiler struct {
	*facts.Profiler
	*fb.Builder
}

// Initializes and returns a cpu facts profiler that utilizes FlatBuffers.
func NewProfiler() (p *Profiler, err error) {
	prof, err := facts.NewProfiler()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: prof, Builder: fb.NewBuilder(0)}, nil
}

// Get returns the current cpuinfo (facts) as Flatbuffer serialized bytes.
func (p *Profiler) Get() ([]byte, error) {
	facts, err := p.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return p.Serialize(facts), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the current cpuinfo (facts) as Flatbuffer serialized bytes
// using the package global profiler.
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

// Serialize serializes Facts using Flatbuffers.
func (p *Profiler) Serialize(fcts *facts.Facts) []byte {
	// ensure the Builder is in a usable state.
	p.Builder.Reset()
	uoffs := make([]fb.UOffsetT, len(fcts.CPU))
	for i, cpu := range fcts.CPU {
		uoffs[i] = p.SerializeFact(&cpu)
	}
	FactsStartCPUVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	cpus := p.Builder.EndVector(len(uoffs))
	FactsStart(p.Builder)
	FactsAddTimestamp(p.Builder, fcts.Timestamp)
	FactsAddCPU(p.Builder, cpus)
	p.Builder.Finish(FactsEnd(p.Builder))
	b := p.Builder.Bytes[p.Builder.Head():]
	// copy them (otherwise gets lost in reset)
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// Serialize serializes a CPU's Fact using flatbuffers and returns the
// resulting UOffsetT.
func (p *Profiler) SerializeFact(f *facts.Fact) fb.UOffsetT {
	vendorID := p.Builder.CreateString(f.VendorID)
	cpuFamily := p.Builder.CreateString(f.CPUFamily)
	model := p.Builder.CreateString(f.Model)
	modelName := p.Builder.CreateString(f.ModelName)
	stepping := p.Builder.CreateString(f.Stepping)
	microcode := p.Builder.CreateString(f.Microcode)
	cacheSize := p.Builder.CreateString(f.CacheSize)
	fpu := p.Builder.CreateString(f.FPU)
	fpuException := p.Builder.CreateString(f.FPUException)
	cpuIDLevel := p.Builder.CreateString(f.CPUIDLevel)
	wp := p.Builder.CreateString(f.WP)
	clFlushSize := p.Builder.CreateString(f.CLFlushSize)
	cacheAlignment := p.Builder.CreateString(f.CacheAlignment)
	addressSize := p.Builder.CreateString(f.AddressSizes)
	powerManagement := p.Builder.CreateString(f.PowerManagement)
	uoffs := make([]fb.UOffsetT, len(f.Flags))
	for i, flag := range f.Flags {
		uoffs[i] = p.Builder.CreateString(flag)
	}
	FactsStartCPUVector(p.Builder, len(uoffs))
	for i := len(uoffs) - 1; i >= 0; i-- {
		p.Builder.PrependUOffsetT(uoffs[i])
	}
	flags := p.Builder.EndVector(len(uoffs))
	FactStart(p.Builder)
	FactAddProcessor(p.Builder, f.Processor)
	FactAddVendorID(p.Builder, vendorID)
	FactAddCPUFamily(p.Builder, cpuFamily)
	FactAddModel(p.Builder, model)
	FactAddModelName(p.Builder, modelName)
	FactAddStepping(p.Builder, stepping)
	FactAddMicrocode(p.Builder, microcode)
	FactAddCPUMHz(p.Builder, f.CPUMHz)
	FactAddCacheSize(p.Builder, cacheSize)
	FactAddPhysicalID(p.Builder, f.PhysicalID)
	FactAddSiblings(p.Builder, f.Siblings)
	FactAddCoreID(p.Builder, f.CoreID)
	FactAddCPUCores(p.Builder, f.CPUCores)
	FactAddApicID(p.Builder, f.ApicID)
	FactAddInitialApicID(p.Builder, f.InitialApicID)
	FactAddFPU(p.Builder, fpu)
	FactAddFPUException(p.Builder, fpuException)
	FactAddCPUIDLevel(p.Builder, cpuIDLevel)
	FactAddWP(p.Builder, wp)
	FactAddBogoMIPS(p.Builder, f.BogoMIPS)
	FactAddCLFlushSize(p.Builder, clFlushSize)
	FactAddCacheAlignment(p.Builder, cacheAlignment)
	FactAddAddressSizes(p.Builder, addressSize)
	FactAddPowerManagement(p.Builder, powerManagement)
	FactAddFlags(p.Builder, flags)
	return FactEnd(p.Builder)
}

// Serialize Facts using the package global profiler.
func Serialize(fcts *facts.Facts) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = NewProfiler()
		if err != nil {
			return nil, err
		}
	}
	return std.Serialize(fcts), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them
// as fact.Facts.
func Deserialize(p []byte) *facts.Facts {
	flatFacts := GetRootAsFacts(p, 0)
	l := flatFacts.CPULength()
	fcts := &facts.Facts{}
	flatFact := &Fact{}
	fct := facts.Fact{}
	fcts.Timestamp = flatFacts.Timestamp()
	for i := 0; i < l; i++ {
		if !flatFacts.CPU(flatFact, i) {
			continue
		}
		fct.Processor = flatFact.Processor()
		fct.VendorID = string(flatFact.VendorID())
		fct.CPUFamily = string(flatFact.CPUFamily())
		fct.Model = string(flatFact.Model())
		fct.ModelName = string(flatFact.ModelName())
		fct.Stepping = string(flatFact.Stepping())
		fct.Microcode = string(flatFact.Microcode())
		fct.CPUMHz = flatFact.CPUMHz()
		fct.CacheSize = string(flatFact.CacheSize())
		fct.PhysicalID = flatFact.PhysicalID()
		fct.Siblings = flatFact.Siblings()
		fct.CoreID = flatFact.CoreID()
		fct.CPUCores = flatFact.CPUCores()
		fct.ApicID = flatFact.ApicID()
		fct.InitialApicID = flatFact.InitialApicID()
		fct.FPU = string(flatFact.FPU())
		fct.FPUException = string(flatFact.FPUException())
		fct.CPUIDLevel = string(flatFact.CPUIDLevel())
		fct.WP = string(flatFact.WP())
		fct.Flags = make([]string, flatFact.FlagsLength())
		for i := 0; i < len(fct.Flags); i++ {
			fct.Flags[i] = string(flatFact.Flags(i))
		}
		fct.BogoMIPS = flatFact.BogoMIPS()
		fct.CLFlushSize = string(flatFact.CLFlushSize())
		fct.CacheAlignment = string(flatFact.CacheAlignment())
		fct.AddressSizes = string(flatFact.AddressSizes())
		fct.PowerManagement = string(flatFact.PowerManagement())
		fcts.CPU = append(fcts.CPU, fct)
	}
	return fcts
}
