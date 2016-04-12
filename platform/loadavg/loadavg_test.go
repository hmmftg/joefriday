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

package loadavg

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	l, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	checkLoad("ticker", l, t)
	t.Logf("%#v\n", l)
}

func TestTicker(t *testing.T) {
	out := make(chan LoadAvg)
	done := make(chan struct{})
	errs := make(chan error)
	go Ticker(time.Duration(400)*time.Millisecond, out, done, errs)
	for i := 0; i < 1; i++ {
		select {
		case l := <-out:
			checkLoad("ticker", l, t)
			t.Logf("%#v\n", l)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	close(done)
}

func checkLoad(n string, l LoadAvg, t *testing.T) {
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
	var l LoadAvg
	b.StopTimer()
	p, err := New()
	if err != nil {
		return
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l, _ = p.Get()
	}
	_ = l
}
