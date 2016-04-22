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

	"github.com/mohae/joefriday/platform/release"
)

func TestSerializeDeserialize(t *testing.T) {
	p, err := Get()
	if err != nil {
		t.Errorf("Get(): got %s, want nil", err)
		return
	}
	r, err := release.Get()
	if err != nil {
		t.Errorf("release.Get(): got %s, want nil", err)
		return
	}
	rD := Deserialize(p)
	if r.ID != rD.ID {
		t.Errorf("ID: got %s; want %s", rD.ID, r.ID)
	}
	if r.IDLike != rD.IDLike {
		t.Errorf("IDLike: got %s; want %s", rD.IDLike, r.IDLike)
	}
	if r.PrettyName != rD.PrettyName {
		t.Errorf("PrettyName: got %s; want %s", rD.PrettyName, r.PrettyName)
	}
	if r.Version != rD.Version {
		t.Errorf("Version: got %s; want %s", rD.Version, r.Version)
	}
	if r.VersionID != rD.VersionID {
		t.Errorf("VersionID: got %s; want %s", rD.VersionID, r.VersionID)
	}
	if r.HomeURL != rD.HomeURL {
		t.Errorf("HomeURL: got %s; want %s", rD.HomeURL, r.HomeURL)
	}
	if r.BugReportURL != rD.BugReportURL {
		t.Errorf("BugReportURL: got %s; want %s", rD.BugReportURL, r.BugReportURL)
	}
}

func BenchmarkGet(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func BenchmarkSerialize(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := New()
	k, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = Serialize(k)
	}
	_ = tmp
}

var r *release.Release

func BenchmarkDeserialize(b *testing.B) {
	b.StopTimer()
	p, _ := New()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r = Deserialize(tmp)
	}
	_ = r
}
