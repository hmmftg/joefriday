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

package utilization

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	time.Sleep(time.Duration(200) * time.Millisecond)
	u, err := p.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	checkUtilization("get", u, t)
}

func TestGetTicker(t *testing.T) {
	results := make(chan *Utilization)
	errs := make(chan error)
	done := make(chan struct{})
	go Ticker(time.Second, results, done, errs)
	var x int
	for {
		if x > 0 {
			close(done)
			break
		}
		select {
		case u, ok := <-results:
			if !ok {
				break
			}
			checkUtilization("ticker", u, t)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}

func checkUtilization(name string, u *Utilization, t *testing.T) {
	if u.Timestamp == 0 {
		t.Errorf("%s: timestamp: expected on-zero", name)
	}
	if u.TimeDelta == 0 {
		t.Errorf("%s: TimeDelta: expected non-zero value, got 0", name)
	}
	if u.CtxtDelta == 0 {
		t.Errorf("%s: CtxtDelta: expected non-zero value, got 0", name)
	}
	if u.BTimeDelta == 0 {
		t.Errorf("%s: BTimeDelta: expected non-zero value, got 0", name)
	}
	if u.Processes == 0 {
		t.Errorf("%s: Processes: expected non-zero value, got 0", name)
	}
	if len(u.CPU) < 2 {
		t.Errorf("%s: cpu: got %d, want at least 2", name, len(u.CPU))
	}
	for i, v := range u.CPU {
		if v.ID == "" {
			t.Errorf("%s: %d: expected ID to have a value, was empty", i, name)
		}
	}
}
