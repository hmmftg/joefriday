package sysinfo

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	u, err := Get()
	if err != nil {
		t.Error(err)
	}
	if u == 0 {
		t.Error("expected uptime to have a non-zero value; was 0")
	}
	t.Logf("%d", u)
}

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
			if v == 0 {
				t.Error("expected uptime to have a non-zero value; was 0")
			}
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
}

func BenchmarkGet(b *testing.B) {
	var u int64
	for i := 0; i < b.N; i++ {
		u, _ = Get()
	}
	_ = u
}
