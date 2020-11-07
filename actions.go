package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/existenzquantor/actions/model"
	"github.com/existenzquantor/actions/reasoning"
)

func main() {
	jsonFile := flag.String("domain", "flipSwitch2", "JSON file that contains a domain description.")
	causalitypath := flag.String("causalitypath", "../causality/", "Path to the executable of causal reasoning, see https://github.com/existenzquantor/causality")
	ontology := flag.String("ontology", "FlipSwitch", "Path to the Ontology to use")
	outputformat := flag.String("outputformat", "types", "types | concepts")
	hermitpath := flag.String("hermitpath", "./ressources/", "Path to the HermiT OWL reasoner")

	flag.Parse()

	m := model.ParseDomainJSON("./ressources/domains/" + *jsonFile + ".json")
	c := model.ToCausalityOutput(m)
	d := model.OutputProgram(m)

	switch *outputformat {
	case "concepts":
		ac := model.ToJSON(reasoning.ActionConcepts(m, c, d, causalitypath))
		fmt.Printf("%v\n", string(ac))
	case "types":
		ac := reasoning.ActionConcepts(m, c, d, causalitypath)
		ap, _ := filepath.Abs("./ressources/ontologies/" + *ontology + ".owl")
		ads := reasoning.ActionDescriptionsFromActionConcepts(ap, *hermitpath, ac, m.ProgramDescription.ActionSequence)
		result := model.ToJSON(ads)
		fmt.Printf("%v\n", string(result))
	}
}
