// automatically generated, do not modify

package flat

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

func (rcv *Stats) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Stats) Devices(obj *Device, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(Device)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Stats) DevicesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func StatsStart(builder *flatbuffers.Builder) { builder.StartObject(2) }
func StatsAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func StatsAddDevices(builder *flatbuffers.Builder, Devices flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(Devices), 0) }
func StatsStartDevicesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func StatsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
