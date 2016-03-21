// automatically generated, do not modify

package cpu

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Procs struct {
	_tab flatbuffers.Table
}

func GetRootAsProcs(buf []byte, offset flatbuffers.UOffsetT) *Procs {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Procs{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Procs) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Procs) Infos(obj *Info, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(Info)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Procs) InfosLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ProcsStart(builder *flatbuffers.Builder) { builder.StartObject(1) }
func ProcsAddInfos(builder *flatbuffers.Builder, Infos flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(Infos), 0) }
func ProcsStartInfosVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func ProcsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
