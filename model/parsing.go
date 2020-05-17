package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

// ParseDomainJSON parses a JSON file to a DomainDescription
func ParseDomainJSON(filename string) DomainDescription {
	jsonFile := openJSON(filename)
	var dd DomainDescription
	json.Unmarshal(jsonFile, &dd)
	return dd
}

// ParsePrologOutput parses the Prolog output
func ParsePrologOutput(s string) Reasons {
	if strings.HasPrefix(s, "[]") {
		return Reasons{}
	}
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	s = strings.ReplaceAll(s, "\n", "")
	sa := strings.Split(s, "),(")
	var la []string
	var lit Literal
	var rea []Reason
	for _, l := range sa {
		la = strings.Split(l, ",")
		if strings.HasPrefix(la[0], "(not(") {
			lit = Literal{Polarity: false, Name: la[0][5 : len(la[0])-1]}
		} else if strings.HasPrefix(la[0], "not(") {
			lit = Literal{Polarity: false, Name: la[0][4 : len(la[0])-1]}
		} else {
			laNo := strings.ReplaceAll(la[0], ")", "")
			laNo = strings.ReplaceAll(la[0], "(", "")
			lit = Literal{Polarity: true, Name: laNo}
		}
		rea = append(rea, Reason{Reason: lit, Witness: strings.Split(la[1], ":")})
	}
	return Reasons{Reasons: rea}
}
