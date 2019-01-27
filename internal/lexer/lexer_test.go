package lexer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	cases := []struct {
		usage    string
		expected []*Token
	}{
		{"OPTIONS", []*Token{{TTOptions, "OPTIONS", 0}}},

		{"XOPTIONS", []*Token{{TTArg, "XOPTIONS", 0}}},
		{"OPTIONSX", []*Token{{TTArg, "OPTIONSX", 0}}},
		{"ARG", []*Token{{TTArg, "ARG", 0}}},
		{"ARG42", []*Token{{TTArg, "ARG42", 0}}},
		{"ARG_EXTRA", []*Token{{TTArg, "ARG_EXTRA", 0}}},

		{"ARG1 ARG2", []*Token{{TTArg, "ARG1", 0}, {TTArg, "ARG2", 5}}},
		{"ARG1  ARG2", []*Token{{TTArg, "ARG1", 0}, {TTArg, "ARG2", 6}}},

		{"(", []*Token{{TTOpenPar, "(", 0}}},
		{")", []*Token{{TTClosePar, ")", 0}}},

		{"(ARG)", []*Token{{TTOpenPar, "(", 0}, {TTArg, "ARG", 1}, {TTClosePar, ")", 4}}},
		{"( ARG )", []*Token{{TTOpenPar, "(", 0}, {TTArg, "ARG", 2}, {TTClosePar, ")", 6}}},

		{"[ARG]", []*Token{{TTOpenSq, "[", 0}, {TTArg, "ARG", 1}, {TTCloseSq, "]", 4}}},
		{"[ ARG ]", []*Token{{TTOpenSq, "[", 0}, {TTArg, "ARG", 2}, {TTCloseSq, "]", 6}}},
		{"ARG [ARG2 ]", []*Token{{TTArg, "ARG", 0}, {TTOpenSq, "[", 4}, {TTArg, "ARG2", 5}, {TTCloseSq, "]", 10}}},
		{"ARG [ ARG2]", []*Token{{TTArg, "ARG", 0}, {TTOpenSq, "[", 4}, {TTArg, "ARG2", 6}, {TTCloseSq, "]", 10}}},

		{"...", []*Token{{TTRep, "...", 0}}},
		{"ARG...", []*Token{{TTArg, "ARG", 0}, {TTRep, "...", 3}}},
		{"ARG ...", []*Token{{TTArg, "ARG", 0}, {TTRep, "...", 4}}},
		{"[ARG...]", []*Token{{TTOpenSq, "[", 0}, {TTArg, "ARG", 1}, {TTRep, "...", 4}, {TTCloseSq, "]", 7}}},

		{"|", []*Token{{TTChoice, "|", 0}}},
		{"ARG|ARG2", []*Token{{TTArg, "ARG", 0}, {TTChoice, "|", 3}, {TTArg, "ARG2", 4}}},
		{"ARG |ARG2", []*Token{{TTArg, "ARG", 0}, {TTChoice, "|", 4}, {TTArg, "ARG2", 5}}},
		{"ARG| ARG2", []*Token{{TTArg, "ARG", 0}, {TTChoice, "|", 3}, {TTArg, "ARG2", 5}}},

		{"[OPTIONS]", []*Token{{TTOpenSq, "[", 0}, {TTOptions, "OPTIONS", 1}, {TTCloseSq, "]", 8}}},

		{"-p", []*Token{{TTShortOpt, "-p", 0}}},
		{"-X", []*Token{{TTShortOpt, "-X", 0}}},

		{"--force", []*Token{{TTLongOpt, "--force", 0}}},
		{"--sig-proxy", []*Token{{TTLongOpt, "--sig-proxy", 0}}},

		{"-aBc", []*Token{{TTOptSeq, "aBc", 0}}},
		{"--", []*Token{{TTDoubleDash, "--", 0}}},
		{"=<bla>", []*Token{{TTOptValue, "=<bla>", 0}}},
		{"=<bla-bla>", []*Token{{TTOptValue, "=<bla-bla>", 0}}},
		{"=<bla--bla>", []*Token{{TTOptValue, "=<bla--bla>", 0}}},
		{"-p=<file-path>", []*Token{{TTShortOpt, "-p", 0}, {TTOptValue, "=<file-path>", 2}}},
		{"--path=<absolute-path>", []*Token{{TTLongOpt, "--path", 0}, {TTOptValue, "=<absolute-path>", 6}}},
	}
	for _, c := range cases {
		t.Run(c.usage, func(t *testing.T) {
			t.Logf("test %s", c.usage)
			tks, err := Tokenize(c.usage)
			if err != nil {
				t.Errorf("[Tokenize '%s']: Unexpected error: %v", c.usage, err)
				return
			}

			t.Logf("actual: %v\n", tks)
			if len(tks) != len(c.expected) {
				t.Errorf("[Tokenize '%s']: token count mismatch:\n\tExpected: %v\n\tActual  : %v", c.usage, c.expected, tks)
				return
			}

			for i, actual := range tks {
				expected := c.expected[i]
				switch {
				case actual.Typ != expected.Typ:
					t.Errorf("[Tokenize '%s']: token type mismatch:\n\tExpected: %v\n\tActual  : %v", c.usage, expected, actual)
				case actual.Val != expected.Val:
					t.Errorf("[Tokenize '%s']: token text mismatch:\n\tExpected: %v\n\tActual  : %v", c.usage, expected, actual)
				case actual.Pos != expected.Pos:
					t.Errorf("[Tokenize '%s']: token pos mismatch:\n\tExpected: %v\n\tActual  : %v", c.usage, expected, actual)
				}
			}
		})
	}
}

func TestTokenizeErrors(t *testing.T) {
	cases := []struct {
		usage string
		pos   int
	}{
		{".", 1},
		{"A.", 2},
		{"A.x", 2},
		{"..", 2},
		{"ARG..", 5},
		{"ARG..x", 5},
		{"-", 1},
		{"---x", 2},
		{"-x-", 2},

		{"=", 1},
		{"=<", 2},
		{"=<dsdf", 6},
		{"=<>", 2},
		{"a", 0},
		{"ARg", 2},
		{"1ARG", 0},
	}

	for _, c := range cases {
		t.Run(c.usage, func(t *testing.T) {
			t.Logf("test case %q", c.usage)
			tks, err := Tokenize(c.usage)

			require.Errorf(t, err, "Tokenize('%s') should have failed, instead got %v", c.usage, tks)

			perr, ok := err.(*ParseError)

			require.True(t, ok, "The returned error should be a *ParseError but instead got %#v", err)

			t.Logf("Got expected error %v", err)
			if perr.Pos != c.pos {
				t.Errorf("[Tokenize '%s']: error pos mismatch:\n\tExpected: %v\n\tActual  : %v", c.usage, c.pos, perr.Pos)

			}
		})
	}
}
