package syntax

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
		if l == "" || l == r {
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

func (e *Environment) MakeLatexExpression(root ASTNode) (string, error) {
	switch node := root.(type) {
	case *ASTUnitOverride:
		return e.formatUnitOverride(node)
	case *ASTLiteral:
		val, unit, err := e.parseLiteral(node)
		if err != nil {
			return "", err
		}
		// only do comment on result
		return e.Formatter.FormatNumber(val, -1, unit, ""), nil
	case *ASTComment:
		res, err := e.Evaluate(root)
		if err != nil {
			return "", err
		}
		comment, precision, err := commentData(node.Content)
		if err != nil {
			return "", err
		}
		return e.Formatter.FormatNumber(res, precision, e.UnitLibrary.GetUnitDisplayName(e.GetUnit(node)), comment), nil
	case *ASTVarSetter:
		return e.MakeLatexExpression(node.Child)
	case *ASTOperator:
		return e.formatOperator(node)
	case *ASTFunction:
		return e.formatFunction(node)
	}
	return "", errors.New("invalid AST node")
}

func (e *Environment) MakeLatexCalculation(root ASTNode) (string, error) {
	// IMPORTANT evaluating first, since it might introduce new variables that could be needed for formatting
	res, err := e.Evaluate(root)
	if err != nil {
		return "", err
	}
	expr, err := e.MakeLatexExpression(root)
	if err != nil {
		return "", err
	}
	node, ok := root.(*ASTComment)
	comment := ""
	precision := 2
	if ok {
		comment, precision, err = commentData(node.Content)
		if err != nil {
			return "", err
		}
	}
	unit := e.UnitLibrary.GetUnitDisplayName(e.GetUnit(root))
	return e.Formatter.FormatLine(expr, e.Formatter.FormatNumber(res, precision, unit, comment)), nil
}

func (e *Environment) MakeMultilineCalculation(root ASTNode) ([]string, error) {
	switch node := root.(type) {
	case *ASTUnitOverride:
		return e.MakeMultilineCalculation(node.Child)
	case *ASTLiteral:
		return nil, nil
	case *ASTComment:
		r, err := e.MakeMultilineCalculation(node.Child)
		if err != nil {
			return nil, err
		}
		line, err := e.MakeLatexCalculation(node.Child)
		if err != nil {
			return nil, err
		}
		if r == nil {
			return []string{line}, nil
		}
		return append(r, line), nil
	case *ASTVarSetter:
		return e.MakeMultilineCalculation(node.Child)
	case *ASTOperator:
		l, err := e.MakeMultilineCalculation(root)
		if err != nil {
			return nil, err
		}
		r, err := e.MakeMultilineCalculation(root)
		if err != nil {
			return nil, err
		}
		if l == nil {
			return r, nil
		} else if r == nil {
			return l, nil
		}
		return append(l, r...), nil
	case *ASTFunction:
		res := make([]string, 0)
		for _, param := range node.Params {
			lines, err := e.MakeMultilineCalculation(param)
			if err != nil {
				return nil, err
			}
			if lines != nil {
				res = append(res, lines...)
			}
		}
		if len(res) == 0 {
			return nil, nil
		}
		return res, nil
	}
	return nil, errors.New("invalid AST node")
}

func (e *Environment) formatFunction(node *ASTFunction) (string, error) {
	fun, err := e.getFunction(node)
	if err != nil {
		return "", err
	}
	replacements := make([]string, 0, len(node.Params)*2)
	for i, param := range node.Params {
		fParam, err := e.MakeLatexExpression(param)
		if err != nil {
			return "", err
		}
		replacements = append(replacements, "@"+fmt.Sprint(i+1))
		replacements = append(replacements, fParam)
	}
	return strings.NewReplacer(replacements...).Replace(fun.Latex), nil
}

func (e *Environment) formatOperator(node *ASTOperator) (string, error) {
	op, ok := e.Operators[node.Operator]
	if !ok {
		return "", fmt.Errorf("operator '%v' is invalid", node.Operator)
	}
	lRes, err := e.MakeLatexExpression(node.Left)
	if err != nil {
		return "", err
	}
	rRes, err := e.MakeLatexExpression(node.Left)
	if err != nil {
		return "", err
	}
	if op.Parenthesis {
		l, r := e.needParenthesis(node)
		if l {
			lRes = e.Formatter.FormatParenthesie(lRes)
		}
		if r {
			rRes = e.Formatter.FormatParenthesie(rRes)
		}
	}
	return strings.ReplaceAll(strings.ReplaceAll(op.Latex, "@l", lRes), "@r", rRes), nil
}

func (e *Environment) formatUnitOverride(node *ASTUnitOverride) (string, error) {
	literal, ok := node.Child.(*ASTLiteral)
	if !ok {
		return e.MakeLatexExpression(node.Child)
	}
	val, _, err := e.parseLiteral(literal)
	if err != nil {
		return "", err
	}
	// only do comment on result
	return e.Formatter.FormatNumber(val, -1, e.UnitLibrary.GetUnitDisplayName(node.Unit), ""), nil
}

func (e *Environment) needParenthesis(op *ASTOperator) (left bool, right bool) {
	lOp, ok := op.Left.(*ASTOperator)
	if ok {
		left = getValue(lOp.Operator, e.OperatorPowers) > getValue(op.Operator, e.OperatorPowers)
	}
	rOp, ok := op.Left.(*ASTOperator)
	if ok {
		right = getValue(rOp.Operator, e.OperatorPowers) >= getValue(op.Operator, e.OperatorPowers)
	}
	return
}

func commentData(comment string) (string, int, error) {
	split := strings.Index(comment, ":")
	if split == -1 {
		return comment, 2, nil
	}
	amt := strings.TrimSpace(comment[:split])
	p, err := strconv.ParseInt(amt, 10, 32)
	if err != nil {
		return "", -1, fmt.Errorf("error while parsing precision '%v'", amt)
	}
	return comment[split+1:], int(p), nil
}
