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

type incidentsAtAGivenHour struct {
	hour      int
	incidents int
}

func main() {
	// processEachElement("AyushKumarAnand")
	//fmt.Println("GeneratedFileName is: ", generateFileName())
	// Result := averageTimeInstanceBetweenTwoGamePingLog()
	// fmt.Println("The result is: ", Result)
	// EmptyHourFile := createEmptyHourFile()
	// mainGameUser()
	fmt.Println("Curr day filename is: ", generateGameLogFileNameFromDate("2024-01-01"))
	fmt.Println("Previous day filename is: ", generatePreviousDayGameLogFileName("2024-01-01"))
}

func main1() {
	EmptyHourFile := createEmptyHourFile()
	Bytes, err := os.ReadFile("/Users/ayushanand/status_page_server/2024-01-01-game-servers-ping.log")
	if err == nil {
		SplittedstringList1 := strings.Split(string(Bytes), "\n")
		// Last line is a white space so I am removing it.
		SplittedstringList := SplittedstringList1[:len(SplittedstringList1)-1]
		Result := test(SplittedstringList, "pokerserv90")
		// test(SplittedstringList, "pokerserv90")
		// fmt.Printf("The result is: %+v\n", Result[1:10])
		process_game_server_log(Result, EmptyHourFile)
	} else {
		fmt.Println("The error encouneterd while reading the file is: ", err)
	}
}

func test(StringSlice []string, ServerNameToProcess string) []serverStatusForATimeStamp {
	Result := make([]serverStatusForATimeStamp, 0)
	for i := 0; i < len(StringSlice); i++ {
		// if i < 1 {
		if true {
			// fmt.Printf("The data at index: %d is: %s. \n", i, StringSlice[i])
			Res, err := processEachElement(StringSlice[i], ServerNameToProcess)
			if len(Res.differentServerStatuses) != 0 {
				// fmt.Println("Ayush: ", Res.differentServerStatuses[0].serverStatus)
			}
			if err == nil {
				if len(Res.differentServerStatuses) == 0 {
					//fmt.Printf("The res of processEachElement is: %+v.\n", Res)
					if len(Res.differentServerStatuses) == 0 {
						Res.differentServerStatuses = []serverStatus{{serverName: ServerNameToProcess, serverStatus: "UNDEFINED "}}
					}
					Result = append(Result, Res)
					//fmt.Println("The new result is: ", Result)
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
	// RegExpPtr, _ := regexp.Compile("[a-zA-Z0-9]+=[A-Z]+\\s{1}")
	//fmt.Println("The regular expression pattern is: ", Pattern)
	RegExpPtr, _ := regexp.Compile(Pattern)
	AllMatchingStringSlice1 := RegExpPtr.FindAllString(StringElement, len(StringElement))
	// fmt.Println("All matching strings are: ", AllMatchingStringSlice1)
	AllMatchingStringSlice := make([]serverStatus, len(AllMatchingStringSlice1))
	for i := 0; i < len(AllMatchingStringSlice1); i++ {
		AllMatchingStringSlice[i] = parseServerStatus(AllMatchingStringSlice1[i])
	}
	AllRemainingPartSlice := RegExpPtr.Split(StringElement, len(StringElement))
	// fmt.Println("The remaining part of the data is: ", strings.Split(AllRemainingPartSlice[0], " ")[0])
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

func process_game_server_log(inputData []serverStatusForATimeStamp, tempOutputData []incidentsAtAGivenHour) {
	for index, incidentData := range tempOutputData {
		hourThis := incidentData.hour
		incidentsAtCurrHour := process_game_server_log_helper(hourThis, inputData)
		tempOutputData[index] = incidentsAtAGivenHour{
			hour:      hourThis,
			incidents: incidentsAtCurrHour,
		}
	}
	fmt.Println("Result is: ", tempOutputData)
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
	// fmt.Printf("The net count fo hour %d is %d.\n", hour, resultCount)
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

func createEmptyHourFile() []incidentsAtAGivenHour {
	Result := make([]incidentsAtAGivenHour, 24)
	for i := 1; i < 25; i++ {
		Result[i-1].hour = i
	}
	fmt.Println("The result is: ", Result)
	return Result
}
