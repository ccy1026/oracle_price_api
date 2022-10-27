package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"oracle_price_api/routers/api"
)

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true

	r.Use(cors.New(corsConfig))
	r.Use(CORSMiddleware())

	r.GET("/lastPrice/:token", api.GetLatestPrice)
	r.GET("/lastPrice/:token/:timestamp", api.GetLatestPriceByTimestamp)
	r.POST("/rangePrice", api.GetAveragePriceFromRange)

	return r
}
