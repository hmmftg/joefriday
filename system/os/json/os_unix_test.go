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

package os

import (
	"testing"

	o "github.com/hmmftg/joefriday/system/os"
)

func TestGet(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	os, err := o.Get()
	if err != nil {
		t.Errorf("release.Get(): got %s, want nil", err)
		return
	}
	osD, err := Deserialize(p)
	if err != nil {
		t.Errorf("deserialize: unexpected error: %s", err)
		return
	}
	if os.Name != osD.Name {
		t.Errorf("Name: got %s; want %s", osD.Name, os.Name)
	}
	if os.ID != osD.ID {
		t.Errorf("ID: got %s; want %s", osD.ID, os.ID)
	}
	if os.IDLike != osD.IDLike {
		t.Errorf("IDLike: got %s; want %s", osD.IDLike, os.IDLike)
	}
	if os.PrettyName != osD.PrettyName {
		t.Errorf("PrettyName: got %s; want %s", osD.PrettyName, os.PrettyName)
	}
	if os.Version != osD.Version {
		t.Errorf("Version: got %s; want %s", osD.Version, os.Version)
	}
	if os.VersionID != osD.VersionID {
		t.Errorf("VersionID: got %s; want %s", osD.VersionID, os.VersionID)
	}
	if os.HomeURL != osD.HomeURL {
		t.Errorf("HomeURL: got %s; want %s", osD.HomeURL, os.HomeURL)
	}
	if os.BugReportURL != osD.BugReportURL {
		t.Errorf("BugReportURL: got %s; want %s", osD.BugReportURL, os.BugReportURL)
	}
}

func BenchmarkGet(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Get()
	}
	_ = jsn
}

func BenchmarkSerialize(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Serialize(v)
	}
	_ = jsn
}

func BenchmarkMarshal(b *testing.B) {
	var jsn []byte
	p, _ := NewProfiler()
	v, _ := p.Profiler.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jsn, _ = p.Marshal(v)
	}
	_ = jsn
}

func BenchmarkDeserialize(b *testing.B) {
	var os *o.OS
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os, _ = Deserialize(tmp)
	}
	_ = os
}

func BenchmarkUnmarshal(b *testing.B) {
	var os *o.OS
	p, _ := NewProfiler()
	tmp, _ := p.Get()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os, _ = Unmarshal(tmp)
	}
	_ = os
}
