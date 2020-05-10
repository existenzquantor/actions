package model

// Literal represents a literal, either positive or negative
type Literal struct {
	Polarity bool
	Name     string
}

// Action represents an action with name, effect, and preconditions
type Action struct {
	Name         string
	Effect       Literal
	Precondition []Literal
}

// Goal represents a goal description
type Goal struct {
	Specification []Literal
}

// InitialState represents the initial state
type InitialState struct {
	State []Literal
}

// Program represents the action sequence
type Program struct {
	ActionSequence []string
}

// DomainDescription represents actions, initial state, goal, and program
type DomainDescription struct {
	ActionDescription       []Action
	InitialStateDescription InitialState
	GoalDescription         Goal
	ProgramDescription      Program
}
