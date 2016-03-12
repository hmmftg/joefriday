joefriday
=========

"All we want are the facts, ma'am"

## Notes:

### Data vs Info 
Anything that ends with Info is either a Go struct or related to a Go struct; e.g. returns a Go struct with the relevant data.  Anything that ends with Data returns is either a Flatbuffers data structure or related to a Flatbuffers data structure; e.g. returns Flatbuffer serialized bytes as `[]byte` for the relevant data.

### /proc/meminfo
`MemAvailable` is not available on pre Linux 3.14 kernels.



