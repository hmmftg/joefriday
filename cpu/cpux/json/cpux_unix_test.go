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
	prof, err := NewProfiler()
	if err != nil {
		t.Fatalf("getting new profiler: %s", err)
	}
	prof.NumCPU = 4
	prof.SystemCPUPath = dir

	p, err := prof.Get()
	if err != nil {
		t.Error(err)
	}

	cpus, err := Deserialize(p)
	if err != nil {
		t.Error(err)
	}

	//compare results w/o cpufreq
	err = testinfo.ValidateCPUX(cpus, false)
	if err != nil {
		t.Error(err)
	}

	js, _ := json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

	// set up test stuff w freq
	dir, err = testinfo.TempSysDevicesSystemCPU(true)
	prof.SystemCPUPath = dir
	p, err = prof.Get()
	if err != nil {
		t.Error(err)
	}

	cpus, err = Deserialize(p)
	if err != nil {
		t.Error(err)
	}

	err = testinfo.ValidateCPUX(cpus, true)
	if err != nil {
		t.Errorf("validate min/max: %s", err)
	}
	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))

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
	var cpus *cpux.CPUs
	p, _ := NewProfiler()
	cpusb, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpus, _ = Deserialize(cpusb)
	}
	_ = cpus
}

func BenchmarkUnmarshal(b *testing.B) {
	var cpus *cpux.CPUs
	p, _ := NewProfiler()
	cpusb, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpus, _ = Unmarshal(cpusb)
	}
	_ = cpus
}
