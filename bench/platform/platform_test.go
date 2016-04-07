package cpu

import (
	"encoding/json"
	"testing"

	"github.com/DataDog/gohai/platform"
	joekernel "github.com/mohae/joefriday/platform/kernel"
	joerelease "github.com/mohae/joefriday/platform/release"
)

func BenchmarkJoeFridayGetKernel(b *testing.B) {
	var fct *joekernel.Kernel
	b.StopTimer()
	p, _ := joekernel.New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		fct, _ = p.Get()
	}
	_ = fct
}

func BenchmarkJoeFridayGetReleases(b *testing.B) {
	var st *joerelease.Release
	b.StopTimer()
	p, _ := joerelease.New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		st, _ = p.Get()
	}
	_ = st
}

func BenchmarkDataDogGohaiplatform(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &platform.Platform{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

// These tests exist to print out the data that is collected; not to test the
// methods themselves.  Run with the -v flag.
func TestJoeFridayGetKernel(t *testing.T) {
	prof, err := joekernel.New()
	if err != nil {
		t.Error(err)
		return
	}
	fct, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(fct, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestJoeFridayGetRelease(t *testing.T) {
	prof, err := joerelease.New()
	if err != nil {
		t.Error(err)
		return
	}
	st, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(st, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestDataDogGohaiPlatform(t *testing.T) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &platform.Platform{}
	c, err := collector.Collect()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}
