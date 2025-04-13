package ast

func (t Type) IsNullable() bool {
	if t.List != nil {
		return t.List.Nullable
	}
	return t.TypeIdent.Nullable
}

