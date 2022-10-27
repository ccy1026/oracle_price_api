package database

import (
	"oracle_price_api/models"
)

func (db *DBProcessor) GetPairAddress() ([]models.PriceFeed, error) {
	var priceFeed []models.PriceFeed
	err := db.DBInstance.Find(&priceFeed).Error
	if err != nil {
		return nil, err
	}
	return priceFeed, nil
}

func (db *DBProcessor) QuerySymbol(symbol string) (*models.PriceFeed, error) {
	var feeds models.PriceFeed
	tx := db.DBInstance.Where("token = ?", symbol).Find(&feeds)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	return &feeds, nil
}

func (db *DBProcessor) CheckSupportToken(symbol string) (bool, error) {
	data, err := db.QuerySymbol(symbol)
	if err != nil {
		return false, err
	}
	if data.Id == 0 {
		return false, err
	}
	return true, nil
}
