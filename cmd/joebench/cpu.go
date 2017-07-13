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
	"github.com/mohae/joefriday/cpu/cpuinfo"
	infofb "github.com/mohae/joefriday/cpu/cpuinfo/flat"
	infojson "github.com/mohae/joefriday/cpu/cpuinfo/json"
	"github.com/mohae/joefriday/cpu/cpustats"
	statsfb "github.com/mohae/joefriday/cpu/cpustats/flat"
	statsjson "github.com/mohae/joefriday/cpu/cpustats/json"
	"github.com/mohae/joefriday/cpu/cpuutil"
	utilfb "github.com/mohae/joefriday/cpu/cpuutil/flat"
	utiljson "github.com/mohae/joefriday/cpu/cpuutil/json"
)

const (
	CPUInfo  = "CPU Info"
	CPUStats = "CPU Stats"
	CPUUtil  = "CPU Utilization"
)

func runCPUBenchmarks(bench benchutil.Benchmarker) {
	b := CPUInfoGet()
	bench.Append(b)

	b = CPUInfoGetFB()
	bench.Append(b)

	b = CPUInfoSerializeFB()
	bench.Append(b)

	b = CPUInfoDeserializeFB()
	bench.Append(b)

	b = CPUInfoGetJSON()
	bench.Append(b)

	b = CPUInfoSerializeJSON()
	bench.Append(b)

	b = CPUInfoDeserializeJSON()
	bench.Append(b)

	b = CPUStatsGet()
	bench.Append(b)

	b = CPUStatsGetFB()
	bench.Append(b)

	b = CPUStatsSerializeFB()
	bench.Append(b)

	b = CPUStatsDeserializeFB()
	bench.Append(b)

	b = CPUStatsGetJSON()
	bench.Append(b)

	b = CPUStatsSerializeJSON()
	bench.Append(b)

	b = CPUStatsDeserializeJSON()
	bench.Append(b)

	b = CPUUtilGet()
	bench.Append(b)

	b = CPUUtilGetFB()
	bench.Append(b)

	b = CPUUtilSerializeFB()
	bench.Append(b)

	b = CPUUtilDeserializeFB()
	bench.Append(b)

	b = CPUUtilGetJSON()
	bench.Append(b)

	b = CPUUtilSerializeJSON()
	bench.Append(b)

	b = CPUUtilDeserializeJSON()
	bench.Append(b)
}

func BenchCPUInfoGet(b *testing.B) {
	var inf *cpuinfo.CPUInfo
	b.StopTimer()
	p, _ := cpuinfo.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func CPUInfoGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUInfo
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoGet))
	return bench
}

func BenchCPUInfoGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := infofb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUInfoGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoGetFB))
	return bench
}

func BenchCPUInfoSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := cpuinfo.NewProfiler()
	fct, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = infofb.Serialize(fct)
	}
	_ = tmp
}

func CPUInfoSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoSerializeFB))
	return bench
}

func BenchCPUInfoDeserializeFB(b *testing.B) {
	var inf *cpuinfo.CPUInfo
	b.StopTimer()
	p, _ := infofb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf = infofb.Deserialize(tmp)
	}
	_ = inf
}

func CPUInfoDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoDeserializeFB))
	return bench
}

func BenchCPUInfoGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := infojson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUInfoGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoGetJSON))
	return bench
}

func BenchCPUInfoSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := infojson.NewProfiler()
	fct, _ := p.Profiler.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = infojson.Serialize(fct)
	}
	_ = tmp
}

func CPUInfoSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoSerializeJSON))
	return bench
}

func BenchCPUInfoDeserializeJSON(b *testing.B) {
	var fct *cpuinfo.CPUInfo
	b.StopTimer()
	p, _ := infojson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = infojson.Deserialize(tmp)
	}
	_ = fct
}

func CPUInfoDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUInfoDeserializeJSON))
	return bench
}

// Stats
func BenchCPUStatsGet(b *testing.B) {
	var sts *cpustats.CPUStats
	b.StopTimer()
	p, _ := cpustats.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sts, _ = p.Get()
	}
	_ = sts
}

func CPUStatsGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUStats
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGet))
	return bench
}

func BenchCPUStatsGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := statsfb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUStatsGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUStats
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGetFB))
	return bench
}

func BenchCPUStatsSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := cpustats.NewProfiler()
	sts, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = statsfb.Serialize(sts)
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
	var sts *cpustats.CPUStats
	b.StopTimer()
	p, _ := statsfb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sts = statsfb.Deserialize(tmp)
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
	p, _ := statsjson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUStatsGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUStats
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUStatsGetJSON))
	return bench
}

func BenchCPUStatsSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := cpustats.NewProfiler()
	sts, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = statsjson.Serialize(sts)
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
	var sts *cpustats.CPUStats
	b.StopTimer()
	p, _ := statsjson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sts, _ = statsjson.Deserialize(tmp)
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
func BenchCPUUtilGet(b *testing.B) {
	var u *cpuutil.CPUUtil
	b.StopTimer()
	p, _ := cpuutil.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func CPUUtilGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = CPUUtil
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilGet))
	return bench
}

func BenchCPUUtilGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := utilfb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUUtilGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilGetFB))
	return bench
}

func BenchCPUUtilSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := cpuutil.NewProfiler()
	u, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = utilfb.Serialize(u)
	}
	_ = tmp
}

func CPUUtilSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilSerializeFB))
	return bench
}

func BenchCPUUtilDeserializeFB(b *testing.B) {
	var u *cpuutil.CPUUtil
	b.StopTimer()
	p, _ := utilfb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u = utilfb.Deserialize(tmp)
	}
	_ = u
}

func CPUUtilDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = CPUUtil
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilDeserializeFB))
	return bench
}

func BenchCPUUtilGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := utiljson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func CPUUtilGetJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilGetJSON))
	return bench
}

func BenchCPUUtilSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := cpuutil.NewProfiler()
	u, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = utiljson.Serialize(u)
	}
	_ = tmp
}

func CPUUtilSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilSerializeJSON))
	return bench
}

func BenchCPUUtilDeserializeJSON(b *testing.B) {
	var u *cpuutil.CPUUtil
	b.StopTimer()
	p, _ := utiljson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = utiljson.Deserialize(tmp)
	}
	_ = u
}

func CPUUtilDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = CPUUtil
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchCPUUtilDeserializeJSON))
	return bench
}
