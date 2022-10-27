package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"log"
	"oracle_price_api/database"
	"oracle_price_api/routers"
	"oracle_price_api/utils"
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
	log.Println("Database Init success")
	err = geth.ConnectEVMClient()
	if err != nil {
		panic("EVM Failed")
	}
	log.Println("EVM Init success")

}

func main() {
	go func() {
		gocron.Every(1).Minute().Do(utils.FetchPriceCronJob)
		<-gocron.Start()
	}()
	log.Println("Running")
	r := routers.SetupRouter()
	r.Run(":1234")
}
