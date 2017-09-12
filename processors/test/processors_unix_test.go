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
	"errors"
	"fmt"
	"testing"

	"github.com/mohae/joefriday"
	"github.com/mohae/joefriday/cpu/testinfo"
	"github.com/mohae/joefriday/processors"
)

func Testi75600u(t *testing.T) {
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
	prof, err := processors.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	prof.Procer = tProc
	prof.NumCPU = int(tCPU.CoresPerPhysicalPackage * tCPU.PhysicalPackageCount)
	prof.SystemCPUPath = tCPU.Dir

	// get the processor info.
	procs, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = ValidateI75600u(procs, tCPU)
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
	procs, err = prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = ValidateI75600u(procs, tCPU)
	if err != nil {
		t.Error(err)
	}

}

func ValidateI75600u(procs *processors.Processors, tCPU testinfo.TempSysDevicesSystemCPU) error {
	if procs.Timestamp <= 0 {
		return errors.New("timestamp: expected a value > 0")
	}
	if procs.Sockets != tCPU.PhysicalPackageCount {
		return fmt.Errorf("sockets: got %d; want %d", procs.Sockets, tCPU.PhysicalPackageCount)
	}

	if int(procs.Sockets) != len(procs.Socket) {
		return fmt.Errorf("socket: got %d; want %d", len(procs.Socket), procs.Sockets)
	}

	if procs.CPUs != int(tCPU.PhysicalPackageCount*tCPU.CoresPerPhysicalPackage) {
		return fmt.Errorf("CPUs: got %d; want %d", procs.CPUs, tCPU.PhysicalPackageCount*tCPU.CoresPerPhysicalPackage)
	}

	for i, proc := range procs.Socket {
		err := testinfo.ValidateI75600uProc(&proc, tCPU.Freq)
		if err != nil {
			return fmt.Errorf("%d: %s", i, err)
		}
	}
	return nil
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
	prof, err := processors.NewProfiler()
	if err != nil {
		t.Error(err)
		return
	}
	prof.Procer = tProc
	prof.NumCPU = int(tCPU.CoresPerPhysicalPackage * tCPU.PhysicalPackageCount)
	prof.SystemCPUPath = tCPU.Dir

	// get the processor info.
	procs, err := prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = ValidateXeonE52690(procs, tCPU)
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
	procs, err = prof.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	// Verify results
	t.Logf("%#v", procs)
	err = ValidateXeonE52690(procs, tCPU)
	if err != nil {
		t.Error(err)
	}

}

func ValidateXeonE52690(procs *processors.Processors, tCPU testinfo.TempSysDevicesSystemCPU) error {
	if procs.Timestamp <= 0 {
		return errors.New("timestamp: expected a value > 0")
	}
	if procs.Sockets != tCPU.PhysicalPackageCount {
		return fmt.Errorf("sockets: got %d; want %d", procs.Sockets, tCPU.PhysicalPackageCount)
	}

	if int(procs.Sockets) != len(procs.Socket) {
		return fmt.Errorf("socket: got %d; want %d", len(procs.Socket), procs.Sockets)
	}

	if procs.CPUs != int(tCPU.PhysicalPackageCount*tCPU.CoresPerPhysicalPackage) {
		return fmt.Errorf("CPUs: got %d; want %d", procs.CPUs, tCPU.PhysicalPackageCount*tCPU.CoresPerPhysicalPackage)
	}

	for i, proc := range procs.Socket {
		err := testinfo.ValidateXeonE52690Proc(&proc, tCPU.Freq)
		if err != nil {
			return fmt.Errorf("%d: %s", i, err)
		}
	}
	return nil
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
