package model

import (
	"reflect"
	"testing"
)

// TestParseJSON tests ParseJSON
func TestParseJSON(t *testing.T) {

	on := Literal{Negation: true, Name: "on"}
	off := Literal{Negation: true, Name: "off"}
	a1 := Action{Name: "FlipSwitch", Effect: on, Precondition: []Literal{off}}
	a2 := Action{Name: "FlipSwitch", Effect: off, Precondition: []Literal{on}}
	init := InitialState{State: []Literal{off}}
	goal := Goal{Specification: []Literal{on}}
	program := Program{ActionSequence: []Action{a1}}
	want := DomainDescription{ActionDescription: []Action{a1, a2}, InitialStateDescription: init, GoalDescription: goal, ProgramDescription: program}

	if got := ParseJSON("../ressources/flipSwitch.json"); reflect.DeepEqual(got, want) {
		t.Errorf("ParseJSON() = %v, want %v", got, want)
	}
}
