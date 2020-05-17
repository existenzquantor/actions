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
	cmd := exec.Command("./causality", "temp.pl", string(d), action, "reason_temporal_empty")
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

func contains(a []string, s string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

func bfs(subs [][]string, path []string, goal string) bool {
	if path[len(path)-1] == goal {
		return true
	}
	if len(path) < len(subs) {
		for _, s := range subs {
			if s[0] == path[len(path)-1] && !contains(path, s[1]) {
				path_new := make([]string, len(path))
				copy(path_new, path)
				path_new = append(path_new, s[1])
				if bfs(subs, path_new, goal) {
					return true
				}
			}
		}
	}
	return false
}

func subsumes(subs [][]string, c string, d string) bool {
	return bfs(subs, []string{c}, d)
}

func callHermiT(p string) [][]string {
	cmd := exec.Command("java", "-jar", "HermiT.jar", "-c", "temp.owl")
	cmd.Dir = p
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	s := string(b)
	s = strings.ReplaceAll(s, "\n", "")
	a := strings.Split(s, ">")
	for i := 0; i < len(a)-1; i++ {
		index := strings.Index(a[i], "#")
		if strings.HasPrefix(a[i], "EquivalentClasses") {
			a[i] = "Eq:" + a[i][index+1:]
		} else {
			a[i] = a[i][index+1:]
		}
	}
	var subsumptions [][]string
	for i := 0; i < len(a[:len(a)-1]); i = i + 2 {
		if strings.HasPrefix(a[:len(a)-1][i], "Eq:") {
			subsumptions = append(subsumptions, []string{a[:len(a)-1][i][3:], a[:len(a)-1][i+1]})
			subsumptions = append(subsumptions, []string{a[:len(a)-1][i+1], a[:len(a)-1][i][3:]})
		} else {
			subsumptions = append(subsumptions, []string{a[:len(a)-1][i], a[:len(a)-1][i+1]})
		}
	}
	return subsumptions
}

func extractAllTypes(subs [][]string) []string {
	var types []string
	for _, i := range subs {
		if !contains(types, i[0]) {
			types = append(types, i[0])
		}
		if !contains(types, i[1]) {
			types = append(types, i[1])
		}
	}
	return types
}

func getDescriptions(subs [][]string, action string) []string {
	types := extractAllTypes(subs)
	var descriptions []string
	for _, t := range types {
		if subsumes(subs, action, t) {
			descriptions = append(descriptions, t)
		}
	}
	return descriptions
}

func UrlToLines(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return LinesFromReader(resp.Body)
}

func LinesFromReader(r io.Reader) ([]string, error) {
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
	jsonFile := flag.String("domain", "./ressources/flipSwitch.json", "JSON file that contains a domain description.")
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
		fmt.Printf("TODO:\nWrite new temp.owl that imports %v\n", *ontology)
		lines, err := UrlToLines(*ontology)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Content: %v\n", lines)
		fmt.Printf("Add new concepts built from %v\n", ac)
		var owlStrings []string
		for i := 0; i < len(ac.Concepts); i++ {
			owlStrings = append(owlStrings, ac.Concepts[i].ToOWLString(i))
		}

		for i, owl := range owlStrings {
			f, err := os.Create(*hermitpath + "/temp.owl")
			if err != nil {
				log.Fatal(err)
			}
			lines[len(lines)-2] = owl
			for _, l := range lines {
				f.WriteString(l + "\n")
			}
			f.Close()

			fmt.Printf("Invoke: %vjava -jar HermiT.jar -c temp.owl\n", *hermitpath)
			s := callHermiT(*hermitpath)
			fmt.Printf("Subs: %v\n", s)
			t := getDescriptions(s, "C"+strconv.Itoa(i))
			fmt.Printf("%v can be described as: %v\n", "C"+strconv.Itoa(i), t)
		}
	}
}
