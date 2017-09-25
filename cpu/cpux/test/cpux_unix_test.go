package cpux

import (
	"encoding/json"
	"testing"

	"github.com/mohae/joefriday/cpu/cpux"
	"github.com/mohae/joefriday/testinfo"
)

func TestCPUX(t *testing.T) {
	// set up test stuff w/o freq
	tSysFS := testinfo.NewTempSysFS()
	tSysFS.PhysicalPackageCount = 1
	tSysFS.CoresPerPhysicalPackage = 2
	tSysFS.ThreadsPerCore = 2
	err := tSysFS.Create()
	if err != nil {
		t.Fatalf("setting up cpux testing info: %s", err)
	}
	prof := &cpux.Profiler{NumCPU: int(tSysFS.CPUs())}
	prof.SysFSSystemPath(tSysFS.Dir)
	cpus, err := prof.Get()
	if err != nil {
		t.Error(err)
	}

	//compare results cpufreq
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
	js, _ := json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// cleanup for next
	err = tSysFS.Clean(false)
	if err != nil {
		t.Error(err)
	}

	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.Create()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto multiSocket
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tSysFS.Clean(false)
	if err != nil {
		t.Error(err)
	}

multiSocket:
	// 2 sockets
	tSysFS.PhysicalPackageCount = 2
	prof.NumCPU = int(tSysFS.CPUs())
	tSysFS.Freq = true
	err = tSysFS.Create()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto noFreq
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	//compare results cpufreq
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// cleanup for next
	err = tSysFS.Clean(false)
	if err != nil {
		t.Error(err)
	}

noFreq:
	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.Create()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto noOffline
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tSysFS.Clean(false)
	if err != nil {
		t.Error(err)
	}

	// no offline file
noOffline:
	tSysFS.OfflineFile = false
	// set up test stuff w/o freq
	err = tSysFS.Create()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto clean
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tSysFS.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

clean:
	// cleanup everything
	err = tSysFS.Clean(true)
	if err != nil {
		t.Error(err)
	}

}
