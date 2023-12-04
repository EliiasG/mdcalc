package setup

import (
	"errors"
	"math"

	"github.com/eliiasg/mdcalc/syntax"
)

func div(l, r float64) (float64, error) {
	if r == 0 {
		return math.NaN(), errors.New("divide by zero")
	}
	return l / r, nil
}

func genOperators() map[string]syntax.Operator {
	return map[string]syntax.Operator{
		"*": syntax.Operator{
			Execute: func(l, r float64) (float64, error) {
				return l * r, nil
			},
			Latex:            "@l*@r",
			ParenthesisLeft:  true,
			ParenthesisRight: true,
			OrderMatters:     false,
		},
		"/": syntax.Operator{
			Execute:          div,
			Latex:            "\\dfrac{@l}{@r}",
			ParenthesisLeft:  false,
			ParenthesisRight: false,
			OrderMatters:     true,
		},
		"//": syntax.Operator{
			Execute:          div,
			Latex:            "@l/@r",
			ParenthesisLeft:  true,
			ParenthesisRight: true,
			OrderMatters:     true,
		},
		"^": syntax.Operator{
			Execute: func(l, r float64) (float64, error) {
				return math.Pow(l, r), nil
			},
			Latex:            "@l^{@r}",
			ParenthesisLeft:  true,
			ParenthesisRight: false,
			OrderMatters:     true,
		},
		"+": syntax.Operator{
			Execute: func(l, r float64) (float64, error) {
				return l + r, nil
			},
			Latex:            "@l+@r",
			ParenthesisLeft:  true,
			ParenthesisRight: true,
			OrderMatters:     false,
		},
		"-": syntax.Operator{
			Execute: func(l, r float64) (float64, error) {
				return l - r, nil
			},
			Latex:            "@l-@r",
			ParenthesisLeft:  true,
			ParenthesisRight: true,
			OrderMatters:     true,
		},
	}
}
