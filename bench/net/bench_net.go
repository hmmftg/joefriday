package net

import (
	"testing"

	"github.com/DataDog/gohai/network"
	"github.com/mohae/benchutil"
	joeinfo "github.com/mohae/joefriday/net/info"
	joestructs "github.com/mohae/joefriday/net/structs"
	joeusage "github.com/mohae/joefriday/net/usage"
	gopsutilnet "github.com/shirou/gopsutil/net"
)

const NetGroup = "Network"

func BenchJoeFridayGetInfo(b *testing.B) {
	var inf *joestructs.Info
	b.StopTimer()
	p, _ := joeinfo.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func JoeFridayGetInfo() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/net/info.Get")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetInfo))
	return bench
}

func BenchJoeFridayGetUsage(b *testing.B) {
	var u *joestructs.Usage
	b.StopTimer()
	p, _ := joeusage.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func JoeFridayGetUsage() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/net/usage.Get")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetUsage))
	return bench
}

func BenchDataDogGohaiNetwork(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &network.Network{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func DataDogGohaiNetwork() benchutil.Bench {
	bench := benchutil.NewBench("DataDog/gohai/network")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDataDogGohaiNetwork))
	return bench
}

func BenchShirouGopsutilNetInterfaces(b *testing.B) {
	var st []gopsutilnet.InterfaceStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilnet.Interfaces()
	}
	_ = st
}

func ShirouGopsutilNetInterfaces() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/net")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilNetInterfaces))
	return bench
}

func BenchShirouGopsutilIOCounters(b *testing.B) {
	var st []gopsutilnet.IOCountersStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilnet.IOCounters(true)
	}
	_ = st
}

func ShirouGopsutilIOCounters() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/net/IOCounters")
	bench.Group = NetGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilIOCounters))
	return bench
}
