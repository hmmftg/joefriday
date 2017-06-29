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

	info "github.com/mohae/joefriday/cpu/cpuinfo"
)

func TestGet(t *testing.T) {
	inf, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	infS, err := Unmarshal(inf)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	if infS.Timestamp == 0 {
		t.Error("expected timestamp to be a non-zero value; got 0")
	}
	if len(infS.CPU) == 0 {
		t.Error("expected CPUs to be a non-zero value; got 0")
	}
	for i, v := range infS.CPU {
		if v.VendorID == "" {
			t.Errorf("%d: expected vendor_id to have a value; it was empty", i)
		}
		if len(v.Flags) == 0 {
			t.Errorf("%d: expected flags to have values; it was empty", i)
		}
	}
	t.Logf("%#v\n", infS)
}

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkSerialize(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Serialize(v)
	}
	_ = jsn
}

func BenchmarkMarshal(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Marshal(v)
	}
	_ = jsn
}

func BenchmarkDeserialize(b *testing.B) {
	var inf *info.Info
	b.StopTimer()
	p, _ := NewProfiler()
	infB, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = Deserialize(infB)
	}
	_ = inf
}

func BenchmarkUnmarshal(b *testing.B) {
	var inf *info.Info
	b.StartTimer()
	p, _ := NewProfiler()
	infB, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = Unmarshal(infB)
	}
	_ = inf
}
