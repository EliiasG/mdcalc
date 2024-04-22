package setup

import (
	"math"

	"github.com/eliiasg/mdcalc/syntax"
)

func genFunctions() map[string]map[int]syntax.Function {
	return map[string]map[int]syntax.Function{
		// functions
		"floor": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Floor(args[0]), nil
				},
				Latex: "\\lfloor@0\\rfloor",
			},
		},
		"ceil": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Ceil(args[0]), nil
				},
				Latex: "\\lceil@0\\rceil",
			},
		},
		"abs": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Abs(args[0]), nil
				},
				Latex: "\\lvert@0\\rvert",
			},
		},
		"sqrt": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Sqrt(args[0]), nil
				},
				Latex: "\\sqrt{@0}",
			},
		},
		"root": {
			2: {
				Execute: func(args []float64) (float64, error) {
					return math.Pow(args[1], 1/args[0]), nil
				},
				Latex: "\\sqrt[@0]{@1}",
			},
		},
		"log10": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Log10(args[0]), nil
				},
				Latex: "\\log @0",
			},
		},
		"log": {
			2: {
				Execute: func(args []float64) (float64, error) {
					return math.Log2(args[1]) / math.Log2(args[0]), nil
				},
				Latex: "\\log_{@0} @1",
			},
		},
		"sin": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Sin(args[0] * math.Pi / 180), nil
				},
				Latex: "\\sin @0",
			},
		},
		"cos": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Cos(args[0] * math.Pi / 180), nil
				},
				Latex: "\\cos @0",
			},
		},
		"tan": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Tan(args[0] * math.Pi / 180), nil
				},
				Latex: "\\tan @0",
			},
		},
		"asin": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Asin(args[0]) / math.Pi * 180, nil
				},
				Latex: "\\arcsin @0",
			},
		},
		"acos": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Acos(args[0]) / math.Pi * 180, nil
				},
				Latex: "\\arccos @0",
			},
		},
		"atan": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return math.Atan(args[0]) / math.Pi * 180, nil
				},
				Latex: "\\arctan @0",
			},
		},
		"mod": {
			2: {
				Execute: func(args []float64) (float64, error) {
					return math.Mod(args[0], args[1]), nil
				},
				Latex: "@0 \\mod @1",
			},
		},
		// symbols
		"pi": {
			0: {
				Execute: func(args []float64) (float64, error) {
					return math.Pi, nil
				},
				Latex: "\\pi",
			},
		},
		"e": {
			0: {
				Execute: func(args []float64) (float64, error) {
					return math.E, nil
				},
				Latex: "\\e",
			},
		},
		// util / formatting
		"par": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return args[0], nil
				},
				Latex: "(@0)",
			},
		},
		"neg": {
			1: {
				Execute: func(args []float64) (float64, error) {
					return -args[0], nil
				},
				Latex: "-@0",
			},
		},
	}
}
