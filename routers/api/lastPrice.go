package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"oracle_price_api/database"
	"oracle_price_api/models"
	"oracle_price_api/models/api"
	"oracle_price_api/utils"
	"oracle_price_api/utils/external"
	"strconv"
	"time"
)

func GetLatestPrice(ctx *gin.Context) {
	if ctx.Param("token") == "" {
		ErrorResponse(ctx, http.StatusBadRequest, errors.New(utils.NullTokenError))
		return
	}

	isValid, err := database.GetDBInstance().CheckSupportToken(ctx.Param("token"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	if !isValid {
		ErrorResponse(ctx, http.StatusBadRequest, errors.New(utils.NotSupportToken))
		return
	}

	response := new(api.PriceResponse)
	lastPriceData, err := database.GetDBInstance().GetLastPriceDataWithTokenName(ctx.Param("token"), 0)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	response.Token = ctx.Param("token")
	response.Price = lastPriceData.Price
	response.LastUpdateTime = lastPriceData.Timestamp
	ctx.JSON(200, response)
}

func GetLatestPriceByTimestamp(ctx *gin.Context) {
	if ctx.Param("token") == "" || ctx.Param("timestamp") == "" {
		ErrorResponse(ctx, http.StatusBadRequest, errors.New(utils.NullTokenTimestamp))
		return
	}

	db := database.GetDBInstance()

	pairFeed, err := db.QuerySymbol(ctx.Param("token"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	response := new(api.PriceResponse)
	stringTimestamp, err := strconv.ParseInt(ctx.Param("timestamp"), 10, 64)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	formattedTime := time.Unix(stringTimestamp, 0).Truncate(time.Minute).Unix()

	if time.Unix(stringTimestamp, 0).Truncate(time.Minute).After(time.Now().Truncate(time.Minute)) {
		ErrorResponse(ctx, http.StatusBadRequest, errors.New(utils.WrongTimestamp))
	}
	lastPriceData, err := database.GetDBInstance().GetLastPriceDataWithTokenName(ctx.Param("token"), formattedTime)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
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
			ErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
		var lastPriceDataArray []models.PriceData
		lastPriceDataArray = append(lastPriceDataArray, *lastPriceData)
		err = database.GetDBInstance().InsertAPIDataToDB(pairFeed.Id, priceRangeData, lastPriceDataArray)
		if err != nil {
			ErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(200, response)
}

func GetAveragePriceFromRange(ctx *gin.Context) {
	var lastPriceRequest = api.AverageRangePriceRequest{}
	err := ctx.ShouldBindBodyWith(&lastPriceRequest, binding.JSON)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = VerifyJsonBodyByStruct(lastPriceRequest)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	db := database.GetDBInstance()
	pairFeed, err := db.QuerySymbol(ctx.Param("token"))
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	lowerTime := time.Unix(lastPriceRequest.FromTimeStamp, 0).Truncate(time.Minute)
	upperTime := time.Unix(lastPriceRequest.ToTimeStamp, 0).Truncate(time.Minute).Add(1 * time.Minute)
	pair := (upperTime.Unix() - lowerTime.Unix()) / 60
	dbPriceRange, err := database.GetDBInstance().GetLatestPriceInTimeRange(lastPriceRequest.Token, lowerTime.Unix(), upperTime.Unix())
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	priceRangeData := new(models.PriceDataModel)

	if len(dbPriceRange) != int(pair) {
		priceRangeData, err = external.GetPriceRangeWithAPI(lastPriceRequest.Token, lastPriceRequest.FromTimeStamp, lastPriceRequest.ToTimeStamp)
		if err != nil {
			ErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}

		err = database.GetDBInstance().InsertAPIDataToDB(pairFeed.Id, priceRangeData, dbPriceRange)
		if err != nil {
			ErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}

	} else {
		for _, val := range dbPriceRange {
			temp := new(models.CustomPriceStruct)
			temp.Close = val.Price
			temp.Time = val.Timestamp.Unix()
			priceRangeData.Data = append(priceRangeData.Data, *temp)
		}
	}
	averagePrice := external.CalculateAveragePrice(priceRangeData)
	ctx.JSON(200, api.AverageRangePriceResponse{
		AverageRangePriceRequest: lastPriceRequest,
		AveragePrice:             averagePrice,
	})

}
