package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorResponse(ctx *gin.Context, status int, err error) {
	ctx.AbortWithStatusJSON(
		status,
		gin.H{"error": err.Error()})
	return
}

func VerifyJsonBodyByStruct(structBody interface{}) error {
	if err := validator.New().Struct(structBody); err != nil {
		return err
	}
	return nil
}
