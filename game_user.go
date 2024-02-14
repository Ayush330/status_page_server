package main

import (
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

func mainGameUser(CurrDate string) []map[int]int {
	CurrFileName := generateGameLogFileNameFromDate(CurrDate)
	PreviousDayFileName := generatePreviousDayGameLogFileName(CurrDate)
	CurrDayData := mainGameUserHelper(CurrFileName)
	PrevDayData := mainGameUserHelper(PreviousDayFileName)
	return []map[int]int{CurrDayData, PrevDayData}
}

func mainGameUserHelper(FileName string) map[int]int {
	FileData, err := os.ReadFile(strings.TrimSpace(FileName))
	if err == nil {
		Res := mainGameUserHelper1(string(FileData))
		return accumulate_results_per_hour(Res)
	} else {
		panic(err.Error())
	}
}

func accumulate_results_per_hour(input []userLog) map[int]int {
	result := make(map[int]int)
	for _, data := range input {
		val, ok := result[data.samplingHour]
		if ok == true {
			result[data.samplingHour] = val + data.averageNumOfUsers
		} else {
			result[data.samplingHour] = data.averageNumOfUsers
		}
	}
	return result
}

func mainGameUserHelper1(input string) []userLog {
	RegExpPtr, err := regexp.Compile(",\\s\\d*\\s,\\s\\d*\\s,\\s\\d*")
	if err == nil {
		DataList1 := RegExpPtr.FindAllString(input, len(input))
		DataList := formatAndConvertToStruct(DataList1)

		return DataList[:len(DataList)-1]

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
