package active

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

func SaveActive(redisClient *redis.Client, userId string) {
	currentTime := time.Now()
	dayTime := currentTime.Format("20060102")
	// dayTime := currentTime.AddDate(0, 0, -11).Format("20060102")
	key := "userActive" + dayTime
	fmt.Println("saveactive", key)
	redisClient.SAdd(key, userId)
}

func GetDayActiveCount(redisClient *redis.Client, searchDay string) int {
	key := "userActive" + searchDay
	count, err := redisClient.SCard(key).Result()
	fmt.Println("getDayActiveCount", key, count)
	if err != nil {
		fmt.Println("getDayActiveCount err: ", err)
		return 0
	}

	res := int(count)
	return res
}

func GetDayRangeCount(redisClient *redis.Client, startOffset int, endOffset int) int {
	var rangeArr []string
	var setRes []string
	currentTime := time.Now()

	for i := startOffset; i >= endOffset; i-- {
		dayTime := currentTime.AddDate(0, 0, i).Format("20060102")
		rangeArr = append(rangeArr, "userActive"+dayTime)
	}

	for i := startOffset - endOffset; i > 0; i-- {
		unionSet, err := redisClient.SUnion(rangeArr[i]).Result()
		if err != nil {
			fmt.Println("GetDayRangeCount err: ", rangeArr[i], err)
			return 0
		}
		setRes = unionSet
	}

	fmt.Println("SCard", setRes)

	// res := len(setRes)
	// return res
	return 0
}
