package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func openJSON(filename string) []byte {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	return byteValue
}

// ParseJSON parses a JSON file to a DomainDescription
func ParseJSON(filename string) DomainDescription {
	jsonFile := openJSON(filename)
	var dd DomainDescription
	json.Unmarshal(jsonFile, &dd)
	return dd
}
