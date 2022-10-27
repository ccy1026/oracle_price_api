package api

type AverageRangePriceRequest struct {
	Token         string `json:"token" validate:"required"`
	FromTimeStamp int64  `json:"from_time_stamp" validate:"required"`
	ToTimeStamp   int64  `json:"to_time_stamp" validate:"required"`
}
type AverageRangePriceResponse struct {
	AverageRangePriceRequest
	AveragePrice float64 `json:"average_price"`
}
