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

package cpu

import (
	//Flat "github.com/google/flatbuffers/go"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	fb "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
	"github.com/mohae/joefriday/cpu/flat"
)

// Facts are a collection of facts, cpuinfo, about the system's cpus.
type Facts struct {
	Timestamp int64
	CPU       []Fact `json:"cpu"`
}

// Fact holds the /proc/cpuinfo for a single cpu
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

// SerializeFlat serializes Facts using Flatbuffers.
func (f *Facts) SerializeFlat() []byte {
	bldr := fb.NewBuilder(0)
	return f.SerializeFlatBuilder(bldr)
}

// SerializeFlatBuilder serializes the Facts using Flatbuffers.  The passed
// builder is used.  It is expected that the builder is ready to use (the
// caller is responsible for either creating a new builder or resetting an
// existing one.)
func (f *Facts) SerializeFlatBuilder(bldr *fb.Builder) []byte {
	flatFacts := make([]fb.UOffsetT, len(f.CPU))
	vendorIDs := make([]fb.UOffsetT, len(f.CPU))
	cpuFamilies := make([]fb.UOffsetT, len(f.CPU))
	models := make([]fb.UOffsetT, len(f.CPU))
	modelNames := make([]fb.UOffsetT, len(f.CPU))
	steppings := make([]fb.UOffsetT, len(f.CPU))
	microcodes := make([]fb.UOffsetT, len(f.CPU))
	cacheSizes := make([]fb.UOffsetT, len(f.CPU))
	fpus := make([]fb.UOffsetT, len(f.CPU))
	fpuExceptions := make([]fb.UOffsetT, len(f.CPU))
	cpuIDLevels := make([]fb.UOffsetT, len(f.CPU))
	wps := make([]fb.UOffsetT, len(f.CPU))
	flags := make([]fb.UOffsetT, len(f.CPU))
	clFlushSizes := make([]fb.UOffsetT, len(f.CPU))
	cacheAlignments := make([]fb.UOffsetT, len(f.CPU))
	addressSizes := make([]fb.UOffsetT, len(f.CPU))
	powerManagements := make([]fb.UOffsetT, len(f.CPU))
	// create the strings
	for i := 0; i < len(f.CPU); i++ {
		vendorIDs[i] = bldr.CreateString(f.CPU[i].VendorID)
		cpuFamilies[i] = bldr.CreateString(f.CPU[i].CPUFamily)
		models[i] = bldr.CreateString(f.CPU[i].Model)
		modelNames[i] = bldr.CreateString(f.CPU[i].ModelName)
		steppings[i] = bldr.CreateString(f.CPU[i].Stepping)
		microcodes[i] = bldr.CreateString(f.CPU[i].Microcode)
		cacheSizes[i] = bldr.CreateString(f.CPU[i].CacheSize)
		fpus[i] = bldr.CreateString(f.CPU[i].FPU)
		fpuExceptions[i] = bldr.CreateString(f.CPU[i].FPUException)
		cpuIDLevels[i] = bldr.CreateString(f.CPU[i].CPUIDLevel)
		wps[i] = bldr.CreateString(f.CPU[i].WP)
		flags[i] = bldr.CreateString(f.CPU[i].Flags)
		clFlushSizes[i] = bldr.CreateString(f.CPU[i].CLFlushSize)
		cacheAlignments[i] = bldr.CreateString(f.CPU[i].CacheAlignment)
		addressSizes[i] = bldr.CreateString(f.CPU[i].AddressSizes)
		powerManagements[i] = bldr.CreateString(f.CPU[i].PowerManagement)
	}
	// create the CPUs
	for i := 0; i < len(f.CPU); i++ {
		flat.FactStart(bldr)
		flat.FactAddProcessor(bldr, f.CPU[i].Processor)
		flat.FactAddVendorID(bldr, vendorIDs[i])
		flat.FactAddCPUFamily(bldr, cpuFamilies[i])
		flat.FactAddModel(bldr, models[i])
		flat.FactAddModelName(bldr, modelNames[i])
		flat.FactAddStepping(bldr, steppings[i])
		flat.FactAddMicrocode(bldr, microcodes[i])
		flat.FactAddCPUMHz(bldr, f.CPU[i].CPUMHz)
		flat.FactAddCacheSize(bldr, cacheSizes[i])
		flat.FactAddPhysicalID(bldr, f.CPU[i].PhysicalID)
		flat.FactAddSiblings(bldr, f.CPU[i].Siblings)
		flat.FactAddCoreID(bldr, f.CPU[i].CoreID)
		flat.FactAddCPUCores(bldr, f.CPU[i].CPUCores)
		flat.FactAddApicID(bldr, f.CPU[i].ApicID)
		flat.FactAddInitialApicID(bldr, f.CPU[i].InitialApicID)
		flat.FactAddFPU(bldr, fpus[i])
		flat.FactAddFPUException(bldr, fpuExceptions[i])
		flat.FactAddCPUIDLevel(bldr, cpuIDLevels[i])
		flat.FactAddWP(bldr, wps[i])
		flat.FactAddFlags(bldr, flags[i])
		flat.FactAddBogoMIPS(bldr, f.CPU[i].BogoMIPS)
		flat.FactAddCLFlushSize(bldr, clFlushSizes[i])
		flat.FactAddCacheAlignment(bldr, cacheAlignments[i])
		flat.FactAddAddressSizes(bldr, addressSizes[i])
		flat.FactAddPowerManagement(bldr, powerManagements[i])
		flatFacts[i] = flat.FactEnd(bldr)
	}
	// Process the flat.Facts vector
	flat.FactsStartCPUVector(bldr, len(flatFacts))
	for i := len(f.CPU) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(flatFacts[i])
	}
	flatFactsV := bldr.EndVector(len(flatFacts))
	flat.FactsStart(bldr)
	flat.FactsAddTimestamp(bldr, f.Timestamp)
	flat.FactsAddCPU(bldr, flatFactsV)
	bldr.Finish(flat.FactsEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeFlat takes some bytes and deserialize's them, using Flatbuffers,
// into the Facts structure.
func DeserializeFlat(p []byte) *Facts {
	flatFacts := flat.GetRootAsFacts(p, 0)
	l := flatFacts.CPULength()
	facts := &Facts{}
	flatFact := &flat.Fact{}
	fact := Fact{}
	facts.Timestamp = flatFacts.Timestamp()
	for i := 0; i < l; i++ {
		if !flatFacts.CPU(flatFact, i) {
			continue
		}
		fact.Processor = flatFact.Processor()
		fact.VendorID = string(flatFact.VendorID())
		fact.CPUFamily = string(flatFact.CPUFamily())
		fact.Model = string(flatFact.Model())
		fact.ModelName = string(flatFact.ModelName())
		fact.Stepping = string(flatFact.Stepping())
		fact.Microcode = string(flatFact.Microcode())
		fact.CPUMHz = flatFact.CPUMHz()
		fact.CacheSize = string(flatFact.CacheSize())
		fact.PhysicalID = flatFact.PhysicalID()
		fact.Siblings = flatFact.Siblings()
		fact.CoreID = flatFact.CoreID()
		fact.CPUCores = flatFact.CPUCores()
		fact.ApicID = flatFact.ApicID()
		fact.InitialApicID = flatFact.InitialApicID()
		fact.FPU = string(flatFact.FPU())
		fact.FPUException = string(flatFact.FPUException())
		fact.CPUIDLevel = string(flatFact.CPUIDLevel())
		fact.WP = string(flatFact.WP())
		fact.Flags = string(flatFact.Flags())
		fact.BogoMIPS = flatFact.BogoMIPS()
		fact.CLFlushSize = string(flatFact.CLFlushSize())
		fact.CacheAlignment = string(flatFact.CacheAlignment())
		fact.AddressSizes = string(flatFact.AddressSizes())
		fact.PowerManagement = string(flatFact.PowerManagement())
		facts.CPU = append(facts.CPU, fact)
	}
	return facts
}

// GetFacts gets the processor information from /proc/cpuinfo
func GetFacts() (*Facts, error) {
	var procCnt, i, pos int
	var v byte
	var name, value string
	t := time.Now().UTC().UnixNano()
	f, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	facts := Facts{Timestamp: t}
	var cpu Fact
	val := make([]byte, 0, 160) // this should be large enough to hold flags...it'll grow if it isn't
	for {
		line, err := buf.ReadSlice('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading output bytes: %s", err)
		}
		// First grab the attribute name; everything up to the ':'.  The key may have
		// spaces and has trailing spaces; that gets trimmed.
		for i, v = range line {
			if v == 0x3A {
				pos = i + 1
				break
			}
			val = append(val, v)
		}
		name = strings.TrimSpace(string(val[:]))
		val = val[:0]
		// if there's anything left, the value is everything else; trim spaces
		if pos < len(line) {
			value = strings.TrimSpace(string(line[pos:]))
		}
		// check to see if this is flat.Facts for a different processor
		if name == "processor" {
			if procCnt > 0 {
				facts.CPU = append(facts.CPU, cpu)
			}
			procCnt++
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "fact: processor", Err: err}
			}
			cpu = Fact{Processor: int16(i)}
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
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: physical id", Err: err}
			}
			cpu.PhysicalID = int16(i)
			continue
		}
		if name == "siblings" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: siblings", Err: err}
			}
			cpu.Siblings = int16(i)
			continue
		}
		if name == "core id" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: core id", Err: err}
			}
			cpu.CoreID = int16(i)
			continue
		}
		if name == "cpu cores" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: cpu cores", Err: err}
			}
			cpu.CPUCores = int16(i)
			continue
		}
		if name == "apicid" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: apicid", Err: err}
			}
			cpu.ApicID = int16(i)
			continue
		}
		if name == "initial apicid" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpu", Op: "facts: initial apicid", Err: err}
			}
			cpu.InitialApicID = int16(i)
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
	return &facts, nil
}
