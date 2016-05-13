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
	runCPUBenches(bench)

	fmt.Println("\ngenerating output...\n")
	err = bench.Out()
	if err != nil {
		fmt.Printf("error generating output: %s\n", err)
	}
}

func runCPUBenches(bench benchutil.Benchmarker) {
	b := JoeFridayGetFacts()
	bench.Add(b)

	b = JoeFridayGetFactsFB()
	bench.Add(b)

	b = JoeFridayFactsSerializeFB()
	bench.Add(b)

	b = JoeFridayFactsDeserializeFB()
	bench.Add(b)

	b = JoeFridayGetFactsJSON()
	bench.Add(b)

	b = JoeFridayFactsSerializeJSON()
	bench.Add(b)

	b = JoeFridayFactsDeserializeJSON()
	bench.Add(b)

	b = JoeFridayGetStats()
	bench.Add(b)

	b = JoeFridayGetStatsFB()
	bench.Add(b)

	b = JoeFridayStatsSerializeFB()
	bench.Add(b)

	b = JoeFridayStatsDeserializeFB()
	bench.Add(b)

	b = JoeFridayGetStatsJSON()
	bench.Add(b)

	b = JoeFridayStatsSerializeJSON()
	bench.Add(b)

	b = JoeFridayStatsDeserializeJSON()
	bench.Add(b)

	b = JoeFridayGetUtilization()
	bench.Add(b)

	b = JoeFridayGetUtilizationFB()
	bench.Add(b)

	b = JoeFridayUtilizationSerializeFB()
	bench.Add(b)

	b = JoeFridayUtilizationDeserializeFB()
	bench.Add(b)

	b = JoeFridayGetUtilizationJSON()
	bench.Add(b)

	b = JoeFridayUtilizationSerializeJSON()
	bench.Add(b)

	b = JoeFridayUtilizationDeserializeJSON()
	bench.Add(b)
}
