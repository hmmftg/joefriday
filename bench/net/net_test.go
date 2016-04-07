package net

import (
	"encoding/json"
	"testing"

	"github.com/DataDog/gohai/network"
	joeinfo "github.com/mohae/joefriday/net/info"
	joestructs "github.com/mohae/joefriday/net/structs"
	joeusage "github.com/mohae/joefriday/net/usage"
	gopsutilnet "github.com/shirou/gopsutil/net"
)

func BenchmarkJoeFridayGetInfo(b *testing.B) {
	var inf *joestructs.Info
	b.StopTimer()
	p, _ := joeinfo.New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		inf, _ = p.Get()
	}
	_ = inf
}

// don't bench usage for now, as it taks 1 second per with current
// implementation.
/*
func BenchmarkJoeFridayGetUsage(b *testing.B) {
	var u *joestructs.Info
	b.StopTimer()
	p, _ := joeusage.New()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		u, _ = p.Get()
	}
	_ = u
}
*/

func BenchmarkDataDogGohaiNetwork(b *testing.B) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &network.Network{}
	var c interface{}
	for i := 0; i < b.N; i++ {
		c, _ = collector.Collect()
	}
	_ = c
}

func BenchmarkShirouGopsutilNetInterfaces(b *testing.B) {
	var st []gopsutilnet.NetInterfaceStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilnet.NetInterfaces()
	}
	_ = st
}

func BenchmarkShirouGopsutilTimeStat(b *testing.B) {
	var st []gopsutilnet.NetIOCountersStat
	for i := 0; i < b.N; i++ {
		st, _ = gopsutilnet.NetIOCounters(true)
	}
	_ = st
}

// These tests exist to print out the data that is collected; not to test the
// methods themselves.  Run with the -v flag.
func TestJoeFridayGetInfo(t *testing.T) {
	prof, err := joeinfo.New()
	if err != nil {
		t.Error(err)
		return
	}
	inf, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(inf, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestJoeFridayGetUsage(t *testing.T) {
	prof, err := joeusage.New()
	if err != nil {
		t.Error(err)
		return
	}
	u, err := prof.Get()
	if err != nil {
		t.Error(err)
		return
	}
	p, err := json.MarshalIndent(u, "", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%s\n", string(p))
}

func TestDataDogGohaiNetwork(t *testing.T) {
	type Collector interface {
		Name() string
		Collect() (interface{}, error)
	}
	var collector = &network.Network{}
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

func TestShirouGopsutilNetInterfaces(t *testing.T) {
	st, err := gopsutilnet.NetInterfaces()
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

func TestShirouGopsutilTimeStat(t *testing.T) {
	st, err := gopsutilnet.NetIOCounters(true)
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
