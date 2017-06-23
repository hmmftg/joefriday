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

import "testing"

func TestProcessors(t *testing.T) {
	procs, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if procs.Timestamp == 0 {
		t.Error("expected timestamp to have a nonzero value, it didn't")
	}
	if procs.Count == 0 {
		t.Errorf("expected the processor count to be a nonzero value, it wasn't")
	}
	if len(procs.CPUs) == 0 {
		t.Error("Expected at least 1 chip entry, got none")
	}
	// spot check some vars
	for i, cpus := range procs.CPUs {
		if cpus.VendorID == "" {
			t.Errorf("%d: expected a vendor id value; it was empty", i)
		}
		if cpus.CPUCores == 0 {
			t.Errorf("%d: expected cpu cores to have a non-zero value; it was 0", i)
		}
		if len(cpus.Flags) == 0 {
			t.Errorf("%d: expected flags to be not be empty; it was", i)
		}
	}
	t.Logf("%#v", procs)
}

func BenchmarkGet(b *testing.B) {
	var procs *Processors
	p, _ := NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		procs, _ = p.Get()
	}
	_ = procs
}
