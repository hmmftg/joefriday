package mem

import (
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	tkr, err := NewTicker(time.Millisecond)
	if err != nil {
		t.Error(err)
		return
	}
	tk := tkr.(*Ticker)
	for i := 0; i < 5; i++ {
		select {
		case <-tk.Done:
			break
		case v, ok := <-tk.Data:
			if !ok {
				break
			}
			checkMemInfo("ticker", v, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func TestGetMemInfo(t *testing.T) {
	m, err := Get()
	if err != nil {
		t.Error(err)
	}
	checkMemInfo("get", m, t)
}

func checkMemInfo(n string, m Info, t *testing.T) {
	if m.Timestamp == 0 {
		t.Errorf("%s: expected the Timestamp to be non-zero, was 0", n)
	}
	if m.TotalRAM == 0 {
		t.Errorf("%s: expected the TotalRAM to be non-zero, was 0", n)
	}
	if m.FreeRAM == 0 {
		t.Errorf("%s: expected the FreeRAM to be non-zero, was 0", n)
	}
	t.Logf("%#v\n", m)
}

func BenchmarkMemInfo(b *testing.B) {
	var tmp Info
	for i := 0; i < b.N; i++ {
		_ = tmp.Get()
	}
	_ = tmp
}
