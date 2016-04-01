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

func (rcv *Info) Interfaces(obj *Interface, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(Interface)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Info) InterfacesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func InfoStart(builder *flatbuffers.Builder) { builder.StartObject(2) }
func InfoAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func InfoAddInterfaces(builder *flatbuffers.Builder, Interfaces flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(Interfaces), 0) }
func InfoStartInterfacesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func InfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
