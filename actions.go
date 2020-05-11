package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/existenzquantor/actions/model"
)

func processPrologOutput(s string) string {
	s = strings.Replace(s, "),(", ");(", -1)
	s = strings.Replace(s, ",", "\",\"", -1)
	s = strings.Replace(s, "(", "[\"", -1)
	s = strings.Replace(s, ")", "\"]", -1)
	s = strings.Replace(s, ";", ",", -1)
	return s
}

func reasonAction(action string, causalitypath *string, c string, d string) model.Reasons {
	tmpFile := *causalitypath + "/temp.pl"
	f, _ := os.Create(tmpFile)
	defer f.Close()
	defer os.Remove(tmpFile)
	f.WriteString(c)
	cmd := exec.Command("./causality", "temp.pl", string(d), action, "reason_temporal_empty")
	cmd.Dir = *causalitypath
	b, _ := cmd.CombinedOutput()
	st := "{\n\"Reasons\": " + processPrologOutput(string(b)) + "}"
	var rea model.ReasonsIntermediate
	err := json.Unmarshal([]byte(st), &rea)
	if err != nil {
		log.Fatal(err)
	}
	var reasons []model.Reason
	for _, r := range rea.Reasons {
		var l model.Literal
		if strings.HasPrefix(r[0], "not(") {
			l = model.Literal{Polarity: false, Name: r[0][4 : len(r[0])-1]}
		} else {
			l = model.Literal{Polarity: true, Name: r[0]}
		}
		as := strings.Split(r[1], ":")
		reasons = append(reasons, model.Reason{Reason: l, ActionSequence: as})
	}
	return model.Reasons{Reasons: reasons}
}

func main() {
	jsonFile := flag.String("json", "./ressources/flipSwitch.json", "JSON file that contains a domain description.")
	causalitypath := flag.String("causalitypath", "../causality/", "Path to executable of causal reasoning")

	flag.Parse()

	m := model.ParseDomainJSON(*jsonFile)
	c := model.ToCausalityOutput(m)
	d := model.OutputProgram(m)

	for _, a := range m.ProgramDescription.ActionSequence {
		o := reasonAction(a, causalitypath, c, d)
		fmt.Printf("%v => %v\n", a, o)
	}
}
