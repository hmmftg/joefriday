// automatically generated, do not modify

package cpu

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type UtilFlat struct {
	_tab flatbuffers.Table
}

func (rcv *UtilFlat) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *UtilFlat) CPU() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *UtilFlat) Usage() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilFlat) User() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilFlat) Nice() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilFlat) System() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilFlat) Idle() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *UtilFlat) IOWait() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func UtilFlatStart(builder *flatbuffers.Builder) { builder.StartObject(7) }
func UtilFlatAddCPU(builder *flatbuffers.Builder, CPU flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(CPU), 0) }
func UtilFlatAddUsage(builder *flatbuffers.Builder, Usage float32) { builder.PrependFloat32Slot(1, Usage, 0) }
func UtilFlatAddUser(builder *flatbuffers.Builder, User float32) { builder.PrependFloat32Slot(2, User, 0) }
func UtilFlatAddNice(builder *flatbuffers.Builder, Nice float32) { builder.PrependFloat32Slot(3, Nice, 0) }
func UtilFlatAddSystem(builder *flatbuffers.Builder, System float32) { builder.PrependFloat32Slot(4, System, 0) }
func UtilFlatAddIdle(builder *flatbuffers.Builder, Idle float32) { builder.PrependFloat32Slot(5, Idle, 0) }
func UtilFlatAddIOWait(builder *flatbuffers.Builder, IOWait float32) { builder.PrependFloat32Slot(6, IOWait, 0) }
func UtilFlatEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
