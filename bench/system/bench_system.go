package system

import (
	"testing"

	"github.com/DataDog/gohai/platform"
	"github.com/cloudfoundry/gosigar"
	"github.com/mohae/benchutil"
	joeloadavg "github.com/mohae/joefriday/system/loadavg"
	joerelease "github.com/mohae/joefriday/system/release"
	joeversion "github.com/mohae/joefriday/system/version"
	joeuptime "github.com/mohae/joefriday/system/uptime"
	sysload "github.com/mohae/joefriday/sysinfo/loadavg"
	sysuptime "github.com/mohae/joefriday/sysinfo/uptime"
	"github.com/shirou/gopsutil/load"
)

const SystemGroup = "System"

func BenchJoeFridayGetVersion(b *testing.B) {
	var fct *joeversion.Kernel
	b.StopTimer()
	p, _ := joeversion.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = p.Get()
	}
	_ = fct
}

func JoeFridayGetVersion() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/system/version.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetVersion))
	return bench
}

func BenchJoeFridayGetRelease(b *testing.B) {
	var os *joerelease.OS
	b.StopTimer()
	p, _ := joerelease.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		os, _ = p.Get()
	}
	_ = os
}

func JoeFridayGetRelease() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/system/release.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetRelease))
	return bench
}

func BenchDataDogGohaiplatform(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &platform.Platform{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func DataDogGohaiplatform() benchutil.Bench {
	bench := benchutil.NewBench("DataDog/gohai/platform.Platform.Collect")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchDataDogGohaiplatform))
	return bench
}

func BenchJoeFridayGetLoadAvg(b *testing.B) {
	var tmp joeloadavg.LoadAvg
	b.StopTimer()
	p, _ := joeloadavg.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetLoadAvg() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/system/loadavg.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetLoadAvg))
	return bench
}

func BenchJoeFridayGetSysinfoLoadAvg(b *testing.B) {
	var tmp sysload.LoadAvg
	for i := 0; i < b.N; i++ {
		_ = tmp.Get()
	}
	_ = tmp
}

func JoeFridayGetSysinfoLoadAvg() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/sysinfo/loadavg.LoadAvg.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetSysinfoLoadAvg))
	return bench
}

func BenchCloudFoundryGoSigarLoadAverage(b *testing.B) {
	var tmp sigar.LoadAverage
	for i := 0; i < b.N; i++ {
		tmp.Get()
	}
	_ = tmp
}

func CloudFoundryGoSigarLoadAverage() benchutil.Bench {
	bench := benchutil.NewBench("cloudfoundry/gosigar.LoadAverage.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCloudFoundryGoSigarLoadAverage))
	return bench
}

func BenchShirouGopsutilLoadAvg(b *testing.B) {
	var tmp *load.AvgStat
	for i := 0; i < b.N; i++ {
		tmp, _ = load.Avg()
	}
	_ = tmp
}

func ShirouGopsutilLoadAvg() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/load.Avg")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilLoadAvg))
	return bench
}

func BenchShirouGopsutilLoadMisc(b *testing.B) {
	var tmp *load.MiscStat
	for i := 0; i < b.N; i++ {
		tmp, _ = load.Misc()
	}
	_ = tmp
}

func ShirouGopsutilLoadMisc() benchutil.Bench {
	bench := benchutil.NewBench("shirou/gopsutil/load.Misc")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchShirouGopsutilLoadMisc))
	return bench
}

func BenchJoeFridayGetUptime(b *testing.B) {
	var tmp joeuptime.Uptime
	b.StopTimer()
	p, _ := joeuptime.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func JoeFridayGetUptime() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/system/uptime.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetUptime))
	return bench
}

func BenchJoeFridayGetSysinfoUptime(b *testing.B) {
	var tmp sysuptime.Uptime
	for i := 0; i < b.N; i++ {
		_ = tmp.Get()
	}
	_ = tmp
}

func JoeFridayGetSysinfoUptime() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/sysinfo/uptime.Uptime.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetSysinfoUptime))
	return bench
}

func BenchCloudFoundryGoSigarUptime(b *testing.B) {
	var tmp sigar.Uptime
	for i := 0; i < b.N; i++ {
		tmp.Get()
	}
	_ = tmp
}

func CloudFoundryGoSigarUptime() benchutil.Bench {
	bench := benchutil.NewBench("cloudfoundry/gosigar.Uptime.Get")
	bench.Group = SystemGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCloudFoundryGoSigarUptime))
	return bench
}
