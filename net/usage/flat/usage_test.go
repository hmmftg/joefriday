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

	"github.com/mohae/joefriday/net/info/flat"
	"github.com/mohae/joefriday/net/structs"
)

func TestGet(t *testing.T) {
	b, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	u := flat.Deserialize(b)
	t.Logf("%#v\n", u)
	checkInfo(u, t)
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
			u := flat.Deserialize(b)
			checkInfo(u, t)
			t.Logf("%#v\n", u)
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		}
		x++
	}
}

func checkInfo(inf *structs.Info, t *testing.T) {
	if inf.Timestamp == 0 {
		t.Errorf("expected timestamp to be a non-zero value; was 0")
	}
	if len(inf.Interfaces) == 0 {
		t.Error("expected interfaces; got none")
		return
	}
	// check name
	for i, v := range inf.Interfaces {
		if v.Name == "" {
			t.Errorf("%d: expected inteface to have a name; was empty", i)
		}
	}
}
