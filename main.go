package main

import (
	"fmt"
	"opnl/config"
	"opnl/data"
	"opnl/pnl"
	"sort"
	"time"
)

func main() {
	config.Init()
	for {
		data.GetData()
		pnlVal, totalValue := pnl.CalculatePNLAndValue()

		output := ""
		if config.Config.DisplayTotalPortfolioValue {
			if config.Config.DisplayWithHTMLColor {
				output += fmt.Sprintf("<span color='%s'>", config.Config.NeutralColor)
			}
			output += fmt.Sprintf("%.2f %s", totalValue, config.Config.UserCurrencySymbol)
			if config.Config.DisplayWithHTMLColor {
				output += "</span>"
			}
			output += " "
		}

		sign := ""
		if pnlVal >= 0 {
			sign = "+"
			if config.Config.DisplayWithHTMLColor {
				output += fmt.Sprintf("<span color='%s'>", config.Config.PositiveColor)
			}
		} else {
			if config.Config.DisplayWithHTMLColor {
				output += fmt.Sprintf("<span color='%s'>", config.Config.NegativeColor)
			}
		}
		output += fmt.Sprintf("%s%.2f %s", sign, pnlVal, config.Config.UserCurrencySymbol)
		if config.Config.DisplayWithHTMLColor {
			output += "</span>"
		}
		if config.Config.DisplayPortfolio {
			qavps := pnl.CalculatePNLQuantityAndValuePerSymbol()
			keys := make([]string, 0, len(qavps))

			for key := range qavps {
				keys = append(keys, key)
			}

			sort.SliceStable(keys, func(i, j int) bool {
				return qavps[keys[i]].Value < qavps[keys[j]].Value
			})
			for key := range qavps {
				pnlVal := qavps[key].PNL
				pnlSign := ""
				pnlHTML := fmt.Sprintf("<span color='%s'>", config.Config.NegativeColor)
				if pnlVal >= 0 {
					pnlSign = "+"
					pnlHTML = fmt.Sprintf("<span color='%s'>", config.Config.PositiveColor)
				}
				val := qavps[key].Value
				if config.Config.DisplayWithHTMLColor {
					output += fmt.Sprintf(" | <span color='%s'>%f %s %.2f %s</span> %s%s%.2f %s</span>", config.Config.NeutralColor, qavps[key].Quantity, key, val, config.Config.UserCurrencySymbol, pnlHTML, pnlSign, pnlVal, config.Config.UserCurrencySymbol)
				} else {
					output += fmt.Sprintf(" | %f %s %.2f %s %s%.2f %s", qavps[key].Quantity, key, val, config.Config.UserCurrencySymbol, pnlSign, pnlVal, config.Config.UserCurrencySymbol)
				}
			}
		}
		fmt.Println(output)
		time.Sleep(time.Duration(config.Config.UpdateFrequency) * time.Second)
	}
}
