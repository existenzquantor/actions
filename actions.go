package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/existenzquantor/actions/model"
	"github.com/existenzquantor/actions/reasoning"
)

func reasonAction(action string, causalitypath *string, c string, d string) model.Reasons {
	tmpFile := *causalitypath + "/temp.pl"
	f, _ := os.Create(tmpFile)
	defer f.Close()
	defer os.Remove(tmpFile)
	f.WriteString(c)
	cmd := exec.Command("./causality", "temp.pl", string(d), action, "reason_temporal_empty")
	cmd.Dir = *causalitypath
	b, _ := cmd.CombinedOutput()
	return model.ParsePrologOutput(string(b))
}

func main() {
	jsonFile := flag.String("json", "./ressources/flipSwitch2.json", "JSON file that contains a domain description.")
	causalitypath := flag.String("causalitypath", "../causality/", "Path to executable of causal reasoning")

	flag.Parse()

	m := model.ParseDomainJSON(*jsonFile)
	c := model.ToCausalityOutput(m)
	d := model.OutputProgram(m)

	for e, a := range m.ProgramDescription.ActionSequence {
		o := reasonAction(a, causalitypath, c, d)
		fmt.Printf("%v: %v => %v\n", e, a, o)
	}

	for i := 0; i < len(m.ProgramDescription.ActionSequence); i++ {
		s := reasoning.StateAt(i, m)
		fmt.Printf("%v: %v => %v\n", i, s, m.ProgramDescription.ActionSequence[i])
	}
}
