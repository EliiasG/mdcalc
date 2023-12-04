package unitlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Operation struct {
	LeftUnit, RightUnit string
	Operator            string
}

const (
	namesFileName      = "units.txt"
	operationsFileName = "operators.txt"
)

type SavedUnitLibrary struct {
	names          map[string]string
	operations     map[Operation]string
	namesPath      string
	operationsPath string
}

func NewSavedUnitLib(dir string) (*SavedUnitLibrary, error) {
	bytes, err := os.ReadFile(dir + "/" + namesFileName)
	var names map[string]string
	if err == nil {
		names, err = loadNames(string(bytes))
		if err != nil {
			return nil, err
		}
	} else {
		names = make(map[string]string)
	}
	bytes, err = os.ReadFile(dir + "/" + operationsFileName)
	var operations map[Operation]string
	if err == nil {
		operations, err = loadOperations(string(bytes))
		if err != nil {
			return nil, err
		}
	} else {
		operations = make(map[Operation]string)
	}
	return &SavedUnitLibrary{
		names:          names,
		operations:     operations,
		namesPath:      dir + "/" + namesFileName,
		operationsPath: dir + "/" + operationsFileName,
	}, nil
}

func loadNames(names string) (map[string]string, error) {
	res := make(map[string]string)
	for i, line := range strings.Split(names, "\n") {
		if line == "" {
			continue
		}
		dat := strings.SplitN(line, " ", 2)
		if len(dat) != 2 {
			return nil, fmt.Errorf("line %v of unit file is invalid", i+1)
		}
		res[dat[0]] = dat[1]
	}
	return res, nil
}

func loadOperations(operations string) (map[Operation]string, error) {
	res := make(map[Operation]string)
	for i, line := range strings.Split(operations, "\n") {
		if line == "" {
			continue
		}
		dat := strings.Split(line, " ")
		if len(dat) != 4 {
			return nil, fmt.Errorf("line %v of operations file is invalid", i+1)
		}
		res[Operation{dat[0], dat[2], dat[1]}] = dat[3]
	}
	return res, nil
}

func (l *SavedUnitLibrary) GetUnitDisplayName(unit string) string {
	if unit == "" {
		return ""
	}
	res, ok := l.names[unit]
	if ok {
		return res
	}
	res = prompt(fmt.Sprintf("name unit '%v': ", unit))
	l.names[unit] = res
	append(l.namesPath, unit+" "+res)
	return res
}

func (l *SavedUnitLibrary) GetOperatorResult(left string, right string, operator string, orderMatters bool) string {
	if operator == "%" {
		operator = "/"
	}
	res, ok := l.operations[Operation{left, right, operator}]
	if ok {
		return res
	}
	if !orderMatters {
		res, ok = l.operations[Operation{right, left, operator}]
		if ok {
			return res
		}
	}
	res = prompt(fmt.Sprintf("determine unit result of '%v' %v '%v': ", left, operator, right))
	l.operations[Operation{left, right, operator}] = res
	append(l.operationsPath, left+" "+operator+" "+right+" "+res)
	return res
}

func append(filePath, line string) {
	f, err := os.OpenFile(filePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error while openening file %v:\n", filePath)
		fmt.Println(err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		fmt.Printf("error while appending to file %v:\n", filePath)
		fmt.Println(err)
		return
	}
}

func prompt(message string) string {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	res, _ := reader.ReadString('\n')
	return strings.TrimSpace(res)
}
