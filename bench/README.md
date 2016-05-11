bench
=====

Benchmarks for various `joefriday` stuff.  Most of the benchmarks reflect current implementations in `joefriday` packages with a few exceptions: ticker and mem implementations.

Ticker implementations are emulated in this package as benchmarking the ticker implementations directly wouldn't provide accurate information.

Most of the memory related benchmarks are either implemented within this package or use 3rd party packages.  Some of these 3rd party versions may not get the same information as the `joefriday` bench versions, but exist for comparative reasons.  The various `joefriday` implementations exist to provide information on the affect of different implementations on CPU usage and memory allocations.  

## Results
####CPU  
Name|Ops|ns/Op|Bytes/Op|Allocs/Op  
:--|--:|--:|--:|--:  
joefriday/cpu/facts.Get|100000|10481|800|18  
joefriday/cpu/stats.Get|200000|8346|264|4  
DataDog/gohai/cpu.Cpu.Collect|5000|306624|999877|822  
shirou/gopsutil/cpu.Info|50000|38845|9656|120  
shirou/gopsutil/cpu.Times|50000|31120|10512|28  

####Memory  
Name|Ops|ns/Op|Bytes/Op|Allocs/Op  
:--|--:|--:|--:|--:  
joefriday/mem.Get|200000|9954|352|1  
joefriday/sysinfo/mem.Info.Get|10000000|189|0|0  
cloudfoundry/gosigar.Mem.Get|50000|33532|13912|108  
DataDog/gohai/memory.Memory.Collect|2000|1037399|3543120|2683  
guillermo/go.procmeminfo.MemInfo.Update|30000|42205|8216|133  
shirou/gopsutil/mem.VirtualMemory|30000|52938|12296|264  

####Network  
Name|Ops|ns/Op|Bytes/Op|Allocs/Op  
:--|--:|--:|--:|--:  
joefriday/net/info.Get|300000|5002|328|4  
joefriday/net/usage.Get|300000|5544|664|6  
DataDog/gohai/network|10000|109043|55109|222  
shirou/gopsutil/net|30000|55134|25568|134  
shirou/gopsutil/net/IOCounters|100000|21182|6200|26  

####Platform  
Name|Ops|ns/Op|Bytes/Op|Allocs/Op  
:--|--:|--:|--:|--:  
joefriday/platform/kernel.Get|1000000|1789|288|7  
joefriday/platform/release.Get|1000000|1771|229|7  
DataDog/gohai/platform.Platform.Collect|1000|1491460|177305|400  
joefriday/platform/loadavg.Get|1000000|1754|12|3  
joefriday/sysinfo/load.LoadAvg.Get|10000000|192|0|0  
cloudfoundry/gosigar.LoadAverage.Get|200000|8245|2440|9  
shirou/gopsutil/load.Avg|200000|9014|2488|11  
shirou/gopsutil/load.Misc|50000|34513|8744|20  
joefriday/platform/uptime.Get|1000000|1438|32|2  
joefriday/sysinfo/uptime.Uptime.Get|10000000|189|0|0  
cloudfoundry/gosigar.Uptime.Get|10000000|136|0|0  
