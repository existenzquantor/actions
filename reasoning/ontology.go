package reasoning

import (
	"log"
	"os/exec"
	"strings"
)

func removeString(t []string, e string) []string {
	var tNew []string
	for _, x := range t {
		if x != e {
			tNew = append(tNew, x)
		}
	}
	return tNew
}

func callHermiT(p string, action string) []string {
	cmd := exec.Command("java", "-jar", "HermiT.jar", "-S:"+action, "temp.owl")
	cmd.Dir = p
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	s := string(b)
	s2 := strings.Split(s, "\n")
	for i := 0; i < len(s2)-1; i++ {
		s2[i] = strings.ReplaceAll(s2[i], " ", "")
		s2[i] = strings.ReplaceAll(s2[i], "\t", "")[1:]
	}
	s2 = removeString(s2, "wl:Thing")
	return s2[1 : len(s2)-1]
}

//GetAllSubsumers calls a reasoner at a particular path to ask for all subsumbers of a given concept
func GetAllSubsumers(path string, concept string) []string {
	return callHermiT(path, concept)
}
