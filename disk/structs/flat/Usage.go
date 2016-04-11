// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Usage struct {
	_tab flatbuffers.Table
}

func GetRootAsUsage(buf []byte, offset flatbuffers.UOffsetT) *Usage {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Usage{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Usage) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Usage) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Usage) TimeDelta() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Usage) Devices(obj *Device, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
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

func (rcv *Usage) DevicesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func UsageStart(builder *flatbuffers.Builder) { builder.StartObject(3) }
func UsageAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func UsageAddTimeDelta(builder *flatbuffers.Builder, TimeDelta int64) { builder.PrependInt64Slot(1, TimeDelta, 0) }
func UsageAddDevices(builder *flatbuffers.Builder, Devices flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(Devices), 0) }
func UsageStartDevicesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func UsageEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
