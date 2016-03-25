// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Utilization struct {
	_tab flatbuffers.Table
}

func GetRootAsUtilization(buf []byte, offset flatbuffers.UOffsetT) *Utilization {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Utilization{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Utilization) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Utilization) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Utilization) BTimeDelta() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Utilization) CtxtDelta() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Utilization) Processes() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Utilization) CPU(obj *Util, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(Util)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Utilization) CPULength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func UtilizationStart(builder *flatbuffers.Builder) { builder.StartObject(5) }
func UtilizationAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func UtilizationAddBTimeDelta(builder *flatbuffers.Builder, BTimeDelta int32) { builder.PrependInt32Slot(1, BTimeDelta, 0) }
func UtilizationAddCtxtDelta(builder *flatbuffers.Builder, CtxtDelta int64) { builder.PrependInt64Slot(2, CtxtDelta, 0) }
func UtilizationAddProcesses(builder *flatbuffers.Builder, Processes int32) { builder.PrependInt32Slot(3, Processes, 0) }
func UtilizationAddCPU(builder *flatbuffers.Builder, CPU flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(CPU), 0) }
func UtilizationStartCPUVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func UtilizationEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
