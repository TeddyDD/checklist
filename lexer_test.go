package checklist

import (
	"fmt"
	"testing"
)

type testcase struct {
	input  string
	output []token
}

func table() []testcase {
	return []testcase{
		{
			"-",
			[]token{
				{tokenHypen, "-"},
				{tokenEOF, ""},
			},
		},
		{
			"- hello",
			[]token{
				{tokenHypen, "-"},
				{tokenText, " hello"},
				{tokenEOF, ""},
			},
		},
		{
			"- -",
			[]token{
				{tokenHypen, "-"},
				{tokenText, " -"},
				{tokenEOF, ""},
			},
		}, {
			"- ",
			[]token{
				{tokenHypen, "-"},
				{tokenText, " "},
				{tokenEOF, ""},
			},
		},
		{
			`- hello
- world`,
			[]token{
				{tokenHypen, "-"},
				{tokenText, " hello\n"},
				{tokenHypen, "-"},
				{tokenText, " world"},
				{tokenEOF, ""},
			},
		},
		{
			`- hello
world`,
			[]token{
				{tokenHypen, "-"},
				{tokenText, " hello\nworld"},
				{tokenEOF, ""},
			},
		},
		{
			"- hello `cmd`",
			[]token{
				{tokenHypen, "-"},
				{tokenText, " hello "},
				{tokenCommand, "`cmd`"},
				{tokenEOF, ""},
			},
		},
		{
			"- hello ``",
			[]token{
				{tokenHypen, "-"},
				{tokenText, " hello "},
				{tokenCommand, "``"},
				{tokenEOF, ""},
			},
		},
		{
			"- hello `cmd`\n- world `cmd2`",
			[]token{
				{tokenHypen, "-"},
				{tokenText, " hello "},
				{tokenCommand, "`cmd`"},
				{tokenHypen, "-"},
				{tokenText, " world "},
				{tokenCommand, "`cmd2`"},
				{tokenEOF, ""},
			},
		}}
}

func errorsTable() []testcase {
	return []testcase{
		{
			"- foo `",
			[]token{
				{tokenHypen, ""},
				{tokenText, ""},
				{tokenError, ""},
			},
		},
	}
}

func assertEqual(t *testing.T, expected, got interface{}, msg string) {
	if expected != got {
		t.Fatalf("%s:\n\t expected: %v\n\tgot: %v\n", msg, expected, got)
	}
}

func assertNotNil(t *testing.T, got interface{}, msg string) {
	if got == nil {
		t.Fatalf("%s", msg)
	}
}

func TestSimpleLexing(t *testing.T) {
	tests := table()
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			l := newLexer(test.input)
			for _, expected := range test.output {
				got := l.nextToken()
				if got.typ != expected.typ {
					t.Fatalf("wrong token type, got %q, expected %q\n token %q\n input: %s ", got.typ, expected.typ, got, test.input)
				}
				if got.val != expected.val {
					t.Fatalf("wrong token value, got %q, expected %q\n token: %q\ninput: %s", got.val, expected.val, got, test.input)
				}
			}
		})
	}
}

func TestLexerErrors(t *testing.T) {
	tests := errorsTable()
	for i, test := range tests {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			l := newLexer(test.input)
			for _, expected := range test.output {
				got := l.nextToken()
				if got.typ != expected.typ {
					t.Fatalf("wrong token type, got %q, expected %q, token: %q", got.typ, expected.typ, got)
				}
			}
		})
	}
}
