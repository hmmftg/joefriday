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

// TempSysDevicesSystemCPU handles the creation of a /sys/devices/system/cpu
// tree in a temp directory for testing purposes. When usage of the temp info
// is done, Clean() should be called to remove everything that was created by
// Create(). By default, the information will be created in its own temp
// directory within the system's temp dir. If the information is to be created
// in a specific location, set the Dir explicitly prior to calling Create. When
// the path is empty when Create is called, Dir will be populated with the path
// to the temp dir within shich the CPUx tree can be found.
//
// PhysicalPackageCount and CoresPerPhysicalPackage should be set if anything
// other than their default values are desired. Each of these values are
// required to be > 0. Create will not check to see if they are > 0; no CPUx
// directories will be created as the product of multiplying by 0 is 0.
//
// The Freq flag is true, cpufreq information will be written. This information
// is not available on all systems so tests should cover both the cpufreq path
// existing and not existing.
type TempSysDevicesSystemCPU struct {
	Dir                     string
	Freq                    bool
	PhysicalPackageCount    int32
	CoresPerPhysicalPackage int32
	ThreadsPerCore          int32
}

// NewTempSysDevicesSystemCPU returns a new TempSysDevicesSystemCPU set to use
// the system's temp dir, populate cpufreq information, and have 4 cores for a
// 1 socket system.
func NewTempSysDevicesSystemCPU() TempSysDevicesSystemCPU {
	return TempSysDevicesSystemCPU{Dir: "", Freq: true, PhysicalPackageCount: 1, CoresPerPhysicalPackage: 2, ThreadsPerCore: 2}
}

// returns the number of CPUs per configuration:
//   PhysicalPackageCount * CoresPerPhysicalPackage * ThreadsPerCore
func (t *TempSysDevicesSystemCPU) CPUs() int32 {
	return t.PhysicalPackageCount * t.CoresPerPhysicalPackage * t.ThreadsPerCore
}

// Create creates the tempdir and cpu info for /sys/devices/system/cpu tests.
// If Dir is an empty string, the information will be written to a randomly
// named subdir within the temp dir and Dir will be set to this path. The
// created dir will be prefixed with cpux. If Dir has a non-empty value, the
// cpuX information will be written to that directory.  The number of cpuX
// entries is the product of PhysicalPackageCount and CoresPerPhysicalPackage.
// If an error occurs that is returned along with an empty string.
func (t *TempSysDevicesSystemCPU) Create() (err error) {
	if t.Dir == "" {
		t.Dir, err = ioutil.TempDir("", "cpux")
		if err != nil {
			return err
		}
	}
	if t.Freq {
		// instead of always checking each cpuX dir for cpufreq information, cpuX
		// processing looks for the existence of the cpufreq path for determining
		// if cpu frequency information is available for processing
		err = os.MkdirAll(filepath.Join(t.Dir, cpux.CPUFreq), 0777)
		if err != nil {
			return err
		}
	}

	// add Possible information:
	err = ioutil.WriteFile(filepath.Join(t.Dir, "possible"), []byte(fmt.Sprintf("%s\n", t.Possible())), 0777)
	if err != nil {
		return err
	}
	// add online info; use the same value as possible.
	err = ioutil.WriteFile(filepath.Join(t.Dir, "online"), []byte(fmt.Sprintf("%s\n", t.Possible())), 0777)
	if err != nil {
		return err
	}

	var x int // tracks current cpu X value

	// Add CPU info for each physical package count
	for i := 0; i < int(t.PhysicalPackageCount); i++ {
		cpusPerSocket := int(t.CoresPerPhysicalPackage * t.ThreadsPerCore)
		for j := 0; j < cpusPerSocket; j++ {
			cpuX := fmt.Sprintf("cpu%d", x)
			x++
			// set the topology core id is in reverse order of cpuX
			tmp := filepath.Join(t.Dir, cpuX, topology)
			err = os.MkdirAll(tmp, 0777)
			if err != nil {
				goto cleanup
			}
			err = ioutil.WriteFile(filepath.Join(tmp, "core_id"), []byte(fmt.Sprintf("%d\n", cpusPerSocket-j)), 0777)
			if err != nil {
				goto cleanup
			}
			err = ioutil.WriteFile(filepath.Join(tmp, "physical_package_id"), []byte(fmt.Sprintf("%d\n", i)), 0777)
			if err != nil {
				goto cleanup
			}

			// cache entries
			for k := range cacheLevels {
				cD := filepath.Join(t.Dir, cpuX, "cache", fmt.Sprintf("index%d", k))
				err = os.MkdirAll(cD, 0777)
				if err != nil {
					goto cleanup
				}
				err = ioutil.WriteFile(filepath.Join(cD, "level"), []byte(cacheLevels[k]+"\n"), 0777)
				if err != nil {
					goto cleanup
				}
				err = ioutil.WriteFile(filepath.Join(cD, "type"), []byte(cacheTypes[k]+"\n"), 0777)
				if err != nil {
					goto cleanup
				}
				err = ioutil.WriteFile(filepath.Join(cD, "size"), []byte(cache[k]+"\n"), 0777)
				if err != nil {
					goto cleanup
				}

			}
			if !t.Freq {
				continue
			}
			tmp = filepath.Join(t.Dir, cpuX, cpux.CPUFreq)
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
	}
	return nil

cleanup:
	os.RemoveAll(t.Dir)
	return err
}

// ValidateCPUX verifies that the info in the struct for cpuX processing is
// consistent with the test data.
func (t *TempSysDevicesSystemCPU) ValidateCPUX(cpus *cpux.CPUs) error {
	if len(cpus.CPU) != int(t.CPUs()) {
		return fmt.Errorf("CPU: got %d; want %d", len(cpus.CPU), t.CPUs())
	}
	if cpus.Sockets != t.PhysicalPackageCount {
		return fmt.Errorf("Sockets: got %d; want %d", cpus.Sockets, t.PhysicalPackageCount)
	}
	if cpus.Possible != t.Possible() {
		return fmt.Errorf("possible: got %q; want %q", cpus.Possible, t.Possible())
	}

	if cpus.Online != t.Possible() {
		return fmt.Errorf("online: got %q; want %q", cpus.Online, t.Possible())
	}
	for i, cpu := range cpus.CPU {
		// find the core_id
		if cpu.CoreID < 0 || cpu.CoreID >= t.CPUs() {
			return fmt.Errorf("%d: core_id: got %d; want [0-%d]", i, cpu.CoreID, t.CPUs())
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
		if t.Freq {
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

// Clean cleans up the information that the struct created during Create. If TRUE
// is passed, the TempSysDevicesSystemCPU.Dir will also be deleted. This can
// also be used to clean up the directory so that Create can be re-run.
func (t *TempSysDevicesSystemCPU) Clean(delDir bool) error {
	if delDir {
		return os.RemoveAll(t.Dir)
	}
	// otherwise get the contents of t.Dir and delete that
	fis, err := ioutil.ReadDir(t.Dir)
	if err != nil {
		return fmt.Errorf("Clean %s: nothing deleted: %s", t.Dir, err)
	}
	for _, fi := range fis {
		p := filepath.Join(t.Dir, fi.Name())
		if fi.IsDir() {
			err = os.RemoveAll(p)
			if err != nil {
				return fmt.Errorf("Clean %s: not all files were deleted: %s", t.Dir, err)
			}
			continue
		}
		err = os.Remove(p)
		if err != nil {
			return fmt.Errorf("Clean %s: not all files were deleted: %s", t.Dir, err)
		}
	}
	return nil
}

// Possible generates the possible string
func (t *TempSysDevicesSystemCPU) Possible() string {
	return fmt.Sprintf("0-%d", (t.PhysicalPackageCount*t.CoresPerPhysicalPackage*t.ThreadsPerCore)-1)
}
