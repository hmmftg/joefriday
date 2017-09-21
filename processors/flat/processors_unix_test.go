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

package processors

import (
	"testing"

	"github.com/mohae/joefriday"
	ps "github.com/mohae/joefriday/processors"
	"github.com/mohae/joefriday/testinfo"
)

func Testi75600u(t *testing.T) {
	// set up the cpuinfo
	tProc, err := joefriday.NewTempFileProc("intel", "i75600", testinfo.I75600uCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()

	// get a new struct for /sys/devices/system/cpux info
	tCPU := testinfo.NewTempSysFSCPU()
	defer tCPU.Clean(true)

	tCPU.Freq = true
	tCPU.PhysicalPackageCount = 1
	tCPU.CoresPerPhysicalPackage = 2
	tCPU.ThreadsPerCore = 2
	// create the /sys/devices/system/cpux info
	err = tCPU.Create()
	if err != nil {
		t.Error(err)
		return
	}

	// get a new profiler and configure it
	prof, err := NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	prof.Procer = tProc
	prof.NumCPU = int(tCPU.CPUs())
	prof.SysFSCPUPath = tCPU.Dir

	// get the processor info.
	p, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	procs := Deserialize(p)

	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateI75600uProc(procs, tCPU.Freq)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tCPU.Clean(false)
	if err != nil {
		t.Error(err)
	}

	// set up test stuff w/o freq
	tCPU.Freq = false
	err = tCPU.Create()
	if err != nil {
		t.Error(err)
	}

	// get the processor info.
	p, err = prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procs = Deserialize(p)

	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateI75600uProc(procs, tCPU.Freq)
	if err != nil {
		t.Error(err)
	}

}

func TestXeonE52690(t *testing.T) {
	// set up the cpuinfo
	tProc, err := joefriday.NewTempFileProc("intel", "e52690", testinfo.XeonE52690CPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()

	// get a new struct for /sys/devices/system/cpux info
	tCPU := testinfo.NewTempSysFSCPU()
	defer tCPU.Clean(true)

	tCPU.Freq = true
	tCPU.PhysicalPackageCount = 2
	tCPU.CoresPerPhysicalPackage = 8
	tCPU.ThreadsPerCore = 2
	// create the /sys/devices/system/cpux info
	err = tCPU.Create()
	if err != nil {
		t.Error(err)
		return
	}

	// get a new profiler and configure it
	prof, err := NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	prof.Procer = tProc
	prof.NumCPU = int(tCPU.CPUs())
	prof.SysFSCPUPath = tCPU.Dir

	// get the processor info.
	p, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procs := Deserialize(p)

	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateXeonE52690Proc(procs, tCPU.Freq)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tCPU.Clean(false)
	if err != nil {
		t.Error(err)
	}

	// set up test stuff w/o freq
	tCPU.Freq = false
	err = tCPU.Create()
	if err != nil {
		t.Error(err)
	}

	// get the processor info.
	p, err = prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procs = Deserialize(p)

	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateXeonE52690Proc(procs, tCPU.Freq)
	if err != nil {
		t.Error(err)
	}

}

func TestR71800x(t *testing.T) {
	// set up the cpuinfo
	tProc, err := joefriday.NewTempFileProc("intel", "r71800x", testinfo.R71800xCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()

	// get a new struct for /sys/devices/system/cpux info
	tCPU := testinfo.NewTempSysFSCPU()
	defer tCPU.Clean(true)

	tCPU.Freq = true
	tCPU.PhysicalPackageCount = 1
	tCPU.CoresPerPhysicalPackage = 8
	tCPU.ThreadsPerCore = 2
	// create the /sys/devices/system/cpux info
	err = tCPU.Create()
	if err != nil {
		t.Error(err)
		return
	}

	// get a new profiler and configure it
	prof, err := NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	prof.Procer = tProc
	prof.NumCPU = int(tCPU.CPUs())
	prof.SysFSCPUPath = tCPU.Dir

	// get the processor info.
	p, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procs := Deserialize(p)

	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateR71800xProc(procs, tCPU.Freq)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tCPU.Clean(false)
	if err != nil {
		t.Error(err)
	}

	// set up test stuff w/o freq
	tCPU.Freq = false
	err = tCPU.Create()
	if err != nil {
		t.Error(err)
	}

	// get the processor info.
	p, err = prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procs = Deserialize(p)

	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateR71800xProc(procs, tCPU.Freq)
	if err != nil {
		t.Error(err)
	}

}

func BenchmarkGet(b *testing.B) {
	var tmp []byte
	p, _ := NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	p, _ := NewProfiler()
	proc, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(proc)
	}
	_ = tmp
}

func BenchmarkDeserialize(b *testing.B) {
	var proc *ps.Processors
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proc = Deserialize(tmp)
	}
	_ = proc
}
