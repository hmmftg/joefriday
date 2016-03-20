// Package bench contains benchmarks for gathering system information.
// These packages are not directly comparable because of the differences
// in what they gather, but I wanted to see some numbers.
//
// This will only work on linux systems due to limitations of
// github.com/mohae/joefriday.
package bench

import (
	"testing"

	"github.com/DataDog/gohai/memory"
	"github.com/EricLagergren/go-gnulib/sysinfo"
	"github.com/cloudfoundry/gosigar"
	gopsutilmem "github.com/shirou/gopsutil/mem"
)

func BenchmarkEricLagergrenMemInfo(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetInfo()
	}
	_ = inf
}

func BenchmarkJoeFridayMemInfoCat(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoCat()
	}
	_ = inf
}

func BenchmarkJoeFridayMemInfoCatToJSON(b *testing.B) {
	var inf []byte
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoCatToJSON()
	}
	_ = inf
}

func BenchmarkJoeFridayMemInfoCatToFlatbuffers(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataCat()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoCatToFlatbuffersReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataCatReuseBldr()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoRead(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoRead()
	}
	_ = inf
}

func BenchmarkJoeFridayMemInfoReadToJSON(b *testing.B) {
	var inf []byte
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoReadToJSON()
	}
	_ = inf
}

func BenchmarkJoeFridayMemInfoReadToFlatbuffers(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataRead()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoReadToFlatbuffersReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReadReuseBldr()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoReadReuseBufioReader(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoReadReuseR()
	}
	_ = inf
}

func BenchmarkJoeFridayMemInfoReadReuseBufioReaderToJSON(b *testing.B) {
	var inf []byte
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoReadReuseRToJSON()
	}
	_ = inf
}

func BenchmarkJoeFridayMemInfoReadReuseBufioReaderToFlatbuffers(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReadReuseR()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoReadReuseBufioReaderToFlatbuffersReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReuseRReuseBldr()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoToFlatbuffersReuseBldrReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemInfoToFlatbuffersReuseBldr()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoToFlatbuffersMinAllocs(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemInfoToFlatbuffersMinAllocs()
	}
	_ = data
}

func BenchmarkGohaiMem(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &memory.Memory{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func BenchmarkGoSigarMem(b *testing.B) {
	var mem sigar.Mem
	for i := 0; i < b.N; i++ {
		mem.Get()
	}
	_ = mem
}

func BenchmarkGopsutilMem(b *testing.B) {
	var mem *gopsutilmem.VirtualMemoryStat
	for i := 0; i < b.N; i++ {
		mem, _ = gopsutilmem.VirtualMemory()
	}
	_ = mem
}

func BenchmarkGnulibSysinfo(b *testing.B) {
	var memA, memT int64
	for i := 0; i < b.N; i++ {
		memA = sysinfo.PhysmemAvailable()
		memT = sysinfo.PhysmemTotal()
	}
	_ = memA
	_ = memT
}
