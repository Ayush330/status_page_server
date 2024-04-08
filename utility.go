package main

import (
	"errors"
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

func generateGameLogFileNameFromDate(date string, gameName string) string {
	Suffix, _ := getSuffixFromGameName(gameName) //"-pokerv2-online.log"
	return fmt.Sprintf("./%s%s\n", date, Suffix)
}

func generatePreviousDayGameLogFileName(currDate string, gameName string) string {
	dateInCorrectFormatCurr, err := time.Parse(time.DateOnly, currDate)
	dateInCorrectFormat := dateInCorrectFormatCurr.AddDate(0, 0, -1)
	year, month, day := dateInCorrectFormat.Date()
	if err == nil {
		suffix, _ := getSuffixFromGameName(gameName)
		return fmt.Sprintf("./%s-%s-%s%s", generateFileNameHelper(year), generateFileNameHelper(int(month)), generateFileNameHelper(day), suffix)
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

func getSuffixFromGameName(gameName string) (string, error) {
	switch gameName {
	case "social_poker":
		return "-pokerv2-online.log", nil
	case "indian_rummy":
		return "-pokerv2-online.log", nil
	case "teenpatti":
		return "-pokerv2-online.log", nil
	case "playrummy":
		return "-pokerv2-online.log", nil
	case "pokerpro":
		return "-pokerv2-online.log", nil
	default:
		return "", errors.New("unknown game")
	}
}
