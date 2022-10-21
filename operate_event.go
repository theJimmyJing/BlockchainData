// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"time"

// 	"github.com/go-redis/redis/v7"
// )

// func connectRedis() *redis.Client {
// 	var redisClient = redis.NewClient(&redis.Options{
// 		Addr:     "127.0.0.1:6379",
// 		Password: "",
// 		DB:       0,
// 	})
// 	return redisClient
// }

// // day格式20220601
// func DayActiveCount(redisClient *redis.Client, day string) int {
// 	eventKey := "active_" + day
// 	var dataCache []string
// 	getInfo, getinfoErr := redisClient.Get(eventKey).Result()
// 	if getinfoErr != nil {
// 		// fmt.Println("DayActiveCount no data", getinfoErr)
// 		return 0
// 	} else {
// 		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
// 		unmarsha1Err := json.Unmarshal([]byte(getInfo), &dataCache)
// 		if unmarsha1Err != nil {
// 			fmt.Println("DayActiveCount failed:", unmarsha1Err)
// 			return 0
// 		} else {
// 			return len(dataCache)
// 		}
// 	}
// }

// // 获取某周的开始和结束时间,0 为当天,-1昨天，1明天以此类推
// func GetDayActiveCount(redisClient *redis.Client, dayOffset int) int {
// 	now := time.Now()
// 	day := now.AddDate(0, 0, dayOffset).Format("2006-01-02")

// 	return DayActiveCount(redisClient, day)
// }

// // xx周的周活
// // 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
// func GetWeekActiveCount(redisClient *redis.Client, week int) int {
// 	weekActiveCount := 0

// 	now := time.Now()
// 	offset := int(time.Monday - now.Weekday())
// 	//周日做特殊判断 因为time.Monday = 0
// 	if offset > 0 {
// 		offset = -6
// 	}
// 	year, month, day := now.Date()
// 	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
// 	// startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02") + " 00:00:00"
// 	// endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02") + " 23:59:59"

// 	for i := 0; i < 7; i++ {
// 		day := thisWeek.AddDate(0, 0, offset+i+7*week).Format("2006-01-02")
// 		weekActiveCount += DayActiveCount(redisClient, day)
// 	}
// 	//
// 	return weekActiveCount
// }

// // 月活数 week为0本月,-1上月，1下月以此类推
// func GetMonthActiveCount(redisClient *redis.Client, monthOffset int) int {
// 	monthActiveCount := 0

// 	year, month, _ := time.Now().Date()
// 	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
// 	monLastDay := thisMonth.AddDate(0, monthOffset+1, -1)
// 	for i := 0; i < monLastDay.Day(); i++ {
// 		day := thisMonth.AddDate(0, monthOffset, i).Format("20060102")
// 		monthActiveCount += DayActiveCount(redisClient, day)
// 	}

// 	return monthActiveCount
// }

// // 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
// func WeekIntervalTime(week int) (startTime, endTime string) {
// 	now := time.Now()
// 	offset := int(time.Monday - now.Weekday())
// 	//周日做特殊判断 因为time.Monday = 0
// 	if offset > 0 {
// 		offset = -6
// 	}

// 	year, month, day := now.Date()
// 	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
// 	// startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02") + " 00:00:00"
// 	// endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02") + " 23:59:59"

// 	startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02") + " 00:00:00"
// 	endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02") + " 23:59:59"

// 	return startTime, endTime
// }

// // 获取某月的开始和结束时间mon为0本月,-1上月，1下月以此类推
// func MonthIntervalTime(mon int) (startTime, endTime string) {
// 	year, month, _ := time.Now().Date()
// 	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
// 	startTime = thisMonth.AddDate(0, mon, 0).Format("2006-01-02")
// 	endTime = thisMonth.AddDate(0, mon+1, -1).Format("2006-01-02")

// 	fmt.Println("month : ", startTime, endTime)
// 	return startTime, endTime
// }
