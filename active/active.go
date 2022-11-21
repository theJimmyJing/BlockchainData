package active

import (
	"fcc/config"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
)

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
	key := "userActive" + dayTime

	_, errCluster := redisClient.Do("CLUSTER KEYSLOT", "userActive"+dayTime)
	if errCluster != nil {
		fmt.Println("SaveActive err: ", errCluster)
	}

	_, err := redisClient.Do("SADD", key, userId)
	if err != nil {
		fmt.Println("SaveActive err: ", err)
	}
}

// 获取打点区间数据
func GetRangeCount(startOffset int, endOffset int) int {
	var rangeArr []string
	redisClient := ConnectRedis()
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
