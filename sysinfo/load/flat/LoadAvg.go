// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type LoadAvg struct {
	_tab flatbuffers.Table
}

func GetRootAsLoadAvg(buf []byte, offset flatbuffers.UOffsetT) *LoadAvg {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &LoadAvg{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *LoadAvg) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *LoadAvg) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LoadAvg) One() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LoadAvg) Five() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LoadAvg) Fifteen() float64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetFloat64(o + rcv._tab.Pos)
	}
	return 0
}

func LoadAvgStart(builder *flatbuffers.Builder) { builder.StartObject(4) }
func LoadAvgAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func LoadAvgAddOne(builder *flatbuffers.Builder, One float64) { builder.PrependFloat64Slot(1, One, 0) }
func LoadAvgAddFive(builder *flatbuffers.Builder, Five float64) { builder.PrependFloat64Slot(2, Five, 0) }
func LoadAvgAddFifteen(builder *flatbuffers.Builder, Fifteen float64) { builder.PrependFloat64Slot(3, Fifteen, 0) }
func LoadAvgEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
