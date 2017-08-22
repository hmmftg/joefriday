package testinfo

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mohae/joefriday/cpu/cpux"
)

const (
	topology = "topology"
)

var (
	cache       = []string{"32K", "32K", "256K", "6144K"}
	cacheIDs    = []string{"L1d cache", "L1i cache", "L2 cache", "L3 cache"}
	cacheTypes  = []string{"Data", "Instruction", "Unified", "Unified"}
	cacheLevels = []string{"1", "1", "2", "3"}
)

// TempSysDevicesSystemCPU creates the tempdir and cpu info for
// /sys/devices/system/cpu tests. There will be 4 cpuX entries, cpu0-cpu3. The
// temp directory path for the cpuX entries dir will be returned. If an error
// occurs that is returned along with empty an empty string.
//
// The freq parm controls whether or not frequency information is created. Some
// systems don't have frequency information.
func TempSysDevicesSystemCPU(freq bool) (dir string, err error) {
	dir, err = ioutil.TempDir("", "jfminmax")
	if err != nil {
		return "", err
	}
	if freq {
		err = os.MkdirAll(filepath.Join(dir, cpux.CPUFreq), 0777)
		if err != nil {
			return "", err
		}
	}
	// set 4 cpus
	for i := 0; i < 4; i++ {
		cpuX := fmt.Sprintf("cpu%d", i)
		// set the topology core id is in reverse order of cpuX
		tmp := filepath.Join(dir, cpuX, topology)
		err = os.MkdirAll(tmp, 0777)
		if err != nil {
			goto cleanup
		}
		err = ioutil.WriteFile(filepath.Join(tmp, "core_id"), []byte(fmt.Sprintf("%d\n", 3-i)), 0777)
		if err != nil {
			goto cleanup
		}
		// make all physical package ids 0; this means testing a multi-socket sytem isn't currently
		// implemented
		err = ioutil.WriteFile(filepath.Join(tmp, "physical_package_id"), []byte("0\n"), 0777)
		if err != nil {
			goto cleanup
		}

		// cache entries
		for j := range cacheLevels {
			cD := filepath.Join(dir, cpuX, "cache", fmt.Sprintf("index%d", j))
			err = os.MkdirAll(cD, 0777)
			if err != nil {
				goto cleanup
			}
			err = ioutil.WriteFile(filepath.Join(cD, "level"), []byte(cacheLevels[j]+"\n"), 0777)
			if err != nil {
				goto cleanup
			}
			err = ioutil.WriteFile(filepath.Join(cD, "type"), []byte(cacheTypes[j]+"\n"), 0777)
			if err != nil {
				goto cleanup
			}
			err = ioutil.WriteFile(filepath.Join(cD, "size"), []byte(cache[j]+"\n"), 0777)
			if err != nil {
				goto cleanup
			}

		}
		if !freq {
			continue
		}
		tmp = filepath.Join(dir, cpuX, cpux.CPUFreq)
		err = os.MkdirAll(tmp, 0777)
		if err != nil {
			goto cleanup
		}
		err = ioutil.WriteFile(filepath.Join(tmp, "cpuinfo_min_freq"), []byte("1600000\n"), 0777)
		if err != nil {
			goto cleanup

		}
		err = ioutil.WriteFile(filepath.Join(tmp, "cpuinfo_max_freq"), []byte("2800000\n"), 0777)
		if err != nil {
			goto cleanup

		}

	}
	return dir, nil

cleanup:
	os.RemoveAll(dir)
	return "", err
}

// ValidateCPUX verifies that the info in the struct for cpuX processing is
// consistent with the test data.
func ValidateCPUX(cpus *cpux.CPUs, freq bool) error {
	if len(cpus.CPU) != 4 {
		return fmt.Errorf("CPU: got %d; want 4", len(cpus.CPU))
	}
	if cpus.Sockets != 1 {
		return fmt.Errorf("Sockets: got %d; want 1", cpus.Sockets)
	}
	for i, cpu := range cpus.CPU {
		if int(cpu.PhysicalPackageID) != 0 {
			return fmt.Errorf("%d: physical package id: got %d; want 0", i, cpu.PhysicalPackageID)
		}
		// find the core_id
		if cpu.CoreID < 0 || cpu.CoreID > 3 {
			return fmt.Errorf("%d: core_id: got %d; want [0-3]", i, cpu.CoreID)
		}
		// get the cache info
		for i, v := range cpu.CacheIDs {
			if v != cacheIDs[i] {
				return fmt.Errorf("%d: got cache %s; want %s", i, v, cacheIDs[i])
			}
			c, ok := cpu.Cache[v]
			if !ok {
				return fmt.Errorf("%d: %s: expected it to exist in the cache map; it didn't", i, v)
			}
			if c != cache[i] {
				return fmt.Errorf("%d: %s: got %s; want %s", i, v, c, cache[i])
			}
		}
		if freq {
			if int(cpu.MHzMax) != 2800 {
				return fmt.Errorf("%d: MHzMax: want 2800.000; got %.3f", i, cpu.MHzMax)
			}
			if int(cpu.MHzMin) != 1600 {
				return fmt.Errorf("%d: MHzMin: want 1600.000; got %.3f", i, cpu.MHzMin)
			}
		} else {
			if int(cpu.MHzMax) != 0 {
				return fmt.Errorf("%d: MHzMax: want 0.000; got %.3f", i, cpu.MHzMax)
			}
			if int(cpu.MHzMin) != 0 {
				return fmt.Errorf("%d: MHzMin: want 0.000; got %.3f", i, cpu.MHzMin)
			}
		}
	}
	return nil
}
