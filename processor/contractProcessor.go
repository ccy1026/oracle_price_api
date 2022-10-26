package processor

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"oracle_price_api/models"
	"oracle_price_api/utils"
	"oracle_price_api/utils/contract"
	"oracle_price_api/utils/database"
	"oracle_price_api/utils/geth"
	"time"
)

type LatestPrice struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

func GetLatestPrice(pairAddress common.Address) (LatestPrice, error) {
	var data LatestPrice

	IChainLink, err := contract.NewChainlinkFeed(pairAddress, geth.GetClient().EthClient)
	if err != nil {
		return data, err
	}
	data, err = IChainLink.LatestRoundData(&bind.CallOpts{})
	return data, nil
}

func FetchPriceCronJob() {
	data, err := GetPairAddress(database.GetDBInstance().DBInstance)
	if err != nil {
		fmt.Println(err)
	}

	for _, val := range data {
		priceTemp := new(models.PriceData)
		priceData, err := GetLatestPrice(common.HexToAddress(val.OracleAddress))
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
			priceTemp.Price = utils.PriceFormat(big.NewFloat(float64(priceData.Answer.Uint64())), big.NewFloat(100000000), 4)
		} else {
			//TODO fetch other source
		}
		priceTemp.PairID = val.Id
		priceTemp.Timestamp = lastMinTime
		priceTemp.CreateTime = lastMinTime

		err = InsertDataToDB(database.GetDBInstance().DBInstance, priceTemp)
		if err != nil {
			log.Println(err)
			continue
		}

	}

}
