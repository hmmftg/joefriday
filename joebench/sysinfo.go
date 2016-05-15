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
	"github.com/mohae/joefriday/sysinfo/load"
	lfb "github.com/mohae/joefriday/sysinfo/load/flat"
	ljson "github.com/mohae/joefriday/sysinfo/load/json"
	"github.com/mohae/joefriday/sysinfo/mem"
	mfb "github.com/mohae/joefriday/sysinfo/mem/flat"
	mjson "github.com/mohae/joefriday/sysinfo/mem/json"
	"github.com/mohae/joefriday/sysinfo/uptime"
	ufb "github.com/mohae/joefriday/sysinfo/uptime/flat"
	ujson "github.com/mohae/joefriday/sysinfo/uptime/json"
)

const (
	SysinfoLoadAvg = "Sysinfo LoadAvg"
	SysinfoMemInfo = "Sysinfo MemInfo"
	SysinfoUptime  = "Sysinfo Uptime"
)

func runSysinfoBenchmarks(bench benchutil.Benchmarker) {
	b := SysinfoLoadAvgGet()
	bench.Add(b)

	b = SysinfoLoadAvgGetFB()
	bench.Add(b)

	b = SysinfoLoadAvgSerializeFB()
	bench.Add(b)

	b = SysinfoLoadAvgDeserializeFB()
	bench.Add(b)

	b = SysinfoLoadAvgGetJSON()
	bench.Add(b)

	b = SysinfoLoadAvgDeserializeJSON()
	bench.Add(b)

	// Mem Info
	b = SysinfoMemInfoGet()
	bench.Add(b)

	b = SysinfoMemInfoGetFB()
	bench.Add(b)

	b = SysinfoMemInfoSerializeFB()
	bench.Add(b)

	b = SysinfoMemInfoDeserializeFB()
	bench.Add(b)

	b = SysinfoMemInfoGetJSON()
	bench.Add(b)

	b = SysinfoMemInfoDeserializeJSON()
	bench.Add(b)

	// Uptime
	b = SysinfoUptimeGet()
	bench.Add(b)

	b = SysinfoUptimeGetFB()
	bench.Add(b)

	b = SysinfoUptimeSerializeFB()
	bench.Add(b)

	b = SysinfoUptimeDeserializeFB()
	bench.Add(b)

	b = SysinfoUptimeGetJSON()
	bench.Add(b)

	b = SysinfoUptimeDeserializeJSON()
	bench.Add(b)
}

// LoadAvg
func BenchSysinfoLoadAvgGet(b *testing.B) {
	l := &load.LoadAvg{}
	for i := 0; i < b.N; i++ {
		l.Get()
	}
	_ = l
}

func SysinfoLoadAvgGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SysinfoLoadAvg
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgGet))
	return bench
}

func BenchSysinfoLoadAvgGetFB(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = ljson.Get()
	}
	_ = tmp
}

func SysinfoLoadAvgGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SysinfoLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgGetFB))
	return bench
}

func BenchSysinfoLoadAvgSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	l := &load.LoadAvg{}
	l.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp = lfb.Serialize(l)
	}
	_ = tmp
}

func SysinfoLoadAvgSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SysinfoLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgSerializeFB))
	return bench
}

func BenchSysinfoLoadAvgDeserializeFB(b *testing.B) {
	var l *load.LoadAvg
	b.StopTimer()
	tmp, _ := lfb.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l = lfb.Deserialize(tmp)
	}
	_ = l
}

func SysinfoLoadAvgDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SysinfoLoadAvg
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgDeserializeFB))
	return bench
}

func BenchSysinfoLoadAvgGetJSON(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = ljson.Get()
	}
	_ = tmp
}

func SysinfoLoadAvgGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SysinfoLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgGetJSON))
	return bench
}

func BenchSysinfoLoadAvgDeserializeJSON(b *testing.B) {
	var l *load.LoadAvg
	b.StopTimer()
	tmp, _ := ljson.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		l, _ = ljson.Deserialize(tmp)
	}
	_ = l
}

func SysinfoLoadAvgDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SysinfoLoadAvg
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoLoadAvgDeserializeJSON))
	return bench
}

// Mem Info
func BenchSysinfoMemInfoGet(b *testing.B) {
	m := &mem.Info{}
	for i := 0; i < b.N; i++ {
		m.Get()
	}
	_ = m
}

func SysinfoMemInfoGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SysinfoMemInfo
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoGet))
	return bench
}

func BenchSysinfoMemInfoGetFB(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = mfb.Get()
	}
	_ = tmp
}

func SysinfoMemInfoGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SysinfoMemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoGetFB))
	return bench
}

func BenchSysinfoMemInfoSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	m := &mem.Info{}
	m.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp = mfb.Serialize(m)
	}
	_ = tmp
}

func SysinfoMemInfoSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SysinfoMemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoSerializeFB))
	return bench
}

func BenchSysinfoMemInfoDeserializeFB(b *testing.B) {
	var m *mem.Info
	b.StopTimer()
	tmp, _ := mfb.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m = mfb.Deserialize(tmp)
	}
	_ = m
}

func SysinfoMemInfoDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SysinfoMemInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoDeserializeFB))
	return bench
}

func BenchSysinfoMemInfoGetJSON(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = mjson.Get()
	}
	_ = tmp
}

func SysinfoMemInfoGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SysinfoMemInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoGetJSON))
	return bench
}

func BenchSysinfoMemInfoDeserializeJSON(b *testing.B) {
	var m *mem.Info
	b.StopTimer()
	tmp, _ := mjson.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m, _ = mjson.Deserialize(tmp)
	}
	_ = m
}

func SysinfoMemInfoDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SysinfoMemInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoMemInfoDeserializeJSON))
	return bench
}

// Uptime
func BenchSysinfoUptimeGet(b *testing.B) {
	u := &uptime.Uptime{}
	for i := 0; i < b.N; i++ {
		u.Get()
	}
	_ = u
}

func SysinfoUptimeGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = SysinfoUptime
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeGet))
	return bench
}

func BenchSysinfoUptimeGetFB(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = ufb.Get()
	}
	_ = tmp
}

func SysinfoUptimeGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = SysinfoUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeGetFB))
	return bench
}

func BenchSysinfoUptimeSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	u := &uptime.Uptime{}
	u.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp = ufb.Serialize(u)
	}
	_ = tmp
}

func SysinfoUptimeSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = SysinfoUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeSerializeFB))
	return bench
}

func BenchSysinfoUptimeDeserializeFB(b *testing.B) {
	var u *uptime.Uptime
	b.StopTimer()
	tmp, _ := ufb.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u = ufb.Deserialize(tmp)
	}
	_ = u
}

func SysinfoUptimeDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = SysinfoUptime
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeDeserializeFB))
	return bench
}

func BenchSysinfoUptimeGetJSON(b *testing.B) {
	var tmp []byte
	for i := 0; i < b.N; i++ {
		tmp, _ = ujson.Get()
	}
	_ = tmp
}

func SysinfoUptimeGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = SysinfoUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeGetJSON))
	return bench
}

func BenchSysinfoUptimeDeserializeJSON(b *testing.B) {
	var u *uptime.Uptime
	b.StopTimer()
	tmp, _ := ujson.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = ujson.Deserialize(tmp)
	}
	_ = u
}

func SysinfoUptimeDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = SysinfoUptime
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchSysinfoUptimeDeserializeJSON))
	return bench
}
