// automatically generated, do not modify

package mem

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type InfoFlat struct {
	_tab flatbuffers.Table
}

func GetRootAsInfoFlat(buf []byte, offset flatbuffers.UOffsetT) *InfoFlat {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &InfoFlat{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *InfoFlat) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *InfoFlat) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) MemTotal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) MemFree() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) MemAvailable() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) Buffers() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) Cached() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) SwapCached() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) Active() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) Inactive() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) SwapTotal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *InfoFlat) SwapFree() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func InfoFlatStart(builder *flatbuffers.Builder) { builder.StartObject(11) }
func InfoFlatAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func InfoFlatAddMemTotal(builder *flatbuffers.Builder, MemTotal int64) { builder.PrependInt64Slot(1, MemTotal, 0) }
func InfoFlatAddMemFree(builder *flatbuffers.Builder, MemFree int64) { builder.PrependInt64Slot(2, MemFree, 0) }
func InfoFlatAddMemAvailable(builder *flatbuffers.Builder, MemAvailable int64) { builder.PrependInt64Slot(3, MemAvailable, 0) }
func InfoFlatAddBuffers(builder *flatbuffers.Builder, Buffers int64) { builder.PrependInt64Slot(4, Buffers, 0) }
func InfoFlatAddCached(builder *flatbuffers.Builder, Cached int64) { builder.PrependInt64Slot(5, Cached, 0) }
func InfoFlatAddSwapCached(builder *flatbuffers.Builder, SwapCached int64) { builder.PrependInt64Slot(6, SwapCached, 0) }
func InfoFlatAddActive(builder *flatbuffers.Builder, Active int64) { builder.PrependInt64Slot(7, Active, 0) }
func InfoFlatAddInactive(builder *flatbuffers.Builder, Inactive int64) { builder.PrependInt64Slot(8, Inactive, 0) }
func InfoFlatAddSwapTotal(builder *flatbuffers.Builder, SwapTotal int64) { builder.PrependInt64Slot(9, SwapTotal, 0) }
func InfoFlatAddSwapFree(builder *flatbuffers.Builder, SwapFree int64) { builder.PrependInt64Slot(10, SwapFree, 0) }
func InfoFlatEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
