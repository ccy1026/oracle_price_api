package routers

import (
	"encoding/json"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io/ioutil"
	"log"
	"net/http"
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

	//Documentation
	r.GET("/docs", ReadDocsJson)
	url := ginSwagger.URL("/docs")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return r
}

func ReadDocsJson(ctx *gin.Context) {
	body, err := ioutil.ReadFile("./docs/swagger.json")
	if err != nil {
		log.Println(err)
	}
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
