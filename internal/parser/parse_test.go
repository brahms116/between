package parser

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseSmoke(t *testing.T) {
	data, err := os.ReadFile("../../testcases/001.bt")
	assert.Nil(t, err)
	defintions, errs := LexAndParse(string(data))
	assert.Equal(t, 0, len(errs))
  for _, err := range errs {
    t.Log(err)
  }
	bytes, err := json.MarshalIndent(defintions, "", "  ")
	println(len(bytes))
}
