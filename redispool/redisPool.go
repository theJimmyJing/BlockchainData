// package main

// import (
// 	"fmt"

// 	"github.com/gomodule/redigo/redis"
// )

// var pool *redis.Pool

// func init() {
// 	pool = &redis.Pool{
// 		MaxIdle:     16,
// 		MaxActive:   0,
// 		IdleTimeout: 300,
// 		Dial: func() (redis.Conn, error) {
// 			return redis.Dial("tcp", "blockchaindata-ro.bllj2c.ng.0001.apne1.cache.amazonaws.com:6379")
// 		},
// 	}
// }

// func main() {
// 	c := pool.Get() // 获取连接
// 	defer c.Close()

// 	_, err := c.Do("PING") // 判断redis是否可用

// 	// _, err := c.Do("Set", "abc", 100)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }
