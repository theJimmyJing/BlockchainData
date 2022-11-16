package main

import (
	"fmt"
	"net/http"
	"strings"
)

func getEquityToken() int {
	// curl 'https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3' \
	//   -H 'authority: api.thegraph.com' \
	//   -H 'accept: application/json, multipart/mixed' \
	//   -H 'accept-language: en,zh-CN;q=0.9,zh;q=0.8' \
	//   -H 'content-type: application/json' \
	//   -H 'origin: https://api.thegraph.com' \
	//   -H 'referer: https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3/graphql?query=%7B%0A++token%28id%3A%220x171b1daefac13a0a3524fcb6beddc7b31e58e079%22%29+%7B%0A++++id%0A++++decimals%0A++++name%0A++++symbol%0A++++poolCount%0A++++totalSupply%0A++++totalValueLocked%0A++++totalValueLockedUSD%0A++++totalValueLockedUSDUntracked%0A++++volume%0A++++volumeUSD%0A++++txCount%0A++%7D%0A++transactions%28%0A++++where%3A+%7Bswaps_%3A+%7Btoken0%3A+%220x171b1daefac13a0a3524fcb6beddc7b31e58e079%22%7D%7D%0A++%29+%7B%0A++++blockNumber%0A++++timestamp%0A++++id%0A++++gasUsed%0A++++gasPrice%0A++++swaps+%7B%0A++++++id%0A++++++amount0%0A++++++amount1%0A++++++token0+%7B%0A++++++++id%0A++++++++name%0A++++++++symbol%0A++++++%7D%0A++++++token1+%7B%0A++++++++decimals%0A++++++++id%0A++++++++name%0A++++++++symbol%0A++++++%7D%0A++++++transaction+%7B%0A++++++++swaps+%7B%0A++++++++++amount0%0A++++++++++amount1%0A++++++++++token0+%7B%0A++++++++++++id%0A++++++++++++name%0A++++++++++%7D%0A++++++++%7D%0A++++++%7D%0A++++%7D%0A++%7D%0A%7D' \
	//   -H 'sec-ch-ua: "Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"' \
	//   -H 'sec-ch-ua-mobile: ?0' \
	//   -H 'sec-ch-ua-platform: "macOS"' \
	//   -H 'sec-fetch-dest: empty' \
	//   -H 'sec-fetch-mode: cors' \
	//   -H 'sec-fetch-site: same-origin' \
	//   -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36' \
	//   --data-raw '{"query":"{\n  token(id: \"0x171b1daefac13a0a3524fcb6beddc7b31e58e079\") {\n    id\n    decimals\n    name\n    symbol\n    poolCount\n    totalSupply\n    totalValueLocked\n    totalValueLockedUSD\n    totalValueLockedUSDUntracked\n    volume\n    volumeUSD\n    txCount\n  }\n  transactions(\n    where: {swaps_: {token0: \"0x171b1daefac13a0a3524fcb6beddc7b31e58e079\"}}\n  ) {\n    blockNumber\n    timestamp\n    id\n    gasUsed\n    gasPrice\n    swaps {\n      id\n      amount0\n      amount1\n      token0 {\n        id\n        name\n        symbol\n      }\n      token1 {\n        decimals\n        id\n        name\n        symbol\n      }\n      transaction {\n        swaps {\n          amount0\n          amount1\n          token0 {\n            id\n            name\n          }\n        }\n      }\n    }\n  }\n}","variables":null,"extensions":{"headers":null}}' \
	//   --compressed

	body := strings.NewReader("{\"query\":\"{n  token(id: \"0x171b1daefac13a0a3524fcb6beddc7b31e58e079\") {n    idn    decimalsn    namen    symboln    poolCountn    totalSupplyn    totalValueLockedn    totalValueLockedUSDn    totalValueLockedUSDUntrackedn    volumen    volumeUSDn    txCountn  }n  transactions(n    where: {swaps_: {token0: \"0x171b1daefac13a0a3524fcb6beddc7b31e58e079\"}}n  ) {n    blockNumbern    timestampn    idn    gasUsedn    gasPricen    swaps {n      idn      amount0n      amount1n      token0 {n        idn        namen        symboln      }n      token1 {n        decimalsn        idn        namen        symboln      }n      transaction {n        swaps {n          amount0n          amount1n          token0 {n            idn            namen          }n        }n      }n    }n  }n}\",\"variables\":null,\"extensions\":{\"headers\":null}}")
	req, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3", body)
	if err != nil {
		fmt.Println("getEquityToken new request error:", err.Error())
		return 0
	}
	req.Header.Set("Authority", "api.thegraph.com")
	req.Header.Set("Accept", "application/json, multipart/mixed")
	req.Header.Set("Accept-Language", "en,zh-CN;q=0.9,zh;q=0.8")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://api.thegraph.com")
	req.Header.Set("Referer", "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3/graphql?query=%7B%0A++token%28id%3A%220x171b1daefac13a0a3524fcb6beddc7b31e58e079%22%29+%7B%0A++++id%0A++++decimals%0A++++name%0A++++symbol%0A++++poolCount%0A++++totalSupply%0A++++totalValueLocked%0A++++totalValueLockedUSD%0A++++totalValueLockedUSDUntracked%0A++++volume%0A++++volumeUSD%0A++++txCount%0A++%7D%0A++transactions%28%0A++++where%3A+%7Bswaps_%3A+%7Btoken0%3A+%220x171b1daefac13a0a3524fcb6beddc7b31e58e079%22%7D%7D%0A++%29+%7B%0A++++blockNumber%0A++++timestamp%0A++++id%0A++++gasUsed%0A++++gasPrice%0A++++swaps+%7B%0A++++++id%0A++++++amount0%0A++++++amount1%0A++++++token0+%7B%0A++++++++id%0A++++++++name%0A++++++++symbol%0A++++++%7D%0A++++++token1+%7B%0A++++++++decimals%0A++++++++id%0A++++++++name%0A++++++++symbol%0A++++++%7D%0A++++++transaction+%7B%0A++++++++swaps+%7B%0A++++++++++amount0%0A++++++++++amount1%0A++++++++++token0+%7B%0A++++++++++++id%0A++++++++++++name%0A++++++++++%7D%0A++++++++%7D%0A++++++%7D%0A++++%7D%0A++%7D%0A%7D")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"106\", \"Google Chrome\";v=\"106\", \"Not;A=Brand\";v=\"99\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("getEquityToken request error:", err.Error())
		return 0
	}
	defer resp.Body.Close()
	fmt.Println("getEquityToken resp:", resp)
	return 1
}
