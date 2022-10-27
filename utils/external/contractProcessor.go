package external

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"oracle_price_api/utils/contract"
	"oracle_price_api/utils/geth"
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
