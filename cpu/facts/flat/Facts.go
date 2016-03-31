// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Facts struct {
	_tab flatbuffers.Table
}

func GetRootAsFacts(buf []byte, offset flatbuffers.UOffsetT) *Facts {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Facts{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Facts) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Facts) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Facts) CPU(obj *Fact, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(Fact)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Facts) CPULength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func FactsStart(builder *flatbuffers.Builder) { builder.StartObject(2) }
func FactsAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func FactsAddCPU(builder *flatbuffers.Builder, CPU flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(CPU), 0) }
func FactsStartCPUVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func FactsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
