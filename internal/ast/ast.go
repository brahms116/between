package ast

// TODO add span location

type Type struct {
	List      *List
	TypeIdent *TypeIdent
}

type TypeIdent struct {
	Id       string
	Nullable bool
}

type List struct {
	Nullable bool
	Type     Type
}

type Field struct {
	Id   string
	Type Type
}

type Product struct {
	Id     string
	Fields []Field
}

type Sum struct {
	Id       string
	Variants []Field
}

type SumStr struct {
	Id       string
	Variants []string
}

type Definition struct {
	Product *Product
	Sum     *Sum
	SumStr  *SumStr
}
