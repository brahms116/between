package translate

type symbolType string

const (
	symbolTypePrimitive symbolType = "primitive"
	symbolTypeProduct   symbolType = "product"
	symbolTypeSum       symbolType = "sum"
	symbolTypeSumString symbolType = "sum_string"
)

type symbolTable map[string]symbolType

func newSymbolTable() symbolTable {
	st := make(symbolTable)
	for t, _ := range PrimitiveTypes {
		st[t] = symbolTypePrimitive
	}
	return st
}

func (st symbolTable) addSymbol(name string, typ symbolType) bool {
  _, ok := st[name]
  if ok {
    return false
  }
	st[name] = typ
  return true
}

func (st symbolTable) getSymbol(name string) (symbolType, bool) {
	typ, ok := st[name]
	return typ, ok
}
