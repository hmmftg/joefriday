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

	procs "github.com/mohae/joefriday/processors"
)

func TestSerialize(t *testing.T) {
	proc, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procD := Deserialize(proc)
	if procD.Timestamp == 0 {
		t.Error("timestamp: expected non-zero timestamp")
	}
	if len(procD.Socket) == 0 {
		t.Error("expected at least 1 CPUs entry; got none")
	}
	for i := 0; i < len(procD.Socket); i++ {
		if procD.Socket[i].VendorID == "" {
			t.Errorf("%d: VendorID: expected Vendor ID to not be empty, it was", i)
		}
		if procD.Socket[i].Model == "" {
			t.Errorf("%d: Model: expected model to not be empty; it was", i)
		}
		if procD.Socket[i].CPUCores == 0 {
			t.Errorf("%d: CPUCores: expected non-zero value; was 0", i)
		}
		if int(procD.Socket[i].BogoMIPS) == 0 {
			t.Errorf("%d: expected a non-zero value for bogomips", i)
		}
	}
	_, err = Serialize(procD)
	if err != nil {
		t.Errorf("unexpected serialization error: %s", err)
		return
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
	var proc *procs.Processors
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		proc = Deserialize(tmp)
	}
	_ = proc
}
