package unitlib

type SimpleUnitLibrary struct{}

func (l *SimpleUnitLibrary) GetUnitDisplayName(unit string) string {
	return unit
}

func (l *SimpleUnitLibrary) GetOperatorResult(left string, right string, operator string, orderMatters bool) string {
	if operator == "%" {
		operator = "/"
	}
	return left + operator + right
}
