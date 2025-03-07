package parser

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/brahms116/between/internal/lex"
	"github.com/stretchr/testify/assert"
)

func TestParseSmoke(t *testing.T) {

	data, err := os.ReadFile("../../testcases/001.bt")
	assert.Nil(t, err)
  tokens, err := lex.Lex(string(data))
  assert.Nil(t, err)
  defintions, err := Parse(tokens)
  assert.Nil(t, err)
  bytes, err := json.Marshal(defintions)
  println(string(bytes))
}
