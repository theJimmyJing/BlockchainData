package active

import (
	"fcc/config"
	"fmt"
	"strconv"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

const ActiveUserKey = "activeUser"

// redigo连接主数据库
func ConnectRedis() redigo.Conn {
	client, err := redigo.Dial("tcp", config.BlockchainDataConfig.Redis.DBAddress[0], redigo.DialPassword(config.BlockchainDataConfig.Redis.DBPassWord))
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
	key := userId + dayTime
	score, _ := strconv.Atoi(dayTime)

	_, err := redisClient.Do("ZADD", ActiveUserKey, score, key)
	if err != nil {
		fmt.Println("SaveActive err: ", err)
	}
}

// 获取打点区间数据
func GetRangeCount(startOffset int, endOffset int) int {
	var rangeArr []int
	redisClient := ConnectRedis()
	defer redisClient.Close()
	currentTime := time.Now()

	for i := startOffset; i >= endOffset; i-- {
		dayTime, _ := strconv.Atoi(currentTime.AddDate(0, 0, i).Format("20060102"))
		rangeArr = append(rangeArr, dayTime-1)
		rangeArr = append(rangeArr, dayTime)
	}
	un, err := redisClient.Do("ZCOUNT", ActiveUserKey, rangeArr[0], rangeArr[1])
	if err != nil {
		fmt.Println("GetDayRangeCount err: ", err)
		return 0
	}

	switch un.(type) {
	case int:
		return int(un.(uint))
		break
	default:
		return 0
		break
	}
	return 0
}
