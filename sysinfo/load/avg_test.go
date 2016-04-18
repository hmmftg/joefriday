package avg

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	l, err := Get()
	if err != nil {
		t.Error(err)
	}
	checkLoadAvg("get", l, t)
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
			checkLoadAvg("ticker", v, t)
		case err := <-tk.Errs:
			t.Errorf("unexpected error: %s", err)
		}
	}
	tk.Stop()
	tk.Close()
}

func checkLoadAvg(n string, l LoadAvg, t *testing.T) {
	if l.One == 0 {
		t.Errorf("%s: expected the 1 minute load avg to be non-zero, was 0", n)
	}
	if l.Five == 0 {
		t.Errorf("%s: expected the 5 minute load avg to be non-zero, was 0", n)
	}
	if l.Fifteen == 0 {
		t.Errorf("%s: expected the 15 minute load avg to be non-zero, was 0", n)
	}
	t.Logf("%#v\n", l)
}

func BenchmarkLoadAvg(b *testing.B) {
	var tmp LoadAvg
	for i := 0; i < b.N; i++ {
		_ = tmp.Get()
	}
	_ = tmp
}
