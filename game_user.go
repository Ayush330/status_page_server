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

// type numOfUsersAtAGivenHour struct {
// 	hour       int
// 	numOfUsers int
// }

func mainGameUser(CurrDate string) []map[int]int {
	CurrFileName := generateGameLogFileNameFromDate(CurrDate)
	PreviousDayFileName := generatePreviousDayGameLogFileName(CurrDate)
	CurrDayData := mainGameUserHelper(CurrFileName)
	PrevDayData := mainGameUserHelper(PreviousDayFileName)
	return []map[int]int{CurrDayData, PrevDayData}
}

func mainGameUserHelper(FileName string) map[int]int {
	// fmt.Println("The filename is: ", FileName)
	FileData, err := os.ReadFile(strings.TrimSpace(FileName))
	if err == nil {
		// fmt.Println("The file data is: ", string(FileData))
		Res := mainGameUserHelper1(string(FileData))
		return accumulate_results_per_hour(Res)
	} else {
		panic(err.Error())
	}
}

func accumulate_results_per_hour(input []userLog) map[int]int {
	result := make(map[int]int)
	for _, data := range input {
		// if data.samplingHour == 0 {
		// 	fmt.Println("Data is: ", data)
		// }
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
		DataList := formatAndConvertToStruct(DataList1) //strings.Split(DataList1[0], ", ")[1:]
		// DataList = format(DataList2)
		// fmt.Printf("The DataList is: %+v\n", DataList[0:8])
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
		// if hour == 0 {
		// 	fmt.Println("Th data regarding 0 hrs is: ", timestamp, " ", hour)
		// }
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
