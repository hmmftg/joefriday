// automatically generated, do not modify

package net

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type IFace struct {
	_tab flatbuffers.Table
}

func (rcv *IFace) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *IFace) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *IFace) RCum(obj *Receive) *Receive {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Receive)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *IFace) TCum(obj *Transmit) *Transmit {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(Transmit)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func IFaceStart(builder *flatbuffers.Builder) { builder.StartObject(3) }
func IFaceAddName(builder *flatbuffers.Builder, Name flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(Name), 0) }
func IFaceAddRCum(builder *flatbuffers.Builder, RCum flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(RCum), 0) }
func IFaceAddTCum(builder *flatbuffers.Builder, TCum flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(TCum), 0) }
func IFaceEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
