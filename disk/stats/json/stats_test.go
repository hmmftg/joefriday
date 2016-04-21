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
	checkStats("get", sts, t)
	t.Logf("%#v\n", sts)
}

func TestTicker(t *testing.T) {
	tkr, err := NewTicker(time.Millisecond)
	if err != nil {
		t.Error(err)
		return
	}
	tk := tkr.(*Ticker)
	for i := 0; i < 5; i++ {
		select {
		case <-tk.Done:
			break
		case v, ok := <-tk.Data:
			if !ok {
				break
			}
			st, err := Deserialize(v)
			if err != nil {
				t.Error(err)
				continue
			}
			checkStats("ticker", st, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkStats(n string, s *structs.Stats, t *testing.T) {
	if s.Timestamp == 0 {
		t.Errorf("%s: Timestamp: wanted non-zero value; got 0", n)
	}
	if len(s.Devices) == 0 {
		t.Errorf("%s: expected there to be devices; didn't get any", n)
	}
	for i := 0; i < len(s.Devices); i++ {
		if s.Devices[i].Major == 0 {
			t.Errorf("%s: Device %d: Major: wanted a non-zero value, was 0", n, i)
		}
		if s.Devices[i].Name == "" {
			t.Errorf("%s: Device %d: Name: wanted a non-empty value; was empty", n, i)
		}
	}
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
	var st *structs.Stats
	b.StopTimer()
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		st, _ = Deserialize(tmp)
	}
	_ = st
}

func BenchmarkUnmarshal(b *testing.B) {
	var st *structs.Stats
	b.StartTimer()
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		st, _ = Unmarshal(tmp)
	}
	_ = st
}
