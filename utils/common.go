package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

func VerifyJsonBodyByStruct(structBody interface{}) error {
	if err := validator.New().Struct(structBody); err != nil {
		return err
	}
	return nil
}

func ErrorResponse(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(
		status,
		gin.H{"error": err.Error()})
	return
}
