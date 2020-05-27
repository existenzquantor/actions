package reasoning

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/existenzquantor/actions/model"
)

func removeString(t []string, e string) []string {
	var tNew []string
	for _, x := range t {
		if x != e {
			tNew = append(tNew, x)
		}
	}
	return tNew
}

func callHermiT(p string, action string) []string {
	cmd := exec.Command("java", "-jar", "HermiT.jar", "-S:"+action, "temp.owl")
	cmd.Dir = p
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	s := string(b)
	s2 := strings.Split(s, "\n")
	for i := 0; i < len(s2)-1; i++ {
		s2[i] = strings.ReplaceAll(s2[i], " ", "")
		s2[i] = strings.ReplaceAll(s2[i], "\t", "")[1:]
	}
	s2 = removeString(s2, "wl:Thing")
	return s2[1 : len(s2)-1]
}

func getAllSubsumers(path string, concept string) []string {
	return callHermiT(path, concept)
}

//ActionDescriptionsFromActionConcepts uses the ontology to infer descriptions (action types) from conceptual action descriptions
func ActionDescriptionsFromActionConcepts(pathOntology string, pathReasoner string, ac model.ActionConcepts, plan []string) model.ActionDescriptions {
	lines := model.ReadOntology(pathOntology)
	lines = lines[0 : len(lines)-1]
	for _, owl := range ac.ToOWLString(plan) {
		lines = append(lines, owl)
	}
	lines = append(lines, ")")
	model.WriteOntology(pathReasoner, lines)
	var ad []model.ActionDescription
	for i := range ac.ToOWLString(plan) {
		t := getAllSubsumers(pathReasoner, "Action"+strconv.Itoa(i))
		ad = append(ad, model.ActionDescription{Step: i, Descriptions: t})
	}
	os.Remove(pathReasoner + "/temp.owl")
	return model.ActionDescriptions{Plan: plan, Descriptions: ad}
}
