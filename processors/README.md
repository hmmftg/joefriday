processors
=======

Package `processors` provides information about a system's nodes, sockets, and physical processors. This information is gathered from both `/proc/cpuinfo` and the `sysfs`. For information gathered from `cpuinfo`, the first core entry for each physical CPU is used.

For `x86/x86-64` systems, the `CPUMHz` field is not reliable as it is the current speed of the first core processed for each physical processor. Modern `x86\x86-64` processor core speeds are dynamic and fall within a range; there may be other cores on the same processor that are at higher or lower speeds than the reported value. For `x86\x86-64` processors, the `MHzMin` and `MHzMax` fields provide information on the processors min and max speeds.

`CPUInfo` from the following processors were used for testing:
* Intel I7 5600u
* Intel Xeon E52690 w 2 sockets
* AMD R7 1800x

For testing, the `sysfs` stuff use generated files and data. Please file an issue or a pull request for any gaps encountered.

Please file an issue or a pull request for any additional processors/architectures.
    
