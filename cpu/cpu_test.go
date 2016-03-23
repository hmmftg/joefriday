package cpu

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	err := Init()
	if err != nil {
		t.Errorf("expected error to be nil; got %s", err)
	}
	if CLK_TCK == 0 {
		t.Errorf("got %d, want a value > 0", CLK_TCK)
	}
}

func TestSerializeDeserializeStats(t *testing.T) {
	Init()
	stats, err := GetStats()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	b := stats.SerializeFlat()
	statsD := DeserializeStatsFlat(b)
	if stats.ClkTck != statsD.ClkTck {
		t.Errorf("ClkTck: got %s; want %s", statsD.ClkTck, stats.ClkTck)
	}
	if stats.Timestamp != statsD.Timestamp {
		t.Errorf("Timestamp: got %s; want %s", statsD.Timestamp, stats.Timestamp)
	}
	if stats.Ctxt != statsD.Ctxt {
		t.Errorf("Ctxt: got %s; want %s", statsD.Ctxt, stats.Ctxt)
	}
	if stats.BTime != statsD.BTime {
		t.Errorf("BTime: got %s; want %s", statsD.BTime, stats.BTime)
	}
	if stats.Processes != statsD.Processes {
		t.Errorf("Processes: got %s; want %s", statsD.Processes, stats.Processes)
	}
	for i := 0; i < len(stats.CPUs); i++ {
		if stats.CPUs[i].CPU != statsD.CPUs[i].CPU {
			t.Errorf("CPU %d: CPU: got %s; want %s", i, statsD.CPUs[i].CPU, stats.CPUs[i].CPU)
		}
		if stats.CPUs[i].User != statsD.CPUs[i].User {
			t.Errorf("CPU %d: User: got %s; want %s", i, statsD.CPUs[i].User, stats.CPUs[i].User)
		}
		if stats.CPUs[i].Nice != statsD.CPUs[i].Nice {
			t.Errorf("CPU %d: Nice: got %s; want %s", i, statsD.CPUs[i].Nice, stats.CPUs[i].Nice)
		}
		if stats.CPUs[i].System != statsD.CPUs[i].System {
			t.Errorf("CPU %d: System: got %s; want %s", i, statsD.CPUs[i].System, stats.CPUs[i].System)
		}
		if stats.CPUs[i].Idle != statsD.CPUs[i].Idle {
			t.Errorf("CPU %d: Idle: got %s; want %s", i, statsD.CPUs[i].Idle, stats.CPUs[i].Idle)
		}
		if stats.CPUs[i].IOWait != statsD.CPUs[i].IOWait {
			t.Errorf("CPU %d: IOWait: got %s; want %s", i, statsD.CPUs[i].IOWait, stats.CPUs[i].IOWait)
		}
		if stats.CPUs[i].IRQ != statsD.CPUs[i].IRQ {
			t.Errorf("CPU %d: IRQ: got %s; want %s", i, statsD.CPUs[i].IRQ, stats.CPUs[i].IRQ)
		}
		if stats.CPUs[i].SoftIRQ != statsD.CPUs[i].SoftIRQ {
			t.Errorf("CPU %d: SoftIRQ: got %s; want %s", i, statsD.CPUs[i].SoftIRQ, stats.CPUs[i].SoftIRQ)
		}
		if stats.CPUs[i].Steal != statsD.CPUs[i].Steal {
			t.Errorf("CPU %d: Steal: got %s; want %s", i, statsD.CPUs[i].Steal, stats.CPUs[i].Steal)
		}
		if stats.CPUs[i].Quest != statsD.CPUs[i].Quest {
			t.Errorf("CPU %d: Quest: got %s; want %s", i, statsD.CPUs[i].Quest, stats.CPUs[i].Quest)
		}
		if stats.CPUs[i].QuestNice != statsD.CPUs[i].QuestNice {
			t.Errorf("CPU %d: QuestNice: got %s; want %s", i, statsD.CPUs[i].QuestNice, stats.CPUs[i].QuestNice)
		}
	}
}

func TestSerializeDeserializeUtilization(t *testing.T) {
	Init()
	u, err := GetUtilization()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	b := u.SerializeFlat()
	uD := DeserializeUtilizationFlat(b)
	if u.Timestamp != uD.Timestamp {
		t.Errorf("Timestamp: got %s; want %s", uD.Timestamp, u.Timestamp)
	}
	if u.CtxtDelta != uD.CtxtDelta {
		t.Errorf("CtxtDelta: got %s; want %s", uD.CtxtDelta, u.CtxtDelta)
	}
	if u.BTimeDelta != uD.BTimeDelta {
		t.Errorf("BTimeDelta: got %s; want %s", uD.BTimeDelta, u.BTimeDelta)
	}
	if u.Processes != uD.Processes {
		t.Errorf("Processes: got %s; want %s", uD.Processes, u.Processes)
	}
	for i := 0; i < len(u.CPUs); i++ {
		if u.CPUs[i].CPU != uD.CPUs[i].CPU {
			t.Errorf("CPU %d: CPU: got %s; want %s", i, uD.CPUs[i].CPU, u.CPUs[i].CPU)
		}
		if u.CPUs[i].Usage != uD.CPUs[i].Usage {
			t.Errorf("CPU %d: Usage: got %s; want %s", i, uD.CPUs[i].Usage, u.CPUs[i].Usage)
		}
		if u.CPUs[i].User != uD.CPUs[i].User {
			t.Errorf("CPU %d: User: got %s; want %s", i, uD.CPUs[i].User, u.CPUs[i].User)
		}
		if u.CPUs[i].Nice != uD.CPUs[i].Nice {
			t.Errorf("CPU %d: Nice: got %s; want %s", i, uD.CPUs[i].Nice, u.CPUs[i].Nice)
		}
		if u.CPUs[i].System != uD.CPUs[i].System {
			t.Errorf("CPU %d: System: got %s; want %s", i, uD.CPUs[i].System, u.CPUs[i].System)
		}
		if u.CPUs[i].Idle != uD.CPUs[i].Idle {
			t.Errorf("CPU %d: Idle: got %s; want %s", i, uD.CPUs[i].Idle, u.CPUs[i].Idle)
		}
		if u.CPUs[i].IOWait != uD.CPUs[i].IOWait {
			t.Errorf("CPU %d: IOWait: got %s; want %s", i, uD.CPUs[i].IOWait, u.CPUs[i].IOWait)
		}
	}
}

func TestUtilizationTicker(t *testing.T) {
	Init()
	out := make(chan Utilization)
	done := make(chan struct{})
	errs := make(chan error)

	ticker := time.NewTicker(time.Duration(2) * time.Second)
	defer ticker.Stop()
	defer close(errs)

	go UtilizationTicker(time.Second, out, done, errs)
testloop:
	for {
		select {
		case <-ticker.C:
			close(done)
			break testloop
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		case u := <-out:
			if u.Timestamp == 0 {
				t.Errorf("expected a timestamp got 0")
				continue
			}
			if len(u.CPUs) == 0 {
				t.Errorf("expected CPU data, got none")
				continue
			}
			for i, v := range u.CPUs {
				if v.CPU == "" {
					t.Errorf("%d: expected CPU to have a CPU ID, got none", i)
					continue
				}
				// only check IDLE: this may fail if on a really busy system
				// but Usage may fail on a non-busy system.
				if fmt.Sprintf("%.1f", v.Idle) == "0.0" {
					t.Errorf("%d: expected Idle to have a non-zero value; it wasn't", i)
				}
			}
		}
	}
}

func TestUtilizationFlatTicker(t *testing.T) {
	out := make(chan []byte)
	done := make(chan struct{})
	errs := make(chan error)

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	defer ticker.Stop()
	defer close(errs)
	go UtilizationFlatTicker(time.Duration(200)*time.Millisecond, out, done, errs)

testloop:
	for {
		select {
		case <-ticker.C:
			done <- struct{}{}
			break testloop
		case err := <-errs:
			t.Errorf("unexpected error: %s", err)
		case b := <-out:
			u := DeserializeUtilizationFlat(b)
			if u.Timestamp == 0 {
				t.Errorf("expected a timestamp got 0")
				continue
			}
			if len(u.CPUs) == 0 {
				t.Errorf("expected CPU data, got none")
				continue
			}
			for i, v := range u.CPUs {
				if v.CPU == "" {
					t.Errorf("%d: expected CPU to have a CPU ID, got none", i)
					continue
				}
				if fmt.Sprintf("%.1f", v.Idle) == "0.0" {
					t.Errorf("%d: expected Idle to have a non-zero value; it wasn't", i)
				}
			}
		}
	}
}
