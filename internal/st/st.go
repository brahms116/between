package st

import "github.com/brahms116/between/internal/lex"

type Type struct {
	List      *List
	TypeIdent *TypeIdent
}

type TypeIdent struct {
	Id       lex.Token
	Nullable *lex.Token
}

type List struct {
	Brackets lex.Token
	Nullable *lex.Token
	Type     Type
}

type Field struct {
	FieldFull  *FieldFull
	FieldShort *FieldShort
}

type FieldFull struct {
	Id        lex.Token
	JsonName  *lex.Token
	Type      Type
	Separator lex.Token
}

type FieldShort struct {
	Id        lex.Token
	Nullable  *lex.Token
	Separator lex.Token
}

type Product struct {
	Keyword    lex.Token
	Id         lex.Token
	LeftBrace  lex.Token
	Fields     []Field
	RightBrace lex.Token
}

type Sum struct {
	Keyword    lex.Token
	Id         lex.Token
	LeftBrace  lex.Token
	Variants   []Field
	RightBrace lex.Token
}

type SumStr struct {
	Keyword    lex.Token
	Id         lex.Token
	LeftBrace  lex.Token
	Variants   []SumStrVariant
	RightBrace lex.Token
}

type SumStrVariant struct {
	Id        lex.Token
	JsonName  *lex.Token
	Separator lex.Token
}

type Definition struct {
	Product *Product
	Sum     *Sum
	SumStr  *SumStr
}
