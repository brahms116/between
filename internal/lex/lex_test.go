package lex

import (
	"encoding/json"
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
					ByteStart: 0,
					ByteEnd:   4,
					Start: Point{
						Row: 0,
						Col: 0,
					},
					End: Point{
						Row: 0,
						Col: 4,
					},
				},
			},
			{
				Type: TOKEN_EOF,
				Loc: Location{
					ByteStart: 4,
					ByteEnd:   4,
					Start: Point{
						Row: 0,
						Col: 4,
					},
					End: Point{
						Row: 0,
						Col: 4,
					},
				},
			},
		},
	},
	{
		input: "User",
		expected: []Token{
			{
				Type:  TOKEN_ID,
				Value: "User",
				Loc: Location{
					ByteStart: 0,
					ByteEnd:   4,
					Start: Point{
						Row: 0,
						Col: 0,
					},
					End: Point{
						Row: 0,
						Col: 4,
					},
				},
			},
			{
				Type: TOKEN_EOF,
				Loc: Location{
					ByteStart: 4,
					ByteEnd:   4,
					Start: Point{
						Row: 0,
						Col: 4,
					},
					End: Point{
						Row: 0,
						Col: 4,
					},
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
	tokens, errs := Lex(string(data))
	assert.Nil(t, errs)
	for _, token := range tokens {
		jsonStr, err := json.MarshalIndent(token, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		log.Print(string(jsonStr))
	}
}
