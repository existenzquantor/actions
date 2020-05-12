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

// Applicable checks if an action is applicable in a state
func (a Action) Applicable(s State) bool {
	for _, p := range a.Precondition {
		if !s.containsStateLiteral(p) {
			return false
		}
	}
	return true
}

// Goal represents a goal description
type Goal struct {
	Specification []Literal
}

// State represents a state
type State struct {
	Time  int
	State []Literal
}

func (s *State) containsStateLiteral(l Literal) bool {
	for _, m := range s.State {
		if m.Polarity == l.Polarity && m.Name == l.Name {
			return true
		}
	}
	return false
}

func (s *State) addLiteral(l Literal) {
	s.State = append(s.State, l)
}

// ApplyAction applies an action to a state
func (s *State) ApplyAction(a Action) {
	if a.Applicable(*s) {
		newState := State{}
		for i := 0; i < len(s.State); i++ {
			if s.State[i].Name == a.Effect.Name {
				newState.addLiteral(Literal{Name: a.Effect.Name, Polarity: a.Effect.Polarity})
			} else {
				newState.addLiteral(s.State[i])
			}
		}
		s.State = newState.State
	}
}

// SetStateTime sets the time field of a state
func (s *State) SetStateTime(t int) {
	s.Time = t
}

// InitialState represents the initial state
type InitialState struct {
	State State
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

// Reason represents a reason
type Reason struct {
	Reason         Literal
	ActionSequence []string
}

// Reasons represents a list of reasons
type Reasons struct {
	Reasons []Reason
}
