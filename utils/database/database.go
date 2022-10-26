package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"oracle_price_api/models"
)

type DBProcessor struct {
	DBInstance *gorm.DB
}

var DBStruct DBProcessor

func GetDBInstance() *DBProcessor {
	return &DBStruct
}

func Setup() error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	temp := GetDBInstance()
	temp.DBInstance = db

	err = temp.Migration()
	if err != nil {
		return err
	}
	return nil
}

func (me *DBProcessor) Migration() error {
	db := me.DBInstance
	if (!db.Migrator().HasTable(&models.PriceFeed{})) {
		db.AutoMigrate(&models.PriceFeed{})
		err := me.InsertChainData()
		if err != nil {
			log.Println(err)
			return err
		}
	}
	if (!db.Migrator().HasTable(&models.PriceData{})) {
		db.AutoMigrate(&models.PriceData{})
	}
	return nil
}

func (me *DBProcessor) InsertChainData() error {
	db := me.DBInstance
	var PriceFeedData = []models.PriceFeed{
		models.PriceFeed{
			Token:         "ETH",
			TokenAddress:  "",
			OracleAddress: "0x5f4ec3df9cbd43714fe2740f5e3616155c5b8419",
		},
		models.PriceFeed{
			Token:         "BTC",
			TokenAddress:  "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
			OracleAddress: "0xf4030086522a5beea4988f8ca5b36dbc97bee88c",
		},
		models.PriceFeed{
			Token:         "MATIC",
			TokenAddress:  "0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0",
			OracleAddress: "0x7bac85a8a13a4bcd8abb3eb7d6b4d632c5a57676",
		},
		models.PriceFeed{
			Token:         "BNB",
			TokenAddress:  "0xb8c77482e45f1f44de1745f52c74426c631bdd52",
			OracleAddress: "0x14e613ac84a31f709eadbdf89c6cc390fdc9540a",
		},
	}
	tx := db.Create(PriceFeedData)
	err := tx.Error
	if err != nil {
		return err
	}
	return nil
}
