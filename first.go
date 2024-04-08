package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type serverStatus struct {
	serverName   string
	serverStatus string
}

type serverStatusForATimeStamp struct {
	differentServerStatuses []serverStatus
	epochTimeStamp          int64
	sampledYear             int
	sampledMonth            time.Month
	sampledDate             int
	samplingHour            int
}

type GameDataStruct struct {
	Hour       int `json:"hour"`
	Quantifier int `json:"quantifier"`
}

type ResponseRawStruct struct {
	GameData []GameDataStruct `json:"gameData"`
	GameName string           `json:"gameName"`
}

func (input ResponseRawStruct) encodeToJSON() ([]byte, error) {
	return json.Marshal(input)
}

func main1(Date string, GameName string) ResponseRawStruct {
	TimeStart := time.Now().UnixNano()
	PingSpecificData := process_ping_file(Date, GameName) //("2024-01-02")
	GameSpecificData := mainGameUser(Date, GameName)      //("2024-01-02")
	TimeEnd := time.Now().UnixNano()
	fmt.Println("The diff to run the above computation in nanoseconds is: ", (TimeEnd - TimeStart))
	fmt.Printf("The PingSpecificData is: %v and the GameSpecidifc Data is: %v.\n", PingSpecificData, GameSpecificData)
	Res := generateResultForSocialGames(PingSpecificData, GameSpecificData)
	fmt.Printf("The actual result is: %+v.\n", Res)
	return Res
}

func generateResultForSocialGames(PingSpecificData map[int]int, GameSpecificData map[int]bool) ResponseRawStruct {
	Response := ResponseRawStruct{
		GameData: make([]GameDataStruct, 24),
		GameName: "testing_social_game",
	}
	IntermediatryResponse := Response.GameData
	for currHour, pingSpecificDataForTheGivenHour := range PingSpecificData {
		userSpecificDataForTheGivenHour := GameSpecificData[currHour]
		if pingSpecificDataForTheGivenHour > 0 && userSpecificDataForTheGivenHour {
			IntermediatryResponse[currHour] = GameDataStruct{
				Hour:       currHour,
				Quantifier: pingSpecificDataForTheGivenHour,
			}
			//pingSpecificDataForTheGivenHour
		} else {
			IntermediatryResponse[currHour] = GameDataStruct{
				Hour:       currHour,
				Quantifier: 0,
			}
		}
	}
	// return Response.gameData
	Response.GameData = IntermediatryResponse
	return Response
}

func process_ping_file(date string, gameName string) map[int]int {
	Bytes, err := os.ReadFile("/Users/ayushanand/status_page_server/2024-01-01-game-servers-ping.log")
	if err == nil {
		SplittedstringList1 := strings.Split(string(Bytes), "\n")
		// Last line is a white space so I am removing it.
		SplittedstringList := SplittedstringList1[:len(SplittedstringList1)-1]
		ResultHelper := test(SplittedstringList, "pokerserv90")
		Result := process_game_server_log(ResultHelper)
		fmt.Println("The result of process_ping_file is:", Result)
		return Result
	} else {
		// Result := make(map[int]int)
		fmt.Println("The error encouneterd while reading the file is: ", err)
		return make(map[int]int)
	}
}

func test(StringSlice []string, ServerNameToProcess string) []serverStatusForATimeStamp {
	Result := make([]serverStatusForATimeStamp, 0)
	for i := 0; i < len(StringSlice); i++ {
		if true {
			Res, err := processEachElement(StringSlice[i], ServerNameToProcess)
			if err == nil {
				if len(Res.differentServerStatuses) == 0 {
					if len(Res.differentServerStatuses) == 0 {
						Res.differentServerStatuses = []serverStatus{{serverName: ServerNameToProcess, serverStatus: "UNDEFINED "}}
					}
					Result = append(Result, Res)
				} else {
					Result = append(Result, Res)
				}
			} else {
				fmt.Println("Error : ", err)
			}
		}
	}
	return Result
}

func processEachElement(StringElement string, ServerNameToProcess string) (serverStatusForATimeStamp, error) {
	Pattern := fmt.Sprintf("%s=[A-Z]+\\s{1}", ServerNameToProcess)
	RegExpPtr, _ := regexp.Compile(Pattern)
	AllMatchingStringSlice1 := RegExpPtr.FindAllString(StringElement, len(StringElement))
	AllMatchingStringSlice := make([]serverStatus, len(AllMatchingStringSlice1))
	for i := 0; i < len(AllMatchingStringSlice1); i++ {
		AllMatchingStringSlice[i] = parseServerStatus(AllMatchingStringSlice1[i])
	}
	AllRemainingPartSlice := RegExpPtr.Split(StringElement, len(StringElement))
	EpochTimeStamp := strings.TrimSpace(strings.Split(AllRemainingPartSlice[0], " ")[0])
	EpochTimeStampNumeric, err1 := (strconv.ParseInt(EpochTimeStamp, 10, 64))
	if err1 == nil {
		ServerStatusForATimeStamp := timeConversion(EpochTimeStampNumeric, AllMatchingStringSlice)
		return ServerStatusForATimeStamp, nil
	} else {
		fmt.Println("Error while converting string to integer: ", err1)
		return serverStatusForATimeStamp{}, errors.New(err1.Error())
	}
}

func parseServerStatus(serverStatusInput string) serverStatus {
	parts := strings.Split(serverStatusInput, "=")
	return serverStatus{
		serverName:   parts[0],
		serverStatus: parts[1],
	}
}

func process_game_server_log(inputData []serverStatusForATimeStamp) map[int]int {
	result := make(map[int]int)
	for hourThis := 0; hourThis < 24; hourThis++ {
		incidentsAtCurrHour := process_game_server_log_helper(hourThis, inputData)
		result[hourThis] = incidentsAtCurrHour
	}
	return result
}

func process_game_server_log_helper(hour int, inputData []serverStatusForATimeStamp) int {
	filteredData := process_game_server_log_helper2(hour, inputData)
	resultCount := 0
	counter := 0
	for _, data := range filteredData {
		if strings.TrimSpace(data.differentServerStatuses[0].serverStatus) != "OK" {
			counter++
		} else {
			if counter == 3 {
				resultCount += 1
			}
			counter = 0
		}
		if counter == 3 {
			resultCount += 1
			counter = 0
		}
	}
	return resultCount
}

func process_game_server_log_helper2(hour int, inputData []serverStatusForATimeStamp) []serverStatusForATimeStamp {
	var filteredData []serverStatusForATimeStamp
	for _, serverDetail := range inputData {
		if serverDetail.samplingHour == hour {
			filteredData = append(filteredData, serverDetail)
		}
	}
	return filteredData
}
