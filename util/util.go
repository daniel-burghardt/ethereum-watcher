package util

import (
	"fmt"
	"math"
	"strconv"
)

func HexToInt(amountHex string) (int64, error) {
	amount, err := strconv.ParseInt(amountHex, 0, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing hex to int: %w", err)
	}

	return amount, nil
}

func WeiToEth(wei int64) float64 {
	return float64(wei) * math.Pow(10, -18)
}
