package bench

import (
	"testing"

	"github.com/mohae/joefriday/platform"
)

func BenchmarkGetPlatformKernel(b *testing.B) {
	var val platform.Kernel
	for i := 0; i < b.N; i++ {
		val, _ = platform.GetKernel()
	}
	_ = val
}

func BenchmarkPlatformKernelSerialize(b *testing.B) {
	var val platform.Kernel
	var p []byte
	b.StopTimer()
	val, _ = platform.GetKernel()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = val.Serialize()
	}
	_ = p
}

func BenchmarkPlatformKernelDeSerialize(b *testing.B) {
	var val platform.Kernel
	var p []byte
	b.StopTimer()
	val, _ = platform.GetKernel()
	p = val.Serialize()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val = platform.DeserializeKernel(p)
	}
	_ = val
}

func BenchmarkGetPlatformRelease(b *testing.B) {
	var val platform.Release
	for i := 0; i < b.N; i++ {
		val, _ = platform.GetRelease()
	}
	_ = val
}

func BenchmarkPlatformReleaseSerialize(b *testing.B) {
	var val platform.Release
	var p []byte
	b.StopTimer()
	val, _ = platform.GetRelease()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		p = val.Serialize()
	}
	_ = p
}

func BenchmarkPlatformReleaseDeSerialize(b *testing.B) {
	var val platform.Release
	var p []byte
	b.StopTimer()
	val, _ = platform.GetRelease()
	p = val.Serialize()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val = platform.DeserializeRelease(p)
	}
	_ = val
}
