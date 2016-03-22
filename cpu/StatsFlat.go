// automatically generated, do not modify

package cpu

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type StatsFlat struct {
	_tab flatbuffers.Table
}

func GetRootAsStatsFlat(buf []byte, offset flatbuffers.UOffsetT) *StatsFlat {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &StatsFlat{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *StatsFlat) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *StatsFlat) ClkTck() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatsFlat) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatsFlat) Ctxt() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatsFlat) BTime() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatsFlat) Processes() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatsFlat) CPUs(obj *StatFlat, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(StatFlat)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *StatsFlat) CPUsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func StatsFlatStart(builder *flatbuffers.Builder) { builder.StartObject(6) }
func StatsFlatAddClkTck(builder *flatbuffers.Builder, ClkTck int16) { builder.PrependInt16Slot(0, ClkTck, 0) }
func StatsFlatAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(1, Timestamp, 0) }
func StatsFlatAddCtxt(builder *flatbuffers.Builder, Ctxt int64) { builder.PrependInt64Slot(2, Ctxt, 0) }
func StatsFlatAddBTime(builder *flatbuffers.Builder, BTime int64) { builder.PrependInt64Slot(3, BTime, 0) }
func StatsFlatAddProcesses(builder *flatbuffers.Builder, Processes int64) { builder.PrependInt64Slot(4, Processes, 0) }
func StatsFlatAddCPUs(builder *flatbuffers.Builder, CPUs flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(CPUs), 0) }
func StatsFlatStartCPUsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func StatsFlatEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
