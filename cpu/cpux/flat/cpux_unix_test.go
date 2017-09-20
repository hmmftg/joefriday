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

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/cpu/cpux"
	"github.com/mohae/joefriday/testinfo"
)

func TestCPUX(t *testing.T) {
	// set up test stuff w/o freq
	tcpu := testinfo.NewTempSysDevicesSystemCPU()
	err := tcpu.Create()
	if err != nil {
		t.Fatalf("setting up cpux testing info: %s", err)
	}
	prof := &Profiler{Builder: fb.NewBuilder(0), Profiler: &cpux.Profiler{SystemCPUPath: tcpu.Dir, NumCPU: int(tcpu.CoresPerPhysicalPackage * tcpu.PhysicalPackageCount)}}
	p, err := prof.Get()
	if err != nil {
		t.Error(err)
	}
	cpus := Deserialize(p)

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
	cpus = Deserialize(p)

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
	cpus = Deserialize(p)

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
	cpus = Deserialize(p)

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
	var f []byte
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f = p.Serialize(v)
	}
	_ = f
}

func BenchmarkDeserialize(b *testing.B) {
	var cpus *cpux.CPUs
	p, _ := NewProfiler()
	cpusb, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpus = Deserialize(cpusb)
	}
	_ = cpus
}
