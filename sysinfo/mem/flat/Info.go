// automatically generated, do not modify

package flat

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

func (rcv *Info) TotalRAM() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) FreeRAM() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) SharedRAM() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) BufferRAM() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) TotalSwap() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) FreeSwap() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func InfoStart(builder *flatbuffers.Builder) { builder.StartObject(7) }
func InfoAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func InfoAddTotalRAM(builder *flatbuffers.Builder, TotalRAM uint64) { builder.PrependUint64Slot(1, TotalRAM, 0) }
func InfoAddFreeRAM(builder *flatbuffers.Builder, FreeRAM uint64) { builder.PrependUint64Slot(2, FreeRAM, 0) }
func InfoAddSharedRAM(builder *flatbuffers.Builder, SharedRAM uint64) { builder.PrependUint64Slot(3, SharedRAM, 0) }
func InfoAddBufferRAM(builder *flatbuffers.Builder, BufferRAM uint64) { builder.PrependUint64Slot(4, BufferRAM, 0) }
func InfoAddTotalSwap(builder *flatbuffers.Builder, TotalSwap uint64) { builder.PrependUint64Slot(5, TotalSwap, 0) }
func InfoAddFreeSwap(builder *flatbuffers.Builder, FreeSwap uint64) { builder.PrependUint64Slot(6, FreeSwap, 0) }
func InfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
