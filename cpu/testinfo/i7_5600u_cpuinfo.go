package testinfo

import (
	"errors"
	"fmt"

	"github.com/mohae/joefriday/cpu/cpufreq"
	"github.com/mohae/joefriday/cpu/cpuinfo"
)

var I75600uCPUInfo = []byte(`processor	: 0
vendor_id	: GenuineIntel
cpu family	: 6
model		: 61
model name	: Intel(R) Core(TM) i7-5600U CPU @ 2.60GHz
stepping	: 4
microcode	: 0x24
cpu MHz		: 2602.062
cache size	: 4096 KB
physical id	: 0
siblings	: 4
core id		: 0
cpu cores	: 2
apicid		: 0
initial apicid	: 0
fpu		: yes
fpu_exception	: yes
cpuid level	: 20
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch epb intel_pt tpr_shadow vnmi flexpriority ept vpid fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm rdseed adx smap xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5187.81
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:

processor	: 1
vendor_id	: GenuineIntel
cpu family	: 6
model		: 61
model name	: Intel(R) Core(TM) i7-5600U CPU @ 2.60GHz
stepping	: 4
microcode	: 0x24
cpu MHz		: 2600.000
cache size	: 4096 KB
physical id	: 0
siblings	: 4
core id		: 0
cpu cores	: 2
apicid		: 1
initial apicid	: 1
fpu		: yes
fpu_exception	: yes
cpuid level	: 20
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch epb intel_pt tpr_shadow vnmi flexpriority ept vpid fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm rdseed adx smap xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5194.01
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:

processor	: 2
vendor_id	: GenuineIntel
cpu family	: 6
model		: 61
model name	: Intel(R) Core(TM) i7-5600U CPU @ 2.60GHz
stepping	: 4
microcode	: 0x24
cpu MHz		: 2939.282
cache size	: 4096 KB
physical id	: 0
siblings	: 4
core id		: 1
cpu cores	: 2
apicid		: 2
initial apicid	: 2
fpu		: yes
fpu_exception	: yes
cpuid level	: 20
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch epb intel_pt tpr_shadow vnmi flexpriority ept vpid fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm rdseed adx smap xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5192.08
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:

processor	: 3
vendor_id	: GenuineIntel
cpu family	: 6
model		: 61
model name	: Intel(R) Core(TM) i7-5600U CPU @ 2.60GHz
stepping	: 4
microcode	: 0x24
cpu MHz		: 2599.682
cache size	: 4096 KB
physical id	: 0
siblings	: 4
core id		: 1
cpu cores	: 2
apicid		: 3
initial apicid	: 3
fpu		: yes
fpu_exception	: yes
cpuid level	: 20
wp		: yes
flags		: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch epb intel_pt tpr_shadow vnmi flexpriority ept vpid fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm rdseed adx smap xsaveopt dtherm ida arat pln pts
bugs		:
bogomips	: 5192.28
clflush size	: 64
cache_alignment	: 64
address sizes	: 39 bits physical, 48 bits virtual
power management:`)

// ValidateI75600uCPUInfo verifies that the info in the struct info is
// consistent with the above data. If everything verifies a nil is returned,
// otherwise an error is returned. This is used for testing.
//
// processor, core id, apicid, initial apicid fields are not checked as they
// vary per entry.
func ValidateI75600uCPUInfo(inf *cpuinfo.CPUInfo) error {
	if inf.Timestamp == 0 {
		return errors.New("expected Timestamp to have a non-zero value; it didn't")
	}
	if len(inf.CPU) != 4 {
		return fmt.Errorf("CPU: got %d; want 4", len(inf.CPU))
	}
	modelName := "Intel(R) Core(TM) i7-5600U CPU @ 2.60GHz"
	for i, cpu := range inf.CPU {
		if cpu.VendorID != "GenuineIntel" {
			return fmt.Errorf("%d: vendor_id: got %q; want \"GenuineIntel\"", i, cpu.VendorID)
		}
		if cpu.CPUFamily != "6" {
			return fmt.Errorf("%d: cpu family: got %q; want \"6\"", i, cpu.CPUFamily)
		}
		if cpu.Model != "61" {
			return fmt.Errorf("%d: model: got %q; want \"61\"", i, cpu.Model)
		}
		if cpu.ModelName != modelName {
			return fmt.Errorf("%d: model name: got %q; want %q", i, cpu.ModelName, modelName)
		}
		if cpu.Stepping != "4" {
			return fmt.Errorf("%d: stepping: got %q; want \"4\"", i, cpu.Stepping)
		}
		if cpu.Microcode != "0x24" {
			return fmt.Errorf("%d: microcode: got %q; want \"0x24\"", i, cpu.Microcode)
		}
		if int(cpu.CPUMHz) < 2500 {
			return fmt.Errorf("%d: cpu MHz: got %.3f; want a value >= 2500", i, cpu.CPUMHz)
		}
		if cpu.CacheSize != "4096 KB" {
			return fmt.Errorf("%d: cache size: got %q; want \"4096 KB\"", i, cpu.CacheSize)
		}
		if cpu.PhysicalID != 0 {
			return fmt.Errorf("%d: physical id: got %d; want 0", i, cpu.PhysicalID)
		}
		if cpu.Siblings != 4 {
			return fmt.Errorf("%d: siblings: got %d; want 4", i, cpu.Siblings)
		}
		if cpu.CPUCores != 2 {
			return fmt.Errorf("%d: cpu cores: got %d; want 2", i, cpu.CPUCores)
		}
		y := "yes"
		if cpu.FPU != y {
			return fmt.Errorf("%d: fpu: got %q; want %q", i, cpu.FPU, y)
		}
		if cpu.FPUException != y {
			return fmt.Errorf("%d: fpu exception: got %q; want %q", i, cpu.FPUException, y)
		}
		if cpu.CPUIDLevel != "20" {
			return fmt.Errorf("%d: cpuid level: got %q; want \"20\"", i, cpu.CPUIDLevel)
		}
		if cpu.WP != y {
			return fmt.Errorf("%d: wp: got %q; want %q", i, cpu.WP, y)
		}
		if len(cpu.Flags) != 99 {
			return fmt.Errorf("%d: flags: got %d; want 99", i, len(cpu.Flags))
		}
		if int(cpu.BogoMIPS) < 5100 {
			return fmt.Errorf("%d: bogomips: got %.3f; want a value >= 5100", i, cpu.BogoMIPS)
		}
		if cpu.CLFlushSize != "64" {
			return fmt.Errorf("%d: clflush size: got %q; want 64", i, cpu.CLFlushSize)
		} 
		if cpu.CacheAlignment != "64" {
			return fmt.Errorf("%d: cache alignment size: got %q; want 64", i, cpu.CacheAlignment)
		} 
		if len(cpu.PowerManagement) != 0 {
			return fmt.Errorf("%d: power management: got %d; wanted 0", i, len(cpu.PowerManagement))
		}
		if len(cpu.AddressSizes) != 2 {
			return fmt.Errorf("%d: address sizes: got %d; want 2", i, len(cpu.AddressSizes))
		}
	}
	return nil
}

// ValidateI75600uCPUFreq verifies that the info in the struct info is
// consistent with relevant parts of the above data. If everything verifies a
// nil is returned, otherwise an error is returned. This is used for testing.
func ValidateI75600uCPUFreq(f *cpufreq.Frequency) error {
	if f.Timestamp == 0 {
		return errors.New("expected Timestamp to have a non-zero value; it didn't")
	}
	if len(f.CPU) != 4 {
		return fmt.Errorf("CPU: got %d; want 4", len(f.CPU))
	}
	for i, cpu := range f.CPU {	
		if int(cpu.CPUMHz) < 2500 {
			return fmt.Errorf("%d: cpu MHz: got %.3f; want a value >= 2500", i, cpu.CPUMHz)
		}
		if int(cpu.PhysicalID) != 0 {
			return fmt.Errorf("%d: physical id: got %d; want 0", i, cpu.PhysicalID)
		}
		if int(cpu.Processor) != i {
			return fmt.Errorf("%d: processor: got %d; want %d", i, cpu.Processor, i)
		}
		// APICID happens to be consistent with i
		if int(cpu.APICID) != i {
			return fmt.Errorf("%d: apicid: got %d; want %d", i, cpu.APICID, i)
		}
		if int(cpu.CoreID) != i/2 {
			return fmt.Errorf("%d: core id: got %d; want %d", i, cpu.CoreID, i/2)
		}
	}
	return nil
}