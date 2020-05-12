package reasoning

import (
	"github.com/existenzquantor/actions/model"
)

func simulateSteps(t int, domain model.DomainDescription) model.State {
	if t == 0 {
		return domain.InitialStateDescription.State
	}
	curr := model.State{Time: 0, State: domain.InitialStateDescription.State.State}
	for i := 0; i < t; i++ {
		curr.SetStateTime(i + 1)
		actName := domain.ProgramDescription.ActionSequence[i]
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
