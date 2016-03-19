// automatically generated, do not modify

package net

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Transmit struct {
	_tab flatbuffers.Table
}

func (rcv *Transmit) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Transmit) Bytes() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Transmit) Packets() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Transmit) Errs() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Transmit) Drop() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Transmit) FIFO() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Transmit) Colls() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Transmit) Carrier() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Transmit) Compressed() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func TransmitStart(builder *flatbuffers.Builder) { builder.StartObject(8) }
func TransmitAddBytes(builder *flatbuffers.Builder, Bytes int64) { builder.PrependInt64Slot(0, Bytes, 0) }
func TransmitAddPackets(builder *flatbuffers.Builder, Packets int64) { builder.PrependInt64Slot(1, Packets, 0) }
func TransmitAddErrs(builder *flatbuffers.Builder, Errs int64) { builder.PrependInt64Slot(2, Errs, 0) }
func TransmitAddDrop(builder *flatbuffers.Builder, Drop int64) { builder.PrependInt64Slot(3, Drop, 0) }
func TransmitAddFIFO(builder *flatbuffers.Builder, FIFO int64) { builder.PrependInt64Slot(4, FIFO, 0) }
func TransmitAddColls(builder *flatbuffers.Builder, Colls int64) { builder.PrependInt64Slot(5, Colls, 0) }
func TransmitAddCarrier(builder *flatbuffers.Builder, Carrier int64) { builder.PrependInt64Slot(6, Carrier, 0) }
func TransmitAddCompressed(builder *flatbuffers.Builder, Compressed int64) { builder.PrependInt64Slot(7, Compressed, 0) }
func TransmitEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
