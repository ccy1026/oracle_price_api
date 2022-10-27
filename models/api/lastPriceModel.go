package api

import "time"

type PriceResponse struct {
	Token          string    `json:"token"`
	Price          float64   `json:"price"`
	LastUpdateTime time.Time `json:"last_update_time"`
}
