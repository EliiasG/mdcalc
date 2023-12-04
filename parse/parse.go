package parse

import (
	"fmt"
	"strings"

	"github.com/eliiasg/mdcalc/setup"
	"github.com/eliiasg/mdcalc/syntax"
)

// Terrible code
// parse mdcalc code: mdc is the code to be parsed, header is the title, and sub is the subproblem name where <n> will be replaced by the index
func Parse(mdc, header, sub string, lib syntax.UnitLibrary) (string, error) {
	env := setup.GenerateEnvironment(lib)
	var sb strings.Builder
	started := false
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
			if !started {
				sb.Reset()
				sb.WriteString("# ")
				sb.WriteString(header + "\n")
				started = true
			}
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
			err := env.WriteCalculation(content, &sb)
			if err != nil {
				fmt.Printf("error on line %v:\n", i+1)
				sb.WriteString(fmt.Sprintf("### <span style=\"color:red\">Error: %v</span>", err.Error()))
				fmt.Println(err.Error())
			}
		}
	}
	return sb.String(), nil
}

func err(ln int, msg string) error {
	return fmt.Errorf("error on line %v: %v", ln+1, msg)
}
