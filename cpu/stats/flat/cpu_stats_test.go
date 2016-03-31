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

	"github.com/mohae/joefriday/cpu/stats"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	statsD := Deserialize(p)
	if int16(stats.CLK_TCK) != statsD.ClkTck {
		t.Errorf("ClkTck: got %s; want %s", statsD.ClkTck, stats.CLK_TCK)
	}
	if statsD.Timestamp == 0 {
		t.Error("Timestamp: wanted non-zero value; got 0")
	}
	if statsD.Ctxt == 0 {
		t.Error("Ctxt: wanted non-zero value; got 0")
	}
	if statsD.BTime == 0 {
		t.Error("BTime: wanted non-zero value; got 0")
	}
	if statsD.Processes == 0 {
		t.Error("Processes: wanted non-zero value; got 0")
	}
	if len(statsD.CPU) < 2 {
		t.Errorf("expected stats for at least 2 CPU entries, got %d", len(statsD.CPU))
	}
	for i := 0; i < len(statsD.CPU); i++ {
		if statsD.CPU[i].ID == "" {
			t.Errorf("CPU %d: ID: wanted a non-empty value; was empty", i)
		}
		if statsD.CPU[i].User == 0 {
			t.Errorf("CPU %d: User: wanted a non-zero value, was 0", i)
		}
		if statsD.CPU[i].System == 0 {
			t.Errorf("CPU %d: System: wanted a non-xero value, was 0", i)
		}
	}
}
