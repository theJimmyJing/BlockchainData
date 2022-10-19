package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func Proxy(c *gin.Context) {
	if c.GetHeader("direct") != "lab" {
		return
	}
	var proxyUrl = new(url.URL)
	proxyUrl.Scheme = "https"
	proxyUrl.Host = "www.okx.com"
	//u.Path = "base" // 这边若是赋值了，做转发的时候，会带上path前缀，例： /hello -> /base/hello

	// proxyUrl.RawQuery = url.QueryEscape("token=" + "VjouhpQHa6wgWvtkPQeDZbQd") // 和如下方式等价
	//var query url.Values
	//query.Add("token", "VjouhpQHa6wgWvtkPQeDZbQd")
	//u.RawQuery = query.Encode()

	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

	//proxy := httputil.ReverseProxy{}
	//proxy.Director = func(req *http.Request) {
	//	fmt.Println(req.URL.String())
	//	req.URL.Scheme = "http"
	//	req.URL.Host = "172.16.60.161"
	//	rawQ := req.URL.Query()
	//	rawQ.Add("token", "VjouhpQHa6wgWvtkPQeDZbQd")
	//	req.URL.RawQuery = rawQ.Encode()
	//}

	// proxy.ErrorHandler // 可以添加错误回调
	// proxy.Transport // 若有需要可以自定义 http.Transport

	proxy.ServeHTTP(c.Writer, c.Request)

	c.Abort()
}
