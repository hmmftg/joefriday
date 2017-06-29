// automatically generated by the FlatBuffers compiler, do not modify

package structs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Stats struct {
	_tab flatbuffers.Table
}

func GetRootAsStats(buf []byte, offset flatbuffers.UOffsetT) *Stats {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Stats{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Stats) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Stats) ClkTck() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Stats) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Stats) Ctxt() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Stats) BTime() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Stats) Processes() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Stats) CPUs(obj *CPU, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
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

func (rcv *Stats) CPUsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func StatsStart(builder *flatbuffers.Builder) { builder.StartObject(6) }
func StatsAddClkTck(builder *flatbuffers.Builder, ClkTck int16) { builder.PrependInt16Slot(0, ClkTck, 0) }
func StatsAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(1, Timestamp, 0) }
func StatsAddCtxt(builder *flatbuffers.Builder, Ctxt int64) { builder.PrependInt64Slot(2, Ctxt, 0) }
func StatsAddBTime(builder *flatbuffers.Builder, BTime int64) { builder.PrependInt64Slot(3, BTime, 0) }
func StatsAddProcesses(builder *flatbuffers.Builder, Processes int64) { builder.PrependInt64Slot(4, Processes, 0) }
func StatsAddCPUs(builder *flatbuffers.Builder, CPUs flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(CPUs), 0) }
func StatsStartCPUsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func StatsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }