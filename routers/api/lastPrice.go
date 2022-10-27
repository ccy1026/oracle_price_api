package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"oracle_price_api/database"
	"oracle_price_api/models"
	"oracle_price_api/models/api"
	"oracle_price_api/utils"
	"oracle_price_api/utils/external"
	"strconv"
	"time"
)

// @BasePath /

// @Summary Get Latest Price
// @Tags Price
// @param  token path string true "token"
// @Success 200
// @Router /lastPrice/{token} [get]
func GetLatestPrice(ctx *gin.Context) {
	if ctx.Param("token") == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New(utils.NullTokenError))
		return
	}

	isValid, err := database.GetDBInstance().CheckSupportToken(ctx.Param("token"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if !isValid {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New(utils.NotSupportToken))
		return
	}

	response := new(api.PriceResponse)
	lastPriceData, err := database.GetDBInstance().GetLastPriceDataWithTokenName(ctx.Param("token"), 0)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	response.Token = ctx.Param("token")
	response.Price = lastPriceData.Price
	response.LastUpdateTime = lastPriceData.Timestamp
	ctx.JSON(200, response)
}

// @BasePath /

// @Summary Get Latest Price by timestamp
// @Tags Price
// @param  token path string true "token"
// @param  timestamp path int true "timestamp"
// @Success 200
// @Router /lastPrice/{token}/{timestamp} [get]
func GetLatestPriceByTimestamp(ctx *gin.Context) {
	if ctx.Param("token") == "" || ctx.Param("timestamp") == "" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New(utils.NullTokenTimestamp))
		return
	}

	db := database.GetDBInstance()

	pairFeed, err := db.QuerySymbol(ctx.Param("token"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	response := new(api.PriceResponse)
	stringTimestamp, err := strconv.ParseInt(ctx.Param("timestamp"), 10, 64)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	formattedTime := time.Unix(stringTimestamp, 0).Truncate(time.Minute).Unix()

	if time.Unix(stringTimestamp, 0).Truncate(time.Minute).After(time.Now().Truncate(time.Minute)) {
		utils.ErrorResponse(ctx, http.StatusBadRequest, errors.New(utils.WrongTimestamp))
	}
	lastPriceData, err := database.GetDBInstance().GetLastPriceDataWithTokenName(ctx.Param("token"), formattedTime)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	response.Token = ctx.Param("token")
	response.LastUpdateTime = time.Unix(formattedTime, 0)
	if lastPriceData.Id != 0 {
		response.Price = lastPriceData.Price
	} else {
		upperTime := time.Unix(stringTimestamp, 0).Truncate(time.Minute).Add(1 * time.Minute).Unix()
		priceRangeData, err := external.GetPriceRangeWithAPI(ctx.Param("token"), formattedTime, upperTime)
		response.Price = priceRangeData.Data[0].Close
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		var lastPriceDataArray []models.PriceData
		lastPriceDataArray = append(lastPriceDataArray, *lastPriceData)
		err = database.GetDBInstance().InsertAPIDataToDB(pairFeed.Id, priceRangeData, lastPriceDataArray)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(200, response)
}
