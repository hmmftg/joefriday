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

package stats

import (
	"testing"
	"time"

	"github.com/mohae/joefriday/disk/structs"
)

func TestGet(t *testing.T) {
	s, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	checkStats(s, t)
}

func TestGetTicker(t *testing.T) {
	results := make(chan *structs.Stats)
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
		case s, ok := <-results:
			if !ok {
				break
			}
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

var stts *structs.Stats

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	p, _ := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		stts, _ = p.Get()
	}
	_ = stts
}
