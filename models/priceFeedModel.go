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
	Id        uint `gorm:"primaryKey"`
	PairID    uint `gorm:"index:,unique,composite:pricePair"`
	Price     float64
	Timestamp time.Time `gorm:"index:,unique,composite:pricePair"`
}
