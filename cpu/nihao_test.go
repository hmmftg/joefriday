package cpu

import "testing"

func TestProcessorsSerialize(t *testing.T) {
	procs, err := NiHao()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	procsD := Deserialize(procs.Serialize())
	if procsD.Timestamp != procs.Timestamp {
		t.Errorf("timestamp: got %d; want %d", procsD.Timestamp, procs.Timestamp)
	}
	for i := 0; i < len(procs.Infos); i++ {
		if procs.Infos[i].Processor != procsD.Infos[i].Processor {
			t.Errorf("Processor: got %v; want %v", procsD.Infos[i].Processor, procs.Infos[i].Processor)
		}
		if procs.Infos[i].VendorID != procsD.Infos[i].VendorID {
			t.Errorf("VendorID: got %v; want %v", procsD.Infos[i].VendorID, procs.Infos[i].VendorID)
		}
		if procs.Infos[i].CPUFamily != procsD.Infos[i].CPUFamily {
			t.Errorf("CPUFamily: got %v; want %v", procsD.Infos[i].CPUFamily, procs.Infos[i].CPUFamily)
		}
		if procs.Infos[i].Model != procsD.Infos[i].Model {
			t.Errorf("Model: got %v; want %v", procsD.Infos[i].Model, procs.Infos[i].Model)
		}
		if procs.Infos[i].ModelName != procsD.Infos[i].ModelName {
			t.Errorf("ModelName: got %v; want %v", procsD.Infos[i].ModelName, procs.Infos[i].ModelName)
		}
		if procs.Infos[i].Stepping != procsD.Infos[i].Stepping {
			t.Errorf("Stepping: got %v; want %v", procsD.Infos[i].Stepping, procs.Infos[i].Stepping)
		}
		if procs.Infos[i].Microcode != procsD.Infos[i].Microcode {
			t.Errorf("Microcode: got %v; want %v", procsD.Infos[i].Microcode, procs.Infos[i].Microcode)
		}
		if procs.Infos[i].CPUMHz != procsD.Infos[i].CPUMHz {
			t.Errorf("CPUMHz: got %v; want %v", procsD.Infos[i].CPUMHz, procs.Infos[i].CPUMHz)
		}
		if procs.Infos[i].CacheSize != procsD.Infos[i].CacheSize {
			t.Errorf("CacheSize: got %v; want %v", procsD.Infos[i].CacheSize, procs.Infos[i].CacheSize)
		}
		if procs.Infos[i].PhysicalID != procsD.Infos[i].PhysicalID {
			t.Errorf("PhysicalID: got %v; want %v", procsD.Infos[i].PhysicalID, procs.Infos[i].PhysicalID)
		}
		if procs.Infos[i].Siblings != procsD.Infos[i].Siblings {
			t.Errorf("Siblings: got %v; want %v", procsD.Infos[i].Siblings, procs.Infos[i].Siblings)
		}
		if procs.Infos[i].CoreID != procsD.Infos[i].CoreID {
			t.Errorf("CoreID: got %v; want %v", procsD.Infos[i].CoreID, procs.Infos[i].CoreID)
		}
		if procs.Infos[i].CPUCores != procsD.Infos[i].CPUCores {
			t.Errorf("CPUCores: got %v; want %v", procsD.Infos[i].CPUCores, procs.Infos[i].CPUCores)
		}
		if procs.Infos[i].ApicID != procsD.Infos[i].ApicID {
			t.Errorf("ApicID: got %v; want %v", procsD.Infos[i].ApicID, procs.Infos[i].ApicID)
		}
		if procs.Infos[i].InitialApicID != procsD.Infos[i].InitialApicID {
			t.Errorf("InitialApicID: got %v; want %v", procsD.Infos[i].InitialApicID, procs.Infos[i].InitialApicID)
		}
		if procs.Infos[i].FPU != procsD.Infos[i].FPU {
			t.Errorf("FPU: got %v; want %v", procsD.Infos[i].FPU, procs.Infos[i].FPU)
		}
		if procs.Infos[i].FPUException != procsD.Infos[i].FPUException {
			t.Errorf("FPUException: got %v; want %v", procsD.Infos[i].FPUException, procs.Infos[i].FPUException)
		}
		if procs.Infos[i].CPUIDLevel != procsD.Infos[i].CPUIDLevel {
			t.Errorf("CPUIDLevel: got %v; want %v", procsD.Infos[i].CPUIDLevel, procs.Infos[i].CPUIDLevel)
		}
		if procs.Infos[i].WP != procsD.Infos[i].WP {
			t.Errorf("WP: got %v; want %v", procsD.Infos[i].WP, procs.Infos[i].WP)
		}
		if procs.Infos[i].Flags != procsD.Infos[i].Flags {
			t.Errorf("Flags: got %v; want %v", procsD.Infos[i].Flags, procs.Infos[i].Flags)
		}
		if procs.Infos[i].BogoMIPS != procsD.Infos[i].BogoMIPS {
			t.Errorf("BogoMIPS: got %v; want %v", procsD.Infos[i].BogoMIPS, procs.Infos[i].BogoMIPS)
		}
		if procs.Infos[i].CLFlushSize != procsD.Infos[i].CLFlushSize {
			t.Errorf("CLFlushSize: got %v; want %v", procsD.Infos[i].CLFlushSize, procs.Infos[i].CLFlushSize)
		}
		if procs.Infos[i].CacheAlignment != procsD.Infos[i].CacheAlignment {
			t.Errorf("CacheAlignment: got %v; want %v", procsD.Infos[i].CacheAlignment, procs.Infos[i].CacheAlignment)
		}
		if procs.Infos[i].AddressSizes != procsD.Infos[i].AddressSizes {
			t.Errorf("AddressSizes: got %v; want %v", procsD.Infos[i].AddressSizes, procs.Infos[i].AddressSizes)
		}
		if procs.Infos[i].PowerManagement != procsD.Infos[i].PowerManagement {
			t.Errorf("PowerManagement: got %v; want %v", procsD.Infos[i].PowerManagement, procs.Infos[i].PowerManagement)
		}
	}
}
