package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
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

func actionConcepts(m model.DomainDescription, c string, d string, causalitypath *string) model.ActionConcepts {
	var concepts []model.ActionConcept
	for i := 0; i < len(m.ProgramDescription.ActionSequence); i++ {
		a := m.ProgramDescription.ActionSequence[i]
		o := reasonAction(a, causalitypath, c, d)
		s := reasoning.StateAt(i, m)
		n := model.ActionConcept{ActionName: a, Context: s, Causes: o}
		concepts = append(concepts, n)
	}
	return model.ActionConcepts{Concepts: concepts}
}

func main() {
	jsonFile := flag.String("domain", "./ressources/flipSwitch.json", "JSON file that contains a domain description.")
	causalitypath := flag.String("causalitypath", "../causality/", "Path to the executable of causal reasoning, see https://github.com/existenzquantor/causality")
	ontology := flag.String("ontology", "https://github.com/existenzquantor/actions/ressources/FlipSwitch.owl", "IRI of the Ontology to use")
	outputformat := flag.String("outputformat", "types", "types | concepts")
	hermitpath := flag.String("hermitpath", "../../reasoner/HermiT/", "Path to the HermiT OWL reasoner")

	flag.Parse()

	m := model.ParseDomainJSON(*jsonFile)
	c := model.ToCausalityOutput(m)
	d := model.OutputProgram(m)

	switch *outputformat {
	case "concepts":
		ac, err := json.Marshal(actionConcepts(m, c, d, causalitypath))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", string(ac))
	case "types":
		fmt.Printf("TODO:\nWrite new temp.owl that imports %v\n", ontology)
		fmt.Printf("Add new concepts built from %v\n", actionConcepts(m, c, d, causalitypath))
		fmt.Printf("Invoke: %vjava -jar HermiT.jar -c temp.owl\n", hermitpath)
		fmt.Printf("Finally parse hermit's output and extract types\n")
	}
}
