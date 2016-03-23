// automatically generated, do not modify

package cpu

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type UtilizationFlat struct {
	_tab flatbuffers.Table
}

func GetRootAsUtilizationFlat(buf []byte, offset flatbuffers.UOffsetT) *UtilizationFlat {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &UtilizationFlat{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *UtilizationFlat) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *UtilizationFlat) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilizationFlat) BTimeDelta() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilizationFlat) CtxtDelta() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilizationFlat) Processes() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilizationFlat) CPUs(obj *UtilFlat, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(UtilFlat)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *UtilizationFlat) CPUsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func UtilizationFlatStart(builder *flatbuffers.Builder) { builder.StartObject(5) }
func UtilizationFlatAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func UtilizationFlatAddBTimeDelta(builder *flatbuffers.Builder, BTimeDelta int32) { builder.PrependInt32Slot(1, BTimeDelta, 0) }
func UtilizationFlatAddCtxtDelta(builder *flatbuffers.Builder, CtxtDelta int64) { builder.PrependInt64Slot(2, CtxtDelta, 0) }
func UtilizationFlatAddProcesses(builder *flatbuffers.Builder, Processes int32) { builder.PrependInt32Slot(3, Processes, 0) }
func UtilizationFlatAddCPUs(builder *flatbuffers.Builder, CPUs flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(CPUs), 0) }
func UtilizationFlatStartCPUsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func UtilizationFlatEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
