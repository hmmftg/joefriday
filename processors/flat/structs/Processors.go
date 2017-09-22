// automatically generated by the FlatBuffers compiler, do not modify

package structs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Processors struct {
	_tab flatbuffers.Table
}

func GetRootAsProcessors(buf []byte, offset flatbuffers.UOffsetT) *Processors {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Processors{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Processors) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Processors) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Processors) Architecture() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) ByteOrder() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) Sockets() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Processors) CPUs() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Processors) Possible() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) Present() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) CoresPerSocket() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Processors) ThreadsPerCore() int8 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetInt8(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Processors) VendorID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) CPUFamily() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) Model() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) ModelName() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) Stepping() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) Microcode() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(32))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) CPUMHz() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(34))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Processors) MHzMin() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(36))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Processors) MHzMax() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(38))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Processors) CacheSize() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(40))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Processors) Cache(obj *CacheInf, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(42))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(CacheInf)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Processors) CacheLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(42))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Processors) BogoMIPS() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(44))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0.0
}

func (rcv *Processors) Flags(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(46))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j * 4))
	}
	return nil
}

func (rcv *Processors) FlagsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(46))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *Processors) OpModes(j int) []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(48))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.ByteVector(a + flatbuffers.UOffsetT(j * 4))
	}
	return nil
}

func (rcv *Processors) OpModesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(48))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func ProcessorsStart(builder *flatbuffers.Builder) { builder.StartObject(23) }
func ProcessorsAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func ProcessorsAddArchitecture(builder *flatbuffers.Builder, Architecture flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(Architecture), 0) }
func ProcessorsAddByteOrder(builder *flatbuffers.Builder, ByteOrder flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(ByteOrder), 0) }
func ProcessorsAddSockets(builder *flatbuffers.Builder, Sockets int32) { builder.PrependInt32Slot(3, Sockets, 0) }
func ProcessorsAddCPUs(builder *flatbuffers.Builder, CPUs int32) { builder.PrependInt32Slot(4, CPUs, 0) }
func ProcessorsAddPossible(builder *flatbuffers.Builder, Possible flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(Possible), 0) }
func ProcessorsAddPresent(builder *flatbuffers.Builder, Present flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(Present), 0) }
func ProcessorsAddCoresPerSocket(builder *flatbuffers.Builder, CoresPerSocket int16) { builder.PrependInt16Slot(7, CoresPerSocket, 0) }
func ProcessorsAddThreadsPerCore(builder *flatbuffers.Builder, ThreadsPerCore int8) { builder.PrependInt8Slot(8, ThreadsPerCore, 0) }
func ProcessorsAddVendorID(builder *flatbuffers.Builder, VendorID flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(9, flatbuffers.UOffsetT(VendorID), 0) }
func ProcessorsAddCPUFamily(builder *flatbuffers.Builder, CPUFamily flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(10, flatbuffers.UOffsetT(CPUFamily), 0) }
func ProcessorsAddModel(builder *flatbuffers.Builder, Model flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(11, flatbuffers.UOffsetT(Model), 0) }
func ProcessorsAddModelName(builder *flatbuffers.Builder, ModelName flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(12, flatbuffers.UOffsetT(ModelName), 0) }
func ProcessorsAddStepping(builder *flatbuffers.Builder, Stepping flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(13, flatbuffers.UOffsetT(Stepping), 0) }
func ProcessorsAddMicrocode(builder *flatbuffers.Builder, Microcode flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(14, flatbuffers.UOffsetT(Microcode), 0) }
func ProcessorsAddCPUMHz(builder *flatbuffers.Builder, CPUMHz float32) { builder.PrependFloat32Slot(15, CPUMHz, 0.0) }
func ProcessorsAddMHzMin(builder *flatbuffers.Builder, MHzMin float32) { builder.PrependFloat32Slot(16, MHzMin, 0.0) }
func ProcessorsAddMHzMax(builder *flatbuffers.Builder, MHzMax float32) { builder.PrependFloat32Slot(17, MHzMax, 0.0) }
func ProcessorsAddCacheSize(builder *flatbuffers.Builder, CacheSize flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(18, flatbuffers.UOffsetT(CacheSize), 0) }
func ProcessorsAddCache(builder *flatbuffers.Builder, Cache flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(19, flatbuffers.UOffsetT(Cache), 0) }
func ProcessorsStartCacheVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func ProcessorsAddBogoMIPS(builder *flatbuffers.Builder, BogoMIPS float32) { builder.PrependFloat32Slot(20, BogoMIPS, 0.0) }
func ProcessorsAddFlags(builder *flatbuffers.Builder, Flags flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(21, flatbuffers.UOffsetT(Flags), 0) }
func ProcessorsStartFlagsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func ProcessorsAddOpModes(builder *flatbuffers.Builder, OpModes flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(22, flatbuffers.UOffsetT(OpModes), 0) }
func ProcessorsStartOpModesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func ProcessorsEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
