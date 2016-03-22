// automatically generated, do not modify

package cpu

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type StatFlat struct {
	_tab flatbuffers.Table
}

func (rcv *StatFlat) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *StatFlat) CPU() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *StatFlat) User() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) Nice() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) System() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) Idle() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) IOWait() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) IRQ() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) SoftIRQ() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) Steal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) Quest() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *StatFlat) QuestNice() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func StatFlatStart(builder *flatbuffers.Builder) { builder.StartObject(11) }
func StatFlatAddCPU(builder *flatbuffers.Builder, CPU flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(CPU), 0) }
func StatFlatAddUser(builder *flatbuffers.Builder, User int64) { builder.PrependInt64Slot(1, User, 0) }
func StatFlatAddNice(builder *flatbuffers.Builder, Nice int64) { builder.PrependInt64Slot(2, Nice, 0) }
func StatFlatAddSystem(builder *flatbuffers.Builder, System int64) { builder.PrependInt64Slot(3, System, 0) }
func StatFlatAddIdle(builder *flatbuffers.Builder, Idle int64) { builder.PrependInt64Slot(4, Idle, 0) }
func StatFlatAddIOWait(builder *flatbuffers.Builder, IOWait int64) { builder.PrependInt64Slot(5, IOWait, 0) }
func StatFlatAddIRQ(builder *flatbuffers.Builder, IRQ int64) { builder.PrependInt64Slot(6, IRQ, 0) }
func StatFlatAddSoftIRQ(builder *flatbuffers.Builder, SoftIRQ int64) { builder.PrependInt64Slot(7, SoftIRQ, 0) }
func StatFlatAddSteal(builder *flatbuffers.Builder, Steal int64) { builder.PrependInt64Slot(8, Steal, 0) }
func StatFlatAddQuest(builder *flatbuffers.Builder, Quest int64) { builder.PrependInt64Slot(9, Quest, 0) }
func StatFlatAddQuestNice(builder *flatbuffers.Builder, QuestNice int64) { builder.PrependInt64Slot(10, QuestNice, 0) }
func StatFlatEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
