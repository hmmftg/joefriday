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

package release

import "testing"

func TestGet(t *testing.T) {
	r, err := Get()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if r.Name == "" {
		t.Error("Name: expected a value; was empty")
	}
	if r.ID == "" {
		t.Error("ID: expected a value; was empty")
	}
	if r.PrettyName == "" {
		t.Error("PrettyName: expected a value; was empty")
	}
	if r.Version == "" {
		t.Error("Version: expected a value; was empty")
	}
	if r.VersionID == "" {
		t.Error("VersionID: expected a value; was empty")
	}
	if r.HomeURL == "" {
		t.Error("HomeURL: expected a value; was empty")
	}
	if r.BugReportURL == "" {
		t.Error("BugReportURL: expected a value; was empty")
	}
	t.Logf("%#v\n", r)
}

func BenchmarkGet(b *testing.B) {
	var r *Release
	for i := 0; i < b.N; i++ {
		r, _ = Get()
	}
	_ = r
}
