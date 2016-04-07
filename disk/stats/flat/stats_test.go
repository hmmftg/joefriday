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

package flat

import (
	"testing"
	"time"

	"github.com/mohae/joefriday/disk/structs"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	statsD := Deserialize(p)
	checkStats(statsD, t)
}

func TestGetTicker(t *testing.T) {
	results := make(chan []byte)
	errs := make(chan error)
	done := make(chan struct{})
	go Ticker(time.Duration(400)*time.Millisecond, results, done, errs)
	var x int
	for {
		if x > 0 {
			close(done)
			break
		}
		select {
		case b, ok := <-results:
			if !ok {
				break
			}
			s := Deserialize(b)
			checkStats(s, t)
			t.Logf("%#v\n", s)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
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
	var tmp []byte
	b.StopTimer()
	p, _ := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := New()
	st, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(st)
	}
	_ = tmp
}

var st *structs.Stats

func BenchmarkDeserialize(b *testing.B) {
	b.StopTimer()
	p, _ := New()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		st = Deserialize(tmp)
	}
	_ = st
}
