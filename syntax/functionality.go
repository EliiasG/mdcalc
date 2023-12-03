package syntax

type Function struct {
	Execute func([]float64) (float64, error)
	// Use @i where 'i' for the formatted parameter staring at i = 0
	Latex string
}

// Only operators that expect 2 arguments are supported.
// For anything else just use a function that formats to an operator.
type Operator struct {
	Execute func(float64, float64) (float64, error)
	// Use @l and @r for the formatted left and right parameters.
	Latex string
	// Add parenthesis if necessary, should only be false for something like a fraction line
	Parenthesis bool
	// Useful for some unit management stuff
	OrderMatters bool
}

type VariableValue struct {
	Value float64
	Unit  string
}

type UnitLibrary interface {
	GetUnitDisplayName(unit string) string
	GetOperatorResult(left, right, operator string, orderMatters bool) string
}

type Formatter interface {
	FormatLine(expr, res string) string
	// precision -1 for number in expression
	FormatNumber(num float64, precision int, unit, comment string) string
	FormatVar(name string) string
	FormatParenthesie(expr string) string
}

type Environment struct {
	Operators map[string]Operator
	// Name, Param amount
	Functions      map[string]map[int]Function
	VariableValues map[string]VariableValue
	OperatorPowers map[string]int
	Formatter      Formatter
	UnitLibrary    UnitLibrary
}
