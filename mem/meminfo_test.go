package mem

import (
	"fmt"
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
	fmt.Printf("%#v\n", inf)
}

func TestGetData(t *testing.T) {
	p, err := GetData()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	inf := Deserialize(p)
	// compare
	data := GetRootAsData(p, 0)
	if inf.Timestamp != data.Timestamp() {
		t.Errorf("got %d; want %d", inf.Timestamp, data.Timestamp())
	}
	if inf.MemTotal != int(data.MemTotal()) {
		t.Errorf("got %d; want %d", inf.MemTotal, data.MemTotal())
	}
	if inf.MemFree != int(data.MemFree()) {
		t.Errorf("got %d; want %d", inf.MemFree, data.MemFree())
	}
	if inf.MemAvailable != int(data.MemAvailable()) {
		t.Errorf("got %d; want %d", inf.MemAvailable, data.MemAvailable())
	}
	if inf.Buffers != int(data.Buffers()) {
		t.Errorf("got %d; want %d", inf.Buffers, data.Buffers())
	}
	if inf.Cached != int(data.Cached()) {
		t.Errorf("got %d; want %d", inf.Cached, data.Cached())
	}
	if inf.SwapCached != int(data.SwapCached()) {
		t.Errorf("got %d; want %d", inf.SwapCached, data.SwapCached())
	}
	if inf.Active != int(data.Active()) {
		t.Errorf("got %d; want %d", inf.Active, data.Active())
	}
	if inf.Inactive != int(data.Inactive()) {
		t.Errorf("got %d; want %d", inf.Inactive, data.Inactive())
	}
	if inf.MemTotal != int(data.MemTotal()) {
		t.Errorf("got %d; want %d", inf.MemTotal, data.MemTotal())
	}
	if inf.SwapTotal != int(data.SwapTotal()) {
		t.Errorf("got %d; want %d", inf.SwapTotal, data.SwapTotal())
	}
	if inf.SwapFree != int(data.SwapFree()) {
		t.Errorf("got %d; want %d", inf.SwapFree, data.SwapFree())
	}
	if inf.SwapFree != int(data.SwapFree()) {
		t.Errorf("got %d; want %d", inf.SwapFree, data.SwapFree())
	}
}
