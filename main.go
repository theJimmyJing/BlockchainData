package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kirinlabs/HttpRequest"
)

func main() {
	fmt.Println("hello")
	origin := analysisJson()
	r := getAbi(&origin)
	fmt.Printf("first result ---------\n %+v", *r)
	a := combileDetails(r)
	fmt.Printf("second result ---------\n %+v", *r)
	file, _ := json.MarshalIndent(a, "", " ")

	_ = ioutil.WriteFile("result.json", file, 0644)
}

func combileDetails(o *ABIReq) AllRsp {
	time.Sleep(1000)
	req := HttpRequest.NewRequest()
	// 设置超时时间，不设置时，默认30s
	req.SetTimeout(5)
	var allRsp AllRsp
	allRsp.AllResult = make([]AllResult, 0)
	for i, t := range o.Tokens {
		time.Sleep(1 * time.Second)
		fmt.Println("i: = ", i)
		var allResult AllResult

		allResult.Abi = t.Abi
		allResult.Address = t.Address
		allResult.ChainID = t.ChainID
		allResult.Decimals = t.Decimals
		allResult.LogoURI = t.LogoURI
		allResult.Name = t.Name

		url := `https://api.etherscan.io/api?module=token&action=tokeninfo&contractaddress=` + t.Address + `&apikey=FZTI57USSADTZ2IZI6TSFY1T98S1IU492M `
		fmt.Println(url)
		// GET 默认调用方法
		resp, err := req.Get(url, nil)
		if err != nil {

		} else {
			fmt.Errorf("%+v", err)
		}
		// resp.Content()
		body, err := ioutil.ReadAll(resp.Response().Body)

		if err != nil {
			panic(err.Error())
		}
		var data DetailRsp
		json.Unmarshal(body, &data)
		fmt.Printf("detail rsp : %v\n-----------------", data)
		allResult.Bitcointalk = data.Result[0].Bitcointalk
		allResult.ContractAddress = data.Result[0].ContractAddress
		allResult.TokenName = data.Result[0].TokenName
		allResult.Symbol = data.Result[0].Symbol
		allResult.Divisor = data.Result[0].Divisor
		allResult.TokenType = data.Result[0].TokenType
		allResult.TotalSupply = data.Result[0].TotalSupply
		allResult.BlueCheckmark = data.Result[0].BlueCheckmark
		allResult.Website = data.Result[0].Website
		allResult.Email = data.Result[0].Email
		allResult.Blog = data.Result[0].Blog
		allResult.Reddit = data.Result[0].Reddit
		allResult.Slack = data.Result[0].Slack
		allResult.Facebook = data.Result[0].Facebook
		allResult.Twitter = data.Result[0].Twitter
		allResult.Github = data.Result[0].Github
		allResult.Telegram = data.Result[0].Telegram
		allResult.Wechat = data.Result[0].Wechat
		allResult.Linkedin = data.Result[0].Linkedin
		allResult.Discord = data.Result[0].Discord
		allResult.Whitepaper = data.Result[0].Whitepaper
		allResult.TokenPriceUSD = data.Result[0].TokenPriceUSD
		allRsp.AllResult = append(allRsp.AllResult, allResult)
	}
	return allRsp
}

// 总结果
type AllRsp struct {
	AllResult []AllResult `json:"result"`
}

type AllResult struct {
	Name     string `json: "name"`
	ChainID  int    `json: "chainId"`
	Decimals int    `json: "decimals"`
	Address  string `json: "address"`
	LogoURI  string `json: "logoURI"`
	Abi      string `json: "abi"`

	// 这里是分割线
	ContractAddress string `json:"contractAddress"`
	TokenName       string `json:"tokenName"`
	Symbol          string `json:"symbol"`
	Divisor         string `json:"divisor"`
	TokenType       string `json:"tokenType"`
	TotalSupply     string `json:"totalSupply"`
	BlueCheckmark   string `json:"blueCheckmark"`
	Description     string `json:"description"`
	Website         string `json:"website"`
	Email           string `json:"email"`
	Blog            string `json:"blog"`
	Reddit          string `json:"reddit"`
	Slack           string `json:"slack"`
	Facebook        string `json:"facebook"`
	Twitter         string `json:"twitter"`
	Bitcointalk     string `json:"bitcointalk"`
	Github          string `json:"github"`
	Telegram        string `json:"telegram"`
	Wechat          string `json:"wechat"`
	Linkedin        string `json:"linkedin"`
	Discord         string `json:"discord"`
	Whitepaper      string `json:"whitepaper"`
	TokenPriceUSD   string `json:"tokenPriceUSD"`
}

// 详情页结果
type DetailRsp struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  []DetailResult `json:"result"`
}
type DetailResult struct {
	ContractAddress string `json:"contractAddress"`
	TokenName       string `json:"tokenName"`
	Symbol          string `json:"symbol"`
	Divisor         string `json:"divisor"`
	TokenType       string `json:"tokenType"`
	TotalSupply     string `json:"totalSupply"`
	BlueCheckmark   string `json:"blueCheckmark"`
	Description     string `json:"description"`
	Website         string `json:"website"`
	Email           string `json:"email"`
	Blog            string `json:"blog"`
	Reddit          string `json:"reddit"`
	Slack           string `json:"slack"`
	Facebook        string `json:"facebook"`
	Twitter         string `json:"twitter"`
	Bitcointalk     string `json:"bitcointalk"`
	Github          string `json:"github"`
	Telegram        string `json:"telegram"`
	Wechat          string `json:"wechat"`
	Linkedin        string `json:"linkedin"`
	Discord         string `json:"discord"`
	Whitepaper      string `json:"whitepaper"`
	TokenPriceUSD   string `json:"tokenPriceUSD"`
}

func analysisJson() ABIReq {
	// 打开json文件
	jsonFile, err := os.Open("test.json")

	// 最好要处理以下错误
	if err != nil {
		fmt.Println(err)
	}

	// 要记得关闭
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var user ABIReq
	json.Unmarshal([]byte(byteValue), &user)

	// fmt.Println(user)
	return user
}

func getAbi(o *ABIReq) *ABIReq {
	req := HttpRequest.NewRequest()
	// 设置超时时间，不设置时，默认30s
	req.SetTimeout(5)
	for i, t := range o.Tokens {
		url := `https://api.etherscan.io/api?module=contract&action=getabi&address=` + t.Address + `&apikey=CXXZYY2UUWW2YWVNXFKDGTG4ZKUUQURYJZ`
		// GET 默认调用方法
		resp, err := req.Get(url, nil)
		if err != nil {

		} else {
			fmt.Errorf("%+v", err)
		}
		// resp.Content()
		body, err := ioutil.ReadAll(resp.Response().Body)

		if err != nil {
			panic(err.Error())
		}

		var data ResultRsp
		json.Unmarshal(body, &data)
		fmt.Printf("Results: %v\n----------------- abi: %v", data, data.Result)
		o.Tokens[i].Abi = data.Result
	}
	return o
}

type ABIReq struct {
	Tokens []Tokens `json: "tokens"`
}
type Tokens struct {
	Name     string `json: "name"`
	ChainID  int    `json: "chainId"`
	Symbol   string `json: "symbol"`
	Decimals int    `json: "decimals"`
	Address  string `json: "address"`
	LogoURI  string `json: "logoURI"`
	Abi      string `json: "abi"`
}

type ResultRsp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}
