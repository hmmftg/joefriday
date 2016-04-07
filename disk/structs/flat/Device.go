// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Device struct {
	_tab flatbuffers.Table
}

func (rcv *Device) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Device) Major() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) Minor() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Device) ReadsCompleted() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) ReadsMerged() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) ReadSectors() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) ReadingTime() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) WritesCompleted() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) WritesMerged() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) WrittenSectors() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) WritingTime() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) IOInProgress() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) IOTime() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Device) WeightedIOTime() uint64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.GetUint64(o + rcv._tab.Pos)
	}
	return 0
}

func DeviceStart(builder *flatbuffers.Builder) { builder.StartObject(14) }
func DeviceAddMajor(builder *flatbuffers.Builder, Major uint32) { builder.PrependUint32Slot(0, Major, 0) }
func DeviceAddMinor(builder *flatbuffers.Builder, Minor uint32) { builder.PrependUint32Slot(1, Minor, 0) }
func DeviceAddName(builder *flatbuffers.Builder, Name flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(Name), 0) }
func DeviceAddReadsCompleted(builder *flatbuffers.Builder, ReadsCompleted uint64) { builder.PrependUint64Slot(3, ReadsCompleted, 0) }
func DeviceAddReadsMerged(builder *flatbuffers.Builder, ReadsMerged uint64) { builder.PrependUint64Slot(4, ReadsMerged, 0) }
func DeviceAddReadSectors(builder *flatbuffers.Builder, ReadSectors uint64) { builder.PrependUint64Slot(5, ReadSectors, 0) }
func DeviceAddReadingTime(builder *flatbuffers.Builder, ReadingTime uint64) { builder.PrependUint64Slot(6, ReadingTime, 0) }
func DeviceAddWritesCompleted(builder *flatbuffers.Builder, WritesCompleted uint64) { builder.PrependUint64Slot(7, WritesCompleted, 0) }
func DeviceAddWritesMerged(builder *flatbuffers.Builder, WritesMerged uint64) { builder.PrependUint64Slot(8, WritesMerged, 0) }
func DeviceAddWrittenSectors(builder *flatbuffers.Builder, WrittenSectors uint64) { builder.PrependUint64Slot(9, WrittenSectors, 0) }
func DeviceAddWritingTime(builder *flatbuffers.Builder, WritingTime uint64) { builder.PrependUint64Slot(10, WritingTime, 0) }
func DeviceAddIOInProgress(builder *flatbuffers.Builder, IOInProgress uint64) { builder.PrependUint64Slot(11, IOInProgress, 0) }
func DeviceAddIOTime(builder *flatbuffers.Builder, IOTime uint64) { builder.PrependUint64Slot(12, IOTime, 0) }
func DeviceAddWeightedIOTime(builder *flatbuffers.Builder, WeightedIOTime uint64) { builder.PrependUint64Slot(13, WeightedIOTime, 0) }
func DeviceEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
