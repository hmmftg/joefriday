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
	"github.com/cloudfoundry/gosigar"
)

func BenchmarkJoeFridayMemInfoCatProcInfo(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoCat()
	}
	_ = inf
}

func BenchmarkJoeFridayMemDataCatProcInfo(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataCat()
	}
	_ = data
}

func BenchmarkJoeFridayMemDataCatProcInfoReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataCatReuseBldr()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoReadProcInfo(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoRead()
	}
	_ = inf
}

func BenchmarkJoeFridayMemDataReadProcInfo(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataRead()
	}
	_ = data
}

func BenchmarkJoeFridayMemDataReadProcInfoReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReadReuseBldr()
	}
	_ = data
}

func BenchmarkJoeFridayMemInfoReadReuseRProcInfo(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoReadReuseR()
	}
	_ = inf
}

func BenchmarkJoeFridayMemDataReadReuseRProcInfo(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReadReuseR()
	}
	_ = data
}

func BenchmarkJoeFridayMemDataReadReuseRProcInfoReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReuseRReuseBldr()
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
