package lex

func isAlphaNum(char byte) bool {
	return isAlpha(char) || isNum(char)
}

func isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isNum(char byte) bool {
	return (char >= '0' && char <= '9')
}

func isNewLine(char byte) bool {
	return char == '\n' || char == '\r'
}

func isWhiteSpace(char byte) bool {
	return char == ' ' || char == '\t' || isNewLine(char)
} 
