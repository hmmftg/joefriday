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

package cpuinfo

import "testing"

func TestFacts(t *testing.T) {
	cpus, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if cpus.Timestamp == 0 {
		t.Error("expected timestamp to have a nonzero value, it didn't")
	}
	if len(cpus.CPU) == 0 {
		t.Error("Expected at least 1 CPU entry, got none")
	}
	// spot check some vars
	for i, cpu := range cpus.CPU {
		if cpu.VendorID == "" {
			t.Errorf("%d: expected a vendor id value; it was empty", i)
		}
		if cpu.CPUCores == 0 {
			t.Errorf("%d: expected cpu cores to have a non-zero value; it was 0", i)
		}
		if len(cpu.Flags) == 0 {
			t.Errorf("%d: expected flags to be not be empty; it was", i)
		}
	}
	t.Logf("%#v", cpus)
}

func BenchmarkGet(b *testing.B) {
	var cpus *CPUs
	b.StopTimer()
	p, _ := NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cpus, _ = p.Get()
	}
	_ = cpus
}
