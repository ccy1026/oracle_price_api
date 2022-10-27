package database

import (
	"gorm.io/gorm"
	"oracle_price_api/models"
	"time"
)

func (db *DBProcessor) GetLastPriceDataWithTokenName(tokenName string, timestamp int64) (*models.PriceData, error) {
	data := new(models.PriceData)
	var tx *gorm.DB
	if timestamp != 0 {
		tx = db.DBInstance.Raw("SELECT * FROM price_data WHERE pair_id IN (SELECT id FROM price_feeds Where token = ?) AND timestamp > CAST(? AS DATE) ORDER BY timestamp DESC Limit 1 ;", tokenName, timestamp).Scan(&data)
	} else {
		tx = db.DBInstance.Raw("SELECT * FROM price_data WHERE pair_id IN (SELECT id FROM price_feeds Where token = ?) ORDER BY timestamp DESC Limit 1 ;", tokenName).Scan(&data)
	}

	err := tx.Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (db *DBProcessor) GetLatestPriceInTimeRange(tokenName string, fromTimestamp int64, toTimestamp int64) ([]models.PriceData, error) {
	data := new([]models.PriceData)
	time.Unix(fromTimestamp, 0).String()

	tx := db.DBInstance.Raw("SELECT * FROM price_data Where timestamp > CAST(? AS DATE) OR timestamp < CAST(? AS DATE) and pair_id IN (SELECT id FROM price_feeds Where token = ?) ORDER BY timestamp DESC;", fromTimestamp, toTimestamp, tokenName).Scan(&data)
	err := tx.Error
	if err != nil {
		return nil, err
	}
	return *data, nil
}

func (db *DBProcessor) InsertDataToDB(data *models.PriceData) error {
	err := db.DBInstance.Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *DBProcessor) InsertAPIDataToDB(pairID uint, priceModelData *models.PriceDataModel, dbData []models.PriceData) error {
	var data []*models.PriceData

	if len(dbData) == 0 {
		for _, val := range priceModelData.Data {
			temp := new(models.PriceData)
			temp.Timestamp = time.Unix(val.Time, 0)
			temp.Price = val.Close
			temp.PairID = pairID
			data = append(data, temp)
		}
	} else {
		for _, val := range priceModelData.Data {
			var exist bool
			for _, dbVal := range dbData {
				if dbVal.Timestamp.Equal(time.Unix(val.Time, 0)) {
					exist = true
				}
			}
			if !exist {
				temp := new(models.PriceData)
				temp.Timestamp = time.Unix(val.Time, 0)
				temp.Price = val.Close
				temp.PairID = pairID
				data = append(data, temp)
			}
		}

	}
	if data != nil {
		err := db.DBInstance.Create(&data).Error
		if err != nil {
			return err
		}
	}

	return nil

}

func removeElementByIndex[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}
