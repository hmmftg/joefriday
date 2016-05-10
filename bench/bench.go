// Copyright 2016 The JoeFriday authors.
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
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mohae/benchutil"
	"github.com/mohae/joefriday/bench/cpu"
	"github.com/mohae/joefriday/bench/mem"
	"github.com/mohae/joefriday/bench/net"
	"github.com/mohae/joefriday/bench/platform"
)

// flags
var (
	output string
	format string
)

func init() {
	flag.StringVar(&output, "output", "stdout", "output destination")
	flag.StringVar(&output, "o", "stdout", "output destination (short)")
	flag.StringVar(&format, "format", "txt", "format of output")
	flag.StringVar(&format, "f", "txt", "format of output")
}

func main() {
	flag.Parse()
	done := make(chan struct{})
	go benchutil.Dot(done)

	// set the output
	var w io.Writer
	var err error
	switch output {
	case "stdout":
		w = os.Stdout
	default:
		w, err = os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer w.(*os.File).Close()
	}
	// get the benchmark for the desired format
	// process the output
	var bench benchutil.Benchmarker
	switch format {
	case "csv":
		bench = benchutil.NewCSVBench(w)
	case "md":
		bench = benchutil.NewMDBench(w)
	default:
		bench = benchutil.NewStringBench(w)
	}

	// CPU
	runCPUBenches(bench)

	// Mem
	runMemBenches(bench)

	// Net
	runNetBenches(bench)

	// Platform
	runPlatformBenches(bench)

	fmt.Println("")
	fmt.Println("generating output...")
	err = bench.Out()
	if err != nil {
		fmt.Printf("error generating output: %s\n", err)
	}
}

func runCPUBenches(bench benchutil.Benchmarker) {
	b := cpu.JoeFridayGetFacts()
	bench.Add(b)

	b = cpu.JoeFridayGetStats()
	bench.Add(b)

	b = cpu.DataDogGohaiCPU()
	bench.Add(b)

	b = cpu.ShirouGopsutilInfoStat()
	bench.Add(b)

	b = cpu.ShirouGopsutilTimeStat()
	bench.Add(b)
}

func runMemBenches(bench benchutil.Benchmarker) {
	b := mem.JoeFridayGetMemInfo()
	bench.Add(b)

	b = mem.JoeFridayGetSysinfoMemInfo()
	bench.Add(b)

	b = mem.CloudFoundryGoSigarMem()
	bench.Add(b)

	b = mem.DataDogGohaiMem()
	bench.Add(b)

	b = mem.GuillermoMemInfo()
	bench.Add(b)

	b = mem.ShirouGopsutilMem()
	bench.Add(b)
}

func runNetBenches(bench benchutil.Benchmarker) {
	b := net.JoeFridayGetInfo()
	bench.Add(b)

	b = net.JoeFridayGetUsage()
	bench.Add(b)

	b = net.DataDogGohaiNetwork()
	bench.Add(b)

	b = net.ShirouGopsutilNetInterfaces()
	bench.Add(b)

	b = net.ShirouGopsutilIOCounters()
	bench.Add(b)
}

func runPlatformBenches(bench benchutil.Benchmarker) {
	b := platform.JoeFridayGetKernel()
	bench.Add(b)

	b = platform.JoeFridayGetRelease()
	bench.Add(b)

	b = platform.DataDogGohaiplatform()
	bench.Add(b)

	b = platform.JoeFridayGetLoadAvg()
	bench.Add(b)

	b = platform.JoeFridayGetSysinfoLoadAvg()
	bench.Add(b)

	b = platform.CloudFoundryGoSigarLoadAverage()
	bench.Add(b)

	b = platform.ShirouGopsutilLoadAvg()
	bench.Add(b)

	b = platform.ShirouGopsutilLoadMisc()
	bench.Add(b)

	b = platform.JoeFridayGetUptime()
	bench.Add(b)

	b = platform.JoeFridayGetSysinfoUptime()
	bench.Add(b)

	b = platform.CloudFoundryGoSigarUptime()
	bench.Add(b)
}
