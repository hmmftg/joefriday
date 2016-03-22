// automatically generated, do not modify

package platform

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type KernelFB struct {
	_tab flatbuffers.Table
}

func GetRootAsKernelFB(buf []byte, offset flatbuffers.UOffsetT) *KernelFB {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &KernelFB{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *KernelFB) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *KernelFB) Version() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KernelFB) CompileUser() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KernelFB) GCC() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KernelFB) OSGCC() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KernelFB) Type() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KernelFB) CompileDate() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *KernelFB) Arch() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func KernelFBStart(builder *flatbuffers.Builder) { builder.StartObject(7) }
func KernelFBAddVersion(builder *flatbuffers.Builder, Version flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(Version), 0) }
func KernelFBAddCompileUser(builder *flatbuffers.Builder, CompileUser flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(CompileUser), 0) }
func KernelFBAddGCC(builder *flatbuffers.Builder, GCC flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(GCC), 0) }
func KernelFBAddOSGCC(builder *flatbuffers.Builder, OSGCC flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(OSGCC), 0) }
func KernelFBAddType(builder *flatbuffers.Builder, Type flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(Type), 0) }
func KernelFBAddCompileDate(builder *flatbuffers.Builder, CompileDate flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(CompileDate), 0) }
func KernelFBAddArch(builder *flatbuffers.Builder, Arch flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(Arch), 0) }
func KernelFBEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
