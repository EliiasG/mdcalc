package parse

import (
	"fmt"
	"strings"

	"github.com/eliiasg/mdcalc/syntax"
)

// Terrible code
// parse mdcalc code: mdc is the code to be parsed, header is the title, and sub is the subproblem name where <n> will be replaced by the index
func Parse(mdc, header, sub string) (string, error) {
	var sb strings.Builder
	sb.WriteString("# ")
	sb.WriteString(header)
	n := 1
	for i, line := range strings.Split(mdc, "\n") {
		sb.WriteString("  \n")
		if len(line) == 0 {
			continue
		}
		var content string
		if len(line) >= 3 {
			content = strings.TrimSpace(line[2:])
		} else {
			content = ""
		}
		switch line[0] {
		default:
			return "", err(i, "Every line must start with either T, C or |")
		case '|':
			sb.WriteString("### ")
			sb.WriteString(strings.ReplaceAll(sub, "<n>", fmt.Sprint(n)))
			n++
		case 'T':
			if len(line) < 3 {
				return "", err(i, "Text lines must start with a T followed by a space followed by text")
			}
			sb.WriteString(content)
		case 'I':
			sb.WriteString(fmt.Sprintf("![Image!](%v)", content))
		case 'C':
			sb.WriteString("$$\\dfrac{\\text{5746 kr. i timen}}{1}=\\text{1 kr. i timen \\textit{(efter AM bidrag)}}$$")
			r := syntax.Tokenize(content)
			tree, err := syntax.GenerateAst(r)
			tree = syntax.ResolveOperatorChains(tree, map[string]int{
				"*": 1,
				"/": 1,
				"^": 2,
			})
			if err != nil {
				panic(err)
			}
			fmt.Println(tree)
		}
	}
	return sb.String(), nil
}

func err(ln int, msg string) error {
	return fmt.Errorf("error on line %v: %v", ln+1, msg)
}
