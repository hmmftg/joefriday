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
	tcpu := testinfo.NewTempSysDevicesSystemCPU()
	err := tcpu.Create()
	if err != nil {
		t.Fatalf("setting up cpux testing info: %s", err)
	}
	prof := &Profiler{Profiler: &cpux.Profiler{SystemCPUPath: tcpu.Dir, NumCPU: int(tcpu.CoresPerPhysicalPackage * tcpu.PhysicalPackageCount)}}
	p, err := prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus, err := Deserialize(p)
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
	
	p, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus, err = Deserialize(p)
	if err != nil {
		t.Error(err)
	}

	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tcpu.ValidateCPUX(cpus)
	if err != nil {
		t.Errorf("validate min/max: %s", err)
	}

	// cleanup for next
	err = tcpu.Clean(false)
	if err != nil {
		t.Error(err)
	}

multiSocket:
	// 2 sockets
	tcpu.PhysicalPackageCount = 2
	prof.NumCPU = int(tcpu.CoresPerPhysicalPackage * tcpu.PhysicalPackageCount)
	tcpu.Freq = true
	err = tcpu.Create()
	if err != nil {
		t.Error("setting up cpux testing info: %s", err)
		goto noFreq
	}
	
	p, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus, err = Deserialize(p)
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
		goto clean
	}
	
	p, err = prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus, err = Deserialize(p)
	if err != nil {
		t.Error(err)
	}

	js, _ = json.MarshalIndent(cpus, "", "\t")
	t.Log(string(js))
	// compare results with frequency
	err = tcpu.ValidateCPUX(cpus)
	if err != nil {
		t.Errorf("validate min/max: %s", err)
	}

clean:
	// cleanup everything
	err = tcpu.Clean(true)
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
