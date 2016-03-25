package mem

import (
	"reflect"
	"testing"

	"github.com/EricLagergren/joefriday/mem/meminfo"
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

func TestGetData(t *testing.T) {
	p, err := GetInfoFlat()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	inf := DeserializeInfoFlat(p)
	// compare
	data := meminfo.GetRootAsInfoFlat(p, 0)
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

var inf Info

func BenchmarkReadMemInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inf, _ = GetInfo()
	}
	// b.Logf("%#v\n", inf)
}
