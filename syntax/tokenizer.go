package syntax

import (
	"strings"

	"github.com/eliiasg/mdcalc/util"
)

type Token interface {
	token()
}

type tokenImpl struct{}

func (t tokenImpl) token() {
	panic("should not be called")
}

type TokenComma struct {
	tokenImpl
}

type TokenParenthesis struct {
	tokenImpl
	// closing if false
	Opening bool
}

// will be followed by start parenthesis, and closed with closing parenthesis
type TokenFunc struct {
	tokenImpl
	Name string
}

// constant number or variable ref
type TokenLiteral struct {
	tokenImpl
	Value string
}

type TokenUnit struct {
	tokenImpl
	Name string
}

type TokenOperator struct {
	tokenImpl
	Operator string
}

type TokenComment struct {
	tokenImpl
	Content string
}

type TokenVarSetter struct {
	tokenImpl
	VarName string
}

type tokenizerState struct {
	res             []Token
	wasNum          bool
	wasAlpha        bool
	readyForUnit    bool
	handlingComment bool
	curRes          strings.Builder
}

func Tokenize(prgm string) []Token {
	state := &tokenizerState{
		res: make([]Token, 0),
	}
	for _, c := range prgm + " " {
		handleNum(state, c)
		handleUnit(state, c)
		handleVarRef(state, c)
		handleOperators(state, c)
		handleVarAssign(state, c)
		handleComma(state, c)
		// single responsibility principle in action
		handleParenthesesAndComments(state, c)
		state.wasNum = util.IsNum(c)
		state.wasAlpha = util.IsAlpha(c)
	}
	return state.res
}

// All the following functions are a mess of spaghetti and side effects

func handleComma(s *tokenizerState, c rune) {
	if !s.handlingComment && c == ',' {
		s.res = append(s.res, TokenComma{})
		s.readyForUnit = false
	}
}

func handleVarAssign(s *tokenizerState, c rune) {
	if c != '=' || len(s.res) == 0 || s.handlingComment {
		return
	}
	elem := s.res[len(s.res)-1]
	switch t := elem.(type) {
	case TokenLiteral:
		s.res[len(s.res)-1] = TokenVarSetter{VarName: t.Value}
		s.readyForUnit = false
	}
}

func handleOperators(s *tokenizerState, c rune) {
	if !s.readyForUnit || c == ' ' || c == '=' || c == ')' || c == ':' || c == ',' || util.IsAlpha(c) {
		return
	}
	switch t := s.res[len(s.res)-1].(type) {
	case TokenUnit, TokenLiteral:
		addOperator(s, c)
	case TokenParenthesis:
		if !t.Opening {
			addOperator(s, c)
		}
	}
}

func addOperator(s *tokenizerState, c rune) {
	s.res = append(s.res, TokenOperator{Operator: string(c)})
	s.readyForUnit = false
}

func handleParenthesesAndComments(s *tokenizerState, c rune) {
	if c == '(' {
		if !s.readyForUnit && !s.handlingComment && s.curRes.Len() > 0 {
			s.res = append(s.res, TokenFunc{Name: s.curRes.String()})
			s.curRes.Reset()
		}
		s.res = append(s.res, TokenParenthesis{Opening: true})
	} else if c == ')' {
		s.readyForUnit = true
		if s.handlingComment {
			s.handlingComment = false
			s.res = append(s.res, TokenComment{Content: strings.TrimSpace(s.curRes.String())})
			s.curRes.Reset()
		}
		s.res = append(s.res, TokenParenthesis{Opening: false})
	} else if c == ':' {
		s.curRes.Reset()
		s.handlingComment = true
		s.readyForUnit = false
		return
	}
	if s.handlingComment {
		s.curRes.WriteRune(c)
	}
}

func handleVarRef(s *tokenizerState, c rune) {
	if s.readyForUnit || s.handlingComment {
		return
	}
	if util.IsAlpha(c) {
		s.curRes.WriteRune(c)
	} else if s.curRes.Len() > 0 && c != '(' && !util.IsNum(c) {
		s.res = append(s.res, TokenLiteral{Value: s.curRes.String()})
		s.curRes.Reset()
		s.readyForUnit = true
	}
}

func handleUnit(s *tokenizerState, c rune) {
	if !s.readyForUnit || s.handlingComment {
		return
	}
	if util.IsAlpha(c) {
		s.curRes.WriteRune(c)
	} else if s.wasAlpha {
		s.res = append(s.res, TokenUnit{Name: s.curRes.String()})
		s.curRes.Reset()
	}
}

func handleNum(s *tokenizerState, c rune) {
	if s.handlingComment {
		return
	}
	num := util.IsNum(c)
	if num {
		s.readyForUnit = false
		s.curRes.WriteRune(c)
	} else if s.wasNum {
		s.res = append(s.res, TokenLiteral{Value: s.curRes.String()})
		s.curRes.Reset()
		s.readyForUnit = true
	}
}

// numeric or .
