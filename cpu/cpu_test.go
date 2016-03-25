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
	for i := 0; i < len(stats.CPU); i++ {
		if stats.CPU[i].ID != statsD.CPU[i].ID {
			t.Errorf("CPU %d: ID: got %s; want %s", i, statsD.CPU[i].ID, stats.CPU[i].ID)
		}
		if stats.CPU[i].User != statsD.CPU[i].User {
			t.Errorf("CPU %d: User: got %s; want %s", i, statsD.CPU[i].User, stats.CPU[i].User)
		}
		if stats.CPU[i].Nice != statsD.CPU[i].Nice {
			t.Errorf("CPU %d: Nice: got %s; want %s", i, statsD.CPU[i].Nice, stats.CPU[i].Nice)
		}
		if stats.CPU[i].System != statsD.CPU[i].System {
			t.Errorf("CPU %d: System: got %s; want %s", i, statsD.CPU[i].System, stats.CPU[i].System)
		}
		if stats.CPU[i].Idle != statsD.CPU[i].Idle {
			t.Errorf("CPU %d: Idle: got %s; want %s", i, statsD.CPU[i].Idle, stats.CPU[i].Idle)
		}
		if stats.CPU[i].IOWait != statsD.CPU[i].IOWait {
			t.Errorf("CPU %d: IOWait: got %s; want %s", i, statsD.CPU[i].IOWait, stats.CPU[i].IOWait)
		}
		if stats.CPU[i].IRQ != statsD.CPU[i].IRQ {
			t.Errorf("CPU %d: IRQ: got %s; want %s", i, statsD.CPU[i].IRQ, stats.CPU[i].IRQ)
		}
		if stats.CPU[i].SoftIRQ != statsD.CPU[i].SoftIRQ {
			t.Errorf("CPU %d: SoftIRQ: got %s; want %s", i, statsD.CPU[i].SoftIRQ, stats.CPU[i].SoftIRQ)
		}
		if stats.CPU[i].Steal != statsD.CPU[i].Steal {
			t.Errorf("CPU %d: Steal: got %s; want %s", i, statsD.CPU[i].Steal, stats.CPU[i].Steal)
		}
		if stats.CPU[i].Quest != statsD.CPU[i].Quest {
			t.Errorf("CPU %d: Quest: got %s; want %s", i, statsD.CPU[i].Quest, stats.CPU[i].Quest)
		}
		if stats.CPU[i].QuestNice != statsD.CPU[i].QuestNice {
			t.Errorf("CPU %d: QuestNice: got %s; want %s", i, statsD.CPU[i].QuestNice, stats.CPU[i].QuestNice)
		}
	}
}

func TestSerializeDeserializeUtilization(t *testing.T) {
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
	for i := 0; i < len(u.CPU); i++ {
		if u.CPU[i].ID != uD.CPU[i].ID {
			t.Errorf("CPU %d: ID: got %s; want %s", i, uD.CPU[i].ID, u.CPU[i].ID)
		}
		if u.CPU[i].Usage != uD.CPU[i].Usage {
			t.Errorf("CPU %d: Usage: got %s; want %s", i, uD.CPU[i].Usage, u.CPU[i].Usage)
		}
		if u.CPU[i].User != uD.CPU[i].User {
			t.Errorf("CPU %d: User: got %s; want %s", i, uD.CPU[i].User, u.CPU[i].User)
		}
		if u.CPU[i].Nice != uD.CPU[i].Nice {
			t.Errorf("CPU %d: Nice: got %s; want %s", i, uD.CPU[i].Nice, u.CPU[i].Nice)
		}
		if u.CPU[i].System != uD.CPU[i].System {
			t.Errorf("CPU %d: System: got %s; want %s", i, uD.CPU[i].System, u.CPU[i].System)
		}
		if u.CPU[i].Idle != uD.CPU[i].Idle {
			t.Errorf("CPU %d: Idle: got %s; want %s", i, uD.CPU[i].Idle, u.CPU[i].Idle)
		}
		if u.CPU[i].IOWait != uD.CPU[i].IOWait {
			t.Errorf("CPU %d: IOWait: got %s; want %s", i, uD.CPU[i].IOWait, u.CPU[i].IOWait)
		}
	}
}

func TestUtilizationTicker(t *testing.T) {
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
			if len(u.CPU) == 0 {
				t.Errorf("expected CPU data, got none")
				continue
			}
			for i, v := range u.CPU {
				if v.ID == "" {
					t.Errorf("%d: expected CPU to have an ID, got none", i)
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

func TestUtilizationTickerFlat(t *testing.T) {
	out := make(chan []byte)
	done := make(chan struct{})
	errs := make(chan error)

	ticker := time.NewTicker(time.Duration(1) * time.Second)
	defer ticker.Stop()
	defer close(errs)
	go UtilizationTickerFlat(time.Duration(200)*time.Millisecond, out, done, errs)

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
			if len(u.CPU) == 0 {
				t.Errorf("expected CPU data, got none")
				continue
			}
			for i, v := range u.CPU {
				if v.ID == "" {
					t.Errorf("%d: expected CPU to have an ID, got none", i)
					continue
				}
				if fmt.Sprintf("%.1f", v.Idle) == "0.0" {
					t.Errorf("%d: expected Idle to have a non-zero value; it wasn't", i)
				}
			}
		}
	}
}
