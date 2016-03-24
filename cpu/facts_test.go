package cpu

import "testing"

func TestFactsSerialize(t *testing.T) {
	facts, err := GetFacts()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	factsD := Deserialize(facts.Serialize())
	if factsD.Timestamp != facts.Timestamp {
		t.Errorf("timestamp: got %d; want %d", factsD.Timestamp, facts.Timestamp)
	}
	for i := 0; i < len(facts.CPUs); i++ {
		if facts.CPUs[i].Processor != factsD.CPUs[i].Processor {
			t.Errorf("Fact: got %v; want %v", factsD.CPUs[i].Processor, facts.CPUs[i].Processor)
		}
		if facts.CPUs[i].VendorID != factsD.CPUs[i].VendorID {
			t.Errorf("VendorID: got %v; want %v", factsD.CPUs[i].VendorID, facts.CPUs[i].VendorID)
		}
		if facts.CPUs[i].CPUFamily != factsD.CPUs[i].CPUFamily {
			t.Errorf("CPUFamily: got %v; want %v", factsD.CPUs[i].CPUFamily, facts.CPUs[i].CPUFamily)
		}
		if facts.CPUs[i].Model != factsD.CPUs[i].Model {
			t.Errorf("Model: got %v; want %v", factsD.CPUs[i].Model, facts.CPUs[i].Model)
		}
		if facts.CPUs[i].ModelName != factsD.CPUs[i].ModelName {
			t.Errorf("ModelName: got %v; want %v", factsD.CPUs[i].ModelName, facts.CPUs[i].ModelName)
		}
		if facts.CPUs[i].Stepping != factsD.CPUs[i].Stepping {
			t.Errorf("Stepping: got %v; want %v", factsD.CPUs[i].Stepping, facts.CPUs[i].Stepping)
		}
		if facts.CPUs[i].Microcode != factsD.CPUs[i].Microcode {
			t.Errorf("Microcode: got %v; want %v", factsD.CPUs[i].Microcode, facts.CPUs[i].Microcode)
		}
		if facts.CPUs[i].CPUMHz != factsD.CPUs[i].CPUMHz {
			t.Errorf("CPUMHz: got %v; want %v", factsD.CPUs[i].CPUMHz, facts.CPUs[i].CPUMHz)
		}
		if facts.CPUs[i].CacheSize != factsD.CPUs[i].CacheSize {
			t.Errorf("CacheSize: got %v; want %v", factsD.CPUs[i].CacheSize, facts.CPUs[i].CacheSize)
		}
		if facts.CPUs[i].PhysicalID != factsD.CPUs[i].PhysicalID {
			t.Errorf("PhysicalID: got %v; want %v", factsD.CPUs[i].PhysicalID, facts.CPUs[i].PhysicalID)
		}
		if facts.CPUs[i].Siblings != factsD.CPUs[i].Siblings {
			t.Errorf("Siblings: got %v; want %v", factsD.CPUs[i].Siblings, facts.CPUs[i].Siblings)
		}
		if facts.CPUs[i].CoreID != factsD.CPUs[i].CoreID {
			t.Errorf("CoreID: got %v; want %v", factsD.CPUs[i].CoreID, facts.CPUs[i].CoreID)
		}
		if facts.CPUs[i].CPUCores != factsD.CPUs[i].CPUCores {
			t.Errorf("CPUCores: got %v; want %v", factsD.CPUs[i].CPUCores, facts.CPUs[i].CPUCores)
		}
		if facts.CPUs[i].ApicID != factsD.CPUs[i].ApicID {
			t.Errorf("ApicID: got %v; want %v", factsD.CPUs[i].ApicID, facts.CPUs[i].ApicID)
		}
		if facts.CPUs[i].InitialApicID != factsD.CPUs[i].InitialApicID {
			t.Errorf("InitialApicID: got %v; want %v", factsD.CPUs[i].InitialApicID, facts.CPUs[i].InitialApicID)
		}
		if facts.CPUs[i].FPU != factsD.CPUs[i].FPU {
			t.Errorf("FPU: got %v; want %v", factsD.CPUs[i].FPU, facts.CPUs[i].FPU)
		}
		if facts.CPUs[i].FPUException != factsD.CPUs[i].FPUException {
			t.Errorf("FPUException: got %v; want %v", factsD.CPUs[i].FPUException, facts.CPUs[i].FPUException)
		}
		if facts.CPUs[i].CPUIDLevel != factsD.CPUs[i].CPUIDLevel {
			t.Errorf("CPUIDLevel: got %v; want %v", factsD.CPUs[i].CPUIDLevel, facts.CPUs[i].CPUIDLevel)
		}
		if facts.CPUs[i].WP != factsD.CPUs[i].WP {
			t.Errorf("WP: got %v; want %v", factsD.CPUs[i].WP, facts.CPUs[i].WP)
		}
		if facts.CPUs[i].Flags != factsD.CPUs[i].Flags {
			t.Errorf("Flags: got %v; want %v", factsD.CPUs[i].Flags, facts.CPUs[i].Flags)
		}
		if facts.CPUs[i].BogoMIPS != factsD.CPUs[i].BogoMIPS {
			t.Errorf("BogoMIPS: got %v; want %v", factsD.CPUs[i].BogoMIPS, facts.CPUs[i].BogoMIPS)
		}
		if facts.CPUs[i].CLFlushSize != factsD.CPUs[i].CLFlushSize {
			t.Errorf("CLFlushSize: got %v; want %v", factsD.CPUs[i].CLFlushSize, facts.CPUs[i].CLFlushSize)
		}
		if facts.CPUs[i].CacheAlignment != factsD.CPUs[i].CacheAlignment {
			t.Errorf("CacheAlignment: got %v; want %v", factsD.CPUs[i].CacheAlignment, facts.CPUs[i].CacheAlignment)
		}
		if facts.CPUs[i].AddressSizes != factsD.CPUs[i].AddressSizes {
			t.Errorf("AddressSizes: got %v; want %v", factsD.CPUs[i].AddressSizes, facts.CPUs[i].AddressSizes)
		}
		if facts.CPUs[i].PowerManagement != factsD.CPUs[i].PowerManagement {
			t.Errorf("PowerManagement: got %v; want %v", factsD.CPUs[i].PowerManagement, facts.CPUs[i].PowerManagement)
		}
	}
}
