package reasoning

import (
	"os"
	"os/exec"
	"strings"

	"github.com/existenzquantor/actions/model"
)

func reasonForAction(action string, causalitypath string, causeProlog string, plan string) model.Reasons {
	tmpFile := causalitypath + "/temp.pl"
	f, _ := os.Create(tmpFile)
	defer f.Close()
	defer os.Remove(tmpFile)
	f.WriteString(causeProlog)
	cmd := exec.Command("./causality", "temp.pl", string(plan), strings.ToLower(action), "reason_temporal_empty_nogoal")
	cmd.Dir = causalitypath
	b, _ := cmd.CombinedOutput()
	return model.ParsePrologOutput(string(b))
}

func plansEqualButNotAtStepI(i int, plan1 []string, plan2 []string) bool {
	if len(plan1) != len(plan2) {
		return false
	}
	for j := 0; j < len(plan1); j++ {
		if i != j && plan1[j] != plan2[j] {
			return false
		}
	}
	return true
}

func keepOnlyReasonsForThisActionToken(reasons model.Reasons, i int, plan string) model.Reasons {
	p := strings.Split(plan, ":")
	var reasons2 []model.Reason
	for _, r := range reasons.Reasons {
		if r.Witness[i] != p[i] {
			if plansEqualButNotAtStepI(i, r.Witness, p) {
				reasons2 = append(reasons2, r)
			}
		}
	}
	return model.Reasons{Reasons: reasons2}
}

func markReasons(goals []model.Literal, reasons model.Reasons) {
	for i := range reasons.Reasons {
		for _, l := range goals {
			if reasons.Reasons[i].Reason.Polarity == l.Polarity && reasons.Reasons[i].Reason.Name == l.Name {
				reasons.Reasons[i].SetIsGoal(true)
			}
		}
	}
}

//ActionConcepts returns the action concepts associated with the actions in the plan
func ActionConcepts(m model.DomainDescription, causeProlog string, plan string, causalitypath string) model.ActionConcepts {
	var concepts []model.ActionConcept
	for i := 0; i < len(m.ProgramDescription); i++ {
		a := m.ProgramDescription[i]
		o := reasonForAction(a, causalitypath, causeProlog, plan)
		markReasons(m.GoalDescription, o)
		o = keepOnlyReasonsForThisActionToken(o, i, plan)
		s := StateAt(i, m)
		n := model.ActionConcept{ActionName: a, Context: s, Causes: o}
		concepts = append(concepts, n)
	}
	return model.ActionConcepts{Concepts: concepts}
}
