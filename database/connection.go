package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"oracle_price_api/models"
	"oracle_price_api/utils/external"
)

type DBProcessor struct {
	DBInstance *gorm.DB
}

var DBStruct DBProcessor

func GetDBInstance() *DBProcessor {
	return &DBStruct
}

func Setup() error {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

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

func (db *DBProcessor) Migration() error {
	var err error
	if (!db.DBInstance.Migrator().HasTable(&models.PriceFeed{})) {
		err = db.DBInstance.AutoMigrate(&models.PriceFeed{})
		if err != nil {
			return err
		}
		err = db.InsertChainData()
		if err != nil {
			return err
		}
	}
	if (!db.DBInstance.Migrator().HasTable(&models.PriceData{})) {
		err = db.DBInstance.AutoMigrate(&models.PriceData{})
		if err != nil {
			return err
		}
		err = db.InsertPriceData()
		if err != nil {
			return err
		}

	}
	return nil
}

func (db *DBProcessor) InsertChainData() error {
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
		models.PriceFeed{
			Token:         "LINK",
			TokenAddress:  "0x514910771AF9Ca656af840dff83E8264EcF986CA",
			OracleAddress: "0x2c1d072e956AFFC0D435Cb7AC38EF18d24d9127c",
		},
	}
	tx := db.DBInstance.Create(PriceFeedData)
	err := tx.Error
	if err != nil {
		return err
	}

	return nil
}

func (db *DBProcessor) InsertPriceData() error {
	ETHData, err := external.GetPriceRangeWithAPI("ETH", 1666850160, 1666856580)
	err = db.InsertAPIDataToDB(1, ETHData, nil)
	if err != nil {
		return err
	}

	BTCData, err := external.GetPriceRangeWithAPI("BTC", 1666850160, 1666856580)
	err = db.InsertAPIDataToDB(2, BTCData, nil)
	if err != nil {
		return err
	}

	MATICData, err := external.GetPriceRangeWithAPI("MATIC", 1666850160, 1666856580)
	err = db.InsertAPIDataToDB(3, MATICData, nil)
	if err != nil {
		return err
	}

	BNBData, err := external.GetPriceRangeWithAPI("BNB", 1666850160, 1666856580)
	err = db.InsertAPIDataToDB(4, BNBData, nil)
	if err != nil {
		return err
	}

	LINKData, err := external.GetPriceRangeWithAPI("LINK", 1666850160, 1666856580)
	err = db.InsertAPIDataToDB(5, LINKData, nil)
	if err != nil {
		return err
	}
	return nil
}
