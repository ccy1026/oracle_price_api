package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"oracle_price_api/database"
	"oracle_price_api/models"
	"oracle_price_api/utils/external"
	"time"
)

func FetchPriceCronJob() {
	data, err := database.GetDBInstance().GetPairAddress()
	if err != nil {
		fmt.Println(err)
	}

	for _, val := range data {
		priceTemp := new(models.PriceData)
		priceData, err := external.GetLatestPrice(common.HexToAddress(val.OracleAddress))
		if err != nil {
			log.Print(err)
			continue
		}

		var now = time.Now()
		lastMinTime := time.Unix(now.Unix(), 0).Truncate(time.Minute)
		updateTimeStamp := time.Unix(priceData.UpdatedAt.Int64(), 0)

		formatTimestamp := updateTimeStamp.Truncate(time.Minute)

		if updateTimeStamp.Second() != 0 {
			formatTimestamp = formatTimestamp.Add(1 * time.Minute)
		}

		if lastMinTime.Equal(formatTimestamp) {
			priceTemp.Price = PriceFormat(big.NewFloat(float64(priceData.Answer.Uint64())), big.NewFloat(100000000), 4)
		} else {
			priceTemp.Price, err = external.GetPriceWithTimestamp(val.Token, lastMinTime.Unix())
			if err != nil {
				log.Println(err)
			}
		}
		priceTemp.PairID = val.Id
		priceTemp.Timestamp = lastMinTime

		err = database.GetDBInstance().InsertDataToDB(priceTemp)
		if err != nil {
			log.Println(err)
			continue
		}

	}

}
