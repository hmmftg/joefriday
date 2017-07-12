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

package cpuutil

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	p, err := NewProfiler()
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
	checkCPUUtil("get", u, t)
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
			checkCPUUtil("ticker", v, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkCPUUtil(name string, u *CPUUtil, t *testing.T) {
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
			t.Errorf("%s: %d: expected ID to have a value, was empty", name, i)
		}
	}
}

func BenchmarkCPUUtil(b *testing.B) {
	var u *CPUUtil
	b.StopTimer()
	p, _ := NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}
