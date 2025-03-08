package lex

import "fmt"

var TokenTypeDisplay map[TokenType]string = map[TokenType]string{
	TOKEN_PRODUCT:   "TOKEN_PRODUCT",
	TOKEN_SUM:       "TOKEN_SUM",
	TOKEN_SUM_STR:   "TOKEN_SUM_STR",
	TOKEN_ID:        "TOKEN_ID",
	TOKEN_LITERAL:   "TOKEN_LITERAL",
	TOKEN_LBRACE:    "TOKEN_LBRACE",
	TOKEN_RBRACE:    "TOKEN_RBRACE",
	TOKEN_LIST:      "TOKEN_LIST",
	TOKEN_SEPARATOR: "TOKEN_SEPARATOR",
	TOKEN_OPTIONAL:  "TOKEN_OPTIONAL",
}

func (t Token) String() string {
	if t.Value == "" {
		return TokenTypeDisplay[t.Type]
	}
	return fmt.Sprintf("%s(%s)", TokenTypeDisplay[t.Type], t.Value)
}
