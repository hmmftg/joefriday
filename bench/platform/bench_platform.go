package platform

import (
	"testing"

	"github.com/DataDog/gohai/platform"
	"github.com/cloudfoundry/gosigar"
	"github.com/mohae/benchutil"
	joekernel "github.com/mohae/joefriday/platform/kernel"
	joeloadavg "github.com/mohae/joefriday/platform/loadavg"
	joerelease "github.com/mohae/joefriday/platform/release"
	joeuptime "github.com/mohae/joefriday/platform/uptime"
	sysload "github.com/mohae/joefriday/sysinfo/load"
	sysuptime "github.com/mohae/joefriday/sysinfo/uptime"
	"github.com/shirou/gopsutil/load"
)

const PlatformGroup = "Platform"

func BenchJoeFridayGetKernel(b *testing.B) {
	var fct *joekernel.Kernel
	b.StopTimer()
	p, _ := joekernel.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = p.Get()
	}
	_ = fct
}

func JoeFridayGetKernel() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/platform/kernel.Get")
	bench.Group = PlatformGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchJoeFridayGetKernel))
	return bench
}

func BenchJoeFridayGetRelease(b *testing.B) {
	var st *joerelease.Release
	b.StopTimer()
	p, _ := joerelease.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		st, _ = p.Get()
	}
	_ = st
}

func JoeFridayGetRelease() benchutil.Bench {
	bench := benchutil.NewBench("joefriday/platform/release.Get")
	bench.Group = PlatformGroup
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
	bench.Group = PlatformGroup
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
	bench := benchutil.NewBench("joefriday/platform/loadavg.Get")
	bench.Group = PlatformGroup
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
	bench := benchutil.NewBench("joefriday/sysinfo/load.LoadAvg.Get")
	bench.Group = PlatformGroup
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
	bench.Group = PlatformGroup
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
	bench.Group = PlatformGroup
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
	bench.Group = PlatformGroup
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
	bench := benchutil.NewBench("joefriday/platform/uptime.Get")
	bench.Group = PlatformGroup
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
	bench.Group = PlatformGroup
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
	bench.Group = PlatformGroup
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCloudFoundryGoSigarUptime))
	return bench
}
