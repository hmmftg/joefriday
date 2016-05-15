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
	"github.com/mohae/joefriday/net/info"
	ifb "github.com/mohae/joefriday/net/info/flat"
	ijson "github.com/mohae/joefriday/net/info/json"
	"github.com/mohae/joefriday/net/structs"
	"github.com/mohae/joefriday/net/usage"
	ufb "github.com/mohae/joefriday/net/usage/flat"
	ujson "github.com/mohae/joefriday/net/usage/json"
)

const (
	NetInfo  = "Net Info"
	NetUsage = "Net Usage"
)

func BenchNetInfoGet(b *testing.B) {
	var inf *structs.Info
	b.StopTimer()
	p, _ := info.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

func NetInfoGet() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = NetInfo
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetInfoGet))
	return bench
}

func BenchNetInfoGetFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ifb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func NetInfoGetFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = NetInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetInfoGetFB))
	return bench
}

func BenchNetInfoSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := info.NewProfiler()
	inf, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ifb.Serialize(inf)
	}
	_ = tmp
}

func NetInfoSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = NetInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetInfoSerializeFB))
	return bench
}

func BenchNetInfoDeserializeFB(b *testing.B) {
	var inf *structs.Info
	b.StopTimer()
	p, _ := ifb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf = ifb.Deserialize(tmp)
	}
	_ = inf
}

func NetInfoDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = NetInfo
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetInfoDeserializeFB))
	return bench
}

func BenchNetInfoGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ijson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func NetInfoGetSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = NetInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetInfoGetJSON))
	return bench
}

func BenchNetInfoSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := info.NewProfiler()
	sts, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ijson.Serialize(sts)
	}
	_ = tmp
}

func NetInfoSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = NetInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetInfoSerializeJSON))
	return bench
}

func BenchNetInfoDeserializeJSON(b *testing.B) {
	var inf *structs.Info
	b.StopTimer()
	p, _ := ijson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = ijson.Deserialize(tmp)
	}
	_ = inf
}

func NetInfoDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = NetInfo
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetInfoDeserializeJSON))
	return bench
}

// Usage
func BenchNetGetUsage(b *testing.B) {
	var u *structs.Usage
	b.StopTimer()
	p, _ := usage.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}

func NetGetUsage() benchutil.Bench {
	bench := benchutil.NewBench("Get")
	bench.Group = NetUsage
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetGetUsage))
	return bench
}

func BenchNetGetUsageFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ufb.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func NetGetUsageFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Get")
	bench.Group = NetUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetGetUsageFB))
	return bench
}

func BenchNetUsageSerializeFB(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := usage.NewProfiler()
	u, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ufb.Serialize(u)
	}
	_ = tmp
}

func NetUsageSerializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Serialize")
	bench.Group = NetUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageSerializeFB))
	return bench
}

func BenchNetUsageDeserializeFB(b *testing.B) {
	var u *structs.Usage
	b.StopTimer()
	p, _ := ufb.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u = ufb.Deserialize(tmp)
	}
	_ = u
}

func NetUsageDeserializeFB() benchutil.Bench {
	bench := benchutil.NewBench("flat.Deserialize")
	bench.Group = NetUsage
	bench.Desc = Flat
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageDeserializeFB))
	return bench
}

func BenchNetUsageGetJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := ujson.NewProfiler()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = p.Get()
	}
	_ = tmp
}

func NetGetUsageJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Get")
	bench.Group = NetUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageGetJSON))
	return bench
}

func BenchNetUsageSerializeJSON(b *testing.B) {
	var tmp []byte
	b.StopTimer()
	p, _ := usage.NewProfiler()
	u, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tmp, _ = ujson.Serialize(u)
	}
	_ = tmp
}

func NetUsageSerializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Serialize")
	bench.Group = NetUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageSerializeJSON))
	return bench
}

func BenchNetUsageDeserializeJSON(b *testing.B) {
	var u *structs.Usage
	b.StopTimer()
	p, _ := ujson.NewProfiler()
	tmp, _ := p.Get()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = ujson.Deserialize(tmp)
	}
	_ = u
}

func NetUsageDeserializeJSON() benchutil.Bench {
	bench := benchutil.NewBench("json.Deserialize")
	bench.Group = NetUsage
	bench.Desc = JSON
	bench.Result = benchutil.ResultFromBenchmarkResult(testing.Benchmark(BenchNetUsageDeserializeJSON))
	return bench
}
