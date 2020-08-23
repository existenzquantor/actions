package model

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func literalToString(l Literal) string {
	if l.Polarity {
		return fmt.Sprintf("%s", strings.ToLower(l.Name))
	}
	return fmt.Sprintf("not(%s)", strings.ToLower(l.Name))
}

// ToCausalityOutput builds a Prolog representation for the causality program
func ToCausalityOutput(domain DomainDescription) string {
	var s string
	// Actions
	var prec string
	var eff string
	for _, a := range domain.ActionDescription {
		eff = literalToString(a.Effect)
		prec = "["
		for _, b := range a.Precondition {
			prec = prec + literalToString(b)
		}
		prec = prec + "]"
		s = s + fmt.Sprintf("effect(%s, %s, %s).\n", strings.ToLower(a.Name), prec, eff)
	}
	// InitialState
	init := ""
	for _, a := range domain.InitialStateDescription.State.State {
		init = init + "," + literalToString(a)
	}
	s = s + fmt.Sprintf("init([%s]).\n", init[1:])
	// Goal
	goal := ""
	for _, a := range domain.GoalDescription.Specification {
		goal = goal + "," + literalToString(a)
	}
	s = s + fmt.Sprintf("goal([%s]).", goal[1:])
	return s
}

// OutputProgram outputs the program
func OutputProgram(domain DomainDescription) string {
	var s string
	for _, a := range domain.ProgramDescription.ActionSequence {
		s = s + ":" + strings.ToLower(a)
	}
	return s[1:]
}

//ToJSON returns a struct as JSON
func ToJSON(astruct interface{}) string {
	ac, err := json.Marshal(astruct)
	if err != nil {
		log.Fatal(err)
	}
	return string(ac)
}

//AddToOntology adds line to ontology in a temp file located at path
func AddToOntology(path string, ontology []string, line string) {
	f, err := os.Create(path + "/temp.owl")
	if err != nil {
		log.Fatal(err)
	}
	ontology[len(ontology)-2] = line
	for _, l := range ontology {
		f.WriteString(l + "\n")
	}
	f.Close()
}

//WriteOntology writes the ontology to temp.owl
func WriteOntology(path string, ontology []string) {
	f, err := os.Create(path + "/temp.owl")
	if err != nil {
		log.Fatal(err)
	}
	for _, l := range ontology {
		f.WriteString(l + "\n")
	}
	f.Close()
}
