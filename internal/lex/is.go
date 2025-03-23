package lex

func isAlphaNum(r rune) bool {
	return isAlpha(r) || isNum(r)
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isNum(r rune) bool {
	return (r >= '0' && r <= '9')
}

func isNewLine(r rune) bool {
	return r == '\n' || r == '\r'
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\t' || isNewLine(r)
}
