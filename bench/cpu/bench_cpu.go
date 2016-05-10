package cpu

import (
	"testing"

	"github.com/DataDog/gohai/cpu"
	"github.com/mohae/benchutil"
	joefacts "github.com/mohae/joefriday/cpu/facts"
	joestats "github.com/mohae/joefriday/cpu/stats"
	gopsutilcpu "github.com/shirou/gopsutil/cpu"
)

func BenchJoeFridayGetFacts(b *testing.B) {
	var fct *joefacts.Facts
	b.StopTimer()
	p, _ := joefacts.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = p.Get()
	}
	_ = fct
}

func JoeFridayGetFacts() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/cpu/facts.Get")
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetFacts))
	return bench
}

func BenchJoeFridayGetStats(b *testing.B) {
	var st *joestats.Stats
	b.StopTimer()
	p, _ := joestats.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		st, _ = p.Get()
	}
	_ = st
}

func JoeFridayGetStats() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/cpu/stats.Get")
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetStats))
	return bench
}

func BenchDataDogGohaiCPU(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &cpu.Cpu{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func DataDogGohaiCPU() benchutil.Bench {
	bench := benchutil.NewBench("DataDog/gohai/cpu.Cpu.Collect")
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDataDogGohaiCPU))
	return bench
}

func BenchShirouGopsutilInfoStat(b *testing.B) {
	var st []gopsutilcpu.InfoStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilcpu.Info()
	}
	_ = st
}

func ShirouGopsutilInfoStat() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/cpu.Info")
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilInfoStat))
	return bench
}

func BenchShirouGopsutilTimeStat(b *testing.B) {
	var st []gopsutilcpu.TimesStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilcpu.Times(true)
	}
	_ = st
}

func ShirouGopsutilTimeStat() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/cpu.Times")
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilTimeStat))
	return bench
}
