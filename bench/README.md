bench
=====

Benchmarks for various `joefriday` stuff.  Most of the benchmarks reflect current implementations in `joefriday` packages with a few exceptions: ticker and mem implementations.

Ticker implementations are emulated in this package as benchmarking the ticker implementations directly wouldn't provide accurate information.

Most of the memory related benchmarks are either implemented within this package or use 3rd party packages.  Some of these 3rd party versions may not get the same information as the `joefriday` bench versions, but exist for comparative reasons.  The various `joefriday` implementations exist to provide information on the affect of different implementations on CPU usage and memory allocations.  

## Results
####CPU  
Name|Ops|ns/Op|Bytes/Op|Allocs/Op  
:--|--:|--:|--:|--:  
joefriday/cpu/stats.Get|200000|8568|264|4  
DataDog/gohai/cpu.Cpu.Collect|5000|365712|999877|822  
shirou/gopsutil/cpu.Info|30000|41024|9656|120  
shirou/gopsutil/cpu.Times|50000|32612|10512|28  

####Memory  
Name|Ops|ns/Op|Bytes/Op|Allocs/Op  
:--|--:|--:|--:|--:  
joefriday/sysinfo/mem.Info.Get|10000000|196|0|0  
cloudfoundry/gosigar.Mem.Get|50000|35624|13912|108  
DataDog/gohai/memory.Memory.Collect|1000|1281577|3543120|2683  
guillermo/go.procmeminfo.MemInfo.Update|30000|44254|8216|133  
shirou/gopsutil/mem.VirtualMemory|30000|55456|12296|264  

####Network  
Name|Ops|ns/Op|Bytes/Op|Allocs/Op  
:--|--:|--:|--:|--:  
joefriday/net/usage.Get|300000|5728|664|6  
DataDog/gohai/network|10000|116706|55109|222  
shirou/gopsutil/net|30000|58493|25568|134  
shirou/gopsutil/net/IOCounters|100000|22491|6200|26  

####Platform  
Name|Ops|ns/Op|Bytes/Op|Allocs/Op  
:--|--:|--:|--:|--:  
joefriday/platform/release.Get|1000000|1765|229|7  
DataDpg/gohai/platform.Platform.Collect|1000|1687743|177303|400  
joefriday/platform/loadavg.Get|1000000|1825|12|3  
joefriday/sysinfo/load.LoadAvg.Get|10000000|202|0|0  
cloudfoundry/gosigar.LoadAverage.Get|200000|8674|2440|9  
shirou/gopsutil/load.Avg|200000|9475|2488|11  
shirou/gopsutil/load.Misc|50000|36010|8744|20  
joefriday/platform/uptime.Get|1000000|1498|32|2  
joefriday/sysinfo/uptime.Uptime.Get|10000000|193|0|0  
cloudfoundry/gosigar.Uptime.Get|10000000|141|0|0  
