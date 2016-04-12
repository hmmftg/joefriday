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

func (rcv *LoadAvg) LastMinute() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LoadAvg) LastFive() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LoadAvg) LastTen() float32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetFloat32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LoadAvg) RunningProcesses() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LoadAvg) TotalProcesses() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *LoadAvg) PID() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func LoadAvgStart(builder *flatbuffers.Builder) { builder.StartObject(6) }
func LoadAvgAddLastMinute(builder *flatbuffers.Builder, LastMinute float32) { builder.PrependFloat32Slot(0, LastMinute, 0) }
func LoadAvgAddLastFive(builder *flatbuffers.Builder, LastFive float32) { builder.PrependFloat32Slot(1, LastFive, 0) }
func LoadAvgAddLastTen(builder *flatbuffers.Builder, LastTen float32) { builder.PrependFloat32Slot(2, LastTen, 0) }
func LoadAvgAddRunningProcesses(builder *flatbuffers.Builder, RunningProcesses int32) { builder.PrependInt32Slot(3, RunningProcesses, 0) }
func LoadAvgAddTotalProcesses(builder *flatbuffers.Builder, TotalProcesses int32) { builder.PrependInt32Slot(4, TotalProcesses, 0) }
func LoadAvgAddPID(builder *flatbuffers.Builder, PID int32) { builder.PrependInt32Slot(5, PID, 0) }
func LoadAvgEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
