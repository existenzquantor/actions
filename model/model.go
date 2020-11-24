package model

import (
	"strconv"
	"strings"
)

// Literal represents a literal, either positive or negative
type Literal struct {
	Polarity bool
	Name     string
}

// Action represents an action with name, effect, and preconditions
type Action struct {
	Name         string
	Effect       []Literal
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

// State represents a state
type State struct {
	Time  int
	State []Literal
}

// Copy returns a copy of the given state
func (s *State) Copy() State {
	return State{Time: s.Time, State: s.State}
}

func (s *State) containsStateLiteral(l Literal) bool {
	for _, m := range s.State {
		if m.Polarity == l.Polarity && m.Name == l.Name {
			return true
		}
	}
	return false
}

func (s *State) containsStateLiteralName(l Literal) bool {
	for _, m := range s.State {
		if m.Name == l.Name {
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
		for _, ac := range a.Effect {
			newState.addLiteral(Literal{Name: ac.Name, Polarity: ac.Polarity})
		}
		for i := 0; i < len(s.State); i++ {
			if !newState.containsStateLiteralName(s.State[i]) {
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

// DomainDescription represents actions, initial state, goal, and program
type DomainDescription struct {
	ActionDescription       []Action
	InitialStateDescription State
	GoalDescription         []Literal
	ProgramDescription      []string
}

// Reason represents a reason
type Reason struct {
	Reason  Literal
	Witness []string
	IsGoal  bool
}

// SetIsGoal sets goal attribute
func (r *Reason) SetIsGoal(b bool) {
	r.IsGoal = b
}

// Reasons represents a list of reasons
type Reasons struct {
	Reasons []Reason
}

// ActionConcept describes an action conceptually
type ActionConcept struct {
	ActionName string
	Context    State
	Causes     Reasons
}

// ToOWLString outputs the action concept as a string in OWL functional syntax
func (c ActionConcept) ToOWLString(planstep int, plan []string) string {
	s := "EquivalentClasses(:Action" + strconv.Itoa(planstep)
	var contextStrings []string
	for _, l := range c.Context.State {
		if l.Polarity == false {
			contextStrings = append(contextStrings, "ObjectComplementOf(:"+l.Name+")")
		} else {
			contextStrings = append(contextStrings, ":"+l.Name)
		}
	}
	sContext := "ObjectSomeValuesFrom(:context "
	sContextLits := ""
	for _, i := range contextStrings {
		sContextLits = sContextLits + " " + i
	}
	if len(contextStrings) > 1 {
		sContextLits = "ObjectIntersectionOf(" + sContextLits + ")"
	}
	sContext = sContext + sContextLits + ")"

	var reasonStrings []string
	for _, l := range c.Causes.Reasons {
		if l.Witness[planstep] != strings.ToLower(plan[planstep]) {
			if l.Reason.Polarity == false {
				reasonStrings = append(reasonStrings, "ObjectComplementOf(:"+l.Reason.Name+")")
			} else {
				reasonStrings = append(reasonStrings, ":"+l.Reason.Name)
			}
		}
	}
	sReasons := "ObjectSomeValuesFrom(:causes "
	sReasonLits := ""
	for _, i := range reasonStrings {
		sReasonLits = sReasonLits + " " + i
	}

	if sReasonLits == "" {
		sReasonLits = " owl:Thing"
	}

	if len(reasonStrings) > 1 {
		sReasonLits = "ObjectIntersectionOf(" + sReasonLits + ")"
	}
	sReasons = sReasons + sReasonLits + ")"

	//

	var reasonStrings2 []string
	for _, l := range c.Causes.Reasons {
		if l.IsGoal {
			if l.Witness[planstep] != strings.ToLower(plan[planstep]) {
				if l.Reason.Polarity == false {
					reasonStrings2 = append(reasonStrings2, "ObjectComplementOf(:"+l.Reason.Name+")")
				} else {
					reasonStrings2 = append(reasonStrings2, ":"+l.Reason.Name)
				}
			}
		}
	}
	sReasons2 := "ObjectSomeValuesFrom(:hasReason "
	sReasonLits2 := ""
	for _, i := range reasonStrings2 {
		sReasonLits2 = sReasonLits2 + " " + i
	}

	if sReasonLits2 == "" {
		sReasonLits2 = " owl:Thing"
	}

	if len(reasonStrings2) > 1 {
		sReasonLits2 = "ObjectIntersectionOf(" + sReasonLits2 + ")"
	}
	sReasons2 = sReasons2 + sReasonLits2 + ")"

	s = s + " ObjectIntersectionOf(:" + c.ActionName + " " + sContext + " " + sReasons + " " + sReasons2 + "))"

	return s
}

//ActionConcepts represents the conceptual descriptions of each action in a sequence
type ActionConcepts struct {
	Concepts []ActionConcept
}

//ToOWLString returns an array of OWL string representations of the action concepts
func (ac ActionConcepts) ToOWLString(plan []string) []string {
	var owlStrings []string
	for i := 0; i < len(ac.Concepts); i++ {
		owlStrings = append(owlStrings, ac.Concepts[i].ToOWLString(i, plan))
	}
	return owlStrings
}

// ActionDescription represents the set of descriptions under which an action falls
type ActionDescription struct {
	Step         int
	Descriptions []string
}

// ActionDescriptions represent a set of ActionDescription structs
type ActionDescriptions struct {
	Plan         []string
	Descriptions []ActionDescription
}
