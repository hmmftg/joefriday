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

type Processors struct {
	Timestamp int64
	Infos     []Inf `json:"cpus"`
}

// Info holds the /proc/cpuinfo for a single cpu
type Inf struct {
	Processor       int16   `json:"processor"`
	VendorID        string  `json:"vendor_id"`
	CPUFamily       string  `json:"cpu_family"`
	Model           string  `json:"model"`
	ModelName       string  `json:"model_name"`
	Stepping        string  `json:"stepping"`
	Microcode       string  `json:"microcode"`
	CPUMHz          string  `json:"cpu_mhz"`
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

// Serialize serializes Processors using Flatbuffers.
func (p *Processors) Serialize() []byte {
	bldr := fb.NewBuilder(0)
	infos := make([]fb.UOffsetT, len(p.Infos))
	vendorIDs := make([]fb.UOffsetT, len(p.Infos))
	cpuFamilies := make([]fb.UOffsetT, len(p.Infos))
	models := make([]fb.UOffsetT, len(p.Infos))
	modelNames := make([]fb.UOffsetT, len(p.Infos))
	steppings := make([]fb.UOffsetT, len(p.Infos))
	microcodes := make([]fb.UOffsetT, len(p.Infos))
	cpuMHzs := make([]fb.UOffsetT, len(p.Infos))
	cacheSizes := make([]fb.UOffsetT, len(p.Infos))
	fpus := make([]fb.UOffsetT, len(p.Infos))
	fpuExceptions := make([]fb.UOffsetT, len(p.Infos))
	cpuIDLevels := make([]fb.UOffsetT, len(p.Infos))
	wps := make([]fb.UOffsetT, len(p.Infos))
	flags := make([]fb.UOffsetT, len(p.Infos))
	clFlushSizes := make([]fb.UOffsetT, len(p.Infos))
	cacheAlignments := make([]fb.UOffsetT, len(p.Infos))
	addressSizes := make([]fb.UOffsetT, len(p.Infos))
	powerManagements := make([]fb.UOffsetT, len(p.Infos))
	// create the strings
	for i := 0; i < len(p.Infos); i++ {
		vendorIDs[i] = bldr.CreateString(p.Infos[i].VendorID)
		cpuFamilies[i] = bldr.CreateString(p.Infos[i].CPUFamily)
		models[i] = bldr.CreateString(p.Infos[i].Model)
		modelNames[i] = bldr.CreateString(p.Infos[i].ModelName)
		steppings[i] = bldr.CreateString(p.Infos[i].Stepping)
		microcodes[i] = bldr.CreateString(p.Infos[i].Microcode)
		cpuMHzs[i] = bldr.CreateString(p.Infos[i].CPUMHz)
		cacheSizes[i] = bldr.CreateString(p.Infos[i].CacheSize)
		fpus[i] = bldr.CreateString(p.Infos[i].FPU)
		fpuExceptions[i] = bldr.CreateString(p.Infos[i].FPUException)
		cpuIDLevels[i] = bldr.CreateString(p.Infos[i].CPUIDLevel)
		wps[i] = bldr.CreateString(p.Infos[i].WP)
		flags[i] = bldr.CreateString(p.Infos[i].Flags)
		clFlushSizes[i] = bldr.CreateString(p.Infos[i].CLFlushSize)
		cacheAlignments[i] = bldr.CreateString(p.Infos[i].CacheAlignment)
		addressSizes[i] = bldr.CreateString(p.Infos[i].AddressSizes)
		powerManagements[i] = bldr.CreateString(p.Infos[i].PowerManagement)
	}
	// create the Infos
	for i := 0; i < len(p.Infos); i++ {
		InfoStart(bldr)
		InfoAddProcessor(bldr, p.Infos[i].Processor)
		InfoAddVendorID(bldr, vendorIDs[i])
		InfoAddCPUFamily(bldr, cpuFamilies[i])
		InfoAddModel(bldr, models[i])
		InfoAddModelName(bldr, modelNames[i])
		InfoAddStepping(bldr, steppings[i])
		InfoAddMicrocode(bldr, microcodes[i])
		InfoAddCPUMHz(bldr, cpuMHzs[i])
		InfoAddCacheSize(bldr, cacheSizes[i])
		InfoAddPhysicalID(bldr, p.Infos[i].PhysicalID)
		InfoAddSiblings(bldr, p.Infos[i].Siblings)
		InfoAddCoreID(bldr, p.Infos[i].CoreID)
		InfoAddCPUCores(bldr, p.Infos[i].CPUCores)
		InfoAddApicID(bldr, p.Infos[i].ApicID)
		InfoAddInitialApicID(bldr, p.Infos[i].InitialApicID)
		InfoAddFPU(bldr, fpus[i])
		InfoAddFPUException(bldr, fpuExceptions[i])
		InfoAddCPUIDLevel(bldr, cpuIDLevels[i])
		InfoAddWP(bldr, wps[i])
		InfoAddFlags(bldr, flags[i])
		InfoAddBogoMIPS(bldr, p.Infos[i].BogoMIPS)
		InfoAddCLFlushSize(bldr, clFlushSizes[i])
		InfoAddCacheAlignment(bldr, cacheAlignments[i])
		InfoAddAddressSizes(bldr, addressSizes[i])
		InfoAddPowerManagement(bldr, powerManagements[i])
		infos[i] = InfoEnd(bldr)
	}
	// Process the Info vector
	ProcsStartInfosVector(bldr, len(infos))
	for i := len(p.Infos) - 1; i >= 0; i-- {
		bldr.PrependUOffsetT(infos[i])
	}
	infosV := bldr.EndVector(len(infos))
	ProcsStart(bldr)
	ProcsAddTimestamp(bldr, p.Timestamp)
	ProcsAddInfos(bldr, infosV)
	bldr.Finish(ProcsEnd(bldr))
	return bldr.Bytes[bldr.Head():]
}

func Deserialize(p []byte) *Processors {
	procs := GetRootAsProcs(p, 0)
	iLen := procs.InfosLength()
	processors := &Processors{}
	info := &Info{}
	inf := Inf{}
	processors.Timestamp = procs.Timestamp()
	for i := 0; i < iLen; i++ {
		if !procs.Infos(info, i) {
			continue
		}
		inf.Processor = info.Processor()
		inf.VendorID = string(info.VendorID())
		inf.CPUFamily = string(info.CPUFamily())
		inf.Model = string(info.Model())
		inf.ModelName = string(info.ModelName())
		inf.Stepping = string(info.Stepping())
		inf.Microcode = string(info.Microcode())
		inf.CPUMHz = string(info.CPUMHz())
		inf.CacheSize = string(info.CacheSize())
		inf.PhysicalID = info.PhysicalID()
		inf.Siblings = info.Siblings()
		inf.CoreID = info.CoreID()
		inf.CPUCores = info.CPUCores()
		inf.ApicID = info.ApicID()
		inf.InitialApicID = info.InitialApicID()
		inf.FPU = string(info.FPU())
		inf.FPUException = string(info.FPUException())
		inf.CPUIDLevel = string(info.CPUIDLevel())
		inf.WP = string(info.WP())
		inf.Flags = string(info.Flags())
		inf.BogoMIPS = info.BogoMIPS()
		inf.CLFlushSize = string(info.CLFlushSize())
		inf.CacheAlignment = string(info.CacheAlignment())
		inf.AddressSizes = string(info.AddressSizes())
		inf.PowerManagement = string(info.PowerManagement())
		processors.Infos = append(processors.Infos, inf)
	}
	return processors
}

// NiHao gets the processor information from /proc/cpuinfo
func NiHao() (*Processors, error) {
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
	procs := Processors{Timestamp: t}
	var cpu Inf
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
		// check to see if this is info for a different processor
		if name == "processor" {
			if procCnt > 0 {
				procs.Infos = append(procs.Infos, cpu)
			}
			procCnt++
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpuinfo", Op: "nihao: processor", Err: err}
			}
			cpu = Inf{Processor: int16(i)}
			continue
		}
		fmt.Println(string(value))
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
			cpu.CPUMHz = value
			continue
		}
		if name == "cache size" {
			cpu.CacheSize = value
			continue
		}
		if name == "physical id" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpuinfo", Op: "nihao: physical id", Err: err}
			}
			cpu.PhysicalID = int16(i)
			continue
		}
		if name == "siblings" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpuinfo", Op: "nihao: siblings", Err: err}
			}
			cpu.Siblings = int16(i)
			continue
		}
		if name == "core id" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpuinfo", Op: "nihao: core id", Err: err}
			}
			cpu.CoreID = int16(i)
			continue
		}
		if name == "cpu cores" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpuinfo", Op: "nihao: cpu cores", Err: err}
			}
			cpu.CPUCores = int16(i)
			continue
		}
		if name == "apicid" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpuinfo", Op: "nihao: apicid", Err: err}
			}
			cpu.ApicID = int16(i)
			continue
		}
		if name == "initial apicid" {
			i, err = strconv.Atoi(value)
			if err != nil {
				return nil, joe.Error{Type: "cpuinfo", Op: "nihao: initial apicid", Err: err}
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
				return nil, joe.Error{Type: "cpuinfo", Op: "nihao: bogomips", Err: err}
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
	procs.Infos = append(procs.Infos, cpu)
	return &procs, nil
}
