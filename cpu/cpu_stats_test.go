package cpu

import "testing"

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
