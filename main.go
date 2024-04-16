package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	_"time"
	"math"
	"github.com/redis/go-redis/v9"
)

type BodyGo struct {
	Date     string `json:"date"`
	GameName string `json:"game_name"`
}

const webPort = ":8080"
var GameNamesList = []string{                                                                                                                 
                "indian_rummy",                                                                                                                   
                "social_poker",                                                                                                                   
                "teenpatti",                                                                                                                      
                "playrummy",                                                                                                                      
                "pokerpro",                                                                                                                      
        }  
func main() {
	go setHourlyDataUpdationRedis()
	// useful data starts here ............
	go redisMain()
	fmt.Println("Starting....")
	http.HandleFunc("/fetchData", fetchData)
	http.HandleFunc("/fetchDataPastIncidents", fetchDataPastIncidents)
	Res := http.ListenAndServe(webPort, nil)
	fmt.Println("Ending....", Res)
	// useful data ends here ................
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
	//time.Sleep(5 * time.Second)
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	cntxt := context.Background()
	//val2temp, _ := RedisClient.Do(cntxt, "LRANGE", "data", "0", "-1").Text()
	//val2Total, _ := RedisClient.Do(cntxt, "LRANGE", "data", "0", "-1").StringSlice()
	val2Total, _ := RedisClient.Do(cntxt, "ZRANGE", "data", "0", "-1").StringSlice()
	fmt.Println("val2Total is : ", val2Total)
	val2 := []string{}
	if len(val2Total) > 0{
		val2 = val2Total[1:int(math.Min(7.0, float64(len(val2Total))))]
	}
	//fmt.Println("The val is: ", val2temp)
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
	fmt.Println("Final Json is: ", FinalJson)
	if errJson == nil{
		//writer.WriteHeader(http.StatusOK)
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
	//time.Sleep(5 * time.Second)
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	BodyPlaceHolder := make([]byte, request.ContentLength)
	request.Body.Read(BodyPlaceHolder)
	Body := &BodyGo{}
	err := json.Unmarshal(BodyPlaceHolder, Body)
	if contains(GameNamesList, Body.GameName){
		if err == nil {
			RedisClient := redis.NewClient(&redis.Options{
				Addr:     "localhost:6379",
				Password: "",
				DB:       0,
			})
			//Resp, _ := main1(Body.Date, Body.GameName).encodeToJSON()
			cntxt := context.Background()          
			Resp, err := RedisClient.Get(cntxt, Body.GameName).Result()
			if err != nil{
				Resp11, err11 := main1(Body.Date, Body.GameName)
				if err11 == nil{
					Resp1, _ := Resp11.encodeToJSON()
					fmt.Fprintf(writer, "%s\n", string(Resp1))
				}else{
					fmt.Println("The error is: ", err11)
					http.Error(writer, "Issues With File", http.StatusBadRequest)
					return 
				}
			}else{				
				fmt.Fprintf(writer, "%s\n", string(Resp))
			}
		} else {
			http.Error(writer, "Error Encountered", http.StatusBadRequest)                                           
		}
	}else{
		http.Error(writer, "Game Not Supported", http.StatusBadRequest)
	}
}
