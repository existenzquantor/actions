package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/existenzquantor/actions/model"
)

func main() {
	jsonFile := flag.String("json", "./ressources/flipSwitch.json", "JSON file that contains a domain description.")
	causalityProgram := flag.String("causalitypath", "../causality/", "Path to executable of causal reasoning")

	flag.Parse()

	m := model.ParseDomainJSON(*jsonFile)
	fmt.Printf("%v\n", m)

	c := model.ToCausalityOutput(m)
	println(c)

	d := model.OutputProgram(m)
	println(d)

	tmpFile := *causalityProgram + "/temp.pl"
	f, _ := os.Create(tmpFile)
	defer f.Close()
	f.WriteString(c)
	println(tmpFile)
	cmd := exec.Command("./causality", "temp.pl", "flipswitch", "on", "temporal_empty")
	cmd.Dir = *causalityProgram
	println(cmd.String())
	b, _ := cmd.CombinedOutput()
	println(string(b))
}
