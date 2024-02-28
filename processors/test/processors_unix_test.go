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

	"github.com/hmmftg/joefriday"
	"github.com/hmmftg/joefriday/processors"
	"github.com/hmmftg/joefriday/testinfo"
)

func TestI75600u(t *testing.T) {
	// set up the cpuinfo
	tProc, err := joefriday.NewTempFileProc("intel", "i75600", testinfo.I75600uCPUInfo)
	if err != nil {
		t.Fatal(err)
	}
	defer tProc.Remove()

	// get a new struct for the sysfs tree
	tSysFS := testinfo.NewTempSysFS()
	err = tSysFS.SetSysFS("")
	if err != nil {
		t.Fatalf("setting up sysfs tree: %s", err)
	}
	defer tSysFS.Clean()

	tSysFS.Freq = true
	tSysFS.PhysicalPackageCount = 1
	tSysFS.CoresPerPhysicalPackage = 2
	tSysFS.ThreadsPerCore = 2
	// create the sysfs cpu tree
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error(err)
		return
	}

	// create the sysfs node tree
	err = tSysFS.CreateNode()
	if err != nil {
		t.Error(err)
		return
	}

	// get a new profiler and configure it
	prof, err := processors.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	prof.CPUProf.SysFSSystemPath(tSysFS.Path())
	prof.Procer = tProc
	prof.CPUProf.NumCPU = int(tSysFS.CPUs())
	prof.NodeProf.SysFSSystemPath(tSysFS.Path())

	// get the processor info.
	procs, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateI75600uProc(procs, tSysFS.Freq)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

	// clean up for next
	err = tSysFS.CleanNode()
	if err != nil {
		t.Error(err)
	}

	// create the sysfs node tree
	err = tSysFS.CreateNode()
	if err != nil {
		t.Error(err)
		return
	}

	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error(err)
	}

	// get the processor info.
	procs, err = prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateI75600uProc(procs, tSysFS.Freq)
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

	// get a new struct for the sysfs tree
	tSysFS := testinfo.NewTempSysFS()
	err = tSysFS.SetSysFS("")
	if err != nil {
		t.Fatalf("setting up sysfs tree: %s", err)
	}
	defer tSysFS.Clean()

	tSysFS.Freq = true
	tSysFS.PhysicalPackageCount = 2
	tSysFS.CoresPerPhysicalPackage = 8
	tSysFS.ThreadsPerCore = 2
	// create the sysfs cpu tree
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error(err)
		return
	}

	// create the sysfs node tree
	err = tSysFS.CreateNode()
	if err != nil {
		t.Error(err)
		return
	}

	// get a new profiler and configure it
	prof, err := processors.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	prof.Procer = tProc
	prof.CPUProf.NumCPU = int(tSysFS.CPUs())
	prof.CPUProf.SysFSSystemPath(tSysFS.Path())
	prof.NodeProf.SysFSSystemPath(tSysFS.Path())

	// get the processor info.
	procs, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateXeonE52690Proc(procs, tSysFS.Freq)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

	// clean up for next
	err = tSysFS.CleanNode()
	if err != nil {
		t.Error(err)
	}

	// create the sysfs node tree
	err = tSysFS.CreateNode()
	if err != nil {
		t.Error(err)
		return
	}

	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error(err)
	}

	// get the processor info.
	procs, err = prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateXeonE52690Proc(procs, tSysFS.Freq)
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

	// get a new struct for the sysfs tree
	tSysFS := testinfo.NewTempSysFS()
	err = tSysFS.SetSysFS("")
	if err != nil {
		t.Fatalf("setting up sysfs tree: %s", err)
	}
	defer tSysFS.Clean()

	tSysFS.Freq = true
	tSysFS.PhysicalPackageCount = 1
	tSysFS.CoresPerPhysicalPackage = 8
	tSysFS.ThreadsPerCore = 2
	// create the sysfs cpu tree
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error(err)
		return
	}

	// create the sysfs node tree
	err = tSysFS.CreateNode()
	if err != nil {
		t.Error(err)
		return
	}

	// get a new profiler and configure it
	prof, err := processors.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	prof.Procer = tProc
	prof.CPUProf.NumCPU = int(tSysFS.CPUs())
	prof.CPUProf.SysFSSystemPath(tSysFS.Path())
	prof.NodeProf.SysFSSystemPath(tSysFS.Path())

	// get the processor info.
	procs, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateR71800xProc(procs, tSysFS.Freq)
	if err != nil {
		t.Error(err)
	}

	// cleanup for next
	err = tSysFS.CleanCPU()
	if err != nil {
		t.Error(err)
	}

	// clean up for next
	err = tSysFS.CleanNode()
	if err != nil {
		t.Error(err)
	}

	// create the sysfs node tree
	err = tSysFS.CreateNode()
	if err != nil {
		t.Error(err)
		return
	}

	// set up test stuff w/o freq
	tSysFS.Freq = false
	err = tSysFS.CreateCPU()
	if err != nil {
		t.Error(err)
	}

	// get the processor info.
	procs, err = prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = testinfo.ValidateR71800xProc(procs, tSysFS.Freq)
	if err != nil {
		t.Error(err)
	}

}

func BenchmarkGet(b *testing.B) {
	var procs *processors.Processors
	p, _ := processors.NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		procs, _ = p.Get()
	}
	_ = procs
}
