package net

import (
	"fmt"
	"testing"
)

func TestGetInfo(t *testing.T) {
	inf, err := GetInfo()
	if err != nil {
		t.Errorf("got %s, want nil", err)
		return
	}
	// test flatbuffers stuff
	infS := inf.Serialize()
	infD := Deserialize(infS)
	// compare
	if inf.Timestamp != infD.Timestamp {
		t.Errorf("got %d; want %d", inf.Timestamp, infD.Timestamp)
	}
	for i := 0; i < len(inf.Interfaces); i++ {
		if inf.Interfaces[i].RBytes != infD.Interfaces[i].RBytes {
			t.Errorf("%d: Rbytes: got %d; want %d", i, infD.Interfaces[i].RBytes, inf.Interfaces[i].RBytes)
		}
		if inf.Interfaces[i].RPackets != infD.Interfaces[i].RPackets {
			t.Errorf("%d: RPackets: got %d; want %d", i, infD.Interfaces[i].RPackets, inf.Interfaces[i].RPackets)
		}
		if inf.Interfaces[i].RErrs != infD.Interfaces[i].RErrs {
			t.Errorf("%d: RErrs: got %d; want %d", i, infD.Interfaces[i].RErrs, inf.Interfaces[i].RErrs)
		}
		if inf.Interfaces[i].RDrop != infD.Interfaces[i].RDrop {
			t.Errorf("%d: RDrop: got %d; want %d", i, infD.Interfaces[i].RDrop, inf.Interfaces[i].RDrop)
		}
		if inf.Interfaces[i].RFIFO != infD.Interfaces[i].RFIFO {
			t.Errorf("%d: RFIFO: got %d; want %d", i, infD.Interfaces[i].RFIFO, inf.Interfaces[i].RFIFO)
		}
		if inf.Interfaces[i].RFrame != infD.Interfaces[i].RFrame {
			t.Errorf("%d: RFrame: got %d; want %d", i, infD.Interfaces[i].RFrame, inf.Interfaces[i].RFrame)
		}
		if inf.Interfaces[i].RCompressed != infD.Interfaces[i].RCompressed {
			t.Errorf("%d: RCompressed: got %d; want %d", i, infD.Interfaces[i].RCompressed, inf.Interfaces[i].RCompressed)
		}
		if inf.Interfaces[i].RMulticast != infD.Interfaces[i].RMulticast {
			t.Errorf("%d: RMulticast: got %d; want %d", i, infD.Interfaces[i].RMulticast, inf.Interfaces[i].RMulticast)
		}
		if inf.Interfaces[i].TBytes != infD.Interfaces[i].TBytes {
			t.Errorf("%d: TBytes: got %d; want %d", i, infD.Interfaces[i].TBytes, inf.Interfaces[i].TBytes)
		}
		if inf.Interfaces[i].TPackets != infD.Interfaces[i].TPackets {
			t.Errorf("%d: TPackets: got %d; want %d", i, infD.Interfaces[i].TPackets, inf.Interfaces[i].TPackets)
		}
		if inf.Interfaces[i].TErrs != infD.Interfaces[i].TErrs {
			t.Errorf("%d: TErrs: got %d; want %d", i, infD.Interfaces[i].TErrs, inf.Interfaces[i].TErrs)
		}
		if inf.Interfaces[i].TDrop != infD.Interfaces[i].TDrop {
			t.Errorf("%d: TDrop: got %d; want %d", i, infD.Interfaces[i].TDrop, inf.Interfaces[i].TDrop)
		}
		if inf.Interfaces[i].TFIFO != infD.Interfaces[i].TFIFO {
			t.Errorf("%d: TFIFO: got %d; want %d", i, infD.Interfaces[i].TFIFO, inf.Interfaces[i].TFIFO)
		}
		if inf.Interfaces[i].TColls != infD.Interfaces[i].TColls {
			t.Errorf("%d: TColls: got %d; want %d", i, infD.Interfaces[i].TColls, inf.Interfaces[i].TColls)
		}
		if inf.Interfaces[i].TCarrier != infD.Interfaces[i].TCarrier {
			t.Errorf("%d: TCarrier: got %d; want %d", i, infD.Interfaces[i].TCarrier, inf.Interfaces[i].TCarrier)
		}
		if inf.Interfaces[i].TCompressed != infD.Interfaces[i].TCompressed {
			t.Errorf("%d: TCompressed: got %d; want %d", i, infD.Interfaces[i].TCompressed, inf.Interfaces[i].TCompressed)
		}
	}
	fmt.Println(infD)

}
