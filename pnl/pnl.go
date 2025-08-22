package pnl

import (
	"opnl/config"
	"opnl/data"
)

func CalculatePNLAndValue() (float64, float64) {
	var pnl float64
	var totalValue float64
	for _, order := range config.Config.Portfolio.Orders {
		price := data.Data.Data[order.Symbol][0].Quote[config.Config.UserCurrencySymbol].Price
		pnl += (price - order.Price) * order.Quantity
		pnl -= order.Fees
		totalValue += price * order.Quantity
	}
	return pnl, totalValue
}

func CalculatePNLQuantityAndValuePerSymbol() map[string]struct {
	PNL      float64
	Quantity float64
	Value    float64
} {
	var output map[string]struct {
		PNL      float64
		Quantity float64
		Value    float64
	}
	output = make(map[string]struct {
		PNL      float64
		Quantity float64
		Value    float64
	})
	for _, order := range config.Config.Portfolio.Orders {
		price := data.Data.Data[order.Symbol][0].Quote[config.Config.UserCurrencySymbol].Price
		output[order.Symbol] = struct {
			PNL      float64
			Quantity float64
			Value    float64
		}{
			PNL:      output[order.Symbol].PNL + (price-order.Price)*order.Quantity - order.Fees,
			Quantity: output[order.Symbol].Quantity + order.Quantity,
			Value:    output[order.Symbol].Value + price*order.Quantity,
		}
	}
	return output
}
