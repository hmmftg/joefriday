// automatically generated, do not modify

package cpu

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type FactFlat struct {
	_tab flatbuffers.Table
}

func (rcv *FactFlat) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *FactFlat) Processor() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) VendorID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) CPUFamily() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) Model() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) ModelName() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) Stepping() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) Microcode() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) CPUMHz() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) CacheSize() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) PhysicalID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) Siblings() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) CoreID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) CPUCores() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) ApicID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) InitialApicID() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(32))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) FPU() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(34))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) FPUException() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(36))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) CPUIDLevel() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(38))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) WP() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(40))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) Flags() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(42))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) BogoMIPS() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(44))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FactFlat) CLFlushSize() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(46))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) CacheAlignment() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(48))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) AddressSizes() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(50))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FactFlat) PowerManagement() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(52))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func FactFlatStart(builder *flatbuffers.Builder) { builder.StartObject(25) }
func FactFlatAddProcessor(builder *flatbuffers.Builder, Processor int16) { builder.PrependInt16Slot(0, Processor, 0) }
func FactFlatAddVendorID(builder *flatbuffers.Builder, VendorID flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(VendorID), 0) }
func FactFlatAddCPUFamily(builder *flatbuffers.Builder, CPUFamily flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(CPUFamily), 0) }
func FactFlatAddModel(builder *flatbuffers.Builder, Model flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(Model), 0) }
func FactFlatAddModelName(builder *flatbuffers.Builder, ModelName flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(ModelName), 0) }
func FactFlatAddStepping(builder *flatbuffers.Builder, Stepping flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(5, flatbuffers.UOffsetT(Stepping), 0) }
func FactFlatAddMicrocode(builder *flatbuffers.Builder, Microcode flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(6, flatbuffers.UOffsetT(Microcode), 0) }
func FactFlatAddCPUMHz(builder *flatbuffers.Builder, CPUMHz float32) { builder.PrependFloat32Slot(7, CPUMHz, 0) }
func FactFlatAddCacheSize(builder *flatbuffers.Builder, CacheSize flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(8, flatbuffers.UOffsetT(CacheSize), 0) }
func FactFlatAddPhysicalID(builder *flatbuffers.Builder, PhysicalID int16) { builder.PrependInt16Slot(9, PhysicalID, 0) }
func FactFlatAddSiblings(builder *flatbuffers.Builder, Siblings int16) { builder.PrependInt16Slot(10, Siblings, 0) }
func FactFlatAddCoreID(builder *flatbuffers.Builder, CoreID int16) { builder.PrependInt16Slot(11, CoreID, 0) }
func FactFlatAddCPUCores(builder *flatbuffers.Builder, CPUCores int16) { builder.PrependInt16Slot(12, CPUCores, 0) }
func FactFlatAddApicID(builder *flatbuffers.Builder, ApicID int16) { builder.PrependInt16Slot(13, ApicID, 0) }
func FactFlatAddInitialApicID(builder *flatbuffers.Builder, InitialApicID int16) { builder.PrependInt16Slot(14, InitialApicID, 0) }
func FactFlatAddFPU(builder *flatbuffers.Builder, FPU flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(15, flatbuffers.UOffsetT(FPU), 0) }
func FactFlatAddFPUException(builder *flatbuffers.Builder, FPUException flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(16, flatbuffers.UOffsetT(FPUException), 0) }
func FactFlatAddCPUIDLevel(builder *flatbuffers.Builder, CPUIDLevel flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(17, flatbuffers.UOffsetT(CPUIDLevel), 0) }
func FactFlatAddWP(builder *flatbuffers.Builder, WP flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(18, flatbuffers.UOffsetT(WP), 0) }
func FactFlatAddFlags(builder *flatbuffers.Builder, Flags flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(19, flatbuffers.UOffsetT(Flags), 0) }
func FactFlatAddBogoMIPS(builder *flatbuffers.Builder, BogoMIPS float32) { builder.PrependFloat32Slot(20, BogoMIPS, 0) }
func FactFlatAddCLFlushSize(builder *flatbuffers.Builder, CLFlushSize flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(21, flatbuffers.UOffsetT(CLFlushSize), 0) }
func FactFlatAddCacheAlignment(builder *flatbuffers.Builder, CacheAlignment flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(22, flatbuffers.UOffsetT(CacheAlignment), 0) }
func FactFlatAddAddressSizes(builder *flatbuffers.Builder, AddressSizes flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(23, flatbuffers.UOffsetT(AddressSizes), 0) }
func FactFlatAddPowerManagement(builder *flatbuffers.Builder, PowerManagement flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(24, flatbuffers.UOffsetT(PowerManagement), 0) }
func FactFlatEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
