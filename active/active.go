package active

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

// redigo连接数据库，后期改pool
func ConnectRedis() redis.Conn {
	client, err := redis.Dial("tcp", "localhost:6379", redis.DialPassword(""))
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
