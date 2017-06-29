// automatically generated by the FlatBuffers compiler, do not modify

package structs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Info struct {
	_tab flatbuffers.Table
}

func GetRootAsInfo(buf []byte, offset flatbuffers.UOffsetT) *Info {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Info{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Info) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Info) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) CPUs(obj *CPU, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(CPU)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Info) CPUsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func InfoStart(builder *flatbuffers.Builder) { builder.StartObject(2) }
func InfoAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func InfoAddCPUs(builder *flatbuffers.Builder, CPUs flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(CPUs), 0) }
func InfoStartCPUsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func InfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }