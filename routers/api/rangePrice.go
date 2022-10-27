package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"oracle_price_api/database"
	"oracle_price_api/models"
	"oracle_price_api/models/api"
	"oracle_price_api/utils"
	"oracle_price_api/utils/external"
	"time"
)

// @BasePath /
// @Summary Get average price for the time range with specific token
// @Description Get average price for the time range. If not exist will add into database for cache
// @Tags Price
// @Success 200
// @produce application/json
// @Param rangePrice body api.AverageRangePriceRequest true "Range Price"
// @Router /rangePrice [post]
func GetAveragePriceFromRange(ctx *gin.Context) {
	var lastPriceRequest = api.AverageRangePriceRequest{}
	err := ctx.ShouldBindBodyWith(&lastPriceRequest, binding.JSON)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = utils.VerifyJsonBodyByStruct(lastPriceRequest)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	db := database.GetDBInstance()
	pairFeed, err := db.QuerySymbol(ctx.Param("token"))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	lowerTime := time.Unix(lastPriceRequest.FromTimeStamp, 0).Truncate(time.Minute)
	upperTime := time.Unix(lastPriceRequest.ToTimeStamp, 0).Truncate(time.Minute).Add(1 * time.Minute)
	pair := (upperTime.Unix() - lowerTime.Unix()) / 60
	dbPriceRange, err := database.GetDBInstance().GetLatestPriceInTimeRange(lastPriceRequest.Token, lowerTime.Unix(), upperTime.Unix())
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	priceRangeData := new(models.PriceDataModel)

	if len(dbPriceRange) != int(pair) {
		priceRangeData, err = external.GetPriceRangeWithAPI(lastPriceRequest.Token, lastPriceRequest.FromTimeStamp, lastPriceRequest.ToTimeStamp)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
			return
		}

		err = database.GetDBInstance().InsertAPIDataToDB(pairFeed.Id, priceRangeData, dbPriceRange)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusInternalServerError, err)
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
