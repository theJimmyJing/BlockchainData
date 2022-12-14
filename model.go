package main

// uniswap token struct
type UniswapTokenResp struct {
	Data UniswapTokenData `json:"data"`
}
type UniswapTokenData struct {
	Token UniswapToken `json:"token"`
}

// uniswap token struct
type UniswapToken struct {
	DerivedETH                   string `json:"derivedETH"`
	FeesUSD                      string `json:"feesUSD"`
	Name                         string `json:"name"`
	PoolCount                    string `json:"poolCount"`
	Symbol                       string `json:"symbol"`
	TotalSupply                  string `json:"totalSupply"`
	TotalValueLocked             string `json:"totalValueLocked"`
	TotalValueLockedUSD          string `json:"totalValueLockedUSD"`
	TotalValueLockedUSDUntracked string `json:"totalValueLockedUSDUntracked"`
	TxCount                      string `json:"txCount"`
	UntrackedVolumeUSD           string `json:"untrackedVolumeUSD"`
	Volume                       string `json:"volume"`
	VolumeUSD                    string `json:"volumeUSD"`
	Decimals                     string `json:"decimals"`
}

// uniswap token struct
type OperateData struct {
	Freechat  Freechat    `json:"freechat"`
	User      UserBigData `json:"user"`
	FPay      FPay        `json:"FPay"`
	ECommerce ECommerce   `json:"eCommerce"`
	Ad        Ad          `json:"ad"`
	NFT       NFT         `json:"NFT"`
	Game      Game        `json:"game"`
}

type EquitytokenData struct {
	Hold        int    `json:"hold"`
	OwnerNum    int    `json:"ownerNum"`
	MarketValue string `json:"marketValue"`
}

type EquitytokenRes struct {
	Code int             `json:"code"`
	Data EquitytokenData `json:"data"`
}

type Freechat struct {
	TotalEarn           string `json:"totalEarn"`
	DayEarn             string `json:"dayEarn"`
	DayEarnIncrease     string `json:"dayEarnIncrease"`
	WeekEarn            string `json:"weekEarn"`
	WeekEarnIncrease    string `json:"weekEarnIncrease"`
	MonthEarn           string `json:"monthEarn"`
	MonthEarnIncrease   string `json:"monthEarnIncrease"`
	NowPrice            string `json:"nowPrice"`
	MarketValue         string `json:"marketValue"`
	MarketValueIncrease string `json:"marketValueIncrease"`
	DayVolume           string `json:"dayVolume"`
	DayVolumeIncrease   string `json:"dayVolumeIncrease"`
	FccUser             string `json:"fccUser"`
	FccUserIncrease     string `json:"fccUserIncrease"`
	TotalProfit         string `json:"totalProfit"`
	WaitProfit          string `json:"waitProfit"`
	PerFccProfit        string `json:"perFccProfit"`
	PledgeProfit        string `json:"pledgeProfit"`
	PledgeRate          string `json:"pledgeRate"`
	// ??????1???
}
type UserBigData struct {
	Total                  int `json:"total"`
	DayIncrease            int `json:"dayIncrease"`
	DayActive              int `json:"dayActive"`
	DayActiveIncrease24H   int `json:"dayActiveIncrease"`
	WeekActive             int `json:"weekActive"`
	WeekActiveIncrease24H  int `json:"weekActiveIncrease"`
	MonthActive            int `json:"monthActive"`
	MonthActiveIncrease24H int `json:"monthActiveIncrease"`
}

type FPay struct {
}
type ECommerce struct {
}
type Ad struct {
}
type NFT struct {
}
type Game struct {
}

// ????????????
type EventData struct {
	UserId  string `json:"userId"`  // ??????
	IP      string `json:"ip"`      // IP
	Device  string `json:"device"`  // ??????
	Os      string `json:"system"`  // ????????????
	Browser string `json:"browser"` // ?????????

	Page    string `json:"page"`    // ??????
	Event   string `json:"event"`   // ??????
	Action  string `json:"action"`  // ??????
	Comment string `json:"comment"` // comment
	Date    int    `json:"date"`    // ??????
}

// ???1inch?????????????????????????????????????????????
type Price1inch struct {
	ToTokenAmount   string `json:"toTokenAmount"`
	FromTokenAmount string `json:"fromTokenAmount"`
}
