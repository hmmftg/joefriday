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
