package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type BodyGo struct {
	Date     string `json:"date"`
	GameName string `json:"game_name"`
}

const webPort = ":8080"

func main() {
	go redisMain()
	http.HandleFunc("/fetchData", fetchData)
	http.HandleFunc("/fetchDataPastIncidents", fetchDataPastIncidents)
	http.ListenAndServe(webPort, nil)
}

func handleCrossOrigin(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodOptions {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		writer.WriteHeader(http.StatusOK)
		return
	}
}

func fetchDataPastIncidents(writer http.ResponseWriter, req *http.Request) {
	fmt.Println("You have a request for fetchDataPastIncidents: ", req)
	handleCrossOrigin(writer, req)
	fetchDataPastIncidentsHelper(writer, req)
}

func fetchDataPastIncidentsHelper(writer http.ResponseWriter, req *http.Request) {
	time.Sleep(5 * time.Second)
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	cntxt := context.Background()
	val2temp, _ := RedisClient.Do(cntxt, "LRANGE", "data", "0", "-1").Text()
	val2, _ := RedisClient.Do(cntxt, "LRANGE", "data", "0", "-1").StringSlice()
	fmt.Println("The val is: ", val2temp)
	val := make([]IncidentData, len(val2))//[]map[string]interface{}
	for index, sampleData := range val2{
		//var incident map[string]interface{}
		err := json.Unmarshal([]byte(sampleData), &val[index])
		if err != nil{
			fmt.Printf("Error Parsing JSON: %v\n", sampleData)
			continue
		}
		//val = append(val, incident)
		val[index] = val[index]
	}
	fmt.Println("Then length of the returned data is: ", len(val2))
	fmt.Println("The request for the response is: ", val2, "\n", val)
	fmt.Printf("\n\n")
	FinalJson, errJson := json.Marshal(val)
	if errJson == nil{
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintf(writer, "%s\n", FinalJson)
	}else{
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func fetchData(writer http.ResponseWriter, request *http.Request) {
	handleCrossOrigin(writer, request)
	fetchDataHelper(writer, request)
}

func fetchDataHelper(writer http.ResponseWriter, request *http.Request) {
	time.Sleep(5 * time.Second)
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	BodyPlaceHolder := make([]byte, request.ContentLength)
	request.Body.Read(BodyPlaceHolder)
	Body := &BodyGo{}
	err := json.Unmarshal(BodyPlaceHolder, Body)
	if err == nil {
		Resp, _ := main1(Body.Date, Body.GameName).encodeToJSON()
		fmt.Println("The Response is: ", Resp)
		fmt.Fprintf(writer, "%s\n", string(Resp))
	} else {
		fmt.Println("Encountered error while unmarshllaling the request body: ", err)
	}
}
