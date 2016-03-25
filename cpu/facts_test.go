package cpu

import "testing"

func TestFactsSerialize(t *testing.T) {
	facts, err := GetFacts()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	factsD := DeserializeFlat(facts.SerializeFlat())
	if factsD.Timestamp != facts.Timestamp {
		t.Errorf("timestamp: got %d; want %d", factsD.Timestamp, facts.Timestamp)
	}
	for i := 0; i < len(facts.CPU); i++ {
		if facts.CPU[i].Processor != factsD.CPU[i].Processor {
			t.Errorf("Fact: got %v; want %v", factsD.CPU[i].Processor, facts.CPU[i].Processor)
		}
		if facts.CPU[i].VendorID != factsD.CPU[i].VendorID {
			t.Errorf("VendorID: got %v; want %v", factsD.CPU[i].VendorID, facts.CPU[i].VendorID)
		}
		if facts.CPU[i].CPUFamily != factsD.CPU[i].CPUFamily {
			t.Errorf("CPUFamily: got %v; want %v", factsD.CPU[i].CPUFamily, facts.CPU[i].CPUFamily)
		}
		if facts.CPU[i].Model != factsD.CPU[i].Model {
			t.Errorf("Model: got %v; want %v", factsD.CPU[i].Model, facts.CPU[i].Model)
		}
		if facts.CPU[i].ModelName != factsD.CPU[i].ModelName {
			t.Errorf("ModelName: got %v; want %v", factsD.CPU[i].ModelName, facts.CPU[i].ModelName)
		}
		if facts.CPU[i].Stepping != factsD.CPU[i].Stepping {
			t.Errorf("Stepping: got %v; want %v", factsD.CPU[i].Stepping, facts.CPU[i].Stepping)
		}
		if facts.CPU[i].Microcode != factsD.CPU[i].Microcode {
			t.Errorf("Microcode: got %v; want %v", factsD.CPU[i].Microcode, facts.CPU[i].Microcode)
		}
		if facts.CPU[i].CPUMHz != factsD.CPU[i].CPUMHz {
			t.Errorf("CPUMHz: got %v; want %v", factsD.CPU[i].CPUMHz, facts.CPU[i].CPUMHz)
		}
		if facts.CPU[i].CacheSize != factsD.CPU[i].CacheSize {
			t.Errorf("CacheSize: got %v; want %v", factsD.CPU[i].CacheSize, facts.CPU[i].CacheSize)
		}
		if facts.CPU[i].PhysicalID != factsD.CPU[i].PhysicalID {
			t.Errorf("PhysicalID: got %v; want %v", factsD.CPU[i].PhysicalID, facts.CPU[i].PhysicalID)
		}
		if facts.CPU[i].Siblings != factsD.CPU[i].Siblings {
			t.Errorf("Siblings: got %v; want %v", factsD.CPU[i].Siblings, facts.CPU[i].Siblings)
		}
		if facts.CPU[i].CoreID != factsD.CPU[i].CoreID {
			t.Errorf("CoreID: got %v; want %v", factsD.CPU[i].CoreID, facts.CPU[i].CoreID)
		}
		if facts.CPU[i].CPUCores != factsD.CPU[i].CPUCores {
			t.Errorf("CPUCores: got %v; want %v", factsD.CPU[i].CPUCores, facts.CPU[i].CPUCores)
		}
		if facts.CPU[i].ApicID != factsD.CPU[i].ApicID {
			t.Errorf("ApicID: got %v; want %v", factsD.CPU[i].ApicID, facts.CPU[i].ApicID)
		}
		if facts.CPU[i].InitialApicID != factsD.CPU[i].InitialApicID {
			t.Errorf("InitialApicID: got %v; want %v", factsD.CPU[i].InitialApicID, facts.CPU[i].InitialApicID)
		}
		if facts.CPU[i].FPU != factsD.CPU[i].FPU {
			t.Errorf("FPU: got %v; want %v", factsD.CPU[i].FPU, facts.CPU[i].FPU)
		}
		if facts.CPU[i].FPUException != factsD.CPU[i].FPUException {
			t.Errorf("FPUException: got %v; want %v", factsD.CPU[i].FPUException, facts.CPU[i].FPUException)
		}
		if facts.CPU[i].CPUIDLevel != factsD.CPU[i].CPUIDLevel {
			t.Errorf("CPUIDLevel: got %v; want %v", factsD.CPU[i].CPUIDLevel, facts.CPU[i].CPUIDLevel)
		}
		if facts.CPU[i].WP != factsD.CPU[i].WP {
			t.Errorf("WP: got %v; want %v", factsD.CPU[i].WP, facts.CPU[i].WP)
		}
		if facts.CPU[i].Flags != factsD.CPU[i].Flags {
			t.Errorf("Flags: got %v; want %v", factsD.CPU[i].Flags, facts.CPU[i].Flags)
		}
		if facts.CPU[i].BogoMIPS != factsD.CPU[i].BogoMIPS {
			t.Errorf("BogoMIPS: got %v; want %v", factsD.CPU[i].BogoMIPS, facts.CPU[i].BogoMIPS)
		}
		if facts.CPU[i].CLFlushSize != factsD.CPU[i].CLFlushSize {
			t.Errorf("CLFlushSize: got %v; want %v", factsD.CPU[i].CLFlushSize, facts.CPU[i].CLFlushSize)
		}
		if facts.CPU[i].CacheAlignment != factsD.CPU[i].CacheAlignment {
			t.Errorf("CacheAlignment: got %v; want %v", factsD.CPU[i].CacheAlignment, facts.CPU[i].CacheAlignment)
		}
		if facts.CPU[i].AddressSizes != factsD.CPU[i].AddressSizes {
			t.Errorf("AddressSizes: got %v; want %v", factsD.CPU[i].AddressSizes, facts.CPU[i].AddressSizes)
		}
		if facts.CPU[i].PowerManagement != factsD.CPU[i].PowerManagement {
			t.Errorf("PowerManagement: got %v; want %v", factsD.CPU[i].PowerManagement, facts.CPU[i].PowerManagement)
		}
	}
}
