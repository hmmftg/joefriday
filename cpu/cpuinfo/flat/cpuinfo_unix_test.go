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

import (
	"testing"

	inf "github.com/mohae/joefriday/cpu/cpuinfo"
)

func TestSerialize(t *testing.T) {
	cpus, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	cpusD := Deserialize(cpus)
	if cpusD.Timestamp == 0 {
		t.Error("timestamp: expected non-zero timestamp")
	}
	if len(cpusD.CPU) == 0 {
		t.Error("expected at least 1 cpu entries; got none")
	}
	for i := 0; i < len(cpusD.CPU); i++ {
		if cpusD.CPU[i].VendorID == "" {
			t.Errorf("%d: VendorID: expected Vendor ID to not be empty, it was", i)
		}
		if cpusD.CPU[i].Model == "" {
			t.Errorf("%d: Model: expected model to not be empty; it was", i)
		}
		if cpusD.CPU[i].CPUCores == 0 {
			t.Errorf("%d: CPUCores: expected non-zero value; was 0", i)
		}
		if len(cpusD.CPU[i].Flags) == 0 {
			t.Errorf("%d: Flags: expected some flags, none were found", i)
		}
	}
	_, err = Serialize(cpusD)
	if err != nil {
		t.Errorf("unexpected serialization error: %s", err)
		return
	}
}

func BenchmarkGet(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := NewProfiler()
	cpus, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(cpus)
	}
	_ = tmp
}

func BenchmarkDeserialize(b *testing.B) {
	var cpus *inf.CPUs
	b.StopTimer()
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cpus = Deserialize(tmp)
	}
	_ = cpus
}
