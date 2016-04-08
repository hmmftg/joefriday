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

package flat

import (
	"testing"
	"time"

	"github.com/mohae/joefriday/disk/structs"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	u := Deserialize(p)
	checkUsage(u, t)
}

func TestGetTicker(t *testing.T) {
	results := make(chan []byte)
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
		case b, ok := <-results:
			if !ok {
				break
			}
			u := Deserialize(b)
			checkUsage(u, t)
			t.Logf("%#v\n", u)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}

func checkUsage(u *structs.Usage, t *testing.T) {
	if u.Timestamp == 0 {
		t.Error("Timestamp: wanted non-zero value; got 0")
	}
	if u.TimeDelta == 0 {
		t.Error("TimeDelta: wanted non-zero value; got 0")
	}
	if len(u.Devices) == 0 {
		t.Errorf("expected there to be devices; didn't get any")
	}
	for i := 0; i < len(u.Devices); i++ {
		if u.Devices[i].Major == 0 {
			t.Errorf("Device %d: Major: wanted a non-zero value, was 0", i)
		}
		if u.Devices[i].Name == "" {
			t.Errorf("Device %d: Name: wanted a non-empty value; was empty", i)
		}
	}
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := New()
	u, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(u)
	}
	_ = tmp
}

var u *structs.Usage

func BenchmarkDeserialize(b *testing.B) {
	b.StopTimer()
	p, _ := New()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u = Deserialize(tmp)
	}
	_ = u
}
