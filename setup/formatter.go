package setup

import (
	"fmt"
	"math"
)

type Formatter struct{}

func (f *Formatter) FormatLine(expr string, res string) string {
	return expr + " &= " + res + "\\\\ \\\\ \n"
}

// precision -1 for number in expression
func (f *Formatter) FormatNumber(num float64, precision int, unit string, comment string) string {
	if precision == -1 {
		precision = 10
	}
	amt := math.Pow10(precision)
	num = math.Round(num*amt)/amt
	if unit != "" {
		unit = " " + unit
	}
	if comment != "" {
		comment = fmt.Sprintf(" \\textit{(%v)}", comment)
	}
	return fmt.Sprintf("\\text{%v%v%v}", num, unit, comment)
}

func (f *Formatter) FormatVar(name string) string {
	panic("not implemented") // TODO: Implement
}

func (f *Formatter) FormatParenthesie(expr string) string {
	return "(" + expr + ")"
}
