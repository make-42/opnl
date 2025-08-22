package config

import (
	"opnl/utils"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"gopkg.in/yaml.v2"
)

type PortfolioS struct {
	Orders []OrderS `yaml:"orders"`
}

type OrderS struct {
	Symbol   string  `yaml:"symbol"`
	Quantity float64 `yaml:"quantity"`
	Price    float64 `yaml:"price"` // In user currency
	Date     string  `yaml:"date"`  // ISO 8601 format
	Fees     float64 `yaml:"fees"`  // In user currency
}

type ConfigS struct {
	CoinMarketCapAPIKey        string     `yaml:"coinmarketcap_api_key"`
	UpdateFrequency            int        `yaml:"update_frequency"` // in seconds (under 60s is useless because of CoinMarketCap API caching)
	UserCurrencySymbol         string     `yaml:"user_currency_symbol"`
	UserCurrencyDisplay        string     `yaml:"user_currency_display"`
	DisplayTotalPortfolioValue bool       `yaml:"display_total_portfolio_value"`
	DisplayPortfolio           bool       `yaml:"display_portfolio"`
	DisplayWithHTMLColor       bool       `yaml:"display_with_html_color"`
	PositiveColor              string     `yaml:"positive_color"`
	NegativeColor              string     `yaml:"negative_color"`
	NeutralColor               string     `yaml:"neutral_color"`
	Portfolio                  PortfolioS `yaml:"portfolio"`
}

var DefaultConfig = ConfigS{
	CoinMarketCapAPIKey:        "",
	UpdateFrequency:            300,
	UserCurrencySymbol:         "EUR",
	DisplayTotalPortfolioValue: true,
	DisplayPortfolio:           true,
	PositiveColor:              "#00FF00",
	NegativeColor:              "#FF0000",
	NeutralColor:               "#FFFF00",
	Portfolio: PortfolioS{
		Orders: []OrderS{
			{Symbol: "BTC", Quantity: 1.0, Price: 50000.0, Date: "2021-01-01T00:00:00Z", Fees: 0.0},
			{Symbol: "BTC", Quantity: 2.0, Price: 60000.0, Date: "2021-02-01T00:00:00Z", Fees: 0.0},
		},
	},
}

var Config ConfigS

func Init() {
	configPath := configdir.LocalConfig("ontake", "opnl")
	err := configdir.MakePath(configPath) // Ensure it exists.
	utils.CheckError(err)

	configFile := filepath.Join(configPath, "config.yml")

	// Does the file not exist?
	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		// Create the new config file.
		fh, err := os.Create(configFile)
		utils.CheckError(err)
		defer fh.Close()

		encoder := yaml.NewEncoder(fh)
		encoder.Encode(&DefaultConfig)
		Config = DefaultConfig
	} else {
		// Load the existing file.
		fh, err := os.Open(configFile)
		utils.CheckError(err)
		defer fh.Close()

		decoder := yaml.NewDecoder(fh)
		decoder.Decode(&Config)
	}
}
