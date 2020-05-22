package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/existenzquantor/actions/model"
	"github.com/existenzquantor/actions/reasoning"
)

func reasonForAction(action string, causalitypath *string, c string, d string) model.Reasons {
	tmpFile := *causalitypath + "/temp.pl"
	f, _ := os.Create(tmpFile)
	defer f.Close()
	defer os.Remove(tmpFile)
	f.WriteString(c)
	cmd := exec.Command("./causality", "temp.pl", string(d), strings.ToLower(action), "reason_temporal_empty_nogoal")
	cmd.Dir = *causalitypath
	b, _ := cmd.CombinedOutput()
	return model.ParsePrologOutput(string(b))
}

func actionConcepts(m model.DomainDescription, c string, d string, causalitypath *string) model.ActionConcepts {
	var concepts []model.ActionConcept
	for i := 0; i < len(m.ProgramDescription.ActionSequence); i++ {
		a := m.ProgramDescription.ActionSequence[i]
		o := reasonForAction(a, causalitypath, c, d)
		s := reasoning.StateAt(i, m)
		n := model.ActionConcept{ActionName: a, Context: s, Causes: o}
		concepts = append(concepts, n)
	}
	return model.ActionConcepts{Concepts: concepts}
}

func urlToLines(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return linesFromReader(resp.Body)
}

func linesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func main() {
	jsonFile := flag.String("domain", "./ressources/flipSwitch2.json", "JSON file that contains a domain description.")
	causalitypath := flag.String("causalitypath", "../causality/", "Path to the executable of causal reasoning, see https://github.com/existenzquantor/causality")
	ontology := flag.String("ontology", "https://raw.githubusercontent.com/existenzquantor/actions/master/ressources/FlipSwitch.owl", "IRI of the Ontology to use")
	outputformat := flag.String("outputformat", "types", "types | concepts")
	hermitpath := flag.String("hermitpath", "./ressources/", "Path to the HermiT OWL reasoner")

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
		ac := actionConcepts(m, c, d, causalitypath)
		lines, err := urlToLines(*ontology)
		if err != nil {
			log.Fatal(err)
		}
		var owlStrings []string
		for i := 0; i < len(ac.Concepts); i++ {
			owlStrings = append(owlStrings, ac.Concepts[i].ToOWLString(i, m.ProgramDescription.ActionSequence))
		}

		var ad []model.ActionDescription
		for i, owl := range owlStrings {
			f, err := os.Create(*hermitpath + "/temp.owl")
			defer f.Close()
			defer os.Remove(*hermitpath + "/temp.owl")
			if err != nil {
				log.Fatal(err)
			}
			lines[len(lines)-2] = owl
			for _, l := range lines {
				f.WriteString(l + "\n")
			}
			f.Close()
			t := reasoning.GetAllSubsumers(*hermitpath, "Action"+strconv.Itoa(i))
			ad = append(ad, model.ActionDescription{Step: i, Descriptions: t})
		}
		ads := model.ActionDescriptions{Plan: m.ProgramDescription.ActionSequence, Descriptions: ad}
		result, err := json.Marshal(ads)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n", string(result))
	}
}
