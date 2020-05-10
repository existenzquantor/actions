package main

import (
	"flag"
	"fmt"

	"github.com/existenzquantor/actions/model"
)

func main() {
	jsonFile := flag.String("json", "./ressources/flipSwitch.json", "JSON file that contains a domain description.")
	flag.Parse()
	m := model.ParseJSON(*jsonFile)
	fmt.Printf("%v\n", m)
}
