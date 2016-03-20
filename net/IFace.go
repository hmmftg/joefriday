// automatically generated, do not modify

package net

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type IFace struct {
	_tab flatbuffers.Table
}

func (rcv *IFace) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *IFace) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *IFace) RBytes() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) RPackets() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) RErrs() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) RDrop() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) RFIFO() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) RFrame() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) RCompressed() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) RMulticast() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) TBytes() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) TPackets() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) TErrs() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) TDrop() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) TFIFO() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) TColls() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(32))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) TCarrier() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(34))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *IFace) TCompressed() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(36))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func IFaceStart(builder *flatbuffers.Builder) { builder.StartObject(17) }
func IFaceAddName(builder *flatbuffers.Builder, Name flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(Name), 0) }
func IFaceAddRBytes(builder *flatbuffers.Builder, RBytes int64) { builder.PrependInt64Slot(1, RBytes, 0) }
func IFaceAddRPackets(builder *flatbuffers.Builder, RPackets int64) { builder.PrependInt64Slot(2, RPackets, 0) }
func IFaceAddRErrs(builder *flatbuffers.Builder, RErrs int64) { builder.PrependInt64Slot(3, RErrs, 0) }
func IFaceAddRDrop(builder *flatbuffers.Builder, RDrop int64) { builder.PrependInt64Slot(4, RDrop, 0) }
func IFaceAddRFIFO(builder *flatbuffers.Builder, RFIFO int64) { builder.PrependInt64Slot(5, RFIFO, 0) }
func IFaceAddRFrame(builder *flatbuffers.Builder, RFrame int64) { builder.PrependInt64Slot(6, RFrame, 0) }
func IFaceAddRCompressed(builder *flatbuffers.Builder, RCompressed int64) { builder.PrependInt64Slot(7, RCompressed, 0) }
func IFaceAddRMulticast(builder *flatbuffers.Builder, RMulticast int64) { builder.PrependInt64Slot(8, RMulticast, 0) }
func IFaceAddTBytes(builder *flatbuffers.Builder, TBytes int64) { builder.PrependInt64Slot(9, TBytes, 0) }
func IFaceAddTPackets(builder *flatbuffers.Builder, TPackets int64) { builder.PrependInt64Slot(10, TPackets, 0) }
func IFaceAddTErrs(builder *flatbuffers.Builder, TErrs int64) { builder.PrependInt64Slot(11, TErrs, 0) }
func IFaceAddTDrop(builder *flatbuffers.Builder, TDrop int64) { builder.PrependInt64Slot(12, TDrop, 0) }
func IFaceAddTFIFO(builder *flatbuffers.Builder, TFIFO int64) { builder.PrependInt64Slot(13, TFIFO, 0) }
func IFaceAddTColls(builder *flatbuffers.Builder, TColls int64) { builder.PrependInt64Slot(14, TColls, 0) }
func IFaceAddTCarrier(builder *flatbuffers.Builder, TCarrier int64) { builder.PrependInt64Slot(15, TCarrier, 0) }
func IFaceAddTCompressed(builder *flatbuffers.Builder, TCompressed int64) { builder.PrependInt64Slot(16, TCompressed, 0) }
func IFaceEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
