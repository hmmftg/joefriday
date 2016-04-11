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

package uptime

import "testing"

func TestGet(t *testing.T) {
	u, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if u.Total == 0 {
		t.Errorf("expected total to be a non-zero value; got 0")
	}
	if u.Idle == 0 {
		t.Errorf("expected idle to be a non-zero value; got 0")
	}
	t.Logf("%#v\n", u)
}

func BenchmarkGet(b *testing.B) {
	var u Uptime
	b.StopTimer()
	p, err := New()
	if err != nil {
		return
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}
