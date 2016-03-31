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

package flat

import (
	"sync"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/cpu/facts"
)

type Profiler struct {
	*facts.Profiler
	*fb.Builder
}

// Initializes and returns a cpu facts profiler that utilizes FlatBuffers.
func New() (prof *Profiler, err error) {
	factsProf, err := facts.New()
	if err != nil {
		return nil, err
	}
	return &Profiler{Profiler: factsProf, Builder: fb.NewBuilder(0)}, nil
}

func (prof *Profiler) reset() {
	prof.Profiler.Lock()
	prof.Builder.Reset()
	prof.Profiler.Unlock()
	prof.Profiler.Reset()
}

// Get returns the current cpuinfo as Flatbuffer serialized bytes.
func (prof *Profiler) Get() ([]byte, error) {
	prof.reset()
	facts, err := prof.Profiler.Get()
	if err != nil {
		return nil, err
	}
	return prof.Serialize(facts), nil
}

var std *Profiler    // global for convenience; lazily instantiated.
var stdMu sync.Mutex // protects access

// Get returns the current cpuinfo as Flatbuffer serialized bytes.
func Get() (p []byte, err error) {
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

// Serialize serializes the Facts using Flatbuffers.
func (prof *Profiler) Serialize(fcts *facts.Facts) []byte {
	prof.Profiler.Lock()
	defer prof.Profiler.Unlock()
	flatFacts := make([]fb.UOffsetT, len(fcts.CPU))
	vendorIDs := make([]fb.UOffsetT, len(fcts.CPU))
	cpuFamilies := make([]fb.UOffsetT, len(fcts.CPU))
	models := make([]fb.UOffsetT, len(fcts.CPU))
	modelNames := make([]fb.UOffsetT, len(fcts.CPU))
	steppings := make([]fb.UOffsetT, len(fcts.CPU))
	microcodes := make([]fb.UOffsetT, len(fcts.CPU))
	cacheSizes := make([]fb.UOffsetT, len(fcts.CPU))
	fpus := make([]fb.UOffsetT, len(fcts.CPU))
	fpuExceptions := make([]fb.UOffsetT, len(fcts.CPU))
	cpuIDLevels := make([]fb.UOffsetT, len(fcts.CPU))
	wps := make([]fb.UOffsetT, len(fcts.CPU))
	flags := make([]fb.UOffsetT, len(fcts.CPU))
	clFlushSizes := make([]fb.UOffsetT, len(fcts.CPU))
	cacheAlignments := make([]fb.UOffsetT, len(fcts.CPU))
	addressSizes := make([]fb.UOffsetT, len(fcts.CPU))
	powerManagements := make([]fb.UOffsetT, len(fcts.CPU))
	// create the strings
	for i := 0; i < len(fcts.CPU); i++ {
		vendorIDs[i] = prof.Builder.CreateString(fcts.CPU[i].VendorID)
		cpuFamilies[i] = prof.Builder.CreateString(fcts.CPU[i].CPUFamily)
		models[i] = prof.Builder.CreateString(fcts.CPU[i].Model)
		modelNames[i] = prof.Builder.CreateString(fcts.CPU[i].ModelName)
		steppings[i] = prof.Builder.CreateString(fcts.CPU[i].Stepping)
		microcodes[i] = prof.Builder.CreateString(fcts.CPU[i].Microcode)
		cacheSizes[i] = prof.Builder.CreateString(fcts.CPU[i].CacheSize)
		fpus[i] = prof.Builder.CreateString(fcts.CPU[i].FPU)
		fpuExceptions[i] = prof.Builder.CreateString(fcts.CPU[i].FPUException)
		cpuIDLevels[i] = prof.Builder.CreateString(fcts.CPU[i].CPUIDLevel)
		wps[i] = prof.Builder.CreateString(fcts.CPU[i].WP)
		flags[i] = prof.Builder.CreateString(fcts.CPU[i].Flags)
		clFlushSizes[i] = prof.Builder.CreateString(fcts.CPU[i].CLFlushSize)
		cacheAlignments[i] = prof.Builder.CreateString(fcts.CPU[i].CacheAlignment)
		addressSizes[i] = prof.Builder.CreateString(fcts.CPU[i].AddressSizes)
		powerManagements[i] = prof.Builder.CreateString(fcts.CPU[i].PowerManagement)
	}
	// create the CPUs
	for i := 0; i < len(fcts.CPU); i++ {
		FactStart(prof.Builder)
		FactAddProcessor(prof.Builder, fcts.CPU[i].Processor)
		FactAddVendorID(prof.Builder, vendorIDs[i])
		FactAddCPUFamily(prof.Builder, cpuFamilies[i])
		FactAddModel(prof.Builder, models[i])
		FactAddModelName(prof.Builder, modelNames[i])
		FactAddStepping(prof.Builder, steppings[i])
		FactAddMicrocode(prof.Builder, microcodes[i])
		FactAddCPUMHz(prof.Builder, fcts.CPU[i].CPUMHz)
		FactAddCacheSize(prof.Builder, cacheSizes[i])
		FactAddPhysicalID(prof.Builder, fcts.CPU[i].PhysicalID)
		FactAddSiblings(prof.Builder, fcts.CPU[i].Siblings)
		FactAddCoreID(prof.Builder, fcts.CPU[i].CoreID)
		FactAddCPUCores(prof.Builder, fcts.CPU[i].CPUCores)
		FactAddApicID(prof.Builder, fcts.CPU[i].ApicID)
		FactAddInitialApicID(prof.Builder, fcts.CPU[i].InitialApicID)
		FactAddFPU(prof.Builder, fpus[i])
		FactAddFPUException(prof.Builder, fpuExceptions[i])
		FactAddCPUIDLevel(prof.Builder, cpuIDLevels[i])
		FactAddWP(prof.Builder, wps[i])
		FactAddFlags(prof.Builder, flags[i])
		FactAddBogoMIPS(prof.Builder, fcts.CPU[i].BogoMIPS)
		FactAddCLFlushSize(prof.Builder, clFlushSizes[i])
		FactAddCacheAlignment(prof.Builder, cacheAlignments[i])
		FactAddAddressSizes(prof.Builder, addressSizes[i])
		FactAddPowerManagement(prof.Builder, powerManagements[i])
		flatFacts[i] = FactEnd(prof.Builder)
	}
	// Process the flat.Facts vector
	FactsStartCPUVector(prof.Builder, len(flatFacts))
	for i := len(flatFacts) - 1; i >= 0; i-- {
		prof.Builder.PrependUOffsetT(flatFacts[i])
	}
	flatFactsV := prof.Builder.EndVector(len(flatFacts))
	FactsStart(prof.Builder)
	FactsAddTimestamp(prof.Builder, fcts.Timestamp)
	FactsAddCPU(prof.Builder, flatFactsV)
	prof.Builder.Finish(FactsEnd(prof.Builder))
	return prof.Builder.Bytes[prof.Builder.Head():]
}

func Serialize(fcts *facts.Facts) (p []byte, err error) {
	stdMu.Lock()
	defer stdMu.Unlock()
	if std == nil {
		std, err = New()
		if err != nil {
			return nil, err
		}
	} else {
		std.reset()
	}
	return std.Serialize(fcts), nil
}

// Deserialize takes some Flatbuffer serialized bytes and deserialize's them.
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
		fct.Flags = string(flatFact.Flags())
		fct.BogoMIPS = flatFact.BogoMIPS()
		fct.CLFlushSize = string(flatFact.CLFlushSize())
		fct.CacheAlignment = string(flatFact.CacheAlignment())
		fct.AddressSizes = string(flatFact.AddressSizes())
		fct.PowerManagement = string(flatFact.PowerManagement())
		fcts.CPU = append(fcts.CPU, fct)
	}
	return fcts
}
