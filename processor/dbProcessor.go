package processor

import (
	"gorm.io/gorm"
	"oracle_price_api/models"
)

func GetPairAddress(db *gorm.DB) ([]models.PriceFeed, error) {
	var priceFeed []models.PriceFeed
	err := db.Find(&priceFeed).Error
	if err != nil {
		return nil, err
	}
	return priceFeed, nil
}

func InsertDataToDB(db *gorm.DB, data *models.PriceData) error {
	err := db.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
