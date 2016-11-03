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
)

func TestClkTck(t *testing.T) {
	err := ClkTck()
	if err != nil {
		t.Errorf("expected error to be nil; got %s", err)
	}
	if CLK_TCK == 0 {
		t.Errorf("got %d, want a value > 0", CLK_TCK)
	}
}

func TestGet(t *testing.T) {
	s, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	checkStats("get", s, t)
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
			checkStats("ticker", v, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkStats(n string, s *Stats, t *testing.T) {
	if int16(CLK_TCK) != s.ClkTck {
		t.Errorf("%s: ClkTck: got %s; want %s", n, s.ClkTck, CLK_TCK)
	}
	if s.Timestamp == 0 {
		t.Errorf("%s: Timestamp: wanted non-zero value; got 0", n)
	}
	if s.Ctxt == 0 {
		t.Errorf("%s: Ctxt: wanted non-zero value; got 0", n)
	}
	if s.BTime == 0 {
		t.Errorf("%s: BTime: wanted non-zero value; got 0", n)
	}
	if s.Processes == 0 {
		t.Errorf("%s: Processes: wanted non-zero value; got 0", n)
	}
	if len(s.CPU) < 2 {
		t.Errorf("%s: expected stats for at least 2 CPU entries, got %d", n, len(s.CPU))
	}
	for i := 0; i < len(s.CPU); i++ {
		if s.CPU[i].ID == "" {
			t.Errorf("%s: CPU %d: ID: wanted a non-empty value; was empty", n, i)
		}
		if s.CPU[i].User == 0 {
			t.Errorf("%s: CPU %d: User: wanted a non-zero value, was 0", n, i)
		}
		if s.CPU[i].System == 0 {
			t.Errorf("%s: CPU %d: System: wanted a non-xero value, was 0", n, i)
		}
	}
}

var stts *Stats

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	p, _ := NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		stts, _ = p.Get()
	}
	_ = stts
}
