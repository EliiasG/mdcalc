package syntax

import "github.com/eliiasg/mdcalc/util"

func (e *Environment) GetUnit(root ASTNode) string {
	switch node := root.(type) {
	case *ASTUnitOverride:
		if node.Unit == "None" {
			return ""
		}
		return node.Unit
	case *ASTLiteral:
		_, unit, err := e.parseLiteral(node)
		if err != nil {
			return ""
		}
		return unit
	case *ASTComment:
		return e.GetUnit(node.Child)
	case *ASTVarSetter:
		return e.GetUnit(node.Child)
	case *ASTOperator:
		l := e.GetUnit(node.Left)
		r := e.GetUnit(node.Left)
		if l == "" {
			return r
		}
		if r == "" {
			return l
		}
		op, ok := e.Operators[node.Operator]
		if !ok {
			return ""
		}
		return e.UnitLibrary.GetOperatorResult(l, r, node.Operator, op.OrderMatters)
	case *ASTFunction:
		return ""
	}
	return ""
}

func (e *Environment) MakeLatexExpression(root ASTNode, function bool) (string, error) {
	switch node := root.(type) {
	case *ASTUnitOverride:
		literal, ok := node.Child.(*ASTLiteral)
		if !ok {
			return e.MakeLatexExpression(node.Child, function)
		}
		val, _, err := e.parseLiteral(literal)
		if function && !util.StrIsNum(literal.Value) {
			return e.Formatter.FormatVar(literal.Value), nil
		}
		if err != nil {
			return "", err
		}
		// only do comment on result
		return e.Formatter.FormatNumber(val, -1, node.Unit, ""), nil
	case *ASTLiteral:
		val, unit, err := e.parseLiteral(node)
		if function && !util.StrIsNum(node.Value) {
			return e.Formatter.FormatVar(node.Value), nil
		}
		if err != nil {
			return "", err
		}
		// only do comment on result
		return e.Formatter.FormatNumber(val, -1, unit, ""), nil
	case *ASTComment:
		return e.MakeLatexExpression(node.Child)
	case *ASTVarSetter:
		return e.MakeLatexExpression(node.Child)
	case *ASTOperator:
		res, ok := e.Operators[node.Operator]
		if !ok {
			return ""
		}
	case *ASTFunction:
		return ""
	}
}

func (e *Environment) formatLiteral(node *ASTLiteral, function, overrideUnit bool, unit string) (string, error) {
	val, unit, err := e.parseLiteral(node)
	if function && !util.StrIsNum(node.Value) {
		return e.Formatter.FormatVar(node.Value), nil
	}
	if err != nil {
		return "", err
	}
	// only do comment on result
	if overrideUnit {

	}
	return e.Formatter.FormatNumber(val, -1, unit, ""), nil
}
