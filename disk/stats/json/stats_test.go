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
	"time"

	"github.com/mohae/joefriday/disk/structs"
)

func TestGet(t *testing.T) {
	st, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	sts, err := Deserialize(st)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	checkStats(sts, t)
	t.Logf("%#v\n", sts)
}

func TestTicker(t *testing.T) {
	out := make(chan []byte)
	done := make(chan struct{})
	errs := make(chan error)
	go Ticker(time.Duration(400)*time.Millisecond, out, done, errs)

	for i := 0; i < 1; i++ {
		select {
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		case b := <-out:
			st, err := Deserialize(b)
			if err != nil {
				t.Errorf("unexpected deserialization error: %s", err)
				continue
			}
			checkStats(st, t)
			t.Logf("%#v\n", st)
		}
	}
}

func checkStats(s *structs.Stats, t *testing.T) {
	if s.Timestamp == 0 {
		t.Error("Timestamp: wanted non-zero value; got 0")
	}
	if len(s.Devices) == 0 {
		t.Errorf("expected there to be devices; didn't get any")
	}
	for i := 0; i < len(s.Devices); i++ {
		if s.Devices[i].Major == 0 {
			t.Errorf("Device %d: Major: wanted a non-zero value, was 0", i)
		}
		if s.Devices[i].Name == "" {
			t.Errorf("Device %d: Name: wanted a non-empty value; was empty", i)
		}
	}
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

var st *structs.Stats

func BenchmarkDeserialize(b *testing.B) {
	b.StopTimer()
	p, _ := New()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		st, _ = Deserialize(tmp)
	}
	_ = st
}

func BenchmarkUnmarshal(b *testing.B) {
	b.StartTimer()
	p, _ := New()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		st, _ = Unmarshal(tmp)
	}
	_ = st
}
