package main

import(
	"fmt"
	"io/ioutil"
	"gopkg.in/yaml.v3"
	"errors"
)

func readFile() (map[string]interface{}, error){
	obj := make(map[string]interface{})
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil{
		fmt.Println("Error encountered while reading the yaml file \"config.yaml\". The error is: ", err)
	}
	err = yaml.Unmarshal(yamlFile, obj)
	if err != nil {
		fmt.Printf("Error encountered while Unmarshalling : %v", err)
		return nil, errors.New("Error Encountered While Reading The YAML file.")
	}
	return obj, nil
}

func findData(key string)(interface{}, error){
	obj, err := readFile()
	if err !=nil{
		return nil, errors.New("Data Does Not Exist")
	}
	rData := obj[key]
	if rData == nil{
		return nil, errors.New("Data Does Not Exist")
	}
	return rData, nil
}

