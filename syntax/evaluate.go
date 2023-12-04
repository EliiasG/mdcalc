package syntax

import (
	"errors"
	"fmt"
	"math"
	"strconv"

	"github.com/eliiasg/mdcalc/util"
)

func (e *Environment) Evaluate(root ASTNode) (float64, error) {
	switch node := root.(type) {
	case *ASTComment:
		return e.Evaluate(node.Child)
	case *ASTUnitOverride:
		return e.Evaluate(node.Child)
	case *ASTVarSetter:
		if util.StrIsNum(node.VarName) {
			return math.NaN(), fmt.Errorf("cannot assign to number '%v'", node.VarName)
		}
		res, err := e.Evaluate(node.Child)
		if err != nil {
			return math.NaN(), err
		}
		e.VariableValues[node.VarName] = VariableValue{Value: res, Unit: e.GetUnit(node.Child)}
		return res, nil
	case *ASTLiteral:
		res, _, err := e.parseLiteral(node)
		if err != nil {
			return math.NaN(), err
		}
		return res, nil
	case *ASTOperator:
		op, ok := e.Operators[node.Operator]
		if !ok {
			return math.NaN(), fmt.Errorf("operator '%v' not defined", node.Operator)
		}
		resL, err := e.Evaluate(node.Left)
		if err != nil {
			return math.NaN(), err
		}
		resR, err := e.Evaluate(node.Right)
		if err != nil {
			return math.NaN(), err
		}
		res, err := op.Execute(resL, resR)
		if err != nil {
			return math.NaN(), fmt.Errorf("error on operator '%v': %v", node.Operator, err.Error())
		}
		return res, nil
	case *ASTFunction:
		fun, err := e.getFunction(node)
		if err != nil {
			return math.NaN(), err
		}
		evalRes := make([]float64, len(node.Params))
		for i, param := range node.Params {
			res, err := e.Evaluate(param)
			if err != nil {
				return math.NaN(), err
			}
			evalRes[i] = res
		}
		res, err := fun.Execute(evalRes)
		if err != nil {
			return math.NaN(), fmt.Errorf("error in function '%v': %v", node.Name, err.Error())
		}
		return res, nil
	}
	return math.NaN(), errors.New("invalid ast node")
}

func (e *Environment) getFunction(node *ASTFunction) (Function, error) {
	funs, ok := e.Functions[node.Name]
	if !ok {
		return Function{}, fmt.Errorf("function '%v' does not exist", node.Name)
	}
	fun, ok := funs[len(node.Params)]
	if !ok {
		return Function{}, fmt.Errorf("function '%v' did not expect %v parameter(s)", node.Name, len(node.Params))
	}
	return fun, nil
}

func (e *Environment) parseLiteral(node *ASTLiteral) (float64, string, error) {
	if util.StrIsNum(node.Value) {
		r, err := strconv.ParseFloat(node.Value, 64)
		if err != nil {
			return math.NaN(), "", fmt.Errorf("error while parsing number '%v'", node.Value)
		}
		return r, "", nil
	}
	val, ok := e.VariableValues[node.Value]
	if !ok {
		return math.NaN(), "", fmt.Errorf("variable '%v' undefined", node.Value)
	}
	return val.Value, val.Unit, nil
}
