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

// Package bench contains benchmarks for gathering system information.
// These packages are not directly comparable because of the differences
// in what they gather, but I wanted to see some numbers.
//
// This will only work on linux systems due to limitations of
// github.com/mohae/joefriday.
package bench

import (
	"fmt"
	"os"
	"testing"

	"github.com/DataDog/gohai/memory"
	"github.com/EricLagergren/go-gnulib/sysinfo"
	"github.com/cloudfoundry/gosigar"
	gopsutilmem "github.com/shirou/gopsutil/mem"
)

func BenchmarkOSExecCatMemInfo(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoCat()
	}
	_ = inf
}

func BenchmarkOSExecCatMemInfoToJSON(b *testing.B) {
	var inf []byte
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoCatToJSON()
	}
	_ = inf
}

func BenchmarkOSExecCatMemInfoToFlatbuffers(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataCat()
	}
	_ = data
}

func BenchmarkOSExecCatMemInfoToFlatbuffersReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataCatReuseBldr()
	}
	_ = data
}

func BenchmarkReadMemInfo(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoRead()
	}
	_ = inf
}

func BenchmarkReadMemInfoToJSON(b *testing.B) {
	var inf []byte
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoReadToJSON()
	}
	_ = inf
}

func BenchmarkReadMemInfoToFlatbuffers(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataRead()
	}
	_ = data
}

func BenchmarkReadMemInfoToFlatbuffersReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReadReuseBldr()
	}
	_ = data
}

func BenchmarkReadMemInfoReuseBufio(b *testing.B) {
	var inf *MemInfo
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoReadReuseR()
	}
	_ = inf
}

func BenchmarkReadMemInfoToJSONReuseBufio(b *testing.B) {
	var inf []byte
	for i := 0; i < b.N; i++ {
		inf, _ = GetMemInfoReadReuseRToJSON()
	}
	_ = inf
}

func BenchmarkReadMemInfoToFlatbuffersReuseBufio(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReadReuseR()
	}
	_ = data
}

func BenchmarkReadMemDataToFlatbuffersReuseBufioReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemDataReuseRReuseBldr()
	}
	_ = data
}

func BenchmarkReadMemInfoToFlatbuffersReuseBufioReuseBuilder(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemInfoToFlatbuffersReuseBldr()
	}
	_ = data
}

func BenchmarkReadMemInfoToFlatbuffersMinAllocs(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = GetMemInfoToFlatbuffersMinAllocs()
	}
	_ = data
}

func BenchmarkReadMemInfoToFlatbuffersMinAllocsSeek(b *testing.B) {
	var data []byte
	b.StopTimer()
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("couldn't open /proc/meminfo")
	}
	defer f.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		data, _ = GetMemInfoToFlatbuffersMinAllocsSeek(f)
	}
	_ = data
}

func BenchmarkReadMemInfoEmulateCurrentFlatTicker(b *testing.B) {
	var data []byte
	b.StopTimer()
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("couldn't open /proc/meminfo")
	}
	defer f.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		data, _ = GetMemInfoEmulateCurrentFlatTicker(f)
	}
	_ = data
}

func BenchmarkGetMemInfoCurrent(b *testing.B) {
	var val MemInfo
	b.StopTimer()
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("couldn't open /proc/meminfo")
	}
	defer f.Close()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val, _ = GetMemInfoCurrent(f)
	}
	_ = val
}

func BenchmarkDataDogGohaiMem(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &memory.Memory{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func BenchmarkCloudFoundryGoSigarMem(b *testing.B) {
	var mem sigar.Mem
	for i := 0; i < b.N; i++ {
		mem.Get()
	}
	_ = mem
}

func BenchmarkShirouGopsutilMem(b *testing.B) {
	var mem *gopsutilmem.VirtualMemoryStat
	for i := 0; i < b.N; i++ {
		mem, _ = gopsutilmem.VirtualMemory()
	}
	_ = mem
}

func BenchmarkEricLagergrenGnulibSysinfoPhysmemAvailable(b *testing.B) {
	var mem int64
	for i := 0; i < b.N; i++ {
		mem = sysinfo.PhysmemAvailable()

	}
	_ = mem
}

func BenchmarkEricLagergrenGnulibSysinfoPhysmemTotal(b *testing.B) {
	var mem int64
	for i := 0; i < b.N; i++ {
		mem = sysinfo.PhysmemTotal()

	}
	_ = mem
}
