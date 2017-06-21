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

package uptime

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	u, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	checkUptime("get", u, t)
	t.Logf("%#v\n", u)
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
			checkUptime("ticker", v, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkUptime(n string, inf Info, t *testing.T) {
	if inf.Timestamp == 0 {
		t.Errorf("expected Timestamp to be a non-zero value; got 0")
	}
	if inf.Total == 0 {
		t.Errorf("expected total to be a non-zero value; got 0")
	}
	if inf.Idle == 0 {
		t.Errorf("expected idle to be a non-zero value; got 0")
	}
}

func BenchmarkGet(b *testing.B) {
	var inf Info
	b.StopTimer()
	p, err := NewProfiler()
	if err != nil {
		return
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}
