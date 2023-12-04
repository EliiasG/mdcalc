package setup

import "github.com/eliiasg/mdcalc/syntax"

func GenerateEnvironment(lib syntax.UnitLibrary) *syntax.Environment {
	return &syntax.Environment{
		Operators:      genOperators(),
		Functions:      genFunctions(),
		VariableValues: map[string]syntax.VariableValue{},
		OperatorPowers: map[string]int{
			"*":  1,
			"/":  1,
			"//": 1,
			"^":  2,
		},
		Formatter:   &formatter{},
		UnitLibrary: lib,
	}
}
