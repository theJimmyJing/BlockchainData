package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

// 运营数据 - 事件埋点相关接口
func connectRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "blockchaindata-ro.bllj2c.ng.0001.apne1.cache.amazonaws.com:6379",
		Password: "",
		DB:       0,
	})
	return redisClient
}

// 连接从服务
func ConnectSlaveRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "blockchaindata-001.bllj2c.0001.apne1.cache.amazonaws.com:6379",
		Password: "",
		DB:       0,
	})
	return redisClient
}

func testKeys(key string) {
	redisClient := connectRedis()
	// 测试批量查询匹配的keys
	keys, _, _ := redisClient.Scan(0, key+"_*", 0).Result()

	fmt.Println(len(keys), keys)
}

func saveEventData(redisClient *redis.Client, newEvent EventData) {
	currentTime := time.Now()
	dayTime := currentTime.Format("20060102")
	eventKey := newEvent.Event + "_" + dayTime // 埋点事件存储格式： KEY（event_日期），data

	var dataCache []string
	getInfo, getinfoErr := redisClient.Get(eventKey).Result()
	if getinfoErr != nil {
		fmt.Println("no data", getinfoErr)
	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &dataCache)
		if unmarsha1Err != nil {
			fmt.Println("反序列化失败:", unmarsha1Err)
		} else {
			fmt.Println(dataCache)
		}
	}

	// 结构体转json
	jsonBytes, err := json.Marshal(newEvent)
	if err != nil {
		fmt.Println("struct to bytes err : ", err)
		return
	}
	// 添加到数组
	dataCache = append(dataCache, string(jsonBytes))

	infoByte, infoError := json.Marshal(dataCache) // 数组转bytes
	if infoError == nil {
		inforString := string(infoByte)                                    //转换成字符串
		infoErrorStatus := redisClient.Set(eventKey, inforString, 0).Err() //设置过期时间- 不过期
		if infoErrorStatus != nil {
			fmt.Println("save failed：", infoErrorStatus)
		} else {
			fmt.Println("save success", eventKey, newEvent.Date)
		}
	}
}

// day格式20220601
func DayActiveCount(redisClient *redis.Client, day string) int {
	eventKey := "active_" + day
	var dataCache []string
	getInfo, getinfoErr := redisClient.Get(eventKey).Result()
	if getinfoErr != nil {
		return 0
	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &dataCache)
		if unmarsha1Err != nil {
			return 0
		} else {
			return len(dataCache)
		}
	}
}

// 获取某周的开始和结束时间,0 为当天,-1昨天，1明天以此类推
func GetDayActiveCount(redisClient *redis.Client, dayOffset int) int {
	now := time.Now()
	day := now.AddDate(0, 0, dayOffset).Format("20060102")

	return DayActiveCount(redisClient, day)
}

// xx周的周活
// 获取某周的开始和结束时间,week为0本周,-1上周，1下周以此类推
func GetWeekActiveCount(redisClient *redis.Client, week int) int {
	weekActiveCount := 0

	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}
	year, month, day := now.Date()
	thisWeek := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	// startTime = thisWeek.AddDate(0, 0, offset+7*week).Format("2006-01-02") + " 00:00:00"
	// endTime = thisWeek.AddDate(0, 0, offset+6+7*week).Format("2006-01-02") + " 23:59:59"

	for i := 0; i < 7; i++ {
		day := thisWeek.AddDate(0, 0, offset+i+7*week).Format("20060102")
		weekActiveCount += DayActiveCount(redisClient, day)
	}

	return weekActiveCount
}

// 月活数 week为0本月,-1上月，1下月以此类推
func GetMonthActiveCount(redisClient *redis.Client, monthOffset int) int {
	monthActiveCount := 0

	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	monLastDay := thisMonth.AddDate(0, monthOffset+1, -1)
	for i := 0; i < monLastDay.Day(); i++ {
		day := thisMonth.AddDate(0, monthOffset, i).Format("20060102")
		monthActiveCount += DayActiveCount(redisClient, day)
	}

	return monthActiveCount
}

// 获取指定日期的注册用户数量，day格式register_20220601
func GetDayRegisteredCount(redisClient *redis.Client, eventDayKey string) int {
	var dataCache []string
	getInfo, getinfoErr := redisClient.Get(eventDayKey).Result()
	if getinfoErr != nil {
		return 0
	} else {
		//获取到json字符串,反序列化,原来是二维数组的,反序列化的时候也要用二维数组接收
		unmarsha1Err := json.Unmarshal([]byte(getInfo), &dataCache)
		if unmarsha1Err != nil {
			return 0
		} else {
			return len(dataCache)
		}
	}
}

// 获取所有注册用户数量
func GetAllRegisteredCount(redisClient *redis.Client) int {
	count := 0
	registerKeys, _, _ := redisClient.Scan(0, "register_*", 0).Result()
	for _, key := range registerKeys {
		count += GetDayRegisteredCount(redisClient, key)
	}
	return count
}
