package main

// uniswap token struct
type UTokenInfo struct {
	id                           string `json:"id"`
	decimals                     string `json:"decimals"`
	name                         string `json:"name"`
	feesUSD                      string `json:"feesUSD"`
	derivedETH                   string `json:"derivedETH"`
	poolCount                    string `json:"poolCount"`
	symbol                       string `json:"symbol"`
	totalSupply                  string `json:"totalSupply"`
	totalValueLocked             string `json:"totalValueLocked"`
	totalValueLockedUSDUntracked string `json:"totalValueLockedUSDUntracked"`
	totalValueLockedUSD          string `json:"totalValueLockedUSD"`
	txCount                      string `json:"txCount"`
	untrackedVolumeUSD           string `json:"untrackedVolumeUSD"`
	volume                       string `json:"volume"`
	volumeUSD                    string `json:"volumeUSD"`
}

// uniswap token#DayData struct
type UTokenDayDataInfo struct {
	date                string `json:"date"`
	feesUSD             string `json:"feesUSD"`
	high                string `json:"high"`
	id                  string `json:"id"`
	low                 string `json:"low"`
	priceUSD            string `json:"priceUSD"`
	totalValueLocked    string `json:"totalValueLocked"`
	totalValueLockedUSD string `json:"totalValueLockedUSD"`
	untrackedVolumeUSD  string `json:"untrackedVolumeUSD"`
	volume              string `json:"volume"`
	volumeUSD           string `json:"volumeUSD"`
}
