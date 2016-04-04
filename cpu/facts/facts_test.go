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

package facts

import "testing"

func TestFacts(t *testing.T) {
	facts, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if facts.Timestamp == 0 {
		t.Error("expected timestamp to have a nonzero value, it didn't")
	}
	if len(facts.CPU) == 0 {
		t.Error("Expected at least 1 CPU entrie, got none")
	}
	// spot check some vars
	for i, fact := range facts.CPU {
		if fact.VendorID == "" {
			t.Errorf("%d: expected a vendor id value; it was empty", i)
		}
		if fact.CPUCores == 0 {
			t.Errorf("%d: expected cpu cores to have a non-zero value; it was 0", i)
		}
		if fact.Flags == "" {
			t.Errorf("%d: expected flags to be not be empty; it was", i)
		}
	}
	t.Logf("%#v", facts)
}

var fact *Facts

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	p, _ := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fact, _ = p.Get()
	}
	_ = fact
}
