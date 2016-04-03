// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Fact struct {
	_tab flatbuffers.Table
}

func (rcv *Fact) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Fact) Processor() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) VendorID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) CPUFamily() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) Model() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) ModelName() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) Stepping() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) Microcode() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) CPUMHz() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) CacheSize() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) PhysicalID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) Siblings() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) CoreID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) CPUCores() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) ApicID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) InitialApicID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(32))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) FPU() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(34))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) FPUException() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(36))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) CPUIDLevel() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(38))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) WP() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(40))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) Flags() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(42))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) BogoMIPS() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(44))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Fact) CLFlushSize() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(46))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) CacheAlignment() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(48))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) AddressSizes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(50))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Fact) PowerManagement() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(52))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func FactStart(builder *flatbuffers.Builder) { builder.StartObject(25) }
func FactAddProcessor(builder *flatbuffers.Builder, Processor int16) { builder.PrependInt16Slot(0, Processor, 0) }
func FactAddVendorID(builder *flatbuffers.Builder, VendorID flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(VendorID), 0) }
func FactAddCPUFamily(builder *flatbuffers.Builder, CPUFamily flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(CPUFamily), 0) }
func FactAddModel(builder *flatbuffers.Builder, Model flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(Model), 0) }
func FactAddModelName(builder *flatbuffers.Builder, ModelName flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(ModelName), 0) }
func FactAddStepping(builder *flatbuffers.Builder, Stepping flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(Stepping), 0) }
func FactAddMicrocode(builder *flatbuffers.Builder, Microcode flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(Microcode), 0) }
func FactAddCPUMHz(builder *flatbuffers.Builder, CPUMHz float32) { builder.PrependFloat32Slot(7, CPUMHz, 0) }
func FactAddCacheSize(builder *flatbuffers.Builder, CacheSize flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(CacheSize), 0) }
func FactAddPhysicalID(builder *flatbuffers.Builder, PhysicalID int16) { builder.PrependInt16Slot(9, PhysicalID, 0) }
func FactAddSiblings(builder *flatbuffers.Builder, Siblings int16) { builder.PrependInt16Slot(10, Siblings, 0) }
func FactAddCoreID(builder *flatbuffers.Builder, CoreID int16) { builder.PrependInt16Slot(11, CoreID, 0) }
func FactAddCPUCores(builder *flatbuffers.Builder, CPUCores int16) { builder.PrependInt16Slot(12, CPUCores, 0) }
func FactAddApicID(builder *flatbuffers.Builder, ApicID int16) { builder.PrependInt16Slot(13, ApicID, 0) }
func FactAddInitialApicID(builder *flatbuffers.Builder, InitialApicID int16) { builder.PrependInt16Slot(14, InitialApicID, 0) }
func FactAddFPU(builder *flatbuffers.Builder, FPU flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(15, flatbuffers.UOffsetT(FPU), 0) }
func FactAddFPUException(builder *flatbuffers.Builder, FPUException flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(16, flatbuffers.UOffsetT(FPUException), 0) }
func FactAddCPUIDLevel(builder *flatbuffers.Builder, CPUIDLevel flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(17, flatbuffers.UOffsetT(CPUIDLevel), 0) }
func FactAddWP(builder *flatbuffers.Builder, WP flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(18, flatbuffers.UOffsetT(WP), 0) }
func FactAddFlags(builder *flatbuffers.Builder, Flags flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(19, flatbuffers.UOffsetT(Flags), 0) }
func FactAddBogoMIPS(builder *flatbuffers.Builder, BogoMIPS float32) { builder.PrependFloat32Slot(20, BogoMIPS, 0) }
func FactAddCLFlushSize(builder *flatbuffers.Builder, CLFlushSize flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(21, flatbuffers.UOffsetT(CLFlushSize), 0) }
func FactAddCacheAlignment(builder *flatbuffers.Builder, CacheAlignment flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(22, flatbuffers.UOffsetT(CacheAlignment), 0) }
func FactAddAddressSizes(builder *flatbuffers.Builder, AddressSizes flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(23, flatbuffers.UOffsetT(AddressSizes), 0) }
func FactAddPowerManagement(builder *flatbuffers.Builder, PowerManagement flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(24, flatbuffers.UOffsetT(PowerManagement), 0) }
func FactEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }