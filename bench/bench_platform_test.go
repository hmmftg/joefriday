// Copyright 2016 Joel Scoble and The JoeFriday authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
