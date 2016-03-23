package bench

import (
	"testing"

	fb "github.com/google/flatbuffers/go"
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

func BenchmarkCPUStats(b *testing.B) {
	var val cpu.Stats
	for i := 0; i < b.N; i++ {
		val, _ = cpu.GetStats()
	}
	_ = val
}

func BenchmarkCPUStatsSerializeFlat(b *testing.B) {
	var val cpu.Stats
	var p []byte
	b.StopTimer()
	val, _ = cpu.GetStats()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = val.SerializeFlat()
	}
	_ = p
}

func BenchmarkCPUStatsSerializeFlatBuilder(b *testing.B) {
	var val cpu.Stats
	var p []byte
	bldr := fb.NewBuilder(0)
	b.StopTimer()
	val, _ = cpu.GetStats()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bldr.Reset()
		p = val.SerializeFlatBuilder(bldr)
	}
	_ = p
}

func BenchmarkCPUDeSerializeStatsFlat(b *testing.B) {
	var val cpu.Stats
	var p []byte
	b.StopTimer()
	val, _ = cpu.GetStats()
	p = val.SerializeFlat()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val = cpu.DeserializeStatsFlat(p)
	}
	_ = val
}
