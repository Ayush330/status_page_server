package main

import (
	"context"
	"fmt"
	"time"

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
	scheduleAndRunEvent("poker", Redisclient)
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
	ListOfIncidents := getListOfIncidents()
	fmt.Println("List of Incidents is: ", ListOfIncidents)
	cntxt := context.Background()
	fmt.Println("The Client Data is: ", RedisClient)
	// for i := 0; i < len(ListOfIncidents); i++ {
	// 	Data, _ := json.Marshal(ListOfIncidents[i])
	// 	client.Do(cntxt, "LPUSH", "data", Data).Result()
	// }

	val2, _ := RedisClient.Do(cntxt, "LRANGE", "data", "0", "-1").StringSlice()
	fmt.Println("foo ", val2)
}

func getListOfIncidents() []IncidentData {
	return []IncidentData{
		{"2024-01-31", 10, "indian_rummy"},
		{"2024-01-31", 0, "indian_rummy"},
		{"2024-01-31", 0, "indian_rummy"},
		{"2024-01-31", 0, "social_poker"},
	}
}
