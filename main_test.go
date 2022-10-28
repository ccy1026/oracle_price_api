package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"oracle_price_api/routers/api"
	"testing"
	"time"
)

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	return r
}

type TestPriceResponse struct {
	Token          string    `json:"token"`
	Price          float64   `json:"price"`
	LastUpdateTime time.Time `json:"last_update_time"`
}

type TestAverageRangePriceRequest struct {
	Token         string `json:"token" validate:"required"`
	FromTimeStamp int64  `json:"from_time_stamp" validate:"required"`
	ToTimeStamp   int64  `json:"to_time_stamp" validate:"required"`
}

type TestAverageRangePriceResponse struct {
	TestAverageRangePriceRequest
	AveragePrice float64 `json:"average_price"`
}

func CurrentTime() time.Time {
	return time.Now()
}

func TestRangePrice(t *testing.T) {
	r := SetUpRouter()
	r.POST("/rangePrice", api.GetAveragePriceFromRange)
	reqBody := TestAverageRangePriceRequest{
		Token:         "ETH",
		FromTimeStamp: time.Now().Unix() - 1800,
		ToTimeStamp:   time.Now().Unix(),
	}
	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/rangePrice", bytes.NewBuffer(jsonBody))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)

	resStruct := new(TestAverageRangePriceResponse)
	json.Unmarshal(responseData, resStruct)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, 0, math.Round(resStruct.AveragePrice))
}

func TestTimestampPrice(t *testing.T) {
	r := SetUpRouter()
	r.GET("/lastPrice/:token/:timestamp", api.GetLatestPriceByTimestamp)
	req, _ := http.NewRequest("GET", "/lastPrice/ETH/"+fmt.Sprint(time.Now().Unix()), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)

	resStruct := new(TestPriceResponse)
	json.Unmarshal(responseData, resStruct)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ETH", resStruct.Token)
	assert.NotEqual(t, 0, resStruct.Price)
}

func TestLastPrice(t *testing.T) {
	r := SetUpRouter()
	r.GET("/lastPrice/:token", api.GetLatestPrice)
	req, _ := http.NewRequest("GET", "/lastPrice/ETH", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	responseData, _ := ioutil.ReadAll(w.Body)

	resStruct := new(TestPriceResponse)
	json.Unmarshal(responseData, resStruct)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ETH", resStruct.Token)
	assert.NotEqual(t, 0, resStruct.Price)
}
