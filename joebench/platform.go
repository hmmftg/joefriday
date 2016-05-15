// Copyright 2016 Joel Scoble and The JoeFriday authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	"github.com/mohae/benchutil"
	"github.com/mohae/joefriday/platform/kernel"
	kfb "github.com/mohae/joefriday/platform/kernel/flat"
	kjson "github.com/mohae/joefriday/platform/kernel/json"
	"github.com/mohae/joefriday/platform/loadavg"
	lfb "github.com/mohae/joefriday/platform/loadavg/flat"
	ljson "github.com/mohae/joefriday/platform/loadavg/json"
	"github.com/mohae/joefriday/platform/release"
	rfb "github.com/mohae/joefriday/platform/release/flat"
	rjson "github.com/mohae/joefriday/platform/release/json"
	"github.com/mohae/joefriday/platform/uptime"
	ufb "github.com/mohae/joefriday/platform/uptime/flat"
	ujson "github.com/mohae/joefriday/platform/uptime/json"
)

const (
	PlatformKernel  = "Platform Kernel"
	PlatformLoadAvg = "Platform LoadAvg"
	PlatformRelease = "Platform Release"
	PlatformUptime  = "Platform Uptime"
)

func runPlatformBenchmarks(bench benchutil.Benchmarker) {
	b := PlatformKernelGet()
	bench.Add(b)

	b = PlatformKernelGetFB()
	bench.Add(b)

	b = PlatformKernelSerializeFB()
	bench.Add(b)

	b = PlatformKernelDeserializeFB()
	bench.Add(b)

	b = PlatformKernelGetJSON()
	bench.Add(b)

	b = PlatformKernelSerializeJSON()
	bench.Add(b)

	b = PlatformKernelDeserializeJSON()
	bench.Add(b)

	b = PlatformLoadAvgGet()
	bench.Add(b)

	b = PlatformLoadAvgGetFB()
	bench.Add(b)

	b = PlatformLoadAvgSerializeFB()
	bench.Add(b)

	b = PlatformLoadAvgDeserializeFB()
	bench.Add(b)

	b = PlatformLoadAvgGetJSON()
	bench.Add(b)

	b = PlatformLoadAvgSerializeJSON()
	bench.Add(b)

	b = PlatformLoadAvgDeserializeJSON()
	bench.Add(b)

	b = PlatformReleaseGet()
	bench.Add(b)

	b = PlatformReleaseGetFB()
	bench.Add(b)

	b = PlatformReleaseSerializeFB()
	bench.Add(b)

	b = PlatformReleaseDeserializeFB()
	bench.Add(b)

	b = PlatformReleaseGetJSON()
	bench.Add(b)

	b = PlatformReleaseSerializeJSON()
	bench.Add(b)

	b = PlatformReleaseDeserializeJSON()
	bench.Add(b)

	b = PlatformUptimeGet()
	bench.Add(b)

	b = PlatformUptimeGetFB()
	bench.Add(b)

	b = PlatformUptimeSerializeFB()
	bench.Add(b)

	b = PlatformUptimeDeserializeFB()
	bench.Add(b)

	b = PlatformUptimeGetJSON()
	bench.Add(b)

	b = PlatformUptimeSerializeJSON()
	bench.Add(b)

	b = PlatformUptimeDeserializeJSON()
	bench.Add(b)
}

func BenchPlatformKernelGet(b *testing.B) {
	var k *kernel.Kernel
	b.StopTimer()
	p, _ := kernel.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		k, _ = p.Get()
	}
	_ = k
}

func PlatformKernelGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = PlatformKernel
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformKernelGet))
	return bench
}

func BenchPlatformKernelGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := kfb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func PlatformKernelGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = PlatformKernel
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformKernelGetFB))
	return bench
}

func BenchPlatformKernelSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := kernel.NewProfiler()
	k, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = kfb.Serialize(k)
	}
	_ = tmp
}

func PlatformKernelSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = PlatformKernel
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformKernelSerializeFB))
	return bench
}

func BenchPlatformKernelDeserializeFB(b *testing.B) {
	var k *kernel.Kernel
	b.StopTimer()
	p, _ := kfb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		k = kfb.Deserialize(tmp)
	}
	_ = k
}

func PlatformKernelDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = PlatformKernel
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformKernelDeserializeFB))
	return bench
}

func BenchPlatformKernelGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := kjson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func PlatformKernelGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = PlatformKernel
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformKernelGetJSON))
	return bench
}

func BenchPlatformKernelSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := kernel.NewProfiler()
	k, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = kjson.Serialize(k)
	}
	_ = tmp
}

func PlatformKernelSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = PlatformKernel
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformKernelSerializeJSON))
	return bench
}

func BenchPlatformKernelDeserializeJSON(b *testing.B) {
	var k *kernel.Kernel
	b.StopTimer()
	p, _ := kjson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		k, _ = kjson.Deserialize(tmp)
	}
	_ = k
}

func PlatformKernelDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = PlatformKernel
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformKernelDeserializeJSON))
	return bench
}

// LoadAvg
func BenchPlatformLoadAvgGet(b *testing.B) {
	var l loadavg.LoadAvg
	b.StopTimer()
	p, _ := loadavg.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l, _ = p.Get()
	}
	_ = l
}

func PlatformLoadAvgGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = PlatformLoadAvg
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformLoadAvgGet))
	return bench
}

func BenchPlatformLoadAvgGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := lfb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func PlatformLoadAvgGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = PlatformLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformLoadAvgGetFB))
	return bench
}

func BenchPlatformLoadAvgSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := loadavg.NewProfiler()
	l, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = lfb.Serialize(l)
	}
	_ = tmp
}

func PlatformLoadAvgSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = PlatformLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformLoadAvgSerializeFB))
	return bench
}

func BenchPlatformLoadAvgDeserializeFB(b *testing.B) {
	var l loadavg.LoadAvg
	b.StopTimer()
	p, _ := lfb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l = lfb.Deserialize(tmp)
	}
	_ = l
}

func PlatformLoadAvgDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = PlatformLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformLoadAvgDeserializeFB))
	return bench
}

func BenchPlatformLoadAvgGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ljson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func PlatformLoadAvgGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = PlatformLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformLoadAvgGetJSON))
	return bench
}

func BenchPlatformLoadAvgSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := loadavg.NewProfiler()
	l, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ljson.Serialize(l)
	}
	_ = tmp
}

func PlatformLoadAvgSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = PlatformLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformLoadAvgSerializeJSON))
	return bench
}

func BenchPlatformLoadAvgDeserializeJSON(b *testing.B) {
	var l loadavg.LoadAvg
	b.StopTimer()
	p, _ := ljson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l, _ = ljson.Deserialize(tmp)
	}
	_ = l
}

func PlatformLoadAvgDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = PlatformLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformLoadAvgDeserializeJSON))
	return bench
}

// release
func BenchPlatformReleaseGet(b *testing.B) {
	var r *release.Release
	b.StopTimer()
	p, _ := release.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r, _ = p.Get()
	}
	_ = r
}

func PlatformReleaseGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = PlatformRelease
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformReleaseGet))
	return bench
}

func BenchPlatformReleaseGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := rfb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func PlatformReleaseGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = PlatformRelease
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformReleaseGetFB))
	return bench
}

func BenchPlatformReleaseSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := release.NewProfiler()
	l, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = rfb.Serialize(l)
	}
	_ = tmp
}

func PlatformReleaseSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = PlatformRelease
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformReleaseSerializeFB))
	return bench
}

func BenchPlatformReleaseDeserializeFB(b *testing.B) {
	var r *release.Release
	b.StopTimer()
	p, _ := rfb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r = rfb.Deserialize(tmp)
	}
	_ = r
}

func PlatformReleaseDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = PlatformRelease
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformReleaseDeserializeFB))
	return bench
}

func BenchPlatformReleaseGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := rjson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func PlatformReleaseGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = PlatformRelease
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformReleaseGetJSON))
	return bench
}

func BenchPlatformReleaseSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := release.NewProfiler()
	l, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = rjson.Serialize(l)
	}
	_ = tmp
}

func PlatformReleaseSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = PlatformRelease
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformReleaseSerializeJSON))
	return bench
}

func BenchPlatformReleaseDeserializeJSON(b *testing.B) {
	var r *release.Release
	b.StopTimer()
	p, _ := rjson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		r, _ = rjson.Deserialize(tmp)
	}
	_ = r
}

func PlatformReleaseDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = PlatformRelease
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformReleaseDeserializeJSON))
	return bench
}

// uptime
func BenchPlatformUptimeGet(b *testing.B) {
	var u uptime.Uptime
	b.StopTimer()
	p, _ := uptime.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func PlatformUptimeGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = PlatformUptime
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformUptimeGet))
	return bench
}

func BenchPlatformUptimeGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ufb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func PlatformUptimeGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = PlatformUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformUptimeGetFB))
	return bench
}

func BenchPlatformUptimeSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := uptime.NewProfiler()
	l, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ufb.Serialize(l)
	}
	_ = tmp
}

func PlatformUptimeSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = PlatformUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformUptimeSerializeFB))
	return bench
}

func BenchPlatformUptimeDeserializeFB(b *testing.B) {
	var u uptime.Uptime
	b.StopTimer()
	p, _ := ufb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u = ufb.Deserialize(tmp)
	}
	_ = u
}

func PlatformUptimeDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = PlatformUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformUptimeDeserializeFB))
	return bench
}

func BenchPlatformUptimeGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ujson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func PlatformUptimeGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = PlatformUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformUptimeGetJSON))
	return bench
}

func BenchPlatformUptimeSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := uptime.NewProfiler()
	u, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ujson.Serialize(u)
	}
	_ = tmp
}

func PlatformUptimeSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = PlatformUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformUptimeSerializeJSON))
	return bench
}

func BenchPlatformUptimeDeserializeJSON(b *testing.B) {
	var u uptime.Uptime
	b.StopTimer()
	p, _ := ujson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = ujson.Deserialize(tmp)
	}
	_ = u
}

func PlatformUptimeDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = PlatformUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchPlatformUptimeDeserializeJSON))
	return bench
}
