### WARNING
This will be going through some changes sometime in the near future.  Utilization and Usage related packages should be considered deprecated.  Usage data will be replaced by code that returns the delta between two snap-shots.  How this will be organized has not yet been determined; which is why the changes haven't been made yet.

Utilization related functionality will be removed and moved to a separate package.  The yet to be determined package will leverage JoeFriday and provide Utilization and other calculated metrics using more statistically sound methods.  These changes will occur after the delta changes.

joefriday
=========
[![GoDoc](https://godoc.org/github.com/mohae/joefriday?status.svg)](https://godoc.org/github.com/mohae/joefriday)[![Build Status](https://travis-ci.org/mohae/joefriday.png)](https://travis-ci.org/mohae/joefriday)
> "All we want are the facts, ma'am"  
>   - Joe Friday

JoeFriday is a library that gathers system information: cpu, disk, memory, network, and platform information.  This information can be returned as Go structs, Flatbuffers serialized bytes, or JSON serialized bytes.  For Flatbuffers and JSON, deserialization convenience methods are provided.  When it makes sense, a Ticker based implementation is provided to enable periodic gathering of information.

JoeFriday seeks to minimize allocations and time spent gathering the information.  For minimal resource usage, use the sysinfo implementations, if appropriate, as the data provided by those implementations use syscalls, which are at least an order of magnitude faster than processing the proc files.

This library only supports Linux.

## Packages

### Comparative Benchmarks
The `bench` package provides comparative benchmarks; showing how JoeFriday implementations compare with some other libraries.  

The information provided by the various libraries are not the same; that should also be factored into the decision as to which library best suits your requirements.

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
Platform provides information about the system's kernel, loadavg, release, and uptime.  For loadavg and uptime, the `sysinfo` packages should be used, unless there is a need for information that this package makes available.

### Processors
Processors provides information about the physical Chips on a system.  This differs from CPU in that it provides less detail about the Chips and it only provides information at the Chip level.  There will be one Chip entry per physical CPU chip on the system; CPU provides detailed information about all CPU cores on a system, physical and logical.

### Sysinfo
Sysinfo provides information about the system via syscalls: load information (loadavg), memory information, and uptime.  Less memory and a lot less CPU cycles are used to obtain the information in this manner.

## Notes:
A big thanks to [Eric Lagergren](https://github.com/EricLagergren) for all of his help and his contributions.

## TODO:
Rename files to reflect GOOS and Arch they support.
