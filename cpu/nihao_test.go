package cpu

import (
	"fmt"
	"testing"
)

func TestNihao(t *testing.T) {
	inf, err := NiHao()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	fmt.Println(*inf)
}
