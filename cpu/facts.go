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

	flat "github.com/google/flatbuffers/go"
	joe "github.com/mohae/joefriday"
)

// Facts are a collection of facts, cpuinfo, about the system's cpus.
type Facts struct {
	Timestamp int64
	CPUs      []Fact `json:"cpus"`
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
	bldr := flat.NewBuilder(0)
	return f.SerializeFlatBuilder(bldr)
}

// SerializeFlatBuilder serializes the Facts using Flatbuffers.  The passed
// builder is used.  It is expected that the builder is ready to use (the
// caller is responsible for either creating a new builder or resetting an
// existing one.)
func (f *Facts) SerializeFlatBuilder(bldr *flat.Builder) []byte {
	factFlats := make([]flat.UOffsetT, len(f.CPUs))
	vendorIDs := make([]flat.UOffsetT, len(f.CPUs))
	cpuFamilies := make([]flat.UOffsetT, len(f.CPUs))
	models := make([]flat.UOffsetT, len(f.CPUs))
	modelNames := make([]flat.UOffsetT, len(f.CPUs))
	steppings := make([]flat.UOffsetT, len(f.CPUs))
	microcodes := make([]flat.UOffsetT, len(f.CPUs))
	cacheSizes := make([]flat.UOffsetT, len(f.CPUs))
	fpus := make([]flat.UOffsetT, len(f.CPUs))
	fpuExceptions := make([]flat.UOffsetT, len(f.CPUs))
	cpuIDLevels := make([]flat.UOffsetT, len(f.CPUs))
	wps := make([]flat.UOffsetT, len(f.CPUs))
	flags := make([]flat.UOffsetT, len(f.CPUs))
	clFlushSizes := make([]flat.UOffsetT, len(f.CPUs))
	cacheAlignments := make([]flat.UOffsetT, len(f.CPUs))
	addressSizes := make([]flat.UOffsetT, len(f.CPUs))
	powerManagements := make([]flat.UOffsetT, len(f.CPUs))
	// create the strings
	for i := 0; i < len(f.CPUs); i++ {
		vendorIDs[i] = bldr.CreateString(f.CPUs[i].VendorID)
		cpuFamilies[i] = bldr.CreateString(f.CPUs[i].CPUFamily)
		models[i] = bldr.CreateString(f.CPUs[i].Model)
		modelNames[i] = bldr.CreateString(f.CPUs[i].ModelName)
		steppings[i] = bldr.CreateString(f.CPUs[i].Stepping)
		microcodes[i] = bldr.CreateString(f.CPUs[i].Microcode)
		cacheSizes[i] = bldr.CreateString(f.CPUs[i].CacheSize)
		fpus[i] = bldr.CreateString(f.CPUs[i].FPU)
		fpuExceptions[i] = bldr.CreateString(f.CPUs[i].FPUException)
		cpuIDLevels[i] = bldr.CreateString(f.CPUs[i].CPUIDLevel)
		wps[i] = bldr.CreateString(f.CPUs[i].WP)
		flags[i] = bldr.CreateString(f.CPUs[i].Flags)
		clFlushSizes[i] = bldr.CreateString(f.CPUs[i].CLFlushSize)
		cacheAlignments[i] = bldr.CreateString(f.CPUs[i].CacheAlignment)
		addressSizes[i] = bldr.CreateString(f.CPUs[i].AddressSizes)
		powerManagements[i] = bldr.CreateString(f.CPUs[i].PowerManagement)
	}
	// create the CPUs
	for i := 0; i < len(f.CPUs); i++ {
		FactFlatStart(bldr)
		FactFlatAddProcessor(bldr, f.CPUs[i].Processor)
		FactFlatAddVendorID(bldr, vendorIDs[i])
		FactFlatAddCPUFamily(bldr, cpuFamilies[i])
		FactFlatAddModel(bldr, models[i])
		FactFlatAddModelName(bldr, modelNames[i])
		FactFlatAddStepping(bldr, steppings[i])
		FactFlatAddMicrocode(bldr, microcodes[i])
		FactFlatAddCPUMHz(bldr, f.CPUs[i].CPUMHz)
		FactFlatAddCacheSize(bldr, cacheSizes[i])
		FactFlatAddPhysicalID(bldr, f.CPUs[i].PhysicalID)
		FactFlatAddSiblings(bldr, f.CPUs[i].Siblings)
		FactFlatAddCoreID(bldr, f.CPUs[i].CoreID)
		FactFlatAddCPUCores(bldr, f.CPUs[i].CPUCores)
		FactFlatAddApicID(bldr, f.CPUs[i].ApicID)
		FactFlatAddInitialApicID(bldr, f.CPUs[i].InitialApicID)
		FactFlatAddFPU(bldr, fpus[i])
		FactFlatAddFPUException(bldr, fpuExceptions[i])
		FactFlatAddCPUIDLevel(bldr, cpuIDLevels[i])
		FactFlatAddWP(bldr, wps[i])
		FactFlatAddFlags(bldr, flags[i])
		FactFlatAddBogoMIPS(bldr, f.CPUs[i].BogoMIPS)
		FactFlatAddCLFlushSize(bldr, clFlushSizes[i])
		FactFlatAddCacheAlignment(bldr, cacheAlignments[i])
		FactFlatAddAddressSizes(bldr, addressSizes[i])
		FactFlatAddPowerManagement(bldr, powerManagements[i])
		factFlats[i] = FactFlatEnd(bldr)
	}
	// Process the FactsFlat vector
	FactsFlatStartCPUsVector(bldr, len(factFlats))
	for i := len(f.CPUs) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(factFlats[i])
	}
	factFlatsV := bldr.EndVector(len(factFlats))
	FactsFlatStart(bldr)
	FactsFlatAddTimestamp(bldr, f.Timestamp)
	FactsFlatAddCPUs(bldr, factFlatsV)
	bldr.Finish(FactsFlatEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

// DeserializeFlat takes some bytes and deserialize's them, using Flatbuffers,
// into the Facts structure.
func DeserializeFlat(p []byte) *Facts {
	factsFlat := GetRootAsFactsFlat(p, 0)
	l := factsFlat.CPUsLength()
	facts := &Facts{}
	factFlat := &FactFlat{}
	fact := Fact{}
	facts.Timestamp = factsFlat.Timestamp()
	for i := 0; i < l; i++ {
		if !factsFlat.CPUs(factFlat, i) {
			continue
		}
		fact.Processor = factFlat.Processor()
		fact.VendorID = string(factFlat.VendorID())
		fact.CPUFamily = string(factFlat.CPUFamily())
		fact.Model = string(factFlat.Model())
		fact.ModelName = string(factFlat.ModelName())
		fact.Stepping = string(factFlat.Stepping())
		fact.Microcode = string(factFlat.Microcode())
		fact.CPUMHz = factFlat.CPUMHz()
		fact.CacheSize = string(factFlat.CacheSize())
		fact.PhysicalID = factFlat.PhysicalID()
		fact.Siblings = factFlat.Siblings()
		fact.CoreID = factFlat.CoreID()
		fact.CPUCores = factFlat.CPUCores()
		fact.ApicID = factFlat.ApicID()
		fact.InitialApicID = factFlat.InitialApicID()
		fact.FPU = string(factFlat.FPU())
		fact.FPUException = string(factFlat.FPUException())
		fact.CPUIDLevel = string(factFlat.CPUIDLevel())
		fact.WP = string(factFlat.WP())
		fact.Flags = string(factFlat.Flags())
		fact.BogoMIPS = factFlat.BogoMIPS()
		fact.CLFlushSize = string(factFlat.CLFlushSize())
		fact.CacheAlignment = string(factFlat.CacheAlignment())
		fact.AddressSizes = string(factFlat.AddressSizes())
		fact.PowerManagement = string(factFlat.PowerManagement())
		facts.CPUs = append(facts.CPUs, fact)
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
		// check to see if this is FactsFlat for a different processor
		if name == "processor" {
			if procCnt > 0 {
				facts.CPUs = append(facts.CPUs, cpu)
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
	facts.CPUs = append(facts.CPUs, cpu)
	return &facts, nil
}
