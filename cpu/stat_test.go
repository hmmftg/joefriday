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

package cpu

import (
	"fmt"
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

func TestGetStats(t *testing.T) {
	stats, err := GetStat()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if stats.ClkTck != int16(CLK_TCK) {
		t.Errorf("CLK_TCK: got %d; wanted %d", stats.ClkTck, CLK_TCK)
	}
	if stats.Ctxt == 0 {
		t.Error("ctck: expected non-zero value, got 0")
	}
	if stats.BTime == 0 {
		t.Error("Btime: expected non-zero value, got 0")
	}
	if stats.Processes == 0 {
		t.Error("Processes: expected non-zero value, got 0")
	}
	if len(stats.CPU) < 2 {
		t.Errorf("cpu: got %d, want at least 2", len(stats.CPU))
	}
	for i, v := range stats.CPU {
		if v.ID == "" {
			t.Errorf("%d: expected ID to have a value, was empty", i)
		}
		if v.System == 0 {
			t.Errorf("%d: expected System to be a non-zero value, got 0", i)
		}
	}
}

func TestUtilizationTicker(t *testing.T) {
	out := make(chan *Utilization)
	done := make(chan struct{})
	errs := make(chan error)

	ticker := time.NewTicker(time.Duration(2) * time.Second)
	defer ticker.Stop()
	defer close(errs)

	go UtilizationTicker(time.Second, out, done, errs)
testloop:
	for {
		select {
		case <-ticker.C:
			close(done)
			break testloop
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		case u := <-out:
			if u.Timestamp == 0 {
				t.Errorf("expected a timestamp got 0")
				continue
			}
			if len(u.CPU) == 0 {
				t.Errorf("expected CPU data, got none")
				continue
			}
			for i, v := range u.CPU {
				if v.ID == "" {
					t.Errorf("%d: expected CPU to have an ID, got none", i)
					continue
				}
				// only check IDLE: this may fail if on a really busy system
				// but Usage may fail on a non-busy system.
				if fmt.Sprintf("%.1f", v.Idle) == "0.0" {
					t.Errorf("%d: expected Idle to have a non-zero value; it wasn't", i)
				}
			}
		}
	}
}
