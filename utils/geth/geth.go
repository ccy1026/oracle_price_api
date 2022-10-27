package geth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
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

	client.EthClient, err = ethclient.Dial(os.Getenv("RPC"))
	if err != nil {
		return err
	}
	return nil
}
