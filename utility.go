package main

import (
	"fmt"
	"time"
)

func timeConversion(TimeStampInSeconds int64, ServerStatusList []serverStatus) serverStatusForATimeStamp {
	Time := time.Unix(TimeStampInSeconds, 0).Local().UTC()
	Year, Month, Day := Time.Date()
	return serverStatusForATimeStamp{
		differentServerStatuses: ServerStatusList,
		epochTimeStamp:          TimeStampInSeconds,
		sampledYear:             Year,
		sampledMonth:            Month,
		sampledDate:             Day,
		samplingHour:            Time.Hour(),
	}
}

func timeConversionUserLogHelper(TimeStampInSeconds int64) (int, time.Month, int, int) {
	Time := time.Unix(TimeStampInSeconds, 0).Local().UTC()
	year, month, day := Time.Date()
	return year, month, day, Time.Hour()
}

func generateGameLogFileNameFromDate(date string) string {
	Suffix := "-pokerv2-online.log"
	return fmt.Sprintf("./%s%s\n", date, Suffix)
}

func generatePreviousDayGameLogFileName(currDate string) string {
	dateInCorrectFormatCurr, err := time.Parse(time.DateOnly, currDate)
	dateInCorrectFormat := dateInCorrectFormatCurr.AddDate(0, 0, -1)
	year, month, day := dateInCorrectFormat.Date()
	if err == nil {
		// fmt.Printf("The parsed date is : %s-%s-%s\n", generateFileNameHelper(year), generateFileNameHelper(int(month)), generateFileNameHelper(day))
		return fmt.Sprintf("./%s-%s-%s-pokerv2-online.log", generateFileNameHelper(year), generateFileNameHelper(int(month)), generateFileNameHelper(day))
	} else {
		fmt.Println("The error in parsing the date is : ", err.Error())
		return ""
	}
}

func generateFileNameHelper(num int) string {
	if int(num/10) == 0 {
		return fmt.Sprintf("0%d", num)
	} else {
		return fmt.Sprintf("%d", num)
	}
}
