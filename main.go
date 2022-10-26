package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"log"
	"oracle_price_api/processor"
	"oracle_price_api/routers"
	"oracle_price_api/utils/database"
	"oracle_price_api/utils/geth"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = database.Setup()
	if err != nil {
		panic("Database init failed.")
	}
	err = geth.ConnectEVMClient()
	if err != nil {
		panic("EVM Failed")
	}

}

func main() {
	go func() {
		gocron.Every(1).Minute().Do(processor.FetchPriceCronJob)
		<-gocron.Start()
	}()
	processor.FetchPriceCronJob()
	r := routers.SetupRouter()
	r.Run(":1234")
}
