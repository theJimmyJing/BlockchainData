package main

// uniswap token struct
type UniswapResp struct {
	Data UniswapData `json:"data"`
}
type UniswapData struct {
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
	Hold        int `json:"hold"`
	OwnerNum    int `json:"ownerNum"`
	MarketValue int `json:"marketValue"`
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
	// 少了1对
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

// 埋点事件
type EventData struct {
	UserId  string `json:"userId"`  // 用户
	IP      string `json:"ip"`      // IP
	Device  string `json:"device"`  // 设备
	Os      string `json:"system"`  // 设备系统
	Browser string `json:"browser"` // 浏览器

	Page    string `json:"page"`    // 页面
	Event   string `json:"event"`   // 事件
	Action  string `json:"action"`  // 动作
	Comment string `json:"comment"` // comment
	Date    string `json:"date"`    // 时间
}

// 从1inch获取币兑换数量，可用于计算币价
type Price1inch struct {
	ToTokenAmount   string `json:"toTokenAmount"`
	FromTokenAmount string `json:"fromTokenAmount"`
}
