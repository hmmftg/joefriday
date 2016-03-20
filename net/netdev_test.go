package net

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
	fmt.Println("info:\n", inf.String())
	for _, v := range inf.Interfaces {
		fmt.Println("bytes", v.RCum.Bytes)
	}

}
