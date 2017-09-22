package cpux

import (
	"encoding/json"
	"testing"

	"github.com/mohae/joefriday/cpu/cpux"
	"github.com/mohae/joefriday/testinfo"
)

func TestCPUX(t *testing.T) {
	// set up test stuff w/o freq
	tcpu := testinfo.NewTempSysFSCPU()
	tcpu.PhysicalPackageCount = 1
	tcpu.CoresPerPhysicalPackage = 2
	tcpu.ThreadsPerCore = 2
	err := tcpu.Create()
	if err != nil {
		t.Fatalf("setting up cpux testing info: %s", err)
	}
	prof := &cpux.Profiler{SysFSCPUPath: tcpu.Dir, NumCPU: int(tcpu.CPUs())}
	cpus, err := prof.Get()
	if err != nil {
		t.Error(err)
	}

	//compare results cpufreq
	err = tcpu.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
	js, _ := json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// cleanup for next
	err = tcpu.Clean(false)
	if err != nil {
		t.Error(err)
	}

	// set up test stuff w/o freq
	tcpu.Freq = false
	err = tcpu.Create()
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
	err = tcpu.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tcpu.Clean(false)
	if err != nil {
		t.Error(err)
	}

multiSocket:
	// 2 sockets
	tcpu.PhysicalPackageCount = 2
	prof.NumCPU = int(tcpu.CPUs())
	tcpu.Freq = true
	err = tcpu.Create()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto noFreq
	}
	cpus, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	//compare results cpufreq
	err = tcpu.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// cleanup for next
	err = tcpu.Clean(false)
	if err != nil {
		t.Error(err)
	}

noFreq:
	// set up test stuff w/o freq
	tcpu.Freq = false
	err = tcpu.Create()
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
	err = tcpu.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tcpu.Clean(false)
	if err != nil {
		t.Error(err)
	}

	// no offline file
noOffline:
	tcpu.OfflineFile = false
	// set up test stuff w/o freq
	err = tcpu.Create()
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
	err = tcpu.ValidateCPUX(cpus)
	if err != nil {
		t.Error(err)
	}

clean:
	// cleanup everything
	err = tcpu.Clean(true)
	if err != nil {
		t.Error(err)
	}

}
