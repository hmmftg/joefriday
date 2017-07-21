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
	k, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if k.OS == "" {
		t.Error("OS: wanted a non-empty value; was empty")
	}
	if k.Version == "" {
		t.Error("Version: wanted a non-empty value; was empty")
	}
	if k.CompileUser == "" {
		t.Error("CompileUser: wanted a non-empty value; was empty")
	}
	if k.GCC == "" {
		t.Error("GCC: wanted a non-empty value; was empty")
	}
	if k.OSGCC == "" {
		t.Error("OSGCC: wanted a non-empty value; was empty")
	}
	if k.Type == "" {
		t.Error("Type: wanted a non-empty value; was empty")
	}
	if k.CompileDate == "" {
		t.Error("CompileDate: wanted a non-empty value; was empty")
	}
	if k.Arch == "" {
		t.Error("Arch: wanted a non-empty value; was empty")
	}
	t.Logf("%#v\n", k)
}

func BenchmarkGet(b *testing.B) {
	var k *Kernel
	p, err := NewProfiler()
	if err != nil {
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k, _ = p.Get()
	}
	_ = k
}
