package reasoning

import (
	"os"
	"os/exec"
	"strings"

	"github.com/existenzquantor/actions/model"
)

func reasonForAction(action string, causalitypath *string, causeProlog string, plan string) model.Reasons {
	tmpFile := *causalitypath + "/temp.pl"
	f, _ := os.Create(tmpFile)
	defer f.Close()
	defer os.Remove(tmpFile)
	f.WriteString(causeProlog)
	cmd := exec.Command("./causality", "temp.pl", string(plan), strings.ToLower(action), "reason_temporal_empty_nogoal")
	cmd.Dir = *causalitypath
	b, _ := cmd.CombinedOutput()
	return model.ParsePrologOutput(string(b))
}

//ActionConcepts returns the action concepts associated with the actions in the plan
func ActionConcepts(m model.DomainDescription, causeProlog string, plan string, causalitypath *string) model.ActionConcepts {
	var concepts []model.ActionConcept
	for i := 0; i < len(m.ProgramDescription.ActionSequence); i++ {
		a := m.ProgramDescription.ActionSequence[i]
		o := reasonForAction(a, causalitypath, causeProlog, plan)
		s := StateAt(i, m)
		n := model.ActionConcept{ActionName: a, Context: s, Causes: o}
		concepts = append(concepts, n)
	}
	return model.ActionConcepts{Concepts: concepts}
}
