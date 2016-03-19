// automatically generated, do not modify

package net

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Receive struct {
	_tab flatbuffers.Table
}

func (rcv *Receive) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Receive) Bytes() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Receive) Packets() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Receive) Errs() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Receive) Drop() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Receive) FIFO() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Receive) Frame() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Receive) Compressed() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Receive) Multicast() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func ReceiveStart(builder *flatbuffers.Builder) { builder.StartObject(8) }
func ReceiveAddBytes(builder *flatbuffers.Builder, Bytes int64) { builder.PrependInt64Slot(0, Bytes, 0) }
func ReceiveAddPackets(builder *flatbuffers.Builder, Packets int64) { builder.PrependInt64Slot(1, Packets, 0) }
func ReceiveAddErrs(builder *flatbuffers.Builder, Errs int64) { builder.PrependInt64Slot(2, Errs, 0) }
func ReceiveAddDrop(builder *flatbuffers.Builder, Drop int64) { builder.PrependInt64Slot(3, Drop, 0) }
func ReceiveAddFIFO(builder *flatbuffers.Builder, FIFO int64) { builder.PrependInt64Slot(4, FIFO, 0) }
func ReceiveAddFrame(builder *flatbuffers.Builder, Frame int64) { builder.PrependInt64Slot(5, Frame, 0) }
func ReceiveAddCompressed(builder *flatbuffers.Builder, Compressed int64) { builder.PrependInt64Slot(6, Compressed, 0) }
func ReceiveAddMulticast(builder *flatbuffers.Builder, Multicast int64) { builder.PrependInt64Slot(7, Multicast, 0) }
func ReceiveEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
