package main

import (
	"context"
	"fmt"
	"time"
	"strconv"
	"encoding/json" 
	"github.com/redis/go-redis/v9"
)

type IncidentData struct {
	Date           string `json:"date"`
	NumOfIncidents int    `json:"number_of_incidents"`
	GameName       string `json:"game_name"`
}

const (
	SOCIAL_POKER = "social_poker"
	INDIAN_RUMMY = "indian_rummy"
	TEENPATTI    = "teenpatti"
	PLAYRUMMY    = "playrummy"
	CASH_POKER   = "pokerpro"
)

func redisMain() {
	var GamesList = []string{
		SOCIAL_POKER,
		INDIAN_RUMMY,
		TEENPATTI,
		PLAYRUMMY,
		CASH_POKER,
	}
	Redisclient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	for _, gameName := range GamesList {
		scheduleAndRunEvent(gameName, Redisclient)
	}
	//scheduleAndRunEvent("poker", Redisclient)
	select {}
}

func scheduleAndRunEvent(GameName string, RedisClient *redis.Client) {
	scheduleNextRun := func() {
		now := time.Now()
		next := now.AddDate(0, 0, 1)
		nextRunTime := time.Date(next.Year(), next.Month(), next.Day(), 12, 0, 0, 0, next.Location())
		durationUntilNextRun, _ := time.ParseDuration("1m")
		fmt.Println("The next run time is: ", nextRunTime)
		time.AfterFunc(durationUntilNextRun, func() {
			fmt.Println("Current Time is: ", time.Now())
			redisMainHelperForever(GameName, RedisClient)
			scheduleAndRunEvent(GameName, RedisClient)
		})
	}
	scheduleNextRun()
}

func redisMainHelperForever(GameName string, RedisClient *redis.Client) { 
	cntxt := context.Background()
	fmt.Println("Here: ...........")
	var Res IncidentData
	if GameName == "pokerpro" || GameName == "playrummy"{
		Res = test_redis_data_formation_cash(GameName)
	}else{
		Res = test_redis_data_formation(GameName)
	}
	fmt.Println("The data is: ", Res)
	if Res.NumOfIncidents > 0{
		ResJson, errJson := json.Marshal(Res)
		if errJson != nil{
			fmt.Println("Error Encountered While Marshalling Data: ", Res, " With Error: ", errJson)
		}else{
			// _, err := RedisClient.Do(cntxt, "LPUSH", "data", ResJson).Text()
			SecondsIn5Days := 5 * 24 * 60 * 60
			ExpiryTime := time.Now().Unix() + int64(SecondsIn5Days)
			err := RedisClient.Do(cntxt, "ZADD", "data", strconv.FormatInt(ExpiryTime, 10), ResJson)
			//err := RedisClient.Do(cntxt, "ZADD", "data", "50", ResJson)
			//err := RedisClient.Do(cntxt, "ZADD", "data", "60", "ayush")
			if err != nil{
				fmt.Println("Failed to push data into Redis because of: ", err, "for the data: ", ResJson)
			}
		}
	}
}

func getListOfIncidents() []IncidentData {
	return []IncidentData{
		{"2024-01-31", 10, "indian_rummy"},
		{"2024-01-31", 0, "indian_rummy"},
		{"2024-01-31", 0, "indian_rummy"},
		{"2024-01-31", 0, "social_poker"},
	}
}

func setHourlyDataUpdationRedis(){
	// loop, that will run at every 1 hrs
	// define the keys with ttl of exactly 1 hrs->
	// 1. indian_rummy
	// 2. social_poker
	// 3. teenpatti
	// 4. playrummy
	// 5. pokerpro
	// GameNamesList := []string{
	// 	"indian_rummy",
	// 	"social_poker",
	// 	"teenpatti",
	// 	"playrummy",
	// 	"pokerpro",
	// }
	for _, gameName := range GameNamesList{
		fmt.Println("GameNames Are: ", gameName)
		setHourlyDataUpdationRedisHelper(gameName)
	}
	select {}
}

func setHourlyDataUpdationRedisHelper(GameName string){
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	nextEvent := func(){
		durationUntilNextTime, _ := time.ParseDuration("1m")
		time.AfterFunc(durationUntilNextTime, func(){
			fmt.Println("In Here")
			setHourlyDataUpdationRedisHelper2(GameName, RedisClient)
			setHourlyDataUpdationRedisHelper(GameName)
		})
	}
	nextEvent()
}

func setHourlyDataUpdationRedisHelper2(GameName string, RedisClient *redis.Client){
	cntxt := context.Background()
	//Date := time.Now().Format(time.DateOnly)
	Date := "2024-01-02"
	Resp1, err := main1(Date, GameName)
	if err == nil{
		Resp, _ := Resp1.encodeToJSON()
		SecondsIn1Min := 60
		err := RedisClient.Do(cntxt, "SET", GameName, Resp, "EX", SecondsIn1Min)
		fmt.Println("The error is: ", err)
	}
}
