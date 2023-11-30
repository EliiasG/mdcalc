package util

func IsNum(c rune) bool {
	return c == '.' || (c >= '0' && c <= '9')
}

func IsAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func StrIsAlpha(s string) bool {
	for _, c := range s {
		if !IsAlpha(c) {
			return false
		}
	}
	return true
}

func StrIsNum(s string) bool {
	for _, c := range s {
		if !IsNum(c) {
			return false
		}
	}
	return true
}
