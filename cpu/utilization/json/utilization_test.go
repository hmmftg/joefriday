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

	"github.com/mohae/joefriday/cpu/utilization"
)

func TestGet(t *testing.T) {
	u, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	ut, err := Unmarshal(u)
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	checkUtilization("get", ut, t)
	t.Logf("%#v\n", ut)
}

func TestTicker(t *testing.T) {
	out := make(chan []byte)
	done := make(chan struct{})
	errs := make(chan error)
	go Ticker(time.Duration(400)*time.Millisecond, out, done, errs)
	var err error
	var ut *utilization.Utilization
	for i := 0; i < 1; i++ {
		select {
		case u := <-out:
			ut, err = Unmarshal(u)
			if err != nil {
				t.Errorf("got %s, want nil", err)
				return
			}
			checkUtilization("ticker", ut, t)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	t.Logf("%#v\n", ut)
}

func checkUtilization(name string, u *utilization.Utilization, t *testing.T) {
	if u.Timestamp == 0 {
		t.Errorf("%s: timestamp: expected on-zero", name)
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
