package lex

import (
	"log"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	input    string
	expected []Token
}

var cases []testCase = []testCase{
	{
		input: "prod",
		expected: []Token{
			{
				Type: TOKEN_PRODUCT,
				Loc: Location{
					FilePos:  0,
					Length: 4,
				},
			},
		},
	},
	{
		input: "User",
		expected: []Token{
			{
				Type: TOKEN_ID,
        Value: "User",
				Loc: Location{
					FilePos:  0,
					Length: 4,
				},
			},
		},
	},
}

func TestLex(t *testing.T) {
	for _, testCase := range cases {
		result, err := Lex(testCase.input)
		assert.Nil(t, err)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestSmokeLex(t *testing.T) {
	data, err := os.ReadFile("../../testcases/001.bt")
	assert.Nil(t, err)
  tokens, err := Lex(string(data))
  assert.Nil(t, err)
  for _, token := range tokens {
    log.Println(token.String())
  }
}
