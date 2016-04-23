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

	"github.com/mohae/joefriday/platform/loadavg"
)

func TestGet(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("Get(): got %s, want nil", err)
		return
	}
	l := Deserialize(p)
	checkLoad("get", l, t)
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
			l := Deserialize(b)
			checkLoad("get", l, t)
			t.Logf("%#v\n", l)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}

func checkLoad(n string, l loadavg.LoadAvg, t *testing.T) {
	if l.LastMinute == 0 {
		t.Errorf("%s: expected LastMinute to be a non-zero value; got 0", n)
	}
	if l.LastFive == 0 {
		t.Errorf("%s: expected LastFive to be a non-zero value; got 0", n)
	}
	if l.LastTen == 0 {
		t.Errorf("%s: expected LastTen to be a non-zero value; got 0", n)
	}
	if l.RunningProcesses == 0 {
		t.Errorf("%s: expected RunningProcesses to be a non-zero value; got 0", n)
	}
	if l.TotalProcesses == 0 {
		t.Errorf("%s: expected TotalProcesses to be a non-zero value; got 0", n)
	}
	if l.PID == 0 {
		t.Errorf("%s: expected PID to be a non-zero value; got 0", n)
	}
}

func BenchmarkGet(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := NewProfiler()
	l, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(l)
	}
	_ = tmp
}

func BenchmarkDeserialize(b *testing.B) {
	var l loadavg.LoadAvg
	b.StopTimer()
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l = Deserialize(tmp)
	}
	_ = l
}
