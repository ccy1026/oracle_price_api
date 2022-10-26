package models

import (
	"time"
)

type PriceFeed struct {
	Id            uint `gorm:"primaryKey"`
	Token         string
	TokenAddress  string
	OracleAddress string
}

type PriceData struct {
	Id         uint `gorm:"primaryKey"`
	PairID     uint
	Price      float64
	Timestamp  time.Time
	CreateTime time.Time
}

func GetTimeWithTimeZone() time.Time {
	loc, _ := time.LoadLocation("Asia/Taipei")
	return time.Now().UTC().In(loc)
}
