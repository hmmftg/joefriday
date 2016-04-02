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

package usage

import (
	"testing"
	"time"

	"github.com/mohae/joefriday/net/structs"
)

func TestGet(t *testing.T) {
	inf, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	checkInfo(inf, t)
	t.Logf("%#v\n", inf)
}

func TestTicker(t *testing.T) {
	results := make(chan *structs.Info)
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
		case inf, ok := <-results:
			if !ok {
				break
			}
			checkInfo(inf, t)
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
