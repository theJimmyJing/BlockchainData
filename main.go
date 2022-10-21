package main

import (
	// "crypto/tls"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/kirinlabs/HttpRequest"
)

var redisClient *redis.Client

func main() {
	// uniswapFCCToken()
	// operateAllData()
	// startServer()
	// now := time.Now().UnixNano()
	// time.Sleep(time.Second)
	// now2 := time.Now().UnixNano()
	// fmt.Println(now2 - now)

	// testEvent()
	startGin()
	// startServerV3()
	// startServerV4()

	// conneRedis()
}

// 测试埋点存储和查询
func TestEvent() {
	var redisClient = connectRedis()

	eventData := EventData{}
	eventData.Event = "active"
	eventData.Date = strconv.FormatInt((time.Now().UnixNano() / 1e6), 10) // 时间戳

	saveEventData(eventData)

	count := GetDayActiveCount(redisClient, 0)
	fmt.Println("count ", count)
}

func onUserEvent() {
	router := gin.Default()
	router.POST("/api/v5/operate/event", func(c *gin.Context) {
		var newEvent EventData
		err := c.ShouldBindJSON(&newEvent)

		if err != nil {
			fmt.Println("event : ", err)
			c.JSON(500, gin.H{
				"Code": 500,
				"Msg":  err.Error(),
			})
			return
		}

		// 存储event
		saveEventData(newEvent)

		c.JSON(http.StatusOK, "OK")
	})

	router.Run(":8081")
}

func saveEventData(newEvent EventData) {
	var redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	currentTime := time.Now()
	dayTime := currentTime.Format("20060102")
	eventKey := newEvent.Event + "_" + dayTime // 埋点事件存储格式： KEY（event_日期），data

	var dataCache []string
	getInfo, getinfoErr := redisClient.Get(eventKey).Result()
	if getinfoErr != nil {
		fmt.Println("没有获取到数据", getinfoErr)
	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &dataCache)
		if unmarsha1Err != nil {
			fmt.Println("反序列化失败:", unmarsha1Err)
		} else {
			fmt.Println(dataCache)
		}
	}

	// 结构体转json
	jsonBytes, err := json.Marshal(newEvent)
	if err != nil {
		fmt.Println("struct to bytes err : ", err)
		return
	}
	// 添加到数组
	dataCache = append(dataCache, string(jsonBytes))

	infoByte, infoError := json.Marshal(dataCache) // 数组转bytes
	if infoError == nil {
		inforString := string(infoByte)                                    //转换成字符串
		infoErrorStatus := redisClient.Set(eventKey, inforString, 0).Err() //设置过期时间- 不过期
		if infoErrorStatus != nil {
			fmt.Println("save failed：", infoErrorStatus)
		} else {
			fmt.Println("save success", eventKey, newEvent.Date)
		}
	}
}

func operateAllData() {
	router := gin.Default()
	router.GET("/api/v5/operate/all", func(c *gin.Context) {
		var data = OperateData{}
		// FCC Token
		fccToken := uniswapFCCToken()
		if fccToken != (UniswapToken{}) {
			fmt.Println("fccToken: ", fccToken)
			// TODO 转换返回值
		}

		// UserBigData
		data.User = getUserBigData()

		jsonBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("struct to bytes err : ", err)
		}

		fmt.Println("resp : ", jsonBytes)
		fmt.Println("resp 2: ", string(jsonBytes))
		c.JSON(http.StatusOK, string(jsonBytes))

	})
	router.Run(":8081")
}

/*
* {"derivedETH":"0","feesUSD":"108.3979488349600340149325490092877","name":"Freechat Coin","poolCount":"0",
* "symbol":"FCC","totalSupply":"28368","totalValueLocked":"199940134.57160359567214293","totalValueLockedUSD":"0",
* "totalValueLockedUSDUntracked":"0","txCount":"108","untrackedVolumeUSD":"18066.32480582667233582209150154796",
* "volume":"456421.018775415938823971","volumeUSD":"36132.6496116533446716441830030959","decimals":"18"}
 */
func uniswapFCCToken() UniswapToken {
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

	var resp = UniswapResp{}
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with", err)
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())

	fmt.Println("outStr : ", outStr)
	fmt.Println("errStr ：", errStr)

	// json转结构体
	uerr := json.Unmarshal(stdout.Bytes(), &resp)
	if uerr != nil {
		log.Fatalf("stdout-> UniswapToken err", uerr)
	}
	fmt.Println("fccToken : ", resp)

	// 结构体转json
	jsonBytes, err := json.Marshal(resp.Data.Token)
	if err != nil {
		fmt.Println("struct to bytes err : ", err)
	}
	fmt.Println("UniswapToken: ", string(jsonBytes))

	return UniswapToken(resp.Data.Token)
	// TODO read data from Stdout and
}

func getUserBigData() UserBigData {

	var redisClient = connectRedis()
	day1 := GetDayActiveCount(redisClient, -1)
	day2 := GetDayActiveCount(redisClient, -1)
	week1 := GetWeekActiveCount(redisClient, -1)
	week2 := GetWeekActiveCount(redisClient, -2)
	month1 := GetMonthActiveCount(redisClient, -1)
	month2 := GetMonthActiveCount(redisClient, -2)

	fmt.Println("count: ", day1, day2, week1, week2, month1, month2)

	var data = UserBigData{}

	data.Total = 10000
	data.DayIncrease = 200
	data.DayActive = day1
	data.DayActiveIncrease24H = day1 - day2
	data.WeekActive = week1
	data.WeekActiveIncrease24H = week1 - week2
	data.MonthActive = month1
	data.MonthActiveIncrease24H = month1 - month2

	return data
}

func connectRedis() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})
	}
	return redisClient
}

// day格式20220601
func DayActiveCount(redisClient *redis.Client, day string) int {
	eventKey := "active_" + day
	var dataCache []string
	getInfo, getinfoErr := redisClient.Get(eventKey).Result()
	if getinfoErr != nil {
		fmt.Println("getDayActive no data", eventKey, getinfoErr)
		return 0
	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &dataCache)
		if unmarsha1Err != nil {
			fmt.Println("getDayActive failed:", eventKey, unmarsha1Err)
			return 0
		} else {
			return len(dataCache)
		}
	}
}

// 获取某周的开始和结束时间,0 为当天,-1昨天，1明天以此类推
func GetDayActiveCount(redisClient *redis.Client, dayOffset int) int {
	now := time.Now()
	day := now.AddDate(0, 0, dayOffset).Format("20060102")

	return DayActiveCount(redisClient, day)
}

// xx周的周活
// 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
func GetWeekActiveCount(redisClient *redis.Client, week int) int {
	weekActiveCount := 0

	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}
	year, month, day := now.Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	// startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02") + " 00:00:00"
	// endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02") + " 23:59:59"

	for i := 0; i < 7; i++ {
		day := thisWeek.AddDate(0, 0, offset+i+7*week).Format("20060102")
		weekActiveCount += DayActiveCount(redisClient, day)
	}

	return weekActiveCount
}

// 月活数 week为0本月,-1上月，1下月以此类推
func GetMonthActiveCount(redisClient *redis.Client, monthOffset int) int {
	monthActiveCount := 0

	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	monLastDay := thisMonth.AddDate(0, monthOffset+1, -1)
	for i := 0; i < monLastDay.Day(); i++ {
		day := thisMonth.AddDate(0, monthOffset, i).Format("20060102")
		monthActiveCount += DayActiveCount(redisClient, day)
	}

	return monthActiveCount
}

func getInstIdTickerInfo(params string) *http.Response {

	req, err := http.NewRequest("GET", "https://www.okex.com/api/v5/market/index-tickers?"+params, nil)
	if err != nil {
		// handle err
		log.Printf("%+v", err)

	}
	req.Header.Set("Authority", "www.okex.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "locale=zh-CN")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		log.Default().Printf("%+v", err)

	}
	r, err := ParseResponse(resp)
	if err != nil {
		log.Default().Printf("%+v", err)
		return nil
	}
	log.Default().Printf("resp %+v", r)
	od.InstIdMap[params] = r
	od.InstTimeMap[params] = uint64(time.Now().UnixNano())
	return resp

}

func getExchangeRate() *http.Response {

	req, err := http.NewRequest("GET", "https://www.okex.com/api/v5/market/exchange-rate", nil)
	if err != nil {
		// handle err
		log.Printf("%+v", err)
	}
	req.Header.Set("Authority", "www.okex.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "locale=zh-CN")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"105\", \"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"105\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"macOS\"")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")

	// proxy, _ := url.Parse("http://127.0.0.1:59726")
	// tr := &http.Transport{
	// 	Proxy:           http.ProxyURL(proxy),
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }

	// client := &http.Client{
	// 	Transport: tr,
	// 	Timeout:   time.Second * 5, //超时时间
	// }

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		log.Default().Printf("%+v", err)

	}
	r, err := ParseResponse(resp)
	if err != nil {
		log.Default().Printf("%+v", err)
		return nil
	}
	log.Default().Printf("resp %+v", r)
	od.Rate = r
	od.RateTime = uint64(time.Now().UnixNano())
	// defer resp.Body.Close()
	return resp

}

func ParseResponse(response *http.Response) (map[string]interface{}, error) {
	var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}

	return result, err
}

func startServerV3() {
	router := gin.Default()
	router.GET("/api/v5/market/index-tickers", func(c *gin.Context) {
		// var i IndexTickers
		// c.ShouldBind(&i)
		rq := c.Request.URL.RawQuery

		req, err := http.NewRequest("GET", "https://www.okx.com/api/v5/market/index-tickers?"+rq, nil)
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}

		// q := req.URL.Query()
		// q.Add("api_key", "key_from_environment_or_flag")
		// q.Add("another_thing", "foo & bar")
		// req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())
		// Output:
		// http://api.themoviedb.org/3/tv/popular?another_thing=foo+%26+bar&api_key=key_from_environment_or_flag
		var resp *http.Response

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Print(err)
		}
		defer resp.Body.Close()
	})
}

type OkData struct {
	InstIdMap map[string]map[string]interface{}
	// InstTime    uint64
	InstTimeMap map[string]uint64
	Rate        map[string]interface{}
	RateTime    uint64
}

var od = OkData{
	InstIdMap: make(map[string]map[string]interface{}),
	Rate:      make(map[string]interface{}),
	// InstTime:  0,
	InstTimeMap: make(map[string]uint64),
	RateTime:    0,
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}

func startGin() {

	router := gin.Default()
	router.Use(Cors()) //开启中间件 允许使用跨域请求
	router.GET("/api/v5/market/index-tickers", func(c *gin.Context) {
		now := time.Now().UnixNano()
		rq := c.Request.URL.RawQuery
		diff := now - int64(od.InstTimeMap[rq])
		log.Default().Println("now = " + strconv.Itoa(int(now)) + " insttime = " + strconv.Itoa(int(od.InstTimeMap[rq])))

		if diff < int64(2*time.Second) {
			log.Default().Println(od.InstIdMap[rq])
			c.JSON(http.StatusOK, od.InstIdMap[rq])
			return
		}
		_ = getInstIdTickerInfo(rq)
		log.Default().Println(od.InstIdMap[rq])
		c.JSON(http.StatusOK, od.InstIdMap[rq])
		return
	})
	router.GET("/api/v5/market/exchange-rate", func(c *gin.Context) {
		now := time.Now().UnixNano()
		diff := now - int64(od.RateTime)
		// rq := c.Request.URL.RawQuery
		if diff < int64(2*time.Second) {
			c.JSON(http.StatusOK, od.Rate)
			return
		}
		_ = getExchangeRate()
		c.JSON(http.StatusOK, od.Rate)
		return
	})

	router.GET("/api/v5/market/tokenlist", func(c *gin.Context) {
		// allResults := make([]AllResult, 0)
		// 打开json文件
		jsonFile, err := os.Open("result.json")

		// 最好要处理以下错误
		if err != nil {
			fmt.Println(err)
		}

		// 要记得关闭
		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)
		// fmt.Println(string(byteValue))
		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)

		c.JSON(http.StatusOK, result)
		return
	})

	// 埋点事件
	router.POST("/api/v5/operate/event", func(c *gin.Context) {
		var data EventData
		err := c.ShouldBindJSON(&data)

		if err != nil {
			fmt.Println("event : ", err)
			c.JSON(500, gin.H{
				"Code": 500,
				"Msg":  err.Error(),
			})
			return
		}

		// TODO 事件埋入redis

		fmt.Println("event : ", data)

		c.JSON(http.StatusOK, "OK")
	})

	// 运营数据
	router.GET("/api/v5/operate/all", func(c *gin.Context) {
		var data = OperateData{}
		// FCC Token
		fccToken := uniswapFCCToken()
		if fccToken != (UniswapToken{}) {
			fmt.Println("fccToken: ", fccToken)

			data.Freechat.TotalEarn = "12312321"
			data.Freechat.DayEarn = "2131"
			data.Freechat.DayEarnIncrease = "+15.4%"
			data.Freechat.WeekEarn = "10232"
			data.Freechat.WeekEarnIncrease = "14.23%"
			data.Freechat.MonthEarn = "31232"
			data.Freechat.MonthEarnIncrease = "12.4%"

			data.Freechat.NowPrice = "0.0635"
			data.Freechat.MarketValue = "100000000"
			data.Freechat.MarketValueIncrease = "+4.4%"
			data.Freechat.DayVolume = "10000"
			data.Freechat.DayVolumeIncrease = "4.23%"
			data.Freechat.FccUser = "1000000"
			data.Freechat.FccUserIncrease = "1.4%"

			data.Freechat.TotalProfit = "5231212"
			data.Freechat.WaitProfit = "131232"
			data.Freechat.PerFccProfit = "5.23%"
			data.Freechat.PledgeProfit = "5.3"
			data.Freechat.PledgeRate = "3.2%"

			// TODO 转换返回值
		}

		// UserBigData - 读取用户日活等数据
		data.User = getUserBigData()

		jsonBytes, err := json.Marshal(data)
		if err != nil {
			fmt.Println("struct to bytes err : ", err)
		}

		fmt.Println("resp : ", jsonBytes)
		fmt.Println("resp 2: ", string(jsonBytes))
		c.JSON(http.StatusOK, string(jsonBytes))

	})
	router.Run(":8080")
}

var tokenList string

func syncData() {
	fmt.Println("start sync token list ")
	origin := analysisJson()
	r := getAbi(&origin)
	fmt.Printf("first result ---------\n %+v", *r)
	a := combileDetails(r)
	fmt.Printf("second result ---------\n %+v", *r)
	file, _ := json.MarshalIndent(a, "", " ")
	tokenList = string(file)
	_ = ioutil.WriteFile("result.json", file, 0644)

}

type IndexTickers struct {
	InstId   string `form:"instId" json:"instId"`
	QuoteCcy string `form:"quoteCcy" json:"quoteCcy"`
}

func startServer() {
	engine := gin.New()
	vi := engine.Group("/api")
	vi.Any("/v5/market/index-tickers", WithHeader)
	// GET /api/v5/market/exchange-rate
	vi.Any("/v5/market/exchange-rate", WithHeader)

	err := engine.Run(":8341")
	if err != nil {
		fmt.Println(err)
	}
}

const Host = "www.okx.com"

var simpleHostProxy = httputil.ReverseProxy{
	Director: func(req *http.Request) {
		req.URL.Scheme = "https"
		req.URL.Host = Host
		req.Host = Host
	},
}

func WithHeader(ctx *gin.Context) {

	// ctx.Request.Header.Add("requester-uid", "id")
	simpleHostProxy.ServeHTTP(ctx.Writer, ctx.Request)
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

// uniswap token struct
type UniswapResp struct {
	Data UniswapData `json:"data"`
}
type UniswapData struct {
	Token Token `json:"token"`
}

type Token struct {
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
