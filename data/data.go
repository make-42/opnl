package data

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"opnl/config"
	"slices"
	"strings"
)

type DataS struct {
	Status Status                `json:"status"`
	Data   map[string][]CoinData `json:"data"`
}

type Status struct {
	Timestamp    string  `json:"timestamp"`
	ErrorCode    int     `json:"error_code"`
	ErrorMessage *string `json:"error_message"`
	Elapsed      int     `json:"elapsed"`
	CreditCount  int     `json:"credit_count"`
	Notice       *string `json:"notice"`
}

type CoinData struct {
	ID                            int                  `json:"id"`
	Name                          string               `json:"name"`
	Symbol                        string               `json:"symbol"`
	Slug                          string               `json:"slug"`
	NumMarketPairs                int                  `json:"num_market_pairs"`
	DateAdded                     string               `json:"date_added"`
	Tags                          []Tag                `json:"tags"`
	MaxSupply                     *int                 `json:"max_supply"`
	CirculatingSupply             float64              `json:"circulating_supply"`
	TotalSupply                   float64              `json:"total_supply"`
	IsActive                      int                  `json:"is_active"`
	InfiniteSupply                bool                 `json:"infinite_supply"`
	Platform                      *Platform            `json:"platform"`
	CmcRank                       int                  `json:"cmc_rank"`
	IsFiat                        int                  `json:"is_fiat"`
	SelfReportedCirculatingSupply *float64             `json:"self_reported_circulating_supply"`
	SelfReportedMarketCap         *float64             `json:"self_reported_market_cap"`
	TvlRatio                      *float64             `json:"tvl_ratio"`
	LastUpdated                   string               `json:"last_updated"`
	Quote                         map[string]QuoteData `json:"quote"`
}

type Tag struct {
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type Platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

type QuoteData struct {
	Price                 float64  `json:"price"`
	Volume24h             float64  `json:"volume_24h"`
	VolumeChange24h       float64  `json:"volume_change_24h"`
	PercentChange1h       float64  `json:"percent_change_1h"`
	PercentChange24h      float64  `json:"percent_change_24h"`
	PercentChange7d       float64  `json:"percent_change_7d"`
	PercentChange30d      float64  `json:"percent_change_30d"`
	PercentChange60d      float64  `json:"percent_change_60d"`
	PercentChange90d      float64  `json:"percent_change_90d"`
	MarketCap             *float64 `json:"market_cap"`
	MarketCapDominance    *float64 `json:"market_cap_dominance"`
	FullyDilutedMarketCap *float64 `json:"fully_diluted_market_cap"`
	Tvl                   *float64 `json:"tvl"`
	LastUpdated           string   `json:"last_updated"`
}

var Data DataS

func GetData() {
	client := &http.Client{}
	symbols := []string{}
	for _, order := range config.Config.Portfolio.Orders {
		symbols = append(symbols, order.Symbol)
	}
	slices.Sort(symbols)
	symbols = slices.Compact(symbols)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?symbol=%s&convert=%s", strings.Join(symbols, ","), config.Config.UserCurrencySymbol), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-CMC_PRO_API_KEY", config.Config.CoinMarketCapAPIKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&Data)
	if err != nil {
		log.Fatal(err)
	}
}
