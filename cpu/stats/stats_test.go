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
	checkStats(s, t)
}

func TestGetTicker(t *testing.T) {
	results := make(chan *Stats)
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

func checkStats(s *Stats, t *testing.T) {
	if int16(CLK_TCK) != s.ClkTck {
		t.Errorf("ClkTck: got %s; want %s", s.ClkTck, CLK_TCK)
	}
	if s.Timestamp == 0 {
		t.Error("Timestamp: wanted non-zero value; got 0")
	}
	if s.Ctxt == 0 {
		t.Error("Ctxt: wanted non-zero value; got 0")
	}
	if s.BTime == 0 {
		t.Error("BTime: wanted non-zero value; got 0")
	}
	if s.Processes == 0 {
		t.Error("Processes: wanted non-zero value; got 0")
	}
	if len(s.CPU) < 2 {
		t.Errorf("expected stats for at least 2 CPU entries, got %d", len(s.CPU))
	}
	for i := 0; i < len(s.CPU); i++ {
		if s.CPU[i].ID == "" {
			t.Errorf("CPU %d: ID: wanted a non-empty value; was empty", i)
		}
		if s.CPU[i].User == 0 {
			t.Errorf("CPU %d: User: wanted a non-zero value, was 0", i)
		}
		if s.CPU[i].System == 0 {
			t.Errorf("CPU %d: System: wanted a non-xero value, was 0", i)
		}
	}
}
