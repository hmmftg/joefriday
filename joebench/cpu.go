package main

import (
	"testing"

	"github.com/mohae/benchutil"
	"github.com/mohae/joefriday/cpu/facts"
	ffb "github.com/mohae/joefriday/cpu/facts/flat"
	fjson "github.com/mohae/joefriday/cpu/facts/json"
	"github.com/mohae/joefriday/cpu/stats"
	sfb "github.com/mohae/joefriday/cpu/stats/flat"
	sjson "github.com/mohae/joefriday/cpu/stats/json"
	"github.com/mohae/joefriday/cpu/utilization"
	ufb "github.com/mohae/joefriday/cpu/utilization/flat"
	ujson "github.com/mohae/joefriday/cpu/utilization/json"
)

const (
	CPUFact  = "CPU Facts"
	CPUStats = "CPU Stats"
	CPUUtil  = "CPU Utilization"
)

func BenchCPUFactsGet(b *testing.B) {
	var fct *facts.Facts
	b.StopTimer()
	p, _ := facts.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = p.Get()
	}
	_ = fct
}

func CPUGetFacts() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts.Get")
	bench.Group = CPUFact
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFactsGet))
	return bench
}

func BenchCPUFactsGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ffb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUGetFactsFB() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/flat.Get")
	bench.Group = CPUFact
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFactsGetFB))
	return bench
}

func BenchCPUFactsSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := facts.NewProfiler()
	fct, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ffb.Serialize(fct)
	}
	_ = tmp
}

func CPUFactsSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/flat.Serialize")
	bench.Group = CPUFact
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFactsSerializeFB))
	return bench
}

func BenchCPUFactsDeserializeFB(b *testing.B) {
	var fct *facts.Facts
	b.StopTimer()
	p, _ := ffb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct = ffb.Deserialize(tmp)
	}
	_ = fct
}

func CPUFactsDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/flat.Deserialize")
	bench.Group = CPUFact
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFactsDeserializeFB))
	return bench
}

func BenchCPUFactsGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := fjson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUGetFactsJSON() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/json.Get")
	bench.Group = CPUFact
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFactsGetJSON))
	return bench
}

func BenchCPUFactsSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := fjson.NewProfiler()
	fct, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = fjson.Serialize(fct)
	}
	_ = tmp
}

func CPUFactsSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/json.Serialize")
	bench.Group = CPUFact
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFactsSerializeJSON))
	return bench
}

func BenchCPUFactsDeserializeJSON(b *testing.B) {
	var fct *facts.Facts
	b.StopTimer()
	p, _ := fjson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = fjson.Deserialize(tmp)
	}
	_ = fct
}

func CPUFactsDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/json.Deserialize")
	bench.Group = CPUFact
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUFactsDeserializeJSON))
	return bench
}

// Stats
func BenchCPUStatsGet(b *testing.B) {
	var sts *stats.Stats
	b.StopTimer()
	p, _ := stats.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sts, _ = p.Get()
	}
	_ = sts
}

func CPUGetStats() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUStats
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGet))
	return bench
}

func BenchCPUStatsGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := sfb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUGetStatsFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGetFB))
	return bench
}

func BenchCPUStatsSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := stats.NewProfiler()
	sts, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = sfb.Serialize(sts)
	}
	_ = tmp
}

func CPUStatsSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsSerializeFB))
	return bench
}

func BenchCPUStatsDeserializeFB(b *testing.B) {
	var sts *stats.Stats
	b.StopTimer()
	p, _ := sfb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sts = sfb.Deserialize(tmp)
	}
	_ = sts
}

func CPUStatsDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsDeserializeFB))
	return bench
}

func BenchCPUStatsGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := sjson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUGetStatsJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGetJSON))
	return bench
}

func BenchCPUStatsSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := stats.NewProfiler()
	sts, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = sjson.Serialize(sts)
	}
	_ = tmp
}

func CPUStatsSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsSerializeJSON))
	return bench
}

func BenchCPUStatsDeserializeJSON(b *testing.B) {
	var sts *stats.Stats
	b.StopTimer()
	p, _ := sjson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sts, _ = sjson.Deserialize(tmp)
	}
	_ = sts
}

func CPUStatsDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsDeserializeJSON))
	return bench
}

// Utilization
func BenchCPUUtilizationGet(b *testing.B) {
	var u *utilization.Utilization
	b.StopTimer()
	p, _ := utilization.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func CPUGetUtilization() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUUtil
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilizationGet))
	return bench
}

func BenchCPUUtilizationGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ufb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUGetUtilizationFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilizationGetFB))
	return bench
}

func BenchCPUUtilizationSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := utilization.NewProfiler()
	u, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ufb.Serialize(u)
	}
	_ = tmp
}

func CPUUtilizationSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilizationSerializeFB))
	return bench
}

func BenchCPUUtilizationDeserializeFB(b *testing.B) {
	var u *utilization.Utilization
	b.StopTimer()
	p, _ := ufb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u = ufb.Deserialize(tmp)
	}
	_ = u
}

func CPUUtilizationDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilizationDeserializeFB))
	return bench
}

func BenchCPUUtilizationGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ujson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUGetUtilizationJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilizationGetJSON))
	return bench
}

func BenchCPUUtilizationSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := utilization.NewProfiler()
	u, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ujson.Serialize(u)
	}
	_ = tmp
}

func CPUUtilizationSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilizationSerializeJSON))
	return bench
}

func BenchCPUUtilizationDeserializeJSON(b *testing.B) {
	var u *utilization.Utilization
	b.StopTimer()
	p, _ := ujson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = ujson.Deserialize(tmp)
	}
	_ = u
}

func CPUUtilizationDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilizationDeserializeJSON))
	return bench
}
