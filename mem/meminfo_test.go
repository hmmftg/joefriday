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

	"github.com/mohae/joefriday/mem/flat"
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

func TestGetInfoSerializeDeserializeFlat(t *testing.T) {
	p, err := GetInfoFlat()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	inf := DeserializeInfoFlat(p)
	// compare
	data := flat.GetRootAsInfo(p, 0)
	if inf.Timestamp != data.Timestamp() {
		t.Errorf("got %d; want %d", inf.Timestamp, data.Timestamp())
	}
	if inf.MemTotal != data.MemTotal() {
		t.Errorf("got %d; want %d", inf.MemTotal, data.MemTotal())
	}
	if inf.MemFree != data.MemFree() {
		t.Errorf("got %d; want %d", inf.MemFree, data.MemFree())
	}
	if inf.MemAvailable != data.MemAvailable() {
		t.Errorf("got %d; want %d", inf.MemAvailable, data.MemAvailable())
	}
	if inf.Buffers != data.Buffers() {
		t.Errorf("got %d; want %d", inf.Buffers, data.Buffers())
	}
	if inf.Cached != data.Cached() {
		t.Errorf("got %d; want %d", inf.Cached, data.Cached())
	}
	if inf.SwapCached != data.SwapCached() {
		t.Errorf("got %d; want %d", inf.SwapCached, data.SwapCached())
	}
	if inf.Active != data.Active() {
		t.Errorf("got %d; want %d", inf.Active, data.Active())
	}
	if inf.Inactive != data.Inactive() {
		t.Errorf("got %d; want %d", inf.Inactive, data.Inactive())
	}
	if inf.MemTotal != data.MemTotal() {
		t.Errorf("got %d; want %d", inf.MemTotal, data.MemTotal())
	}
	if inf.SwapTotal != data.SwapTotal() {
		t.Errorf("got %d; want %d", inf.SwapTotal, data.SwapTotal())
	}
	if inf.SwapFree != data.SwapFree() {
		t.Errorf("got %d; want %d", inf.SwapFree, data.SwapFree())
	}
	if inf.SwapFree != data.SwapFree() {
		t.Errorf("got %d; want %d", inf.SwapFree, data.SwapFree())
	}
}

var inf *Info

func BenchmarkGetMemInfo(b *testing.B) {
	p, _ := NewInfoProfiler()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func BenchmarkGetMemInfoFlat(b *testing.B) {
	var infF []byte
	p, _ := NewInfoProfiler()
	for i := 0; i < b.N; i++ {
		infF, _ = p.GetFlat()
	}
	_ = infF
}

func BenchmarkGetMemInfoJSON(b *testing.B) {
	var infF []byte
	p, _ := NewInfoProfiler()
	for i := 0; i < b.N; i++ {
		infF, _ = p.GetJSON()
	}
	_ = infF
}

func BenchmarkDeserializeInfoFlat(b *testing.B) {
	var inf *Info
	p, _ := NewInfoProfiler()
	infB, _ := p.GetFlat()
	for i := 0; i < b.N; i++ {
		inf = DeserializeInfoFlat(infB)
	}
	_ = inf
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
