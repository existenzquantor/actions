package reasoning

import (
	"github.com/existenzquantor/actions/model"
)

func simulateSteps(t int, domain model.DomainDescription) model.State {
	if t == 0 {
		return domain.InitialStateDescription
	}
	curr := model.State{Time: 0, State: domain.InitialStateDescription.State}
	for i := 0; i < t; i++ {
		curr.SetStateTime(i + 1)
		actName := domain.ProgramDescription[i]
		var applicable []model.Action
		for _, a := range domain.ActionDescription {
			if a.Applicable(curr) {
				applicable = append(applicable, a)
			}
		}
		for _, a := range applicable {
			if actName == a.Name {
				curr.ApplyAction(a)
			}
		}
	}
	return curr
}

// StateAt returns the state at time t
func StateAt(t int, domain model.DomainDescription) model.State {
	return simulateSteps(t, domain)
}

//UsedFacts returns the set of literals that where used by the t-th action
func UsedFacts(t int, domain model.DomainDescription) []model.Literal {
	var facts []model.Literal
	s := StateAt(t, domain)
	act := domain.ProgramDescription[t]
	for _, a := range domain.ActionDescription {
		if a.Name == act {
			if a.Applicable(s) {
				facts = append(facts, a.Precondition...)
			}
		}
	}
	return facts
}
