package parser

import (
	"fmt"
	"strings"

	"github.com/brahms116/between/internal/lex"
	"github.com/brahms116/between/internal/st"
)

type parser struct {
	input []lex.Token
	pos   int
}

const NUM_TYPE_STRING = "Num"
const STR_TYPE_STRING = "Str"
const OBJ_TYPE_STRING = "Obj"

func UnexpectedTokenError(token lex.Token) error {
	return fmt.Errorf("Unexpected token %s at %s", token.String(), token.Loc.Start.String())
}

func ExpectedTokenError(expected lex.TokenType, token lex.Token) error {
	return fmt.Errorf("Expected %s at %s, got %s", expected.String(), token.Loc.Start.String(), token.String())
}

func LexAndParse(input string) ([]st.Definition, error) {
	tokens, err := lex.Lex(input)
	if err != nil {
		return nil, err
	}
	return Parse(tokens)
}

func Parse(input []lex.Token) ([]st.Definition, error) {
	p := &parser{input: input}
	return p.parse()
}

func (p *parser) parse() ([]st.Definition, error) {
	definitions := []st.Definition{}

	for p.pos < len(p.input) {
		definition, err := p.parseDefinition()
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, definition)
	}
	return definitions, nil
}

func (p *parser) expectToken(tokenType lex.TokenType) (lex.Token, error) {
  currToken := p.currToken()
	if currToken.Type != tokenType {
		return lex.Token{}, ExpectedTokenError(tokenType, currToken)
	}
	p.pos++
	return currToken, nil
}

func (p *parser) optionalToken(tokenType lex.TokenType) (lex.Token, bool) {
  currToken := p.currToken()
	if currToken.Type == tokenType {
		p.pos++
		return currToken, true
	}
	return lex.Token{}, false
}

func (p *parser) parseDefinition() (st.Definition, error) {
	switch p.currToken().Type {
	case lex.TOKEN_PRODUCT:
		prod, err := p.parseProduct()
		return st.Definition{Product: &prod}, err
	case lex.TOKEN_SUM:
		sum, err := p.parseSum()
		return st.Definition{Sum: &sum}, err
	case lex.TOKEN_SUM_STR:
		sumStr, err := p.parseSumStr()
		return st.Definition{SumStr: &sumStr}, err
	}
	return st.Definition{}, UnexpectedTokenError(p.currToken())
}

func (p *parser) parseSumStr() (st.SumStr, error) {
	keyword, err := p.expectToken(lex.TOKEN_SUM_STR)
	if err != nil {
		return st.SumStr{}, err
	}

	id, err := p.expectToken(lex.TOKEN_ID)
	if err != nil {
		return st.SumStr{}, err
	}

	lBrace, err := p.expectToken(lex.TOKEN_LBRACE)
	if err != nil {
		return st.SumStr{}, err
	}

	variants, err := p.parseSumStrVariants()
	if err != nil {
		return st.SumStr{}, err
	}

	rBrace, err := p.expectToken(lex.TOKEN_RBRACE)
	if err != nil {
		return st.SumStr{}, err
	}

	return st.SumStr{
		Keyword:    keyword,
		Id:         id,
		LeftBrace:  lBrace,
		Variants:   variants,
		RightBrace: rBrace,
	}, nil

}

func (p *parser) parseSumStrVariants() ([]st.SumStrVariant, error) {
	variants := []st.SumStrVariant{}
	for {
		id, ok := p.optionalToken(lex.TOKEN_ID)
		if !ok {
			break
		}

		jsonName := p.parseJsonRename()

		separator, err := p.expectToken(lex.TOKEN_SEPARATOR)
		if err != nil {
			return nil, err
		}

		variants = append(variants, st.SumStrVariant{
			Id:        id,
			JsonName:  jsonName,
			Separator: separator,
		})

	}
	return variants, nil
}

func (p *parser) parseSum() (st.Sum, error) {
	keyword, err := p.expectToken(lex.TOKEN_SUM)
	if err != nil {
		return st.Sum{}, err
	}

	id, err := p.expectToken(lex.TOKEN_ID)
	if err != nil {
		return st.Sum{}, err
	}

	lBrace, err := p.expectToken(lex.TOKEN_LBRACE)
	if err != nil {
		return st.Sum{}, err
	}

	fields, err := p.parseFields()
	if err != nil {
		return st.Sum{}, err
	}

	rBrace, err := p.expectToken(lex.TOKEN_RBRACE)
	if err != nil {
		return st.Sum{}, err
	}

	return st.Sum{
		Keyword:    keyword,
		Id:         id,
		LeftBrace:  lBrace,
		Variants:   fields,
		RightBrace: rBrace,
	}, nil
}

func (p *parser) parseProduct() (st.Product, error) {
	keyword, err := p.expectToken(lex.TOKEN_PRODUCT)
	if err != nil {
		return st.Product{}, err
	}

	id, err := p.expectToken(lex.TOKEN_ID)
	if err != nil {
		return st.Product{}, err
	}

	lBrace, err := p.expectToken(lex.TOKEN_LBRACE)
	if err != nil {
		return st.Product{}, err
	}

	fields, err := p.parseFields()
	if err != nil {
		return st.Product{}, err
	}

	rBrace, err := p.expectToken(lex.TOKEN_RBRACE)
	if err != nil {
		return st.Product{}, err
	}

	return st.Product{
		Keyword:    keyword,
		LeftBrace:  lBrace,
		Id:         id,
		Fields:     fields,
		RightBrace: rBrace,
	}, nil
}

func (p *parser) parseFields() ([]st.Field, error) {
	fields := []st.Field{}
	for {
		currToken := p.currToken()
		if currToken.Type != lex.TOKEN_ID {
			break
		}
		field, err := p.parseField()
		if err != nil {
			return nil, err
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func (p *parser) parseField() (st.Field, error) {
	id, err := p.expectToken(lex.TOKEN_ID)
	if err != nil {
		return st.Field{}, err
	}

	currToken := p.currToken()
	if currToken.Type == lex.TOKEN_ID || currToken.Type == lex.TOKEN_LIST || currToken.Type == lex.TOKEN_LITERAL {
		jsonName := p.parseJsonRename()
		fieldType, err := p.parseType()
		if err != nil {
			return st.Field{}, err
		}

		separator, err := p.expectToken(lex.TOKEN_SEPARATOR)
		if err != nil {
			return st.Field{}, err
		}
		return st.Field{
			FieldFull: &st.FieldFull{
				Id:        id,
				JsonName:  jsonName,
				Type:      fieldType,
				Separator: separator,
			},
		}, nil
	}
	if currToken.Type == lex.TOKEN_OPTIONAL || currToken.Type == lex.TOKEN_SEPARATOR {
		fieldNullable := p.parseNullability()

		separator, err := p.expectToken(lex.TOKEN_SEPARATOR)
		if err != nil {
			return st.Field{}, err
		}
		return st.Field{
			FieldShort: &st.FieldShort{
				Id:        id,
				Nullable:  fieldNullable,
				Separator: separator,
			},
		}, nil
	}

	return st.Field{}, ExpectedTokenError(lex.TOKEN_SEPARATOR, p.currToken())
}

func (p *parser) parseTypeIdent() (st.TypeIdent, error) {
	id, err := p.expectToken(lex.TOKEN_ID)
	if err != nil {
		return st.TypeIdent{}, err
	}
	nullable := p.parseNullability()

	return st.TypeIdent{
		Id:       id,
		Nullable: nullable,
	}, nil
}

func (p *parser) parseList() (st.List, error) {
	brackets, err := p.expectToken(lex.TOKEN_LIST)
	if err != nil {
	}
	nullable := p.parseNullability()
	listType, err := p.parseType()
	if err != nil {
		return st.List{}, err
	}

	return st.List{
		Brackets: brackets,
		Nullable: nullable,
		Type:     listType,
	}, nil
}

func (p *parser) parseType() (st.Type, error) {
	switch p.currToken().Type {
	case lex.TOKEN_ID:
		typeIdent, err := p.parseTypeIdent()
		if err != nil {
			return st.Type{}, err
		}
		return st.Type{
			TypeIdent: &typeIdent,
		}, nil
	case lex.TOKEN_LIST:
		typeList, err := p.parseList()
		if err != nil {
			return st.Type{}, err
		}
		return st.Type{
			List: &typeList,
		}, nil
	default:
	}
	return st.Type{}, ExpectedTokenError(lex.TOKEN_ID, p.currToken())
}

func (p *parser) parseNullability() *lex.Token {
	nullable, ok := p.optionalToken(lex.TOKEN_OPTIONAL)
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
	return p.input[p.pos]
}

func lowerCaseFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}
