package cpu

import (
	//fb "github.com/google/flatbuffers/go"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	fb "github.com/google/flatbuffers/go"
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

// Serialize serializes Facts using Flatbuffers.
func (p *Facts) Serialize() []byte {
	bldr := fb.NewBuilder(0)
	factFBs := make([]fb.UOffsetT, len(p.CPUs))
	vendorIDs := make([]fb.UOffsetT, len(p.CPUs))
	cpuFamilies := make([]fb.UOffsetT, len(p.CPUs))
	models := make([]fb.UOffsetT, len(p.CPUs))
	modelNames := make([]fb.UOffsetT, len(p.CPUs))
	steppings := make([]fb.UOffsetT, len(p.CPUs))
	microcodes := make([]fb.UOffsetT, len(p.CPUs))
	cacheSizes := make([]fb.UOffsetT, len(p.CPUs))
	fpus := make([]fb.UOffsetT, len(p.CPUs))
	fpuExceptions := make([]fb.UOffsetT, len(p.CPUs))
	cpuIDLevels := make([]fb.UOffsetT, len(p.CPUs))
	wps := make([]fb.UOffsetT, len(p.CPUs))
	flags := make([]fb.UOffsetT, len(p.CPUs))
	clFlushSizes := make([]fb.UOffsetT, len(p.CPUs))
	cacheAlignments := make([]fb.UOffsetT, len(p.CPUs))
	addressSizes := make([]fb.UOffsetT, len(p.CPUs))
	powerManagements := make([]fb.UOffsetT, len(p.CPUs))
	// create the strings
	for i := 0; i < len(p.CPUs); i++ {
		vendorIDs[i] = bldr.CreateString(p.CPUs[i].VendorID)
		cpuFamilies[i] = bldr.CreateString(p.CPUs[i].CPUFamily)
		models[i] = bldr.CreateString(p.CPUs[i].Model)
		modelNames[i] = bldr.CreateString(p.CPUs[i].ModelName)
		steppings[i] = bldr.CreateString(p.CPUs[i].Stepping)
		microcodes[i] = bldr.CreateString(p.CPUs[i].Microcode)
		cacheSizes[i] = bldr.CreateString(p.CPUs[i].CacheSize)
		fpus[i] = bldr.CreateString(p.CPUs[i].FPU)
		fpuExceptions[i] = bldr.CreateString(p.CPUs[i].FPUException)
		cpuIDLevels[i] = bldr.CreateString(p.CPUs[i].CPUIDLevel)
		wps[i] = bldr.CreateString(p.CPUs[i].WP)
		flags[i] = bldr.CreateString(p.CPUs[i].Flags)
		clFlushSizes[i] = bldr.CreateString(p.CPUs[i].CLFlushSize)
		cacheAlignments[i] = bldr.CreateString(p.CPUs[i].CacheAlignment)
		addressSizes[i] = bldr.CreateString(p.CPUs[i].AddressSizes)
		powerManagements[i] = bldr.CreateString(p.CPUs[i].PowerManagement)
	}
	// create the CPUs
	for i := 0; i < len(p.CPUs); i++ {
		FactFBStart(bldr)
		FactFBAddProcessor(bldr, p.CPUs[i].Processor)
		FactFBAddVendorID(bldr, vendorIDs[i])
		FactFBAddCPUFamily(bldr, cpuFamilies[i])
		FactFBAddModel(bldr, models[i])
		FactFBAddModelName(bldr, modelNames[i])
		FactFBAddStepping(bldr, steppings[i])
		FactFBAddMicrocode(bldr, microcodes[i])
		FactFBAddCPUMHz(bldr, p.CPUs[i].CPUMHz)
		FactFBAddCacheSize(bldr, cacheSizes[i])
		FactFBAddPhysicalID(bldr, p.CPUs[i].PhysicalID)
		FactFBAddSiblings(bldr, p.CPUs[i].Siblings)
		FactFBAddCoreID(bldr, p.CPUs[i].CoreID)
		FactFBAddCPUCores(bldr, p.CPUs[i].CPUCores)
		FactFBAddApicID(bldr, p.CPUs[i].ApicID)
		FactFBAddInitialApicID(bldr, p.CPUs[i].InitialApicID)
		FactFBAddFPU(bldr, fpus[i])
		FactFBAddFPUException(bldr, fpuExceptions[i])
		FactFBAddCPUIDLevel(bldr, cpuIDLevels[i])
		FactFBAddWP(bldr, wps[i])
		FactFBAddFlags(bldr, flags[i])
		FactFBAddBogoMIPS(bldr, p.CPUs[i].BogoMIPS)
		FactFBAddCLFlushSize(bldr, clFlushSizes[i])
		FactFBAddCacheAlignment(bldr, cacheAlignments[i])
		FactFBAddAddressSizes(bldr, addressSizes[i])
		FactFBAddPowerManagement(bldr, powerManagements[i])
		factFBs[i] = FactFBEnd(bldr)
	}
	// Process the FactsFB vector
	FactsFBStartCPUsVector(bldr, len(factFBs))
	for i := len(p.CPUs) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(factFBs[i])
	}
	factFBsV := bldr.EndVector(len(factFBs))
	FactsFBStart(bldr)
	FactsFBAddTimestamp(bldr, p.Timestamp)
	FactsFBAddCPUs(bldr, factFBsV)
	bldr.Finish(FactsFBEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

func Deserialize(p []byte) *Facts {
	factsFB := GetRootAsFactsFB(p, 0)
	l := factsFB.CPUsLength()
	facts := &Facts{}
	factFB := &FactFB{}
	fact := Fact{}
	facts.Timestamp = factsFB.Timestamp()
	for i := 0; i < l; i++ {
		if !factsFB.CPUs(factFB, i) {
			continue
		}
		fact.Processor = factFB.Processor()
		fact.VendorID = string(factFB.VendorID())
		fact.CPUFamily = string(factFB.CPUFamily())
		fact.Model = string(factFB.Model())
		fact.ModelName = string(factFB.ModelName())
		fact.Stepping = string(factFB.Stepping())
		fact.Microcode = string(factFB.Microcode())
		fact.CPUMHz = factFB.CPUMHz()
		fact.CacheSize = string(factFB.CacheSize())
		fact.PhysicalID = factFB.PhysicalID()
		fact.Siblings = factFB.Siblings()
		fact.CoreID = factFB.CoreID()
		fact.CPUCores = factFB.CPUCores()
		fact.ApicID = factFB.ApicID()
		fact.InitialApicID = factFB.InitialApicID()
		fact.FPU = string(factFB.FPU())
		fact.FPUException = string(factFB.FPUException())
		fact.CPUIDLevel = string(factFB.CPUIDLevel())
		fact.WP = string(factFB.WP())
		fact.Flags = string(factFB.Flags())
		fact.BogoMIPS = factFB.BogoMIPS()
		fact.CLFlushSize = string(factFB.CLFlushSize())
		fact.CacheAlignment = string(factFB.CacheAlignment())
		fact.AddressSizes = string(factFB.AddressSizes())
		fact.PowerManagement = string(factFB.PowerManagement())
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
		// check to see if this is FactsFB for a different processor
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
