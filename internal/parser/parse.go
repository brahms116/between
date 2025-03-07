package parser

import (
	"fmt"
	"strings"

	"github.com/brahms116/between/internal/ast"
	"github.com/brahms116/between/internal/lex"
)

type parser struct {
	input []lex.Token
	pos   int
}

const NUM_TYPE_STRING = "Num"
const STR_TYPE_STRING = "Str"
const OBJ_TYPE_STRING = "Obj"

func Parse(input []lex.Token) ([]ast.Definition, error) {
	p := &parser{input: input}
	return p.parse()
}

func (p *parser) parse() ([]ast.Definition, error) {
	definitions := []ast.Definition{}

	for p.pos < len(p.input) {
		definition, err := p.parseDefinition()
		if err != nil {
			return nil, err
		}
		definitions = append(definitions, definition)
	}
	return definitions, nil
}

func (p *parser) parseDefinition() (ast.Definition, error) {
	switch p.currToken().Type {
	case lex.TOKEN_PRODUCT:
		prod, err := p.parseProduct()
		return ast.Definition{Product: &prod}, err
	case lex.TOKEN_SUM:
		sum, err := p.parseSum()
		return ast.Definition{Sum: &sum}, err
	case lex.TOKEN_SUM_STR:
		sumStr, err := p.parseSumStr()
		return ast.Definition{SumStr: &sumStr}, err
	}
	return ast.Definition{}, fmt.Errorf("Unexpected token %s, at %d", p.currToken().String(), p.currToken().Loc.Length)
}

func (p *parser) parseSumStr() (ast.SumStr, error) {
	if p.currToken().Type != lex.TOKEN_SUM_STR {
		return ast.SumStr{}, fmt.Errorf("Expected keyword sumstr at %d", p.currToken().Loc.FilePos)
	}
	p.pos++

	if p.currToken().Type != lex.TOKEN_ID {
		return ast.SumStr{}, fmt.Errorf("Expected identifier at %d", p.currToken().Loc.FilePos)
	}
	name := p.currToken().Value
	p.pos++

	if p.currToken().Type != lex.TOKEN_LBRACE {
		return ast.SumStr{}, fmt.Errorf("Expected { at %d", p.currToken().Loc.FilePos)
	}
	p.pos++

	variants, err := p.parseSumStrVariants()
	if err != nil {
		return ast.SumStr{}, err
	}

	if p.currToken().Type != lex.TOKEN_RBRACE {
		return ast.SumStr{}, fmt.Errorf("Expected } at %d", p.currToken().Loc.FilePos)
	}
	p.pos++

	return ast.SumStr{
		Id:       name,
		Variants: variants,
	}, nil

}

func (p *parser) parseSumStrVariants() ([]string, error) {
	variants := []string{}
	for {
		if p.currToken().Type != lex.TOKEN_LITERAL {
			break
		}

		variants = append(variants, p.currToken().Value)
		p.pos++

		if p.currToken().Type != lex.TOKEN_SEPARATOR {
			return nil, fmt.Errorf("Expected SEPARATOR at %d", p.currToken().Loc.FilePos)
		}
		p.pos++

	}
	return variants, nil
}

func (p *parser) parseSum() (ast.Sum, error) {
	if p.currToken().Type != lex.TOKEN_SUM {
		return ast.Sum{}, fmt.Errorf("Expected keyword sum at %d", p.currToken().Loc.FilePos)
	}
	p.pos++

	if p.currToken().Type != lex.TOKEN_ID {
		return ast.Sum{}, fmt.Errorf("Expected identifier at %d", p.currToken().Loc.FilePos)
	}
	name := p.currToken().Value
	p.pos++

	if p.currToken().Type != lex.TOKEN_LBRACE {
		return ast.Sum{}, fmt.Errorf("Expected { at %d", p.currToken().Loc.FilePos)
	}
	p.pos++

	fields, err := p.parseFields()
	if err != nil {
		return ast.Sum{}, err
	}

	if p.currToken().Type != lex.TOKEN_RBRACE {
		return ast.Sum{}, fmt.Errorf("Expected } at %d", p.currToken().Loc.FilePos)
	}
	p.pos++

	return ast.Sum{
		Id:       name,
		Variants: fields,
	}, nil
}

func (p *parser) parseProduct() (ast.Product, error) {
	if p.currToken().Type != lex.TOKEN_PRODUCT {
		return ast.Product{}, fmt.Errorf("Expected keyword prod at %d", p.currToken().Loc.FilePos)
	}
	p.pos++

	if p.currToken().Type != lex.TOKEN_ID {
		return ast.Product{}, fmt.Errorf("Expected identifier at %d", p.currToken().Loc.FilePos)
	}
	name := p.currToken().Value
	p.pos++

	if p.currToken().Type != lex.TOKEN_LBRACE {
		return ast.Product{}, fmt.Errorf("Expected { at %d", p.currToken().Loc.FilePos)
	}

	p.pos++
	fields, err := p.parseFields()
	if err != nil {
		return ast.Product{}, err
	}

	if p.currToken().Type != lex.TOKEN_RBRACE {
		return ast.Product{}, fmt.Errorf("Expected } at %d", p.currToken().Loc.FilePos)
	}
	p.pos++

	return ast.Product{
		Id:     name,
		Fields: fields,
	}, nil
}

func (p *parser) parseFields() ([]ast.Field, error) {
	fields := []ast.Field{}
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

func (p *parser) parseField() (ast.Field, error) {
	currToken := p.currToken()
	if currToken.Type != lex.TOKEN_ID {
		return ast.Field{}, fmt.Errorf("Unexpect ID at %d", currToken.Loc.FilePos)
	}
	name := currToken.Value
	p.pos++

	currToken = p.currToken()
	if currToken.Type == lex.TOKEN_ID ||
		currToken.Type == lex.TOKEN_STR ||
		currToken.Type == lex.TOKEN_NUM ||
		currToken.Type == lex.TOKEN_OBJ {
		fieldType, err := p.parseType()
		if err != nil {
			return ast.Field{}, err
		}

		currToken = p.currToken()
		if currToken.Type != lex.TOKEN_SEPARATOR {
			return ast.Field{}, fmt.Errorf("Expected SEPARATOR at %d", currToken.Loc.FilePos)
		}
		p.pos++

		return ast.Field{
			Id:   name,
			Type: fieldType,
		}, nil
	}
	if currToken.Type == lex.TOKEN_OPTIONAL || currToken.Type == lex.TOKEN_SEPARATOR {
		fieldNullable := p.parseNullability()

		currToken = p.currToken()
		if currToken.Type != lex.TOKEN_SEPARATOR {
			return ast.Field{}, fmt.Errorf("Expected SEPARATOR at %d", currToken.Loc.FilePos)
		}
		p.pos++

		return ast.Field{
			Id: lowerCaseFirstLetter(name),
			Type: ast.Type{
				Name:     name,
				Nullable: fieldNullable,
			},
		}, nil
	}

	return ast.Field{}, fmt.Errorf("Unexpected token %s at %d", currToken.String(), currToken.Loc.FilePos)
}

func (p *parser) parseType() (ast.Type, error) {
	var name string
	nameToken := p.currToken()
	switch nameToken.Type {
	case lex.TOKEN_ID:
		name = nameToken.Value
	case lex.TOKEN_NUM:
		name = NUM_TYPE_STRING
	case lex.TOKEN_STR:
		name = STR_TYPE_STRING
  case lex.TOKEN_OBJ:
    name = OBJ_TYPE_STRING
	default:
		return ast.Type{}, fmt.Errorf("Unexpected token at %d, expected name of a type", nameToken.Loc.FilePos)
	}
	p.pos++
	typeNullable := p.parseNullability()
	typeList := p.parseList()
	return ast.Type{
		Name:     name,
		Nullable: typeNullable,
		List:     typeList,
	}, nil
}

func (p *parser) parseList() *ast.List {
	if p.currToken().Type == lex.TOKEN_LIST {
		p.pos++
		isNullable := p.parseNullability()
		return &ast.List{
			Nullable: isNullable,
		}
	}
	return nil
}

func (p *parser) parseNullability() bool {
	if p.currToken().Type == lex.TOKEN_OPTIONAL {
		p.pos++
		return true
	}
	return false
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
