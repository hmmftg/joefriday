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

	"github.com/mohae/joefriday/net"
)

func BenchmarkGetNetDevInfo(b *testing.B) {
	var inf *net.Info
	for i := 0; i < b.N; i++ {
		inf, _ = net.GetInfo()
	}
	_ = inf
}

func BenchmarkGetNetDevData(b *testing.B) {
	var inf []byte
	for i := 0; i < b.N; i++ {
		inf, _ = net.GetData()
	}
	_ = inf
}

func BenchmarkEmulateNetDevDataTicker(b *testing.B) {
	var inf []byte
	for i := 0; i < b.N; i++ {
		inf, _ = EmulateNetDevDataTicker()
	}
	_ = inf

}
