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

package json

import (
	"testing"

	"github.com/mohae/joefriday/mem"
)

func TestGet(t *testing.T) {
	nf, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	info, err := Unmarshal(nf)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	if info.Timestamp == 0 {
		t.Error("expected timestamp to be a non-zero value; got 0")
	}
	if info.MemTotal == 0 {
		t.Error("expected mem_total to be a non-zero value; got 0")
	}
	t.Logf("%#v\n", info)
}

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkSerialize(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := New()
	v, _ := p.Prof.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Serialize(v)
	}
	_ = jsn
}

func BenchmarkMarshal(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := New()
	v, _ := p.Prof.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Marshal(v)
	}
	_ = jsn
}

var inf *mem.Info

func BenchmarkDeserialize(b *testing.B) {
	b.StopTimer()
	p, _ := New()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = Deserialize(tmp)
	}
	_ = inf
}

func BenchmarkUnmarshal(b *testing.B) {
	b.StartTimer()
	p, _ := New()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = Unmarshal(tmp)
	}
	_ = inf
}
