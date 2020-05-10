package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func openJSON(filename string) *os.File {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	return jsonFile
}

// ParseJSON parses a JSON file to a DomainDescription
func ParseJSON(filename string) DomainDescription {
	jsonFile := openJSON(filename)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dd DomainDescription
	json.Unmarshal(byteValue, &dd)
	return dd
}
