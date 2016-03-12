// Package bench contains benchmarks for gathering system information.
// These packages are not directly comparable because of the differences
// in what they gather, but I wanted to see some numbers.
//
// This will only work on linux systems due to limitations of
// github.com/mohae/joefriday.
package bench

import (
	"testing"

	"github.com/DataDog/gohai/memory"
	joemem "github.com/mohae/joefriday/mem"
)

func BenchmarkJoeFridayMemInfo(b *testing.B) {
	var inf *joemem.Info
	for i := 0; i < b.N; i++ {
		inf, _ = joemem.GetInfo()
	}
	_ = inf
}

func BenchmarkJoeFridayMemData(b *testing.B) {
	var data []byte
	for i := 0; i < b.N; i++ {
		data, _ = joemem.GetData()
	}
	_ = data
}

func BenchmarkGohaiMem(b *testing.B) {
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
