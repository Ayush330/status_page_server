package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type userLog struct {
	averageNumOfUsers int
	epochTimeStamp    int64
	sampledYear       int
	sampledMonth      time.Month
	sampledDate       int
	samplingHour      int
}

// type numOfUsersAtAGivenHour struct {
// 	hour       int
// 	numOfUsers int
// }

func mainGameUser() {
	FileData, err := os.ReadFile("./2024-01-01-pokerv2-online.log")
	if err == nil {
		// fmt.Println("The file data is: ", string(FileData))
		mainGameUserHelper(string(FileData))
	} else {
		panic(err.Error())
	}
}

func mainGameUserHelper(input string) {
	RegExpPtr, err := regexp.Compile(",\\s\\d*\\s,\\s\\d*\\s,\\s\\d*")
	if err == nil {
		DataList1 := RegExpPtr.FindAllString(input, len(input))
		DataList := formatAndConvertToStruct(DataList1) //strings.Split(DataList1[0], ", ")[1:]
		// DataList = format(DataList2)
		fmt.Printf("The DataList is: %+v\n", DataList[0:8])
	} else {
		panic(err.Error())
	}
}

func formatAndConvertToStruct(input []string) []userLog {
	outputData := make([]userLog, len(input))
	for index, StringEl := range input {
		tempData := strings.Split(StringEl, ", ")[1:]
		numUsers, _ := strconv.ParseInt(strings.TrimSpace(tempData[1]), 10, 64)
		timestamp, _ := strconv.ParseInt(strings.TrimSpace(tempData[0]), 10, 64)
		year, month, date, hour := timeConversionUserLogHelper(timestamp)
		UserLog := userLog{
			averageNumOfUsers: int(numUsers),
			epochTimeStamp:    timestamp,
			sampledYear:       year,
			sampledMonth:      month,
			sampledDate:       date,
			samplingHour:      hour,
		}
		outputData[index] = UserLog
	}
	return outputData
}

func format(input []string) []string {
	for index, Element := range input {
		ModfEl1 := strings.TrimSpace(Element)
		input[index] = ModfEl1
	}
	return input
}
