package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mohae/benchutil"
)

// flags
var (
	output         string
	format         string
	section        bool
	sectionHeaders bool
	nameSections   bool
)

func init() {
	flag.StringVar(&output, "output", "stdout", "output destination")
	flag.StringVar(&output, "o", "stdout", "output destination (short)")
	flag.StringVar(&format, "format", "txt", "format of output")
	flag.StringVar(&format, "f", "txt", "format of output")
	flag.BoolVar(&nameSections, "namesections", false, "use group as section name: some restrictions apply")
	flag.BoolVar(&nameSections, "n", false, "use group as section name: some restrictions apply")
	flag.BoolVar(&section, "sections", false, "don't separate groups of tests into sections")
	flag.BoolVar(&section, "s", false, "don't separate groups of tests into sections")
	flag.BoolVar(&sectionHeaders, "sectionheader", false, "if there are sections, add a section header row")
	flag.BoolVar(&sectionHeaders, "h", false, "if there are sections, add a section header row")
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
		bench.(*benchutil.MDBench).GroupAsSectionName = nameSections
	default:
		bench = benchutil.NewStringBench(w)
	}
	bench.SectionPerGroup(section)
	bench.SectionHeaders(sectionHeaders)
	// CPU
	runCPUBenchmarks(bench)

	fmt.Println("\ngenerating output...\n")
	err = bench.Out()
	if err != nil {
		fmt.Printf("error generating output: %s\n", err)
	}
}

func runCPUBenchmarks(bench benchutil.Benchmarker) {
	b := CPUGetFacts()
	bench.Add(b)

	b = CPUGetFactsFB()
	bench.Add(b)

	b = CPUFactsSerializeFB()
	bench.Add(b)

	b = CPUFactsDeserializeFB()
	bench.Add(b)

	b = CPUGetFactsJSON()
	bench.Add(b)

	b = CPUFactsSerializeJSON()
	bench.Add(b)

	b = CPUFactsDeserializeJSON()
	bench.Add(b)

	b = CPUGetStats()
	bench.Add(b)

	b = CPUGetStatsFB()
	bench.Add(b)

	b = CPUStatsSerializeFB()
	bench.Add(b)

	b = CPUStatsDeserializeFB()
	bench.Add(b)

	b = CPUGetStatsJSON()
	bench.Add(b)

	b = CPUStatsSerializeJSON()
	bench.Add(b)

	b = CPUStatsDeserializeJSON()
	bench.Add(b)

	b = CPUGetUtilization()
	bench.Add(b)

	b = CPUGetUtilizationFB()
	bench.Add(b)

	b = CPUUtilizationSerializeFB()
	bench.Add(b)

	b = CPUUtilizationDeserializeFB()
	bench.Add(b)

	b = CPUGetUtilizationJSON()
	bench.Add(b)

	b = CPUUtilizationSerializeJSON()
	bench.Add(b)

	b = CPUUtilizationDeserializeJSON()
	bench.Add(b)
}
