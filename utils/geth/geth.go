package geth

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type EVMClients struct {
	EthClient *ethclient.Client
}

var ClientStruct EVMClients

func GetClient() *EVMClients {
	return &ClientStruct
}

func ConnectEVMClient() error {
	client := GetClient()
	var err error

	client.EthClient, err = ethclient.Dial("https://mainnet.infura.io/v3/3630895f60c94b159c58e16c0680b93a")
	if err != nil {
		return err
	}
	return nil
}
