package setup

import "github.com/eliiasg/mdcalc/syntax"

func GetOperatorPowers() map[string]int {
	return map[string]int{
		"*": 1,
		"/": 1,
		"^": 2,
	}
}

func GenerateEnvironment() *syntax.Environment {
	
}
