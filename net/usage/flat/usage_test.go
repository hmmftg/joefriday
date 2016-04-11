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

	"github.com/mohae/joefriday/net/structs"
)

func TestGet(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	b, err := p.Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	u := Deserialize(b)
	t.Logf("%#v\n", u)
	checkUsage("get", u, t)
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
			checkUsage("ticker", u, t)
			t.Logf("%#v\n", u)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}

func checkUsage(n string, u *structs.Usage, t *testing.T) {
	if u.Timestamp == 0 {
		t.Errorf("%s: expected timestamp to be a non-zero value; was 0", n)
	}
	if u.TimeDelta == 0 {
		t.Errorf("%s: expected TimeDelta to be a non-zero value; was 0", n)
	}
	if len(u.Interfaces) == 0 {
		t.Error("%s: expected interfaces; got none", n)
		return
	}
	// check name
	for i, v := range u.Interfaces {
		if v.Name == "" {
			t.Errorf("%s: %d: expected inteface to have a name; was empty", n, i)
		}
	}
}
