package main

import (
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

func main() {
	process_ping_file("2024-01-02")
}

func main2() {
	TimeStart := time.Now().UnixNano()
	Res := mainGameUser("2024-01-02")
	TimeEnd := time.Now().UnixNano()
	fmt.Println("The diff to run the above computation in nanoseconds is: ", (TimeEnd - TimeStart))
	fmt.Println("Result is: ", Res)
}

func process_ping_file(date string) {
	Bytes, err := os.ReadFile("/Users/ayushanand/status_page_server/2024-01-01-game-servers-ping.log")
	if err == nil {
		SplittedstringList1 := strings.Split(string(Bytes), "\n")
		// Last line is a white space so I am removing it.
		SplittedstringList := SplittedstringList1[:len(SplittedstringList1)-1]
		ResultHelper := test(SplittedstringList, "pokerserv90")
		Result := process_game_server_log(ResultHelper)
		fmt.Println("The result of process_ping_file is:", Result)
	} else {
		fmt.Println("The error encouneterd while reading the file is: ", err)
	}
}

func test(StringSlice []string, ServerNameToProcess string) []serverStatusForATimeStamp {
	Result := make([]serverStatusForATimeStamp, 0)
	for i := 0; i < len(StringSlice); i++ {
		if true {
			Res, err := processEachElement(StringSlice[i], ServerNameToProcess)
			if len(Res.differentServerStatuses) != 0 {
			}
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

// Code for getting the file name

func generateFileName() string {
	Suffix := "-game-servers-ping.log"
	Year, Month, Day := time.Now().UTC().Date()
	MonthInteger := int(Month)
	return fmt.Sprintf("%d-%s-%s%s\n", Year, generateFileNameHelper(MonthInteger), generateFileNameHelper(Day), Suffix)
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
