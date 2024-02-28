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

package diskusage

import (
	"testing"
	"time"

	"github.com/hmmftg/joefriday/disk/structs"
)

func TestGet(t *testing.T) {
	p, err := NewProfiler()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	time.Sleep(time.Duration(300) * time.Millisecond)
	st, err := p.Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	checkUsage("get", st, t)
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
			checkUsage("ticker", v, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkUsage(n string, s *structs.DiskUsage, t *testing.T) {
	if s.Timestamp == 0 {
		t.Errorf("%s: Timestamp: wanted non-zero value; got 0", n)
	}
	if s.TimeDelta == 0 {
		t.Errorf("%s: TimeDelta: wanted non-zero value; got 0", n)
	}
	if len(s.Device) == 0 {
		t.Errorf("%s: expected there to be devices; didn't get any", n)
	}
	for i := 0; i < len(s.Device); i++ {
		if s.Device[i].Major == 0 {
			t.Errorf("%s: Device %d: Major: wanted a non-zero value, was 0", n, i)
		}
		if s.Device[i].Name == "" {
			t.Errorf("%s: Device %d: Name: wanted a non-empty value; was empty", n, i)
		}
	}
	t.Logf("%#v\n", s)
}
