package setup

import (
	"math"

	"github.com/eliiasg/mdcalc/syntax"
)

func genFunctions() map[string]map[int]syntax.Function {
	return map[string]map[int]syntax.Function{
		"floor": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Floor(args[0]), nil
				},
				Latex: "\\lfloor@0\\rfloor",
			},
		},
	}
}
