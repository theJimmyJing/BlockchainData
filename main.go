package main

import (
	// "crypto/tls"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kirinlabs/HttpRequest"
)

func main() {
	// syncData()
	// startServer()
	// now := time.Now().UnixNano()
	// time.Sleep(time.Second)
	// now2 := time.Now().UnixNano()
	// fmt.Println(now2 - now)

	startGin()
	// startServerV3()
	// startServerV4()
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
	od.InstTime = uint64(time.Now().UnixNano())
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
	InstTime  uint64
	Rate      map[string]interface{}
	RateTime  uint64
}

var od = OkData{
	InstIdMap: make(map[string]map[string]interface{}),
	Rate:      make(map[string]interface{}),
	InstTime:  0,
	RateTime:  0,
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
		diff := now - int64(od.InstTime)
		log.Default().Println("now = " + strconv.Itoa(int(now)) + " insttime = " + strconv.Itoa(int(od.InstTime)))
		rq := c.Request.URL.RawQuery
		if diff < int64(2*time.Second) {
			c.JSON(http.StatusOK, od.InstIdMap[rq])
			return
		}
		_ = getInstIdTickerInfo(rq)
		// reader := response.Body
		// contentLength := response.ContentLength
		// contentType := response.Header.Get("Content-Type")
		// extraHeaders := map[string]string{
		// 	//"Content-Disposition": `attachment; filename="gopher.png"`,
		// }
		// c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
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
		// reader := response.Body
		// contentLength := response.ContentLength
		// contentType := response.Header.Get("Content-Type")
		// extraHeaders := map[string]string{
		// 	//"Content-Disposition": `attachment; filename="gopher.png"`,
		// }
		// c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
		c.JSON(http.StatusOK, od.Rate)
		return
	})
	router.Run(":8080")
}

func syncData() {
	fmt.Println("hello")
	origin := analysisJson()
	r := getAbi(&origin)
	fmt.Printf("first result ---------\n %+v", *r)
	a := combileDetails(r)
	fmt.Printf("second result ---------\n %+v", *r)
	file, _ := json.MarshalIndent(a, "", " ")

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
