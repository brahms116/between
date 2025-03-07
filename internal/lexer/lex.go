package lex

import (
	"fmt"
)

type TokenType int

const (
	TOKEN_PRODUCT TokenType = iota
	TOKEN_SUM
	TOKEN_STR_SUM
	TOKEN_ID
	TOKEN_LITERAL
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_LIST
	TOKEN_SEPARATOR

	TOKEN_STR
	TOKEN_NUM
	TOKEN_OBJ

	TOKEN_OPTIONAL
)

var StringToToken map[string]TokenType = map[string]TokenType{
	"prod":   TOKEN_PRODUCT,
	"sum":    TOKEN_SUM,
	"strsum": TOKEN_STR_SUM,
	"Str":    TOKEN_STR,
	"Num":    TOKEN_STR_SUM,
	"Obj":    TOKEN_OBJ,
}

type Location struct {
	FilePos int
	Length  int
}

type Token struct {
	Type  TokenType
	Value string
	Loc   Location
}

type Lexer struct {
	input    string
	startPos int
	currPos  int
	tokens   []Token
}

func Lex(input string) ([]Token, error) {
	lexer := &Lexer{input: input}
	return lexer.Lex()
}

func (l *Lexer) Lex() ([]Token, error) {

	for {
		currChar := l.next()
		if currChar == nil {
			break
		}

		if isWhiteSpace(*currChar) {
			l.eatWhile(isWhiteSpace)
			l.startPos = l.currPos
			continue
		}

		switch *currChar {
		case ',':
			l.acceptToken(TOKEN_SEPARATOR)
			continue
		case '{':
			l.acceptToken(TOKEN_LBRACE)
			continue
		case '}':
			l.acceptToken(TOKEN_RBRACE)
			continue
		case '?':
			l.acceptToken(TOKEN_OPTIONAL)
			continue
		case '[':
			{
				currChar = l.next()
				if currChar == nil {
					return nil, fmt.Errorf("Unexpected EOF")
				}
				if *currChar != ']' {
					return nil, fmt.Errorf("Expected ']' at pos %d", l.currPos)
				}
				l.acceptToken(TOKEN_LIST)
				continue
			}
		case '"':
			err := l.lexLiteral()
			if err != nil {
				return nil, err
			}
			continue
		default:
		}

		if isAlpha(*currChar) {
			l.lexAlphaNum()
			continue
		}

		return nil, fmt.Errorf("Unexpected char %s at %d", string(*currChar), l.currPos)
	}
	return l.tokens, nil
}

func (l *Lexer) lexLiteral() error {
	l.eatWhile(func(b byte) bool {
		return b != '"'
	})
	next := l.next()
	if next == nil {
		return fmt.Errorf("Unexpected EOF, expected '\"' at %d", l.currPos)
	}
	str := l.currString()
	l.acceptTokenWithValue(TOKEN_LITERAL, str[1:len(str)-1])
	return nil
}

func (l *Lexer) lexWhitespace() {
	l.eatWhile(isWhiteSpace)
}

func (l *Lexer) lexAlphaNum() {
	l.eatWhile(isAlphaNum)
	str := l.currString()

	token, ok := StringToToken[str]
	if ok {
		l.acceptToken(token)
		return
	}
	l.acceptTokenWithValue(TOKEN_ID, str)
}

func (l *Lexer) acceptToken(tokenType TokenType) {
	l.acceptTokenWithValue(tokenType, "")
}

func (l *Lexer) acceptTokenWithValue(tokenType TokenType, value string) {
	length := l.currPos - l.startPos
	start := l.startPos
	token := Token{
		Type:  tokenType,
		Value: value,
		Loc: Location{
			FilePos: start,
			Length:  length,
		},
	}
	l.tokens = append(l.tokens, token)
	l.startPos = l.currPos
}

func (l *Lexer) next() *byte {
	if l.currPos >= len(l.input) {
		l.currPos++
		return nil
	}
	char := l.input[l.currPos]
	l.currPos++
	return &char
}

func (l *Lexer) backup() {
	if l.currPos > l.startPos {
		l.currPos--
	}
}

func (l *Lexer) currString() string {
	return l.input[l.startPos:l.currPos]
}

func (l *Lexer) eatWhile(fn func(byte) bool) {
	for {
		next := l.next()
		if next == nil {
			l.backup()
			return
		}

		if !fn(*next) {
			l.backup()
			return
		}
	}
}
