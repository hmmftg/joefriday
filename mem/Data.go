// automatically generated, do not modify

package mem

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Data struct {
	_tab flatbuffers.Table
}

func GetRootAsData(buf []byte, offset flatbuffers.UOffsetT) *Data {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Data{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Data) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Data) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) MemTotal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) MemFree() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) MemAvailable() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) Buffers() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) Cached() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) SwapCached() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) Active() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) Inactive() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) SwapTotal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) SwapFree() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func DataStart(builder *flatbuffers.Builder) { builder.StartObject(11) }
func DataAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func DataAddMemTotal(builder *flatbuffers.Builder, MemTotal int64) { builder.PrependInt64Slot(1, MemTotal, 0) }
func DataAddMemFree(builder *flatbuffers.Builder, MemFree int64) { builder.PrependInt64Slot(2, MemFree, 0) }
func DataAddMemAvailable(builder *flatbuffers.Builder, MemAvailable int64) { builder.PrependInt64Slot(3, MemAvailable, 0) }
func DataAddBuffers(builder *flatbuffers.Builder, Buffers int64) { builder.PrependInt64Slot(4, Buffers, 0) }
func DataAddCached(builder *flatbuffers.Builder, Cached int64) { builder.PrependInt64Slot(5, Cached, 0) }
func DataAddSwapCached(builder *flatbuffers.Builder, SwapCached int64) { builder.PrependInt64Slot(6, SwapCached, 0) }
func DataAddActive(builder *flatbuffers.Builder, Active int64) { builder.PrependInt64Slot(7, Active, 0) }
func DataAddInactive(builder *flatbuffers.Builder, Inactive int64) { builder.PrependInt64Slot(8, Inactive, 0) }
func DataAddSwapTotal(builder *flatbuffers.Builder, SwapTotal int64) { builder.PrependInt64Slot(9, SwapTotal, 0) }
func DataAddSwapFree(builder *flatbuffers.Builder, SwapFree int64) { builder.PrependInt64Slot(10, SwapFree, 0) }
func DataEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
