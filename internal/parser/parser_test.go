package parser

import (
	"testing"

	"strings"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/jawher/mow.cli/internal/fsm/fsmtest"
	"github.com/jawher/mow.cli/internal/lexer"
	"github.com/stretchr/testify/require"
)

var (
	optACon = &container.Container{
		Name:  "-a --all",
		Names: []string{"-a", "--all"},
	}
	optBCon = &container.Container{
		Name:  "-b --ball",
		Names: []string{"-b", "--ball"},
	}
	optsIndex = map[string]*container.Container{
		"-a":     optACon,
		"--all":  optACon,
		"-b":     optBCon,
		"--ball": optBCon,
	}

	argCon = &container.Container{
		Name: "ARG",
	}
	argsIndex = map[string]*container.Container{
		"ARG": argCon,
	}
)

func TestParse(t *testing.T) {
	cases := []struct {
		spec        string
		expectedFsm string
	}{
		{
			spec:        "",
			expectedFsm: "",
		},
		{
			spec:        "-a",
			expectedFsm: "S1 -a (S2)",
		},
		{
			spec:        "--all",
			expectedFsm: "S1 -a (S2)",
		},
		{
			spec: "-a -b",
			expectedFsm: `
						S1 -a S2
						S2 -b (S3)
			`,
		},
		{
			spec: "-a | -b",
			expectedFsm: `
						S1 -a (S2)
						S1 -b (S3)
			`,
		},
		{
			spec: "[ -a ]",
			expectedFsm: `
						(S1) -a (S2)
			`,
		},
		{
			spec: "-a...",
			expectedFsm: `
						S1 -a (S2)
						(S2) -a (S2)
			`,
		},
		{
			spec: "[-a...]",
			expectedFsm: `
						(S1) -a (S2)
						(S2) -a (S2)
			`,
		},
		{
			spec: "[-a]...",
			expectedFsm: `
						(S1) -a (S2)
						(S2) -a (S2)
			`,
		},
		{
			spec: "-a -b | ARG",
			expectedFsm: `
						S1 -a S2
						S2 -b (S3)
						S2 ARG (S4)
			`,
		},
		{
			spec: "-a (-b | ARG)",
			expectedFsm: `
						S1 -a S2
						S2 -b (S3)
						S2 ARG (S4)
			`,
		},
		{
			spec: "( -a -b ) | ARG",
			expectedFsm: `
						S1 -a S2
						S1 ARG (S4)
						S2 -b (S3)
			`,
		},
		{
			spec: "( -a -b ) | ARG...",
			expectedFsm: `
						S1 -a S2
						S1 ARG (S4)
						S2 -b (S3)
						(S4) ARG (S4)
			`,
		},
		{
			spec: "-ab",
			expectedFsm: `
						S1 -ab (S2)
			`,
		},
		{
			spec: "-a -- ARG",
			expectedFsm: `
						S1 -a S2
						S2 -- S3
						S3 ARG (S4)
			`,
		},
		{
			spec: "[OPTIONS]",
			expectedFsm: `
						(S1) -ab (S2)
			`,
		},
	}

	for _, cas := range cases {

		t.Run(cas.spec, func(t *testing.T) {
			t.Logf("Testing spec %q", cas.spec)

			tokens, lerr := lexer.Tokenize(cas.spec)
			require.NoErrorf(t, lerr, "Lexing error %v", lerr)

			s, perr := Parse(tokens, Params{
				Args:       []*container.Container{argCon},
				ArgsIdx:    argsIndex,
				Options:    []*container.Container{optACon, optBCon},
				OptionsIdx: optsIndex,
				Spec:       cas.spec,
			})
			require.NoErrorf(t, perr, "Parsing error %v", perr)

			actualFsm := fsmtest.FsmStr(s)
			require.Equal(t, cleanFsmStr(cas.expectedFsm), actualFsm)
		})
	}
}

func TestParseErrors(t *testing.T) {
	cases := []struct {
		spec string
		msg  string
		pos  int
	}{
		{
			spec: "-c",
			msg:  "Undeclared option -c",
			pos:  0,
		},
		{
			spec: "--close",
			msg:  "Undeclared option --close",
			pos:  0,
		},
		{
			spec: "NOPE",
			msg:  "Undeclared arg NOPE",
			pos:  0,
		},
		{
			spec: "ARG -- -a",
			msg:  "No options after --",
			pos:  7,
		},
		{
			spec: "ARG -- --all",
			msg:  "No options after --",
			pos:  7,
		},
		{
			spec: "ARG -- -ab",
			msg:  "No options after --",
			pos:  7,
		},
		{
			spec: "ARG -- [OPTIONS]",
			msg:  "No options after --",
			pos:  8,
		},
		{
			spec: "-ac",
			msg:  "Undeclared option -c",
			pos:  0,
		},
		{
			spec: ")",
			msg:  "Unexpected input",
			pos:  0,
		},
		{
			spec: "]",
			msg:  "Unexpected input",
			pos:  0,
		},
		{
			spec: "|",
			msg:  "Unexpected input",
			pos:  0,
		},
		{
			spec: "-a |",
			msg:  "Unexpected end of input",
			pos:  4,
		},
		{
			spec: "( -a",
			msg:  "Was expecting ClosePar",
			pos:  4,
		},
	}

	for _, cas := range cases {

		t.Run(cas.spec, func(t *testing.T) {
			t.Logf("Testing spec %q", cas.spec)

			tokens, lerr := lexer.Tokenize(cas.spec)
			require.NoErrorf(t, lerr, "Lexing error %v", lerr)

			_, err := Parse(tokens, Params{
				Args:       []*container.Container{argCon},
				ArgsIdx:    argsIndex,
				Options:    []*container.Container{optACon, optBCon},
				OptionsIdx: optsIndex,
				Spec:       cas.spec,
			})
			require.Error(t, err, "Parsing should have failed")

			perr, ok := err.(*lexer.ParseError)

			require.Truef(t, ok, "the returned error should be of type ParseError, but instead was %T", err)

			t.Logf("Got expected parse error %v", perr)
			require.Contains(t, perr.Msg, cas.msg)
			require.Equal(t, cas.pos, perr.Pos)
		})
	}
}

func cleanFsmStr(s string) string {
	lines := strings.FieldsFunc(s, func(r rune) bool { return r == '\n' })
	clines := make([]string, 0, len(lines))
	for _, line := range lines {
		cline := strings.TrimSpace(line)
		if cline == "" {
			continue
		}
		clines = append(clines, cline)
	}
	//fmt.Printf("%#v\n", lines)
	return strings.Join(clines, "\n")
}
