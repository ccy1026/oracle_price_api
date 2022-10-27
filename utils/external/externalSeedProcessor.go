package external

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"oracle_price_api/models"
	"os"
	"strconv"
	"time"
)

func GetPriceWithTimestamp(tokenName string, timestamp int64) (float64, error) {

	request, err := http.Get(os.Getenv("API_URL") + "/histominute?fsym=" + tokenName + "&tsym=USD&toTS=" + strconv.FormatInt(timestamp, 10) + "&limit=1&api_key=" + os.Getenv("API_KEY"))
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return 0, err
	}
	defer request.Body.Close()
	priceData := new(models.APIPriceModel)
	err = json.Unmarshal(body, &priceData)
	if err != nil {
		return 0, err
	}
	return priceData.Data.Data[1].Close, nil
}
func GetPriceRangeWithAPI(tokenName string, fromTimestamp int64, toTimestamp int64) (*models.PriceDataModel, error) {
	limit := strconv.FormatInt((time.Now().Unix()-fromTimestamp)/60, 10)
	request, err := http.Get(os.Getenv("API_URL") + "/histominute?fsym=" + tokenName + "&tsym=USD&toTS=" + strconv.FormatInt(toTimestamp, 10) + "&limit=" + limit + "&api_key=" + os.Getenv("API_KEY"))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()
	priceData := new(models.APIPriceModel)
	err = json.Unmarshal(body, &priceData)
	if err != nil {
		return nil, err
	}

	resStruct := new(models.PriceDataModel)

	for _, val := range priceData.Data.Data {
		if val.Time >= fromTimestamp && val.Time <= toTimestamp {
			resStruct.Data = append(resStruct.Data, val)
		}
	}

	return resStruct, nil
}
func CalculateAveragePrice(model *models.PriceDataModel) float64 {
	var temp float64
	for _, val := range model.Data {
		temp += val.Close
	}
	return temp / float64(len(model.Data))
}
