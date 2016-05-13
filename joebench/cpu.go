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
	Flat     = "FlatBuffers"
	JSON     = "JSON"
)

func BenchJoeFridayFactsGet(b *testing.B) {
	var fct *facts.Facts
	b.StopTimer()
	p, _ := facts.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = p.Get()
	}
	_ = fct
}

func JoeFridayGetFacts() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts.Get")
	bench.Group = CPUFact
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayFactsGet))
	return bench
}

func BenchJoeFridayFactsGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ffb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetFactsFB() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/flat.Get")
	bench.Group = CPUFact
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayFactsGetFB))
	return bench
}

func BenchJoeFridayFactsSerializeFB(b *testing.B) {
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

func JoeFridayFactsSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/flat.Serialize")
	bench.Group = CPUFact
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayFactsSerializeFB))
	return bench
}

func BenchJoeFridayFactsDeserializeFB(b *testing.B) {
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

func JoeFridayFactsDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/flat.Deserialize")
	bench.Group = CPUFact
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayFactsDeserializeFB))
	return bench
}

func BenchJoeFridayFactsGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := fjson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetFactsJSON() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/json.Get")
	bench.Group = CPUFact
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayFactsGetJSON))
	return bench
}

func BenchJoeFridayFactsSerializeJSON(b *testing.B) {
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

func JoeFridayFactsSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/json.Serialize")
	bench.Group = CPUFact
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayFactsSerializeJSON))
	return bench
}

func BenchJoeFridayFactsDeserializeJSON(b *testing.B) {
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

func JoeFridayFactsDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("cpu/facts/json.Deserialize")
	bench.Group = CPUFact
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayFactsDeserializeJSON))
	return bench
}

// Stats
func BenchJoeFridayStatsGet(b *testing.B) {
	var sts *stats.Stats
	b.StopTimer()
	p, _ := stats.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sts, _ = p.Get()
	}
	_ = sts
}

func JoeFridayGetStats() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUStats
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayStatsGet))
	return bench
}

func BenchJoeFridayStatsGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := sfb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetStatsFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayStatsGetFB))
	return bench
}

func BenchJoeFridayStatsSerializeFB(b *testing.B) {
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

func JoeFridayStatsSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayStatsSerializeFB))
	return bench
}

func BenchJoeFridayStatsDeserializeFB(b *testing.B) {
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

func JoeFridayStatsDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayStatsDeserializeFB))
	return bench
}

func BenchJoeFridayStatsGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := sjson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetStatsJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayStatsGetJSON))
	return bench
}

func BenchJoeFridayStatsSerializeJSON(b *testing.B) {
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

func JoeFridayStatsSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayStatsSerializeJSON))
	return bench
}

func BenchJoeFridayStatsDeserializeJSON(b *testing.B) {
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

func JoeFridayStatsDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayStatsDeserializeJSON))
	return bench
}

// Utilization
func BenchJoeFridayUtilizationGet(b *testing.B) {
	var u *utilization.Utilization
	b.StopTimer()
	p, _ := utilization.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func JoeFridayGetUtilization() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUUtil
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayUtilizationGet))
	return bench
}

func BenchJoeFridayUtilizationGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ufb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetUtilizationFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayUtilizationGetFB))
	return bench
}

func BenchJoeFridayUtilizationSerializeFB(b *testing.B) {
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

func JoeFridayUtilizationSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayUtilizationSerializeFB))
	return bench
}

func BenchJoeFridayUtilizationDeserializeFB(b *testing.B) {
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

func JoeFridayUtilizationDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayUtilizationDeserializeFB))
	return bench
}

func BenchJoeFridayUtilizationGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ujson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetUtilizationJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayUtilizationGetJSON))
	return bench
}

func BenchJoeFridayUtilizationSerializeJSON(b *testing.B) {
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

func JoeFridayUtilizationSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayUtilizationSerializeJSON))
	return bench
}

func BenchJoeFridayUtilizationDeserializeJSON(b *testing.B) {
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

func JoeFridayUtilizationDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayUtilizationDeserializeJSON))
	return bench
}
