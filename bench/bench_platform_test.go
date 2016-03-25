package bench

import (
	"testing"

	"github.com/mohae/joefriday/platform"
)

func BenchmarkGetPlatformKernel(b *testing.B) {
	var val *platform.Kernel
	for i := 0; i < b.N; i++ {
		val, _ = platform.GetKernel()
	}
	_ = val
}

func BenchmarkPlatformKernelSerializeFlat(b *testing.B) {
	var val *platform.Kernel
	var p []byte
	b.StopTimer()
	val, _ = platform.GetKernel()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = val.SerializeFlat()
	}
	_ = p
}

func BenchmarkPlatformKernelDeserializeFlat(b *testing.B) {
	var val *platform.Kernel
	var p []byte
	b.StopTimer()
	val, _ = platform.GetKernel()
	p = val.SerializeFlat()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val = platform.DeserializeKernelFlat(p)
	}
	_ = val
}

func BenchmarkGetPlatformRelease(b *testing.B) {
	var val *platform.Release
	for i := 0; i < b.N; i++ {
		val, _ = platform.GetRelease()
	}
	_ = val
}

func BenchmarkPlatformReleaseSerializeFlat(b *testing.B) {
	var val *platform.Release
	var p []byte
	b.StopTimer()
	val, _ = platform.GetRelease()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = val.SerializeFlat()
	}
	_ = p
}

func BenchmarkPlatformReleaseDeserializeFlat(b *testing.B) {
	var val *platform.Release
	var p []byte
	b.StopTimer()
	val, _ = platform.GetRelease()
	p = val.SerializeFlat()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val = platform.DeserializeReleaseFlat(p)
	}
	_ = val
}
