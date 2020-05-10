package model

import (
	"reflect"
	"testing"
)

func getDomainDescription() DomainDescription {
	on := Literal{Polarity: true, Name: "on"}
	off := Literal{Polarity: false, Name: "on"}
	a1 := Action{Name: "FlipSwitch", Effect: on, Precondition: []Literal{off}}
	a2 := Action{Name: "FlipSwitch", Effect: off, Precondition: []Literal{on}}
	init := InitialState{State: []Literal{off}}
	goal := Goal{Specification: []Literal{on}}
	program := Program{ActionSequence: []string{a1.Name, a1.Name}}
	want := DomainDescription{ActionDescription: []Action{a1, a2}, InitialStateDescription: init, GoalDescription: goal, ProgramDescription: program}
	return want
}

// TestParseJSON tests ParseJSON
func TestParseJSON(t *testing.T) {

	want := getDomainDescription()
	got := ParseDomainJSON("../ressources/flipSwitch.json")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\nw: %v\ng: %v", want, got)
	}
}

func TestOutputProgram(t *testing.T) {
	want := "flipswitch;flipswitch"
	got := OutputProgram(ParseDomainJSON("../ressources/flipSwitch.json"))

	if want != got {
		t.Errorf("\nw: %v\ng: %v", want, got)
	}
}
