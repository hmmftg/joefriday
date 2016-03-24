// automatically generated, do not modify

package cpu

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type FactsFB struct {
	_tab flatbuffers.Table
}

func GetRootAsFactsFB(buf []byte, offset flatbuffers.UOffsetT) *FactsFB {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &FactsFB{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *FactsFB) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *FactsFB) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactsFB) CPUs(obj *FactFB, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(FactFB)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *FactsFB) CPUsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func FactsFBStart(builder *flatbuffers.Builder) { builder.StartObject(2) }
func FactsFBAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func FactsFBAddCPUs(builder *flatbuffers.Builder, CPUs flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(CPUs), 0) }
func FactsFBStartCPUsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func FactsFBEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
