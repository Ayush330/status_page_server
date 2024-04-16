package main

import(
	"fmt"
	"regexp"
	"os"
	"bufio"
	"log"
	"strings"
	"time"
)

var pattern_mapping = map[string]string{
	"teenpatti": "serv[0-9]+=[A-Z a-z]*",
	"social_poker": "pokerserv[0-9]+=[A-Z a-z]*",
	"playrummy": "",
	"indian_rummy": "rummyServ[0-9]+=[A-Z a-z]*",
	"pokerpro": "",
}


func test_redis_data_formation_cash(gameName string) IncidentData{
	Date := "2024-01-02"
	Res, _ := mainGameUser(Date, gameName)
	fmt.Println("The Result Is: ", Res)
	Count := 0
	for _, failover := range Res{
		if failover == true{
			Count++
		}
	}
	fmt.Println("The count of occurence of incidents is: ", Count)
	return IncidentData{Date, Count, gameName}
}

func test_redis_data_formation(gameName string) IncidentData{
	CurrDate := time.Now().Format("2006-01-02")
	Pattern :=  pattern_mapping[gameName]
	RegexCompliantPattern, _:= regexp.Compile(Pattern)
	//fmt.Println("Hello! Ayush Kumar Anand.......................", RegexCompliantPattern)
	file, err := os.Open("./2024-01-01-game-servers-ping.log")
	if err!=nil{
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	Count := 0
	for scanner.Scan(){
		//fmt.Println(scanner.Text())
		Res := test_redis_data_formation_helper(scanner.Text(), RegexCompliantPattern, gameName)
		if Res{
			Count++
		}
	}
	return IncidentData{CurrDate, Count, gameName}
}

func test_redis_data_formation_helper(data string, pattern *regexp.Regexp, gameName string) bool{
	MatchingString := pattern.FindAllString(data, len(data))
	Count := 0
	for _, dataSplit := range MatchingString{
		if (dataSplit != "\n" && dataSplit != "\t"){
			Status := strings.TrimSpace(strings.Split(dataSplit, "=")[1])
			if Status == "OK"{
				Count ++
			}
		}
	}
	//fmt.Println("The Count Is: ", Count, "And the Length Of Server List is: ", len(MatchingString))
	return len(MatchingString) ==  Count
}

