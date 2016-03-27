bench
=====

Benchmarks for various `joefriday` stuff.  Most of the benchmarks reflect current implementations in `joefriday` packages with a few exceptions: ticker and mem implementations.

Ticker implementations are emulated in this package as benchmarking the ticker implementations directly wouldn't provide accurate information.

Most of the memory related benchmarks are either implemented within this package or use 3rd party packages.  Some of these 3rd party versions may not get the same information as the `joefriday` bench versions, but exist for comparative reasons.  The various `joefriday` implementations exist to provide information on the affect of different implementations on CPU usage and memory allocations.  
