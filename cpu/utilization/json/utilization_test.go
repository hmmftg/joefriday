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
	checkUtilization(ut, t)
	t.Logf("%#v\n", ut)
}

func TestTicker(t *testing.T) {
	out := make(chan []byte)
	done := make(chan struct{})
	errs := make(chan error)
	go Ticker(time.Duration(400)*time.Millisecond, out, done, errs)
	for i := 0; i < 1; i++ {
		select {
		case u := <-out:
			ut, err := Unmarshal(u)
			if err != nil {
				t.Errorf("got %s, want nil", err)
				return
			}
			checkUtilization(ut, t)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	t.Logf("%#v\n", ut)
}

func checkUtilization(ut *utilization.Utilization, t *testing.T) {
	if ut.Timestamp == 0 {
		t.Error("expected timestamp to be a non-zero value; got 0")
	}
	if ut.BTimeDelta == 0 {
		t.Error("expected btime delta to be a non-zero value; got 0")
	}
	if len(ut.CPU) == 0 {
		t.Error("expected CPUs to be a non-zero value; got 0")
	}
	for i, v := range ut.CPU {
		if v.ID == "" {
			t.Errorf("%d: expected id to have a value; it was empty", i)
		}
	}
}

var ut *utilization.Utilization

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	p, _ := New()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkUnmarshal(b *testing.B) {
	p, _ := New()
	utB, _ := p.Get()
	for i := 0; i < b.N; i++ {
		ut, _ = Unmarshal(utB)
	}
	_ = ut
}
