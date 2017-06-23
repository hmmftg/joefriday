package uptime

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	u, err := Get()
	if err != nil {
		t.Error(err)
	}
	if u.Timestamp == 0 {
		t.Error("expected timestamp to have a non-zero value; was 0")
	}
	if u.Uptime == 0 {
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
		case u, ok := <-tk.Data:
			if !ok {
				break
			}
			if u.Timestamp == 0 {
				t.Error("expected timestamp to have a non-zero value; was 0")
			}
			if u.Uptime == 0 {
				t.Error("expected uptime to have a non-zero value; was 0")
			}
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func BenchmarkGet(b *testing.B) {
	var u Info
	for i := 0; i < b.N; i++ {
		u, _ = Get()
	}
	_ = u
}
