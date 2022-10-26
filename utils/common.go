package utils

import (
	"math"
	"math/big"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func PriceFormat(price *big.Float, decimal *big.Float, dp int) float64 {
	formattedPrice, _ := new(big.Float).Quo(price, decimal).Float64()
	output := math.Pow(10, float64(dp))
	return float64(round(formattedPrice*output)) / output
}
