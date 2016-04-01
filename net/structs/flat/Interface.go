// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Interface struct {
	_tab flatbuffers.Table
}

func (rcv *Interface) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Interface) Name() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Interface) RBytes() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) RPackets() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) RErrs() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) RDrop() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) RFIFO() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) RFrame() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) RCompressed() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) RMulticast() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) TBytes() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) TPackets() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) TErrs() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) TDrop() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) TFIFO() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) TColls() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(32))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) TCarrier() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(34))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Interface) TCompressed() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(36))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func InterfaceStart(builder *flatbuffers.Builder) { builder.StartObject(17) }
func InterfaceAddName(builder *flatbuffers.Builder, Name flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(Name), 0) }
func InterfaceAddRBytes(builder *flatbuffers.Builder, RBytes int64) { builder.PrependInt64Slot(1, RBytes, 0) }
func InterfaceAddRPackets(builder *flatbuffers.Builder, RPackets int64) { builder.PrependInt64Slot(2, RPackets, 0) }
func InterfaceAddRErrs(builder *flatbuffers.Builder, RErrs int64) { builder.PrependInt64Slot(3, RErrs, 0) }
func InterfaceAddRDrop(builder *flatbuffers.Builder, RDrop int64) { builder.PrependInt64Slot(4, RDrop, 0) }
func InterfaceAddRFIFO(builder *flatbuffers.Builder, RFIFO int64) { builder.PrependInt64Slot(5, RFIFO, 0) }
func InterfaceAddRFrame(builder *flatbuffers.Builder, RFrame int64) { builder.PrependInt64Slot(6, RFrame, 0) }
func InterfaceAddRCompressed(builder *flatbuffers.Builder, RCompressed int64) { builder.PrependInt64Slot(7, RCompressed, 0) }
func InterfaceAddRMulticast(builder *flatbuffers.Builder, RMulticast int64) { builder.PrependInt64Slot(8, RMulticast, 0) }
func InterfaceAddTBytes(builder *flatbuffers.Builder, TBytes int64) { builder.PrependInt64Slot(9, TBytes, 0) }
func InterfaceAddTPackets(builder *flatbuffers.Builder, TPackets int64) { builder.PrependInt64Slot(10, TPackets, 0) }
func InterfaceAddTErrs(builder *flatbuffers.Builder, TErrs int64) { builder.PrependInt64Slot(11, TErrs, 0) }
func InterfaceAddTDrop(builder *flatbuffers.Builder, TDrop int64) { builder.PrependInt64Slot(12, TDrop, 0) }
func InterfaceAddTFIFO(builder *flatbuffers.Builder, TFIFO int64) { builder.PrependInt64Slot(13, TFIFO, 0) }
func InterfaceAddTColls(builder *flatbuffers.Builder, TColls int64) { builder.PrependInt64Slot(14, TColls, 0) }
func InterfaceAddTCarrier(builder *flatbuffers.Builder, TCarrier int64) { builder.PrependInt64Slot(15, TCarrier, 0) }
func InterfaceAddTCompressed(builder *flatbuffers.Builder, TCompressed int64) { builder.PrependInt64Slot(16, TCompressed, 0) }
func InterfaceEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
