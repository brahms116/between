package parser

import (
	"encoding/json"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseSmoke(t *testing.T) {
	data, err := os.ReadFile("../../testcases/001.bt")
	assert.Nil(t, err)
	defintions, err := LexAndParse(string(data))
	assert.Nil(t, err)
	bytes, err := json.MarshalIndent(defintions, "", "  ")
	println(string(bytes))
}
