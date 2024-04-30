package util

import (
	"fmt"
	"math"
	"strconv"
)

func HexToInt(number string) (int64, error) {
	amount, err := strconv.ParseInt(number, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing hex to int: %w", err)
	}

	return amount, nil
}

func IntToHex(number int64) string {
	return fmt.Sprintf("0x%x", number)
}

func WeiToEth(wei int64) float64 {
	return float64(wei) * math.Pow(10, -18)
}
