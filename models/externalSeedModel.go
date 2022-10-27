package models

type CustomPriceStruct struct {
	Time  int64   `json:"time"`
	Close float64 `json:"close"`
}

type PriceDataModel struct {
	Data []struct {
		Time  int64   `json:"time"`
		Close float64 `json:"close"`
	} `json:"Data"`
}

type APIPriceModel struct {
	Data struct {
		TimeFrom int `json:"TimeFrom"`
		TimeTo   int `json:"TimeTo"`
		PriceDataModel
	} `json:"Data"`
}
