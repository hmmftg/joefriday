joefriday
=========

> "All we want are the facts, ma'am"  
>   - Joe Friday

JoeFriday is a library that gathers system information: cpu, disk, memory, network, and platform information.  This information can be returned as Go structs, Flatbuffers serialized bytes, or JSON serialized bytes.  For Flatbuffers and JSON, deserialization convenience methods are provided.  When it makes sense, a Ticker based implementation is provided to enable periodic gathering of information.

JoeFriday seeks to minimize allocations and time spent gathering the information.  For minimal resource usage, use the sysinfo implementations, if appropriate, as the data provided by those implementations use syscalls, which are at least an order of magnitude faster than processing the proc files.

This library only supports Linux.

## Packages

### Comparative Benchmarks
The `bench` package provides comparative benchmarks; showing how JoeFriday implementations compare with some other libraries.  

Sometimes the data provided by other libraries is not the same as what is provided by JoeFriday: these benchmarks are not necessarily apples to apples comparisons, data wise.

### CPU
CPU provides facts, stats, and utilization information.  Facts provide information about the system's CPUs: e.g. model, make, speed, flags, etc.  Stats provides a snapshot of the current CPU state.  Utilization provides information about CPU utilization, which is the difference between two stats snapshots.

### Disk
Disk provides stats and usage information about the system's block devices.  The device stats provide read and write information, io information, and time spent on IO.  The device usage information is the difference between two stats snapshots.

Currently, no disk usage information is provided.  This is a TODO.

### Mem
Mem provides all available memory information: `sysinfo/mem` should be used instead, unless there is a need for information that this package makes available.

`MemAvailable` is not available on pre Linux 3.14 kernels.

### Net
Net provides information about the system's network information.  Info is the current information about the interfaces, this data is cumulative.  Usage is the difference between two info snapshots and provides information about the usage, for a given time delta, of the system's network interfaces.

### Platform
Platform provides information about the system's kernel, loadavg, release, and uptime.  For loadavg and uptime, the `sysinfo` packages should be used, unless ther is a need for information that this package makes available.

### Sysinfo
Sysinfo provides information about the system via syscalls: load information (loadavg), memory information, and uptime.  Less memory and a lot less CPU cycles are used to obtain the information in this manner.

## Notes:
A big thanks to [Eric Lagergren](https://github.com/EricLagergren) for all of his help and his contributions.
