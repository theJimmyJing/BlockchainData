package active

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// redigo连接主数据库
func ConnectRedis() redis.Conn {
	client, err := redis.Dial("tcp", "blockchaindata-ro.bllj2c.ng.0001.apne1.cache.amazonaws.com:6379", redis.DialPassword(""))
	if err != nil {
		panic(err)
	}
	return client
}

// redigo连接从数据库
func ConnectSlaveRedis() redis.Conn {
	client, err := redis.Dial("tcp", "blockchaindata-001.bllj2c.0001.apne1.cache.amazonaws.com:6379", redis.DialPassword(""))
	if err != nil {
		panic(err)
	}
	return client
}

// 保存event事件
func SaveActive(userId string) {
	redisClient := ConnectRedis()
	defer redisClient.Close()
	currentTime := time.Now()
	dayTime := currentTime.Format("20060102")
	key := "userActive" + dayTime

	_, err := redisClient.Do("SADD", key, userId)
	if err != nil {
		fmt.Println("SaveActive err: ", err)
	}
}

// 获取打点区间数据
func GetRangeCount(startOffset int, endOffset int) int {
	var rangeArr []string
	redisClient := ConnectSlaveRedis()
	defer redisClient.Close()
	currentTime := time.Now()

	for i := startOffset; i >= endOffset; i-- {
		dayTime := currentTime.AddDate(0, 0, i).Format("20060102")
		rangeArr = append(rangeArr, "userActive"+dayTime)
	}

	union, err := redis.Strings(redisClient.Do("SUNION", redis.Args{}.AddFlat(rangeArr)...))
	if err != nil {
		fmt.Println("GetDayRangeCount err: ", err)
		return 0
	}

	return len(union)
}
