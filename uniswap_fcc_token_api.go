package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

// uniswap token struct
type FccTokenResp struct {
	Data FccTokenData `json:"data"`
}
type FccTokenData struct {
	Token FccToken `json:"token"`
}

// uniswap token struct
type FccToken struct {
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

/*
* 查询uniswap token接口
* {"derivedETH":"0","feesUSD":"108.3979488349600340149325490092877","name":"Freechat Coin","poolCount":"0",
* "symbol":"FCC","totalSupply":"28368","totalValueLocked":"199940134.57160359567214293","totalValueLockedUSD":"0",
* "totalValueLockedUSDUntracked":"0","txCount":"108","untrackedVolumeUSD":"18066.32480582667233582209150154796",
* "volume":"456421.018775415938823971","volumeUSD":"36132.6496116533446716441830030959","decimals":"18"}
 */
func uniswapFCCToken() FccToken {
	// scurl := `'https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3' \
	//   -H 'authority: api.thegraph.com' \
	//   -H 'accept: application/json, multipart/mixed' \
	//   -H 'accept-language: zh-CN,zh;q=0.9' \
	//   -H 'content-type: application/json' \
	//   -H 'origin: https://api.thegraph.com' \
	//   -H 'referer: https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3/graphql?query=%7B%0A++token%28id%3A+%220x171b1daefac13a0a3524fcb6beddc7b31e58e079%22%29+%7B%0A++++decimals%0A++++derivedETH%0A++++feesUSD%0A++++name%0A++++poolCount%0A++%7D%0A%7D' \
	//   -H 'sec-ch-ua: "Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"' \
	//   -H 'sec-ch-ua-mobile: ?0' \
	//   -H 'sec-ch-ua-platform: "macOS"' \
	//   -H 'sec-fetch-dest: empty' \
	//   -H 'sec-fetch-mode: cors' \
	//   -H 'sec-fetch-site: same-origin' \
	//   -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36' \
	//   --data '{"query":"{\n  token(id: \"0x171b1daefac13a0a3524fcb6beddc7b31e58e079\") {\n    decimals\n    derivedETH\n    feesUSD\n    name\n    poolCount\n  }\n}","variables":null,"extensions":{"headers":null}}' \
	//   --compressed`

	// 	scurl := `curl 'https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3' \
	//   -H 'authority: api.thegraph.com' \
	//   -H 'accept: application/json, multipart/mixed' \
	//   -H 'accept-language: zh-CN,zh;q=0.9' \
	//   -H 'content-type: application/json' \
	//   -H 'origin: https://api.thegraph.com' \
	//   -H 'referer: https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3/graphql?query=%7B%0A++token%28id%3A+%220x171b1daefac13a0a3524fcb6beddc7b31e58e079%22%29+%7B%0A++++decimals%0A++++derivedETH%0A++++feesUSD%0A++++name%0A++++poolCount%0A++%7D%0A%7D' \
	//   -H 'sec-ch-ua: "Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"' \
	//   -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36' \
	//   --data-raw '{"query":"{\n  token(id: \"0x171b1daefac13a0a3524fcb6beddc7b31e58e079\") {\n    decimals\n    derivedETH\n    feesUSD\n    name\n    poolCount\n  }\n}","variables":null,"extensions":{"headers":null}}' \
	//   --compressed`

	cmd := exec.Command("sh", "./uniswap_fcc_token.sh") // chmod -R 777 Tokenlist (Permission failed)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误

	var resp = FccTokenResp{}
	err := cmd.Run()
	if err != nil {
		log.Fatalf("FCCToken cmd.Run() failed with", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	fmt.Println("FCCToken outStr : ", outStr)
	fmt.Println("FCCToken errStr ：", errStr)

	// json转结构体
	uerr := json.Unmarshal(stdout.Bytes(), &resp)
	if uerr != nil {
		log.Fatalf("FCCTokenstdout-> UniswapToken err", uerr)
	}
	fmt.Println("FCCToken : ", resp)

	// 结构体转json
	jsonBytes, err := json.Marshal(resp.Data.Token)
	if err != nil {
		fmt.Println("FCCToken struct to bytes err : ", err)
	}
	fmt.Println("FCCToken : ", string(jsonBytes))

	return FccToken(resp.Data.Token)
	// TODO read data from Stdout and
}
