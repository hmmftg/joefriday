// automatically generated by the FlatBuffers compiler, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type DiskUsage struct {
	_tab flatbuffers.Table
}

func GetRootAsDiskUsage(buf []byte, offset flatbuffers.UOffsetT) *DiskUsage {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &DiskUsage{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *DiskUsage) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *DiskUsage) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *DiskUsage) TimeDelta() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *DiskUsage) Devices(obj *Device, j int) bool {
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

func (rcv *DiskUsage) DevicesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func DiskUsageStart(builder *flatbuffers.Builder) { builder.StartObject(3) }
func DiskUsageAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func DiskUsageAddTimeDelta(builder *flatbuffers.Builder, TimeDelta int64) { builder.PrependInt64Slot(1, TimeDelta, 0) }
func DiskUsageAddDevices(builder *flatbuffers.Builder, Devices flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(Devices), 0) }
func DiskUsageStartDevicesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func DiskUsageEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }