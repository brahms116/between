package lex

import (
	"fmt"
	"unicode/utf8"
)

type TokenType int

const (
	TOKEN_PRODUCT TokenType = iota
	TOKEN_SUM
	TOKEN_SUM_STR
	TOKEN_ID
	TOKEN_LITERAL
	TOKEN_LBRACE
	TOKEN_RBRACE
	TOKEN_LIST
	TOKEN_SEPARATOR
	TOKEN_OPTIONAL
	TOKEN_EOF
)

var stringToToken map[string]TokenType = map[string]TokenType{
	"prod":   TOKEN_PRODUCT,
	"sum":    TOKEN_SUM,
	"sumstr": TOKEN_SUM_STR,
}

type Location struct {
	ByteStart int
	ByteEnd   int
	Start     Point
	End       Point
}

type Point struct {
	Row int
	Col int
}

func (p Point) String() string {
	return fmt.Sprintf("(Row %d, Col %d)", p.Row+1, p.Col+1)
}

type Token struct {
	Type  TokenType
	Value string
	Loc   Location
	IsErr bool
}

type lexer struct {
	input    string
	startPos int
	currPos  int

	startPt Point
	currPt  Point

	// Should be okay because we don't backup more than 1 rune ever
	lastRowEnd int

	tokens []Token
	errs   []error
}

func Lex(input string) ([]Token, []error) {
	lexer := &lexer{input: input}
	lexer.Lex()
	return lexer.tokens, lexer.errs
}

func (l *lexer) err(err error) {
	l.errs = append(l.errs, err)
}

func (l *lexer) Lex() {

	for {
		currChar := l.next()
		if currChar == nil {
			break
		}

		if isWhiteSpace(*currChar) {
			l.eatWhile(isWhiteSpace)
			l.updateStart()
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
					l.err(fmt.Errorf("Unexpected EOF, expected ']' at %d", l.currPos))
					continue
				}
				if *currChar != ']' {
					l.err(fmt.Errorf("Unexpected char %s, expected ']' at %d", string(*currChar), l.currPos))
					continue
				}
				l.acceptToken(TOKEN_LIST)
				continue
			}
		case '"':
			l.lexLiteral()
			continue
		default:
		}

		if isAlpha(*currChar) {
			l.lexAlphaNum()
			continue
		}
		l.err(fmt.Errorf("Unexpected char %s at %d", string(*currChar), l.currPos))
	}

	l.acceptToken(TOKEN_EOF)
}

func (l *lexer) lexLiteral() {
	l.eatWhile(func(b rune) bool {
		return b != '"'
	})
	next := l.next()
	if next == nil {
		l.err(fmt.Errorf("Unexpected EOF, expected '\"' at %d", l.currPos))
	}
	str := l.currString()
	l.acceptTokenWithValue(TOKEN_LITERAL, str[1:len(str)-1])
}

func (l *lexer) lexWhitespace() {
	l.eatWhile(isWhiteSpace)
}

func (l *lexer) lexAlphaNum() {
	l.eatWhile(isAlphaNum)
	str := l.currString()

	token, ok := stringToToken[str]
	if ok {
		l.acceptToken(token)
		return
	}
	l.acceptTokenWithValue(TOKEN_ID, str)
}

func (l *lexer) acceptToken(tokenType TokenType) {
	l.acceptTokenWithValue(tokenType, "")
}

func (l *lexer) acceptTokenWithValue(tokenType TokenType, value string) {
	token := Token{
		Type:  tokenType,
		Value: value,
		Loc: Location{
			ByteStart: l.startPos,
			ByteEnd:   l.currPos,
			Start:     l.startPt,
			End:       l.currPt,
		},
	}
	l.tokens = append(l.tokens, token)
	l.updateStart()
}

func (l *lexer) updateStart() {
	l.startPos = l.currPos
	l.startPt = l.currPt
}

func (l *lexer) next() *rune {
	if l.currPos >= len(l.input) {
		return nil
	}

	r, w := utf8.DecodeRuneInString(l.input[l.currPos:])
	if w == 0 {
		return nil
	}
	l.currPos += w

	if r == '\n' {
		l.lastRowEnd = l.currPt.Col
		l.currPt.Col = 0
		l.currPt.Row++
	} else {
		l.currPt.Col++
	}

	return &r
}

func (l *lexer) backup() {
	if l.currPos > l.startPos {
		l.currPos--
		if l.currPt.Col == 0 {
			l.currPt.Row--
			l.currPt.Col = l.lastRowEnd
		} else {
			l.currPt.Col--
		}
	}
}

func (l *lexer) currString() string {
	return l.input[l.startPos:l.currPos]
}

func (l *lexer) eatWhile(fn func(rune) bool) {
	for {
		next := l.next()
		if next == nil {
			return
		}

		if !fn(*next) {
			l.backup()
			return
		}
	}
}
