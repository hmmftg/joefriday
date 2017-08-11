package cpux

import (
	"encoding/json"
	"testing"

	"github.com/mohae/joefriday/cpu/cpux"
	"github.com/mohae/joefriday/cpu/testinfo"
)

func TestCPUX(t *testing.T) {
	// set up test stuff w/o freq
	dir, err := testinfo.TempSysDevicesSystemCPU(false)
	if err != nil {
		t.Fatalf("setting up tempdir: %s", err)
	}
	prof := &cpux.Profiler{SystemCPUPath: dir, NumCPU: 4}
	cpus, err := prof.Get()
	if err != nil {
		t.Error(err)
	}
	//compare results w/o cpufreq
	err = testinfo.ValidateCPUX(cpus, false)
	if err != nil {
		t.Error(err)
	}
	// set up test stuff w freq
	dir, err = testinfo.TempSysDevicesSystemCPU(true)
	prof.SystemCPUPath = dir
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	js, _ := json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = testinfo.ValidateCPUX(cpus, true)
	if err != nil {
		t.Errorf("validate min/max: %s", err)
	}
}
