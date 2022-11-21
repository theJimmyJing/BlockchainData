package active

import (
	"context"
	"fcc/config"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"
)

const ActiveUserKey = "activeUser"

var ctx = context.Background()

func ConnectRedis() redis.UniversalClient {
	var rdb redis.UniversalClient
	if config.BlockchainDataConfig.Redis.EnableCluster {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    config.BlockchainDataConfig.Redis.DBAddress,
			PoolSize: 50,
		})
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:     config.BlockchainDataConfig.Redis.DBAddress[0],
			Password: config.BlockchainDataConfig.Redis.DBPassWord, // no password set
			DB:       0,                                            // use default DB
			PoolSize: 100,                                          // 连接池大小
		})
	}

	return rdb
}

// 保存event事件
func SaveActive(userId string) {
	redisClient := ConnectRedis()
	defer redisClient.Close()
	currentTime := time.Now()
	dayTime := currentTime.Format("20060102")
	key := userId + dayTime
	score, _ := strconv.ParseFloat(dayTime, 64)

	_, err := redisClient.ZAdd(ctx, ActiveUserKey, redis.Z{Score: score, Member: key}).Result()
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
	min := strconv.Itoa(rangeArr[0])
	max := strconv.Itoa(rangeArr[1])
	un, err := redisClient.ZCount(ctx, ActiveUserKey, min, max).Result()
	// un, err := redisClient.Do("ZCOUNT", ActiveUserKey, rangeArr[0], rangeArr[1])
	if err != nil {
		fmt.Println("GetDayRangeCount err: ", err)
		return 0
	}

	return int(un)
}
