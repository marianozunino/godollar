package dto

type RateDto struct {
	Date           string  `json:"date"`
	BuyPrice       float64 `json:"buyPrice"`
	SellPrice      float64 `json:"sellPrice"`
	EbourSellPrice float64 `json:"ebourSellPrice"`
	EbourBuyPrice  float64 `json:"ebourBuyPrice"`
}
