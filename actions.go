package main

import (
	"encoding/json"
	"fmt"

	"github.com/existenzquantor/actions/model"
)

func main() {
	on := model.Literal{Negation: true, Name: "on"}
	off := model.Literal{Negation: true, Name: "off"}
	a1 := model.Action{Name: "FlipSwitch", Effect: on, Precondition: []model.Literal{off}}
	a2 := model.Action{Name: "FlipSwitch", Effect: off, Precondition: []model.Literal{on}}
	init := model.InitialState{State: []model.Literal{off}}
	goal := model.Goal{Specification: []model.Literal{on}}
	program := model.Program{ActionSequence: []model.Action{a1}}
	domain := model.DomainDescription{ActionDescription: []model.Action{a1, a2}, InitialStateDescription: init, GoalDescription: goal, ProgramDescription: program}
	m, _ := json.Marshal(domain)
	fmt.Printf("%v\n", string(m))

	var n model.DomainDescription
	_ = json.Unmarshal(m, &n)
	fmt.Printf("%v\n", n)
}
