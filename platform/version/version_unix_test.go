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

package version

import "testing"

func TestGet(t *testing.T) {
	inf, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if inf.OS == "" {
		t.Error("OS: wanted a non-empty value; was empty")
	}
	if inf.Version == "" {
		t.Error("Version: wanted a non-empty value; was empty")
	}
	if inf.CompileUser == "" {
		t.Error("CompileUser: wanted a non-empty value; was empty")
	}
	if inf.GCC == "" {
		t.Error("GCC: wanted a non-empty value; was empty")
	}
	if inf.OSGCC == "" {
		t.Error("OSGCC: wanted a non-empty value; was empty")
	}
	if inf.Type == "" {
		t.Error("Type: wanted a non-empty value; was empty")
	}
	if inf.CompileDate == "" {
		t.Error("CompileDate: wanted a non-empty value; was empty")
	}
	if inf.Arch == "" {
		t.Error("Arch: wanted a non-empty value; was empty")
	}
	t.Logf("%#v\n", inf)
}

func BenchmarkGet(b *testing.B) {
	var inf *Info
	b.StopTimer()
	p, err := NewProfiler()
	if err != nil {
		return
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}
