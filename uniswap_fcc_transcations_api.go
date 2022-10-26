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

var FCC_TOKEN_ID = "0x171b1daefac13a0a3524fcb6beddc7b31e58e079"

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

	// linux 服务器版本curl不识别--data-raw, 换成--data即可
	cmd := exec.Command("sh", "./uniswap_fcc_transactions.sh") // chmod -R 777 Tokenlist (Permission failed)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	// var resp = FccTranscationsResp{}
	err := cmd.Run()
	if err != nil {
		log.Fatalf("FccTransactions cmd.Run() failed with ", err)
	}
	outStr, _ := string(stdout.Bytes()), string(stderr.Bytes())

	// outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	// fmt.Println("------------------------")
	// fmt.Println("FccTransactions outStr : ", outStr)
	// fmt.Println("------------------------")
	// fmt.Println("FccTransactions errStr ：", errStr)
	// fmt.Println("------------------------")

	redisClient := connectRedis()
	saveFccTransactionsData(redisClient, outStr)
}

// save FccTransactions to redis
func saveFccTransactionsData(redisClient *redis.Client, data string) {
	Fcc_Transactions_KEY := "fcc_transactions"
	//转换成字符串
	infoErrorStatus := redisClient.Set(Fcc_Transactions_KEY, data, 0).Err() //设置过期时间- 不过期
	if infoErrorStatus != nil {
		fmt.Println("saveFccTransactionsData failed：", infoErrorStatus)
		// } else {
		// 	fmt.Println("Fcc Transactions Data updated")
	}
}

func getCachedFccTransactionsData(redisClient *redis.Client) FccTranscationsResp {
	Fcc_Transactions_KEY := "fcc_transactions"
	getInfo, getinfoErr := redisClient.Get(Fcc_Transactions_KEY).Result()

	dataCached := FccTranscationsResp{}
	if getinfoErr != nil {
		fmt.Println("getCachedFccTransactionsData no data", getinfoErr)
	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &dataCached)
		if unmarsha1Err != nil {
			fmt.Println("反序列化失败:", unmarsha1Err)
			// } else {
			// 	fmt.Println("cached", dataCache)
		}
	}

	return dataCached
}

// 获取FCC与指定token(id)的交易
func getTransactionsWithToken(redisClient *redis.Client, tokenId string) FccTranscationsResp {
	dataCached := getCachedFccTransactionsData(redisClient)

	transactions := []FccTransactions{}

	for i := 0; i < len(dataCached.Data.Transactions); i++ {

		swaps := dataCached.Data.Transactions[i].Swaps
		if swaps != nil {
		TAG:
			for j := 0; j < len(swaps); j++ {
				if swaps[j].Token0.ID == tokenId || swaps[j].Token1.ID == tokenId {
					transactions = append(transactions, dataCached.Data.Transactions[i])
					break TAG
				} else if swaps[j].Transaction.Swaps != nil {
					swapsInner := swaps[j].Transaction.Swaps
					for k := 0; k < len(swapsInner); k++ {
						if swapsInner[k].Token0.ID == tokenId || swapsInner[k].Token1.ID == tokenId {
							transactions = append(transactions, dataCached.Data.Transactions[i])
							break TAG
						}
					}
				}
			}
		}
	}

	dataCached.Data.Transactions = transactions

	return dataCached
}

func getCachedFccTransactions(redisClient *redis.Client, symbol string, id string) FccTranscationsResp {
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
		} else {
			fmt.Println(dataCache)
		}
	}

	return dataCache
}

type FccTranscationsResp struct {
	Data FccTranscationsData `json:"data"`
}

type FccTranscationsData struct {
	Token        FccTranscationsToken `json:"token,omitempty,omitempty"`
	Transactions []FccTransactions    `json:"transactions,omitempty,omitempty"`
}

type FccTransactions struct {
	BlockNumber string                 `json:"blockNumber,omitempty"`
	Timestamp   string                 `json:"timestamp,omitempty"`
	ID          string                 `json:"id,omitempty"`
	GasUsed     string                 `json:"gasUsed,omitempty"`
	GasPrice    string                 `json:"gasPrice,omitempty"`
	Swaps       []FccTranscationsSwaps `json:"swaps,omitempty"`
}

type FccTranscationsSwaps struct {
	ID          string              `json:"id,omitempty"`
	Amount0     string              `json:"amount0,omitempty"`
	Amount1     string              `json:"amount1,omitempty"`
	Token0      FccTranscationToken `json:"token0,omitempty"`
	Token1      FccTranscationToken `json:"token1,omitempty"`
	Transaction FccTransaction      `json:"transaction,omitempty"`
}

type FccTranscationsToken struct {
	Decimals                     string `json:"decimals,omitempty"`
	DerivedETH                   string `json:"derivedETH,omitempty"`
	FeesUSD                      string `json:"feesUSD,omitempty"`
	ID                           string `json:"id,omitempty"`
	Name                         string `json:"name,omitempty"`
	Symbol                       string `json:"symbol,omitempty"`
	TotalSupply                  string `json:"totalSupply,omitempty"`
	TotalValueLocked             string `json:"totalValueLocked,omitempty"`
	TotalValueLockedUSD          string `json:"totalValueLockedUSD,omitempty"`
	TotalValueLockedUSDUntracked string `json:"totalValueLockedUSDUntracked,omitempty"`
	TxCount                      string `json:"txCount,omitempty"`
}

type FccTranscationToken struct {
	Decimals string `json:"decimals,omitempty"`
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
}

type FccTransaction struct {
	Swaps []FccTranscationsSwaps `json:"swaps,omitempty"`
}
