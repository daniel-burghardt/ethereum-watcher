package util

import "testing"

func TestHexToInt(t *testing.T) {
	value, err := HexToInt("0x6f05b59d3b20000")
	if err != nil {
		t.Errorf("Error = %v; want nil", err)
	}
	if value != 500000000000000000 {
		t.Errorf("Int value = %v; want 500000000000000000", value)
	}
}

func TestWeiToEth(t *testing.T) {
	value := WeiToEth(500000000000000000)
	if value != 0.5 {
		t.Errorf("Eth value = %f; want 0.5", value)
	}
}
