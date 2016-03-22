// automatically generated, do not modify

package platform

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type ReleaseFB struct {
	_tab flatbuffers.Table
}

func GetRootAsReleaseFB(buf []byte, offset flatbuffers.UOffsetT) *ReleaseFB {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &ReleaseFB{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *ReleaseFB) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ReleaseFB) ID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReleaseFB) IDLike() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReleaseFB) PrettyName() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReleaseFB) Version() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReleaseFB) VersionID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReleaseFB) HomeURL() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *ReleaseFB) BugReportURL() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func ReleaseFBStart(builder *flatbuffers.Builder) { builder.StartObject(7) }
func ReleaseFBAddID(builder *flatbuffers.Builder, ID flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(ID), 0) }
func ReleaseFBAddIDLike(builder *flatbuffers.Builder, IDLike flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(IDLike), 0) }
func ReleaseFBAddPrettyName(builder *flatbuffers.Builder, PrettyName flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(PrettyName), 0) }
func ReleaseFBAddVersion(builder *flatbuffers.Builder, Version flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(Version), 0) }
func ReleaseFBAddVersionID(builder *flatbuffers.Builder, VersionID flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(VersionID), 0) }
func ReleaseFBAddHomeURL(builder *flatbuffers.Builder, HomeURL flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(HomeURL), 0) }
func ReleaseFBAddBugReportURL(builder *flatbuffers.Builder, BugReportURL flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(BugReportURL), 0) }
func ReleaseFBEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
