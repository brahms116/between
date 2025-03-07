package lex

import (
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
				Loc: Loc{
					Start:  0,
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
