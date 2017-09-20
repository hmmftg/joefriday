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

func TestI75600u(t *testing.T) {
	// set up the cpuinfo
	tProc, err := joefriday.NewTempFileProc("intel", "i75600", testinfo.I75600uCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()

	// get a new struct for /sys/devices/system/cpux info
	tCPU := testinfo.NewTempSysDevicesSystemCPU()
	defer tCPU.Clean(true)

	tCPU.Freq = true
	tCPU.PhysicalPackageCount = 1
	tCPU.CoresPerPhysicalPackage = 4
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
	prof.NumCPU = int(tCPU.CoresPerPhysicalPackage * tCPU.PhysicalPackageCount)
	prof.SystemCPUPath = tCPU.Dir

	// get the processor info.
	p, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	procs, err := Unmarshal(p)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
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
	procs, err = Unmarshal(p)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}

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
	tCPU := testinfo.NewTempSysDevicesSystemCPU()
	defer tCPU.Clean(true)

	tCPU.Freq = true
	tCPU.PhysicalPackageCount = 2
	tCPU.CoresPerPhysicalPackage = 16
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
	prof.NumCPU = int(tCPU.CoresPerPhysicalPackage * tCPU.PhysicalPackageCount)
	prof.SystemCPUPath = tCPU.Dir

	// get the processor info.
	p, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procs, err := Unmarshal(p)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}

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
	procs, err = Unmarshal(p)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}

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
	tCPU := testinfo.NewTempSysDevicesSystemCPU()
	defer tCPU.Clean(true)

	tCPU.Freq = true
	tCPU.PhysicalPackageCount = 1
	tCPU.CoresPerPhysicalPackage = 8
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
	prof.NumCPU = int(tCPU.CoresPerPhysicalPackage * tCPU.PhysicalPackageCount)
	prof.SystemCPUPath = tCPU.Dir

	// get the processor info.
	p, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procs, err := Unmarshal(p)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}

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
	procs, err = Unmarshal(p)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}

	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateR71800xProc(procs, tCPU.Freq)
	if err != nil {
		t.Error(err)
	}

}

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkSerialize(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Serialize(v)
	}
	_ = jsn
}

func BenchmarkMarshal(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Marshal(v)
	}
	_ = jsn
}

func BenchmarkDeserialize(b *testing.B) {
	var proc *ps.Processors
	p, _ := NewProfiler()
	pB, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proc, _ = Deserialize(pB)
	}
	_ = proc
}

func BenchmarkUnmarshal(b *testing.B) {
	var proc *ps.Processors
	p, _ := NewProfiler()
	procB, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proc, _ = Unmarshal(procB)
	}
	_ = proc
}
