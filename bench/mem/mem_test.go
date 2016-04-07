package mem

import (
	"testing"

	"github.com/DataDog/gohai/memory"
	"github.com/cloudfoundry/gosigar"
	joe "github.com/mohae/joefriday/mem"
	gopsutilmem "github.com/shirou/gopsutil/mem"
)

func BenchmarkJoeFridayGet(b *testing.B) {
	var mem *joe.Info
	b.StopTimer()
	p, _ := joe.New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mem, _ = p.Get()
	}
	_ = mem
}

func BenchmarkCloudFoundryGoSigarMem(b *testing.B) {
	var mem sigar.Mem
	for i := 0; i < b.N; i++ {
		mem.Get()
	}
	_ = mem
}

func BenchmarkDataDogGohaiMem(b *testing.B) {
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

func BenchmarkShirouGopsutilMem(b *testing.B) {
	var mem *gopsutilmem.VirtualMemoryStat
	for i := 0; i < b.N; i++ {
		mem, _ = gopsutilmem.VirtualMemory()
	}
	_ = mem
}
