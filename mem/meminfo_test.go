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

package mem

import (
	"reflect"
	"testing"
)

func TestGetInfo(t *testing.T) {
	inf, err := GetInfo()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	// just test to make sure the returned value != the zero value of Info.
	if reflect.DeepEqual(inf, Info{}) {
		t.Errorf("expected %v to not be equal to the zero value of Info, it was", inf)
	}
	t.Logf("%#v\n", inf)
}

var inf *Info

func BenchmarkGetMemInfo(b *testing.B) {
	p, _ := NewInfoProfiler()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func BenchmarkGetMemInfoJSON(b *testing.B) {
	var infF []byte
	p, _ := NewInfoProfiler()
	for i := 0; i < b.N; i++ {
		infF, _ = p.GetJSON()
	}
	_ = infF
}

func BenchmarkUnmarshalInfoJSON(b *testing.B) {
	var inf *Info
	p, _ := NewInfoProfiler()
	infB, _ := p.GetJSON()
	for i := 0; i < b.N; i++ {
		inf, _ = UnmarshalInfoJSON(infB)
	}
	_ = inf
}
