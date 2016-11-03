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

	"github.com/mohae/joefriday/platform/uptime"
)

func TestGet(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("Get(): got %s, want nil", err)
		return
	}
	u, err := Deserialize(p)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	checkUptime("get", u, t)
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
			u, err := Deserialize(v)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
				continue
			}
			checkUptime("ticker", u, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkUptime(n string, u uptime.Uptime, t *testing.T) {
	if u.Timestamp == 0 {
		t.Errorf("expected Timestamp to be a non-zero value; got 0")
	}
	if u.Total == 0 {
		t.Errorf("expected total to be a non-zero value; got 0")
	}
	if u.Idle == 0 {
		t.Errorf("expected idle to be a non-zero value; got 0")
	}
}

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkSerialize(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Serialize(v)
	}
	_ = jsn
}

func BenchmarkMarshal(b *testing.B) {
	var jsn []byte
	b.StopTimer()
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Marshal(v)
	}
	_ = jsn
}

func BenchmarkDeserialize(b *testing.B) {
	var u uptime.Uptime
	b.StopTimer()
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = Deserialize(tmp)
	}
	_ = u
}

func BenchmarkUnmarshal(b *testing.B) {
	var u uptime.Uptime
	b.StartTimer()
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = Unmarshal(tmp)
	}
	_ = u
}
