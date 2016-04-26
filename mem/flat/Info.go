// automatically generated, do not modify

package flat

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Info struct {
	_tab flatbuffers.Table
}

func GetRootAsInfo(buf []byte, offset flatbuffers.UOffsetT) *Info {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Info{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Info) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Info) Timestamp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Active() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) ActiveAnon() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) ActiveFile() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) AnonHugePages() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) AnonPages() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(14))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Bounce() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(16))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Buffers() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(18))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Cached() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(20))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) CommitLimit() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(22))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) CommittedAS() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(24))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) DirectMap4K() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(26))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) DirectMap2M() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(28))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Dirty() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(30))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) HardwareCorrupted() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(32))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) HugePagesFree() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(34))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) HugePagesRsvd() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(36))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) HugePagesSize() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(38))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) HugePagesSurp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(40))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) HugePagesTotal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(42))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Inactive() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(44))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) InactiveAnon() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(46))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) InactiveFile() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(48))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) KernelStack() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(50))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Mapped() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(52))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) MemAvailable() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(54))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) MemFree() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(56))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) MemTotal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(58))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Mlocked() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(60))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) NFSUnstable() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(62))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) PageTables() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(64))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Shmem() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(66))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Slab() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(68))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) SReclaimable() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(70))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) SUnreclaim() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(72))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) SwapCached() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(74))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) SwapFree() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(76))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) SwapTotal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(78))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Unevictable() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(80))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) VmallocChunk() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(82))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) VmallocTotal() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(84))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) VmallocUsed() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(86))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) Writeback() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(88))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Info) WritebackTmp() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(90))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func InfoStart(builder *flatbuffers.Builder) { builder.StartObject(44) }
func InfoAddTimestamp(builder *flatbuffers.Builder, Timestamp int64) { builder.PrependInt64Slot(0, Timestamp, 0) }
func InfoAddActive(builder *flatbuffers.Builder, Active int64) { builder.PrependInt64Slot(1, Active, 0) }
func InfoAddActiveAnon(builder *flatbuffers.Builder, ActiveAnon int64) { builder.PrependInt64Slot(2, ActiveAnon, 0) }
func InfoAddActiveFile(builder *flatbuffers.Builder, ActiveFile int64) { builder.PrependInt64Slot(3, ActiveFile, 0) }
func InfoAddAnonHugePages(builder *flatbuffers.Builder, AnonHugePages int64) { builder.PrependInt64Slot(4, AnonHugePages, 0) }
func InfoAddAnonPages(builder *flatbuffers.Builder, AnonPages int64) { builder.PrependInt64Slot(5, AnonPages, 0) }
func InfoAddBounce(builder *flatbuffers.Builder, Bounce int64) { builder.PrependInt64Slot(6, Bounce, 0) }
func InfoAddBuffers(builder *flatbuffers.Builder, Buffers int64) { builder.PrependInt64Slot(7, Buffers, 0) }
func InfoAddCached(builder *flatbuffers.Builder, Cached int64) { builder.PrependInt64Slot(8, Cached, 0) }
func InfoAddCommitLimit(builder *flatbuffers.Builder, CommitLimit int64) { builder.PrependInt64Slot(9, CommitLimit, 0) }
func InfoAddCommittedAS(builder *flatbuffers.Builder, CommittedAS int64) { builder.PrependInt64Slot(10, CommittedAS, 0) }
func InfoAddDirectMap4K(builder *flatbuffers.Builder, DirectMap4K int64) { builder.PrependInt64Slot(11, DirectMap4K, 0) }
func InfoAddDirectMap2M(builder *flatbuffers.Builder, DirectMap2M int64) { builder.PrependInt64Slot(12, DirectMap2M, 0) }
func InfoAddDirty(builder *flatbuffers.Builder, Dirty int64) { builder.PrependInt64Slot(13, Dirty, 0) }
func InfoAddHardwareCorrupted(builder *flatbuffers.Builder, HardwareCorrupted int64) { builder.PrependInt64Slot(14, HardwareCorrupted, 0) }
func InfoAddHugePagesFree(builder *flatbuffers.Builder, HugePagesFree int64) { builder.PrependInt64Slot(15, HugePagesFree, 0) }
func InfoAddHugePagesRsvd(builder *flatbuffers.Builder, HugePagesRsvd int64) { builder.PrependInt64Slot(16, HugePagesRsvd, 0) }
func InfoAddHugePagesSize(builder *flatbuffers.Builder, HugePagesSize int64) { builder.PrependInt64Slot(17, HugePagesSize, 0) }
func InfoAddHugePagesSurp(builder *flatbuffers.Builder, HugePagesSurp int64) { builder.PrependInt64Slot(18, HugePagesSurp, 0) }
func InfoAddHugePagesTotal(builder *flatbuffers.Builder, HugePagesTotal int64) { builder.PrependInt64Slot(19, HugePagesTotal, 0) }
func InfoAddInactive(builder *flatbuffers.Builder, Inactive int64) { builder.PrependInt64Slot(20, Inactive, 0) }
func InfoAddInactiveAnon(builder *flatbuffers.Builder, InactiveAnon int64) { builder.PrependInt64Slot(21, InactiveAnon, 0) }
func InfoAddInactiveFile(builder *flatbuffers.Builder, InactiveFile int64) { builder.PrependInt64Slot(22, InactiveFile, 0) }
func InfoAddKernelStack(builder *flatbuffers.Builder, KernelStack int64) { builder.PrependInt64Slot(23, KernelStack, 0) }
func InfoAddMapped(builder *flatbuffers.Builder, Mapped int64) { builder.PrependInt64Slot(24, Mapped, 0) }
func InfoAddMemAvailable(builder *flatbuffers.Builder, MemAvailable int64) { builder.PrependInt64Slot(25, MemAvailable, 0) }
func InfoAddMemFree(builder *flatbuffers.Builder, MemFree int64) { builder.PrependInt64Slot(26, MemFree, 0) }
func InfoAddMemTotal(builder *flatbuffers.Builder, MemTotal int64) { builder.PrependInt64Slot(27, MemTotal, 0) }
func InfoAddMlocked(builder *flatbuffers.Builder, Mlocked int64) { builder.PrependInt64Slot(28, Mlocked, 0) }
func InfoAddNFSUnstable(builder *flatbuffers.Builder, NFSUnstable int64) { builder.PrependInt64Slot(29, NFSUnstable, 0) }
func InfoAddPageTables(builder *flatbuffers.Builder, PageTables int64) { builder.PrependInt64Slot(30, PageTables, 0) }
func InfoAddShmem(builder *flatbuffers.Builder, Shmem int64) { builder.PrependInt64Slot(31, Shmem, 0) }
func InfoAddSlab(builder *flatbuffers.Builder, Slab int64) { builder.PrependInt64Slot(32, Slab, 0) }
func InfoAddSReclaimable(builder *flatbuffers.Builder, SReclaimable int64) { builder.PrependInt64Slot(33, SReclaimable, 0) }
func InfoAddSUnreclaim(builder *flatbuffers.Builder, SUnreclaim int64) { builder.PrependInt64Slot(34, SUnreclaim, 0) }
func InfoAddSwapCached(builder *flatbuffers.Builder, SwapCached int64) { builder.PrependInt64Slot(35, SwapCached, 0) }
func InfoAddSwapFree(builder *flatbuffers.Builder, SwapFree int64) { builder.PrependInt64Slot(36, SwapFree, 0) }
func InfoAddSwapTotal(builder *flatbuffers.Builder, SwapTotal int64) { builder.PrependInt64Slot(37, SwapTotal, 0) }
func InfoAddUnevictable(builder *flatbuffers.Builder, Unevictable int64) { builder.PrependInt64Slot(38, Unevictable, 0) }
func InfoAddVmallocChunk(builder *flatbuffers.Builder, VmallocChunk int64) { builder.PrependInt64Slot(39, VmallocChunk, 0) }
func InfoAddVmallocTotal(builder *flatbuffers.Builder, VmallocTotal int64) { builder.PrependInt64Slot(40, VmallocTotal, 0) }
func InfoAddVmallocUsed(builder *flatbuffers.Builder, VmallocUsed int64) { builder.PrependInt64Slot(41, VmallocUsed, 0) }
func InfoAddWriteback(builder *flatbuffers.Builder, Writeback int64) { builder.PrependInt64Slot(42, Writeback, 0) }
func InfoAddWritebackTmp(builder *flatbuffers.Builder, WritebackTmp int64) { builder.PrependInt64Slot(43, WritebackTmp, 0) }
func InfoEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
