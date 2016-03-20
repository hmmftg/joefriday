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
