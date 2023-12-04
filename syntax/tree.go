package syntax

import (
	"errors"
	"fmt"
)

type ASTNode interface {
	astVal()
}

func GenerateAst(code []Token) (ASTNode, error) {
	err := checkError(code)
	if err != nil {
		return nil, err
	}
	res, err := strip(code)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	return resolveExpression(code)
}

// modifies original
func ResolveOperatorChains(n ASTNode, values map[string]int) ASTNode {
	switch node := n.(type) {
	case *ASTUnitOverride:
		node.Child = ResolveOperatorChains(node.Child, values)
	case *ASTComment:
		node.Child = ResolveOperatorChains(node.Child, values)
	case *ASTVarSetter:
		node.Child = ResolveOperatorChains(node.Child, values)
	case *ASTOperator:
		node.Left = ResolveOperatorChains(node.Left, values)
		node.Right = ResolveOperatorChains(node.Right, values)
	case *ASTFunction:
		for i, param := range node.Params {
			node.Params[i] = ResolveOperatorChains(param, values)
		}
	case *ASTOperatorChain:
		r := generateInitalOperator(node)
		r = sortOperators(r, values)
		return ResolveOperatorChains(r, values)
	}
	return n
}

type astValImpl struct{}

func (a *astValImpl) astVal() {
	panic("should not be called")
}

// Number constant or variable reference
type ASTLiteral struct {
	astValImpl
	Value string
}

type ASTUnitOverride struct {
	astValImpl
	Unit  string
	Child ASTNode
}

type ASTComment struct {
	astValImpl
	Content string
	Child   ASTNode
}

type ASTVarSetter struct {
	astValImpl
	VarName string
	Child   ASTNode
}

// Only allows operators with a left and right (so no negate or not operator)
type ASTOperator struct {
	astValImpl
	Operator string
	Left     ASTNode
	Right    ASTNode
}

// For generating an initial tree where operator priorities are unknown
type ASTOperatorChain struct {
	astValImpl
	Operators []string
	// Should always be 1 longer than Operators
	Values []ASTNode
}

type ASTFunction struct {
	astValImpl
	Name   string
	Params []ASTNode
}

func resolveExpression(code []Token) (ASTNode, error) {
	res := &ASTOperatorChain{
		Operators: make([]string, 0),
		Values:    make([]ASTNode, 0),
	}
	var expr ASTNode
	i := 0
	for i < len(code) {
		switch tok := code[i].(type) {
		case TokenUnit:
			if expr == nil {
				return nil, fmt.Errorf("cannot have unit without expression, got unit '%v'", tok.Name)
			}
			expr = &ASTUnitOverride{
				Unit:  tok.Name,
				Child: expr,
			}
		case TokenLiteral:
			if expr != nil {
				return nil, fmt.Errorf("expected operator or ) after expression, got literal '%v'", tok.Value)
			}
			expr = &ASTLiteral{Value: tok.Value}
		case TokenParenthesis:
			if !tok.Opening {
				return nil, fmt.Errorf("unexpected ')'")
			}
			if expr != nil {
				return nil, errors.New("expected operator or ) after expression, got (")
			}
			next := closingIdx(code, i)
			if next == -1 {
				return nil, errors.New("expected )")
			}
			var err error
			expr, err = GenerateAst(code[i+1 : next])
			if err != nil {
				return nil, err
			}
			i = next
		case TokenFunc:
			if expr != nil {
				return nil, errors.New("expected operator or ) after expression, got (")
			}
			next := closingIdx(code, i+1)
			if next == -1 {
				return nil, errors.New("expected )")
			}
			var err error
			expr, err = resolveFunc(code[i : next+1])
			if err != nil {
				return nil, err
			}
			i = next
		case TokenOperator:
			res.Values = append(res.Values, expr)
			expr = nil
			res.Operators = append(res.Operators, tok.Operator)
		case TokenComment:
			return nil, errors.New("unexpected ':' comments can only be at end of code or parenthesis")
		case TokenVarSetter:
			return nil, errors.New("unexpected =")
		case TokenComma:
			return nil, errors.New("unexpected ,")
		}
		i++
	}
	if expr == nil {
		return nil, errors.New("expected expression")
	}
	if len(res.Values) == 0 {
		return expr, nil
	}
	res.Values = append(res.Values, expr)
	return res, nil
}

func resolveFunc(code []Token) (ASTNode, error) {
	fun, _ := code[0].(TokenFunc)
	if p, ok := code[len(code)-1].(TokenParenthesis); !ok || p.Opening {
		return nil, fmt.Errorf("missing close parenthesis for function '%v'", fun.Name)
	}
	funNode := &ASTFunction{
		Name:   fun.Name,
		Params: make([]ASTNode, 0),
	}
	if len(code) == 3 {
		return funNode, nil
	}
	startidx := 2
	for i, t := range code[1:] {
		if _, ok := t.(TokenComma); !ok {
			continue
		}
		ast, err := GenerateAst(code[startidx : i+1])
		if err != nil {
			return nil, err
		}
		startidx = i + 1
		funNode.Params = append(funNode.Params, ast)
	}
	return funNode, nil
}

// Only finds length error
func checkError(code []Token) error {
	if len(code) == 0 {
		return errors.New("code section must have a nonzero length")
	}
	return nil
}

func strip(code []Token) (ASTNode, error) {
	res, err := resolveVarSetter(code)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	res, err = resolveComment(code)
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res, nil
	}
	return nil, nil
}

// returns nil if no var setter
func resolveVarSetter(code []Token) (ASTNode, error) {
	if setter, ok := code[0].(TokenVarSetter); ok {
		res, err := GenerateAst(code[1:])
		if err != nil {
			return nil, err
		}
		return &ASTVarSetter{
			VarName: setter.VarName,
			Child:   res,
		}, nil
	}
	return nil, nil
}

// returns nil if no comment
func resolveComment(code []Token) (ASTNode, error) {
	idx := len(code) - 1
	if comment, ok := code[idx].(TokenComment); ok {
		res, err := GenerateAst(code[:idx])
		if err != nil {
			return nil, err
		}
		return &ASTComment{
			Content: comment.Content,
			Child:   res,
		}, nil
	}
	return nil, nil
}

func closingIdx(code []Token, startIdx int) int {
	val := 0
	for i := startIdx; i < len(code); i++ {
		tkn := code[i]
		switch p := tkn.(type) {
		case TokenParenthesis:
			if p.Opening {
				val++
			} else {
				val--
			}
		}
		if val == 0 {
			return i
		}
	}
	return -1
}

func sortOperators(n ASTNode, values map[string]int) ASTNode {
	opNode, ok := n.(*ASTOperator)
	if !ok {
		return n
	}
	opNext, ok := opNode.Left.(*ASTOperator)
	if !ok {
		return n
	}
	op := sortOperators(opNext, values)
	opNext = op.(*ASTOperator)
	if getValue(opNode.Operator, values) > getValue(opNext.Operator, values) {
		opNode.Left, opNext.Right = opNext.Right, opNode
		opNode, opNext = opNext, opNode
	} else {
		opNode.Left = opNext
	}
	return opNode
}

func getValue(op string, values map[string]int) int {
	val, ok := values[op]
	if !ok {
		return 0
	}
	return val
}

func generateInitalOperator(n *ASTOperatorChain) ASTNode {
	if len(n.Values) == 1 {
		return n.Values[0]
	}
	var root *ASTOperator
	var last *ASTOperator
	for i := range n.Operators {
		opIdx := len(n.Operators) - i - 1
		cur := &ASTOperator{
			Operator: n.Operators[opIdx],
			Right:    n.Values[opIdx+1],
		}
		if root == nil {
			root = cur
		} else {
			last.Left = cur
		}
		last = cur
	}
	last.Left = n.Values[0]
	return root
}
