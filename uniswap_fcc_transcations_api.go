package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/go-redis/redis/v7"
)

var transactionsData FccTranscationsData

// 定时更新交易记录
func uniswap_fcc_transactions_timer() {
	for {
		uniswap_fcc_transactions()
		time.Sleep(time.Minute) // 每分钟刷新
	}
}

// 获取fcc的交易记录
func uniswap_fcc_transactions() {
	// func uniswap_fcc_transactions() FccTransactions {

	cmd := exec.Command("sh", "./uniswap_fcc_transactions.sh") // chmod -R 777 Tokenlist (Permission failed)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	// var resp = FccTranscationsResp{}
	err := cmd.Run()
	if err != nil {
		log.Fatalf("FccTransactions cmd.Run() failed with", err)
	}
	// outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	// fmt.Println("FccTransactions outStr : ", outStr)
	// fmt.Println("FccTransactions errStr ：", errStr)

	redisClient := connectRedis()
	saveFccTransactionsData(redisClient, stdout.Bytes())

	// // json转结构体
	// uerr := json.Unmarshal(stdout.Bytes(), &resp)
	// if uerr != nil {
	// 	log.Fatalf("stdout-> UniswapToken err", uerr)
	// }
	// fmt.Println("FccTransactions resp : ", resp)

	// // 结构体转json
	// jsonBytes, err := json.Marshal(resp.Data.Transactions)
	// if err != nil {
	// 	fmt.Println("FccTransactions struct to bytes err : ", err)
	// }
	// fmt.Println("FccTransactions: ", string(jsonBytes))

	// return UniswapToken(resp.Data.Token)
	// TODO read data from Stdout and
}

// save FccTransactions to redis
func saveFccTransactionsData(redisClient *redis.Client, data []byte) {
	Fcc_Transactions_KEY := "fcc_transactions"

	transactionsValue := string(data)                                                    //转换成字符串
	infoErrorStatus := redisClient.Set(Fcc_Transactions_KEY, transactionsValue, 0).Err() //设置过期时间- 不过期
	if infoErrorStatus != nil {
		fmt.Println("saveFccTransactionsData failed：", infoErrorStatus)
	} else {
		fmt.Println("saveFccTransactionsData update success")
	}
}

func getCachedFccTransactionsData(redisClient *redis.Client) FccTranscationsResp {
	Fcc_Transactions_KEY := "fcc_transactions"
	getInfo, getinfoErr := redisClient.Get(Fcc_Transactions_KEY).Result()

	dataCache := FccTranscationsResp{}
	if getinfoErr != nil {
		fmt.Println("getCachedFccTransactionsData no data", getinfoErr)
	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &dataCache)
		if unmarsha1Err != nil {
			fmt.Println("反序列化失败:", unmarsha1Err)
			// } else {
			// 	fmt.Println(dataCache)
		}
	}

	return dataCache
}

type FccTranscationsResp struct {
	Data FccTranscationsData `json:"data"`
}

type FccTranscationsData struct {
	Token        FccTranscationsToken   `json:"token"`
	Transactions []FccTranscationsSwaps `json:"transactions"`
}

type FccTranscationsToken struct {
	Decimals                     string `json:"decimals"`
	DerivedETH                   string `json:"derivedETH"`
	FeesUSD                      string `json:"feesUSD"`
	ID                           string `json:"id"`
	Name                         string `json:"name"`
	Symbol                       string `json:"symbol"`
	TotalSupply                  string `json:"totalSupply"`
	TotalValueLocked             string `json:"totalValueLocked"`
	TotalValueLockedUSD          string `json:"totalValueLockedUSD"`
	TotalValueLockedUSDUntracked string `json:"totalValueLockedUSDUntracked"`
	TxCount                      string `json:"txCount"`
}
type FccTranscationsToken0 struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
type FccTranscationsToken1 struct {
	Decimals string `json:"decimals"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
}

type FccTransaction struct {
	Swaps []FccTranscationsSwaps `json:"swaps"`
}
type FccTranscationsSwaps struct {
	ID          string                `json:"id"`
	Amount0     string                `json:"amount0"`
	Amount1     string                `json:"amount1"`
	Token0      FccTranscationsToken0 `json:"token0"`
	Token1      FccTranscationsToken1 `json:"token1"`
	Transaction FccTransaction        `json:"transaction"`
}
type FccTransactions struct {
	BlockNumber string                 `json:"blockNumber"`
	Timestamp   string                 `json:"timestamp"`
	ID          string                 `json:"id"`
	GasUsed     string                 `json:"gasUsed"`
	GasPrice    string                 `json:"gasPrice"`
	Swaps       []FccTranscationsSwaps `json:"swaps"`
}
