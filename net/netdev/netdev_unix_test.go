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

package netdev

import (
	"testing"
	"time"

	"github.com/hmmftg/joefriday/net/structs"
)

func TestGet(t *testing.T) {
	inf, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	checkInfo("get", inf, t)
	t.Logf("%#v\n", inf)
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
			checkInfo("ticker", v, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkInfo(n string, dev *structs.DevInfo, t *testing.T) {
	if dev.Timestamp == 0 {
		t.Errorf("%s: expected timestamp to be a non-zero value; was 0", n)
	}
	if len(dev.Device) == 0 {
		t.Errorf("%s: expected interfaces; got none", n)
		return
	}
	// check name
	for i, v := range dev.Device {
		if v.Name == "" {
			t.Errorf("%s: %d: expected inteface to have a name; was empty", n, i)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	var inf *structs.DevInfo
	p, _ := NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}
