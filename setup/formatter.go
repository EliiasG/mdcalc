package setup

import (
	"fmt"
	"math"
	"strings"
)

type formatter struct{}

func (f *formatter) FormatLine(expr string, res string) string {
	return expr + " &= " + res + "\\\\ \\\\ \n"
}

// precision -1 for number in expression
func (f *formatter) FormatNumber(num float64, precision int, unit string, comment string) string {
	if precision == -1 {
		precision = 10
	}
	amt := math.Pow10(precision)
	num = math.Round(num*amt) / amt
	if unit != "" {
		unit = " " + unit
	}
	if comment != "" {
		comment = fmt.Sprintf("\\textit{ (%v)}", comment)
	}
	fNum := strings.ReplaceAll(fmt.Sprintf("%v", num), ".", ",")
	return fmt.Sprintf("\\textbf{%v}\\text{\\scriptsize{%v}}%v", fNum, unit, comment)
}

func (f *formatter) FormatVar(name string) string {
	panic("not implemented") // TODO: Implement
}

func (f *formatter) FormatParenthesie(expr string) string {
	return "(" + expr + ")"
}
