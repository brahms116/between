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

type UnexpectedTokenError struct {
	Expected []lex.TokenType
	Actual   lex.Token
}

func newUnexpectedTokenError(expected []lex.TokenType, actual lex.Token) UnexpectedTokenError {
	return UnexpectedTokenError{
		Expected: expected,
		Actual:   actual,
	}
}

func (e UnexpectedTokenError) Error() string {
	expectedStr := ""
	for i, t := range e.Expected {
		if i > 0 {
			expectedStr += " or "
		}
		expectedStr += t.String()
	}

	return fmt.Sprintf("Expected %s at %s, got %s", expectedStr, e.Actual.Loc.Start.String(), e.Actual.String())
}

func (e UnexpectedTokenError) LspMessage() string {
	expectedStr := ""
	for i, t := range e.Expected {
		if i > 0 {
			expectedStr += " or "
		}
		expectedStr += t.String()
	}

	return fmt.Sprintf("Expected %s got %s", expectedStr, e.Actual.String())
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
	definitions := p.parseDefinitions()
	return definitions, p.errors
}

func (p *parser) appendErr(err error) {
	p.errors = append(p.errors, err)
}

func (p *parser) expect(tokenType lex.TokenType, follow []lex.TokenType) lex.Token {
	currToken := p.currToken()
	if currToken.Type == tokenType {
		p.pos++
		return currToken
	}
	p.errorUntil([]lex.TokenType{tokenType}, follow)
	return lex.Token{IsErr: true}
}

func (p *parser) errorUntil(expectedTokens []lex.TokenType, follow []lex.TokenType) {
	if !p.isEofError {
		p.appendErr(newUnexpectedTokenError(expectedTokens, p.currToken()))
	}

outer:
	for {
		currToken := p.currToken()
		for _, tokenType := range follow {
			if currToken.Type == tokenType {
				break outer
			}
		}
		if currToken.Type == lex.TOKEN_EOF {
			p.isEofError = true
			break outer
		}
		p.pos++
	}
}

func (p *parser) optionalNextToken(tokenType lex.TokenType) (lex.Token, bool) {
	currToken := p.currToken()
	if currToken.Type == tokenType {
		p.pos++
		return currToken, true
	}
	return lex.Token{}, false
}

func (p *parser) parseDefinitions() []st.Definition {
	definitions := []st.Definition{}
	for {
		if p.currToken().Type == lex.TOKEN_EOF {
			return definitions
		}
		definitions = append(definitions, p.parseDefinition())
	}
}

var definitionFirsts = []lex.TokenType{
	productFirst,
	sumFirst,
  sumStrFirst,
}

var definitionFollows = append(definitionFirsts, []lex.TokenType{
	lex.TOKEN_EOF,
}...)

func (p *parser) parseDefinition() st.Definition {
	switch p.currToken().Type {
	case lex.TOKEN_PRODUCT:
		prod := p.parseProduct()
		return st.Definition{Product: &prod}
	case lex.TOKEN_SUM:
		sum := p.parseSum()
		return st.Definition{Sum: &sum}
	case lex.TOKEN_SUM_STR:
		sumStr := p.parseSumStr()
		return st.Definition{SumStr: &sumStr}
	default:
		p.errorUntil(definitionFirsts, definitionFollows)
		return st.Definition{}
	}
}

var sumStrFirst = lex.TOKEN_SUM_STR
var sumStrFollows = definitionFollows

func (p *parser) parseSumStr() st.SumStr {
	keyword := p.expect(lex.TOKEN_SUM_STR, []lex.TokenType{lex.TOKEN_ID})
	id := p.expect(lex.TOKEN_ID, []lex.TokenType{lex.TOKEN_LBRACE})
	lBrace := p.expect(lex.TOKEN_LBRACE, []lex.TokenType{sumStrVariantsFirst, sumStrVariantsFollow})
	variants := p.parseSumStrVariants()
	rBrace := p.expect(lex.TOKEN_RBRACE, sumStrFollows)

	return st.SumStr{
		Keyword:    keyword,
		Id:         id,
		LeftBrace:  lBrace,
		Variants:   variants,
		RightBrace: rBrace,
	}
}

var sumStrVariantsFirst = sumStrVariantFirst
var sumStrVariantsFollow = lex.TOKEN_RBRACE

var sumStrVariantFirst = lex.TOKEN_ID
var sumStrVariantFollows = []lex.TokenType{
	sumStrVariantsFollow,
	sumStrVariantFirst,
}

func (p *parser) parseSumStrVariants() []st.SumStrVariant {
	variants := []st.SumStrVariant{}
	for {
		switch p.currToken().Type {
		case lex.TOKEN_ID:
			id := p.expect(lex.TOKEN_ID, []lex.TokenType{lex.TOKEN_SEPARATOR, lex.TOKEN_SEPARATOR})
			jsonName := p.parseJsonRename()
			separator := p.expect(lex.TOKEN_SEPARATOR, sumStrVariantFollows)
			variants = append(variants, st.SumStrVariant{
				Id:        id,
				JsonName:  jsonName,
				Separator: separator,
			})
		case sumStrVariantsFollow:
			return variants
		default:
			p.errorUntil([]lex.TokenType{sumStrVariantsFirst, sumStrVariantsFollow}, []lex.TokenType{sumStrVariantsFollow})
			return variants
		}
	}
}

var sumFirst = lex.TOKEN_SUM
var sumFollows = definitionFollows

func (p *parser) parseSum() st.Sum {
	keyword := p.expect(lex.TOKEN_SUM, []lex.TokenType{lex.TOKEN_ID})
	id := p.expect(lex.TOKEN_ID, []lex.TokenType{lex.TOKEN_LBRACE})
	lBrace := p.expect(lex.TOKEN_LBRACE, []lex.TokenType{fieldsFirst, lex.TOKEN_RBRACE})
	fields := p.parseFields()
	rBrace := p.expect(lex.TOKEN_RBRACE, sumFollows)
	return st.Sum{
		Keyword:    keyword,
		Id:         id,
		LeftBrace:  lBrace,
		Variants:   fields,
		RightBrace: rBrace,
	}
}

var productFirst = lex.TOKEN_PRODUCT
var productFollows = definitionFollows

func (p *parser) parseProduct() st.Product {
	keyword := p.expect(lex.TOKEN_PRODUCT, []lex.TokenType{lex.TOKEN_ID})
	id := p.expect(lex.TOKEN_ID, []lex.TokenType{lex.TOKEN_LBRACE})
	lBrace := p.expect(lex.TOKEN_LBRACE, []lex.TokenType{lex.TOKEN_ID})
	fields := p.parseFields()
	rBrace := p.expect(lex.TOKEN_RBRACE, definitionFollows)
	return st.Product{
		Keyword:    keyword,
		LeftBrace:  lBrace,
		Id:         id,
		Fields:     fields,
		RightBrace: rBrace,
	}
}

var fieldsFirst = fieldFirst
var fieldsFollow = lex.TOKEN_RBRACE

func (p *parser) parseFields() []st.Field {
	fields := []st.Field{}
	for {
		switch p.currToken().Type {
		case fieldsFirst:
			fields = append(fields, p.parseField())
		case fieldsFollow:
			return fields
		default:
			p.errorUntil([]lex.TokenType{fieldsFirst, fieldsFollow}, []lex.TokenType{fieldsFollow})
			return fields
		}
	}
}

var fieldFirst = lex.TOKEN_ID
var fieldFollows = []lex.TokenType{
	fieldsFollow,
	fieldFirst,
}

func (p *parser) parseField() st.Field {
	id := p.expect(lex.TOKEN_ID, []lex.TokenType{
		lex.TOKEN_ID,
		lex.TOKEN_LIST,
		lex.TOKEN_LITERAL,
		lex.TOKEN_OPTIONAL,
		lex.TOKEN_SEPARATOR,
	})

	currToken := p.currToken()
	if currToken.Type == lex.TOKEN_ID || currToken.Type == lex.TOKEN_LIST || currToken.Type == lex.TOKEN_LITERAL {
		jsonName := p.parseJsonRename()
		fieldType := p.parseType()
		separator := p.expect(lex.TOKEN_SEPARATOR, fieldFollows)
		return st.Field{
			FieldFull: &st.FieldFull{
				Id:        id,
				JsonName:  jsonName,
				Type:      fieldType,
				Separator: separator,
			},
		}
	}
	if currToken.Type == lex.TOKEN_OPTIONAL || currToken.Type == lex.TOKEN_SEPARATOR {
		fieldNullable := p.parseNullability()
		separator := p.expect(lex.TOKEN_SEPARATOR, fieldFollows)
		return st.Field{
			FieldShort: &st.FieldShort{
				Id:        id,
				Nullable:  fieldNullable,
				Separator: separator,
			},
		}
	}
	p.errorUntil([]lex.TokenType{
		lex.TOKEN_ID,
		lex.TOKEN_LIST,
		lex.TOKEN_LITERAL,
		lex.TOKEN_OPTIONAL,
		lex.TOKEN_SEPARATOR,
	}, fieldFollows)
	return st.Field{}
}

var typeIdentFirst = lex.TOKEN_ID
var typeIdentFollow = typeFollow

func (p *parser) parseTypeIdent() st.TypeIdent {
	id := p.expect(lex.TOKEN_ID, []lex.TokenType{
		typeIdentFollow,
		lex.TOKEN_OPTIONAL,
	})
	nullable := p.parseNullability()
	return st.TypeIdent{
		Id:       id,
		Nullable: nullable,
	}
}

var listFirst = lex.TOKEN_LIST
var listFollow = typeFollow

func (p *parser) parseList() st.List {
	brackets := p.expect(lex.TOKEN_LIST, append(typeFirsts, lex.TOKEN_OPTIONAL))
	nullable := p.parseNullability()
	listType := p.parseType()
	return st.List{
		Brackets: brackets,
		Nullable: nullable,
		Type:     listType,
	}
}

var typeFirsts = []lex.TokenType{
	typeIdentFirst,
	listFirst,
}
var typeFollow = lex.TOKEN_SEPARATOR

func (p *parser) parseType() st.Type {
	switch p.currToken().Type {
	case lex.TOKEN_ID:
		typeIdent := p.parseTypeIdent()
		return st.Type{
			TypeIdent: &typeIdent,
		}
	case lex.TOKEN_LIST:
		typeList := p.parseList()
		return st.Type{
			List: &typeList,
		}
	default:
		p.errorUntil(typeFirsts, []lex.TokenType{typeFollow})
		return st.Type{}
	}
}

var nullabilityFirst = lex.TOKEN_OPTIONAL

func (p *parser) parseNullability() *lex.Token {
	nullable, ok := p.optionalNextToken(lex.TOKEN_OPTIONAL)
	if ok {
		return &nullable
	}
	return nil
}

var jsonRenameFirst = lex.TOKEN_LITERAL

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
