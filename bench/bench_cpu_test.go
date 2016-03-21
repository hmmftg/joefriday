package bench

import (
	"testing"

	"github.com/mohae/joefriday/cpu"
)

func BenchmarkCPUNiHao(b *testing.B) {
	var procs *cpu.Processors
	for i := 0; i < b.N; i++ {
		procs, _ = cpu.NiHao()
	}
	_ = procs
}

func BenchmarkCPUProcessorsSerialize(b *testing.B) {
	var procs *cpu.Processors
	var p []byte
	b.StopTimer()
	procs, _ = cpu.NiHao()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = procs.Serialize()
	}
	_ = p
}

func BenchmarkCPUProcessorsDeSerialize(b *testing.B) {
	var procs *cpu.Processors
	var p []byte
	b.StopTimer()
	procs, _ = cpu.NiHao()
	p = procs.Serialize()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		procs = cpu.Deserialize(p)
	}
	_ = procs
}
