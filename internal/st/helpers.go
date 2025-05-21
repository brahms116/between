package st

import "github.com/brahms116/between/internal/lex"

func (f Field) Id() lex.Token {
  if f.FieldFull != nil {
    return f.FieldFull.Id
  } else if f.FieldShort != nil {
    return f.FieldShort.Id
  }
  panic("unreachable")
}
