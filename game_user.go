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

func mainGameUser(CurrDate string, GameName string) map[int]bool {
	CurrFileName := generateGameLogFileNameFromDate(CurrDate, GameName)
	PreviousDayFileName := generatePreviousDayGameLogFileName(CurrDate, GameName)
	CurrDayData1 := mainGameUserHelper(CurrFileName)
	PrevDayData1 := mainGameUserHelper(PreviousDayFileName)
	CurrDayData := sanitizeData(CurrDayData1)
	PrevDayData := sanitizeData(PrevDayData1)
	return generateResult(CurrDayData, PrevDayData)
	//return []map[int]int{CurrDayData, PrevDayData}
}

func generateResult(CurrDayData map[int]int, PrevDayData map[int]int) map[int]bool {
	Result := make(map[int]bool, 24)
	for key, CurrDayDataForTheGivenKey := range CurrDayData {
		if CurrDayDataForTheGivenKey != 0 {
			PrevDayDataForTheGivenKey := PrevDayData[key]
			DeltaChangePercentage := (float64(PrevDayDataForTheGivenKey-CurrDayDataForTheGivenKey) / float64(PrevDayDataForTheGivenKey)) * 100.00
			if DeltaChangePercentage > 10 {
				Result[key] = true
			} else {
				Result[key] = false
			}
		} else {
			Result[key] = false
		}
	}
	return Result
}

func sanitizeData(data map[int]int) map[int]int {
	for i := 0; i < 24; i++ {
		_, ok := data[i]
		if !ok {
			data[i] = 0
		}
	}
	return data
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
		if ok  {
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
