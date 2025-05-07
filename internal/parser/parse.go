package parser

import (
	"fmt"

	"github.com/brahms116/between/internal/lex"
	"github.com/brahms116/between/internal/st"
)

type parser struct {
	input  []lex.Token
	pos    int
	errors []error

	isEofError bool
}

const NUM_TYPE_STRING = "Num"
const STR_TYPE_STRING = "Str"
const OBJ_TYPE_STRING = "Obj"

func UnexpectedTokenError(token lex.Token) error {
	return fmt.Errorf("Unexpected token %s at %s", token.String(), token.Loc.Start.String())
}

func ExpectedTokenError(expected []lex.TokenType, token lex.Token) error {
	expectedStr := ""
	for i, t := range expected {
		if i > 0 {
			expectedStr += " or "
		}
		expectedStr += t.String()
	}

	return fmt.Errorf("Expected %s at %s, got %s", expectedStr, token.Loc.Start.String(), token.String())
}

func EOFError() error {
	return fmt.Errorf("Unexpected end of file")
}

func LexAndParse(input string) ([]st.Definition, []error) {
	tokens, lexErrs := lex.Lex(input)
	d, parseErrs := Parse(tokens)
	var errs []error
	errs = append(errs, lexErrs...)
	errs = append(errs, parseErrs...)
	return d, errs
}

func Parse(input []lex.Token) ([]st.Definition, []error) {
	p := &parser{input: input}
	definitions := p.parse()
	return definitions, p.errors
}

func (p *parser) parse() []st.Definition {
	definitions := []st.Definition{}
	lastDefinitionOk := true
	for {
		if lastDefinitionOk {
			currToken, ok := p.eatUntilOneOf([]lex.TokenType{
				lex.TOKEN_PRODUCT,
				lex.TOKEN_SUM,
				lex.TOKEN_SUM_STR,
				lex.TOKEN_EOF,
			}, true)

			if !ok {
				lastDefinitionOk = false
				continue
			}

			if currToken.Type == lex.TOKEN_EOF {
				return definitions
			}
		} else {
			currToken, _ := p.eatUntilOneOf([]lex.TokenType{
				lex.TOKEN_PRODUCT,
				lex.TOKEN_SUM,
				lex.TOKEN_SUM_STR,
				lex.TOKEN_EOF,
			}, false)

			if currToken.Type == lex.TOKEN_EOF {
				return definitions
			}
		}
		definition, ok := p.parseDefinition()
		lastDefinitionOk = ok
		definitions = append(definitions, definition)
	}
}

func (p *parser) appendErr(err error) {
	p.errors = append(p.errors, err)
}

func (p *parser) advance(tokenTypes []lex.TokenType, shouldFailImmediately, shouldConsumeToken bool) (lex.Token, bool) {
	currToken := p.currToken()
	found := false
	for _, tokenType := range tokenTypes {
		if currToken.Type == tokenType {
			found = true
			break
		}
	}

	if found {
		if shouldConsumeToken {
			p.pos++
		}
		return currToken, true
	}

	if currToken.Type == lex.TOKEN_EOF {
		p.appendErr(EOFError())
		return lex.Token{IsErr: true}, false
	}

	if shouldFailImmediately {
		p.appendErr(ExpectedTokenError(tokenTypes, currToken))
		return lex.Token{IsErr: true}, false
	}
	p.pos++
	return p.advance(tokenTypes, shouldFailImmediately, shouldConsumeToken)
}

func (p *parser) expect(tokenType lex.TokenType, follow []lex.TokenType) lex.Token {
	token, _ := p.expectOneOf([]lex.TokenType{tokenType}, follow)
	return token
}

func (p *parser) expectOneOf(tokenTypes []lex.TokenType, follow []lex.TokenType) (lex.Token, bool) {
	currToken := p.currToken()
	for _, tokenType := range tokenTypes {
		if currToken.Type == tokenType {
			p.pos++
			return currToken, true
		}
	}

	if !p.isEofError {
		p.appendErr(ExpectedTokenError(tokenTypes, currToken))
	}

	for {
		for _, tokenType := range follow {
			if currToken.Type == tokenType {
				return lex.Token{IsErr: true}, false
			}
		}
		if currToken.Type == lex.TOKEN_EOF {
			p.isEofError = true
			return lex.Token{IsErr: true}, false
		}
		p.pos++
	}
}

func (p *parser) expectToken(tokenType lex.TokenType, shouldImmediatelyFail bool) (lex.Token, bool) {
	return p.advance([]lex.TokenType{tokenType}, shouldImmediatelyFail, true)
}

func (p *parser) eatUntilOneOf(tokenTypes []lex.TokenType, shouldImmediatelyFail bool) (lex.Token, bool) {
	return p.advance(tokenTypes, shouldImmediatelyFail, false)
}

func (p *parser) optionalNextToken(tokenType lex.TokenType) (lex.Token, bool) {
	currToken := p.currToken()
	if currToken.Type == tokenType {
		p.pos++
		return currToken, true
	}
	return lex.Token{}, false
}

func (p *parser) parseDefinition() (st.Definition, bool) {

  currToken, ok := p.expectOneOf([]lex.TokenType{
    lex.TOKEN_PRODUCT,
    lex.TOKEN_SUM,
    lex.TOKEN_SUM_STR,
  }, []lex.TokenType{
    lex.TOKEN_EOF,
    lex.TOKEN_PRODUCT,
    lex.TOKEN_SUM,
    lex.TOKEN_SUM_STR,
  })

	if !ok {
		return st.Definition{}, false
	}

	switch currToken.Type {
	case lex.TOKEN_PRODUCT:
		prod, ok := p.parseProduct()
		return st.Definition{Product: &prod}, ok
	case lex.TOKEN_SUM:
		sum, ok := p.parseSum()
		return st.Definition{Sum: &sum}, ok
	case lex.TOKEN_SUM_STR:
		sumStr, ok := p.parseSumStr()
		return st.Definition{SumStr: &sumStr}, ok
	default:
		panic("Unreachable")
	}
}

func (p *parser) parseSumStr() (st.SumStr, bool) {
	keyword, ok := p.expectToken(lex.TOKEN_SUM_STR, true)
	id, ok := p.expectToken(lex.TOKEN_ID, ok)
	lBrace, ok := p.expectToken(lex.TOKEN_LBRACE, ok)
	variants, ok := p.parseSumStrVariants()
	rBrace, ok := p.expectToken(lex.TOKEN_RBRACE, ok)

	return st.SumStr{
		Keyword:    keyword,
		Id:         id,
		LeftBrace:  lBrace,
		Variants:   variants,
		RightBrace: rBrace,
	}, ok

}

func (p *parser) parseSumStrVariants() ([]st.SumStrVariant, bool) {
	variants := []st.SumStrVariant{}
	lastVariantOk := false
	for {
		if lastVariantOk {
			if p.currToken().Type != lex.TOKEN_ID {
				return variants, p.currToken().Type == lex.TOKEN_RBRACE
			}
		} else {
			currentToken, ok := p.eatUntilOneOf([]lex.TokenType{
				lex.TOKEN_ID,
				lex.TOKEN_LBRACE,
			}, false)
			if !ok {
				return variants, false
			}
			if currentToken.Type == lex.TOKEN_LBRACE {
				return variants, true
			}
		}

		id, ok := p.expectToken(lex.TOKEN_ID, true)
		jsonName := p.parseJsonRename()

		separator, ok := p.expectToken(lex.TOKEN_SEPARATOR, ok)
		lastVariantOk = ok

		variants = append(variants, st.SumStrVariant{
			Id:        id,
			JsonName:  jsonName,
			Separator: separator,
		})
	}
}

func (p *parser) parseSum() (st.Sum, bool) {
	keyword, ok := p.expectToken(lex.TOKEN_SUM, true)
	id, ok := p.expectToken(lex.TOKEN_ID, ok)
	lBrace, ok := p.expectToken(lex.TOKEN_LBRACE, ok)
	fields, ok := p.parseFields()
	rBrace, ok := p.expectToken(lex.TOKEN_RBRACE, ok)
	return st.Sum{
		Keyword:    keyword,
		Id:         id,
		LeftBrace:  lBrace,
		Variants:   fields,
		RightBrace: rBrace,
	}, ok
}

func (p *parser) parseProduct() (st.Product, bool) {
	keyword, ok := p.expectToken(lex.TOKEN_PRODUCT, true)
	id, ok := p.expectToken(lex.TOKEN_ID, ok)
	lBrace, ok := p.expectToken(lex.TOKEN_LBRACE, ok)
	fields, ok := p.parseFields()
	rBrace, ok := p.expectToken(lex.TOKEN_RBRACE, ok)
	return st.Product{
		Keyword:    keyword,
		LeftBrace:  lBrace,
		Id:         id,
		Fields:     fields,
		RightBrace: rBrace,
	}, ok
}

func (p *parser) parseFields() ([]st.Field, bool) {
	fields := []st.Field{}
	lastFieldOk := true
	for {
		if lastFieldOk {
			if p.currToken().Type != lex.TOKEN_ID {
				return fields, p.currToken().Type == lex.TOKEN_RBRACE
			}
		} else {
			currentToken, ok := p.eatUntilOneOf([]lex.TokenType{
				lex.TOKEN_LBRACE,
				lex.TOKEN_ID,
			}, false)
			if !ok {
				return fields, false
			}
			if currentToken.Type == lex.TOKEN_LBRACE {
				return fields, true
			}
		}
		field, ok := p.parseField()
		lastFieldOk = ok
		fields = append(fields, field)
	}
}

func (p *parser) parseField() (st.Field, bool) {
	id, ok := p.expectToken(lex.TOKEN_ID, true)

	_, ok = p.eatUntilOneOf([]lex.TokenType{
		lex.TOKEN_ID,
		lex.TOKEN_LIST,
		lex.TOKEN_LITERAL,
		lex.TOKEN_OPTIONAL,
		lex.TOKEN_SEPARATOR,
	}, ok)

	if !ok {
		return st.Field{}, false
	}

	currToken := p.currToken()

	if currToken.Type == lex.TOKEN_ID || currToken.Type == lex.TOKEN_LIST || currToken.Type == lex.TOKEN_LITERAL {
		jsonName := p.parseJsonRename()
		fieldType, ok := p.parseType()
		separator, ok := p.expectToken(lex.TOKEN_SEPARATOR, ok)
		return st.Field{
			FieldFull: &st.FieldFull{
				Id:        id,
				JsonName:  jsonName,
				Type:      fieldType,
				Separator: separator,
			},
		}, ok
	}
	if currToken.Type == lex.TOKEN_OPTIONAL || currToken.Type == lex.TOKEN_SEPARATOR {
		fieldNullable := p.parseNullability()
		separator, ok := p.expectToken(lex.TOKEN_SEPARATOR, true)
		return st.Field{
			FieldShort: &st.FieldShort{
				Id:        id,
				Nullable:  fieldNullable,
				Separator: separator,
			},
		}, ok
	}
	panic("Unreachable")
}

func (p *parser) parseTypeIdent() (st.TypeIdent, bool) {
	id, ok := p.expectToken(lex.TOKEN_ID, true)
	nullable := p.parseNullability()
	return st.TypeIdent{
		Id:       id,
		Nullable: nullable,
	}, ok
}

func (p *parser) parseList() (st.List, bool) {
	brackets, ok := p.expectToken(lex.TOKEN_LIST, true)
	nullable := p.parseNullability()
	listType, ok := p.parseType()

	return st.List{
		Brackets: brackets,
		Nullable: nullable,
		Type:     listType,
	}, ok
}

func (p *parser) parseType() (st.Type, bool) {
	switch p.currToken().Type {
	case lex.TOKEN_ID:
		typeIdent, ok := p.parseTypeIdent()
		return st.Type{
			TypeIdent: &typeIdent,
		}, ok
	case lex.TOKEN_LIST:
		typeList, ok := p.parseList()
		return st.Type{
			List: &typeList,
		}, ok
	default:
	}
	return st.Type{}, false
}

func (p *parser) parseNullability() *lex.Token {
	nullable, ok := p.optionalNextToken(lex.TOKEN_OPTIONAL)
	if ok {
		return &nullable
	}
	return nil
}

func (p *parser) parseJsonRename() *lex.Token {
	jsonRename := p.currToken()
	if jsonRename.Type == lex.TOKEN_LITERAL {
		p.pos++
		return &jsonRename
	}
	return nil
}

func (p *parser) currToken() lex.Token {
	if p.pos >= len(p.input) {
		return p.input[len(p.input)-1]
	}
	return p.input[p.pos]
}
