package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func averageTimeInstanceBetweenTwoGamePingLog() int {
	var FileName string = "./test.log"
	FileDataSliceOfBytes, err := os.ReadFile(FileName)
	// FileDataSliceOfBytes := TempFileDataSliceOfBytes[1:100]
	TempFileDataSliceAsStringSlice := strings.Split(string(FileDataSliceOfBytes), "\n")
	FileDataAsSliceOfString := TempFileDataSliceAsStringSlice[:len(TempFileDataSliceAsStringSlice)-1]
	handle_error(err)
	// fmt.Printf("The data read from the file %s is %s.\n", FileName, FileDataAsSliceOfString[0])
	return processSlice(FileDataAsSliceOfString)
}

// Internal Functions

func processSlice(inputData []string) int {
	//diffSlice := make([]int, len(inputData)-1)
	TimeSum := 0
	for i := 1; i < len(inputData); i++ {
		// fmt.Println("input is: ", inputData[i])
		input1, err1 := strconv.ParseInt(strings.TrimSpace(inputData[i]), 10, 64)
		input2, err2 := strconv.ParseInt(strings.TrimSpace(inputData[i-1]), 10, 64)
		handle_error(err1)
		handle_error(err2)
		TimeSum += (int(input1) - int(input2))
		fmt.Println("TimeSum is: ", TimeSum)
	}
	fmt.Println("TimeSum ....... is: ", TimeSum)
	return TimeSum / (len(inputData) - 1)
}

func handle_error(errorMsg error) {
	if errorMsg == nil {
		//do nothing as of now
	} else {
		//panic(fmt.Sprintf("Panicking because of error: %v", errorMsg))
		fmt.Println("The error message is: ", errorMsg)
	}
}
