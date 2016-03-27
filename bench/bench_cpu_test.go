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

package bench

import (
	"testing"

	fb "github.com/google/flatbuffers/go"
	"github.com/mohae/joefriday/cpu"
)

func BenchmarkCPUGetFacts(b *testing.B) {
	var val *cpu.Facts
	for i := 0; i < b.N; i++ {
		val, _ = cpu.GetFacts()
	}
	_ = val
}

func BenchmarkCPUFactsSerializeSlat(b *testing.B) {
	var val *cpu.Facts
	var p []byte
	b.StopTimer()
	val, _ = cpu.GetFacts()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = val.SerializeFlat()
	}
	_ = p
}

func BenchmarkCPUFactsDeSerialize(b *testing.B) {
	var val *cpu.Facts
	var p []byte
	b.StopTimer()
	val, _ = cpu.GetFacts()
	p = val.SerializeFlat()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val = cpu.DeserializeFlat(p)
	}
	_ = val
}

func BenchmarkCPUStats(b *testing.B) {
	var val cpu.Stats
	for i := 0; i < b.N; i++ {
		val, _ = cpu.GetStats()
	}
	_ = val
}

func BenchmarkCPUStatsSerializeFlat(b *testing.B) {
	var val cpu.Stats
	var p []byte
	b.StopTimer()
	val, _ = cpu.GetStats()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = val.SerializeFlat()
	}
	_ = p
}

func BenchmarkCPUStatsSerializeFlatBuilder(b *testing.B) {
	var val cpu.Stats
	var p []byte
	bldr := fb.NewBuilder(0)
	b.StopTimer()
	val, _ = cpu.GetStats()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bldr.Reset()
		p = val.SerializeFlatBuilder(bldr)
	}
	_ = p
}

func BenchmarkCPUDeSerializeStatsFlat(b *testing.B) {
	var val cpu.Stats
	var p []byte
	b.StopTimer()
	val, _ = cpu.GetStats()
	p = val.SerializeFlat()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val = cpu.DeserializeStatsFlat(p)
	}
	_ = val
}

func BenchmarkCPUUtilizationSerializeFlat(b *testing.B) {
	var val cpu.Utilization
	var p []byte
	b.StopTimer()
	val, _ = cpu.GetUtilization()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = val.SerializeFlat()
	}
	_ = p
}

func BenchmarkCPUUtilizationSerializeFlatBuilder(b *testing.B) {
	var val cpu.Utilization
	var p []byte
	bldr := fb.NewBuilder(0)
	b.StopTimer()
	val, _ = cpu.GetUtilization()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bldr.Reset()
		p = val.SerializeFlatBuilder(bldr)
	}
	_ = p
}

func BenchmarkCPUDeSerializeUtilizationFlat(b *testing.B) {
	var val cpu.Utilization
	var p []byte
	b.StopTimer()
	val, _ = cpu.GetUtilization()
	p = val.SerializeFlat()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val = cpu.DeserializeUtilizationFlat(p)
	}
	_ = val
}
