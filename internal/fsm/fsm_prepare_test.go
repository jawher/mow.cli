package fsm_test

import (
	"testing"

	"github.com/jawher/mow.cli/internal/fsm/fsmtest"
	"github.com/jawher/mow.cli/internal/matcher"
	"github.com/jawher/mow.cli/internal/matcher/matchertest"
	"github.com/stretchr/testify/require"
)

var (
	testMatchers = map[string]matcher.Matcher{
		"*":    matcher.NewShortcut(),
		"!":    fsmtest.NopeMatcher{},
		"-a":   matchertest.NewOpt("-a"),
		"ARG":  matchertest.NewArg("ARG"),
		"--":   matcher.NewOptsEnd(),
		"-abc": matchertest.NewOptions("-abc"),
	}
)

func TestSimplify(t *testing.T) {

	cases := []struct {
		original   string
		simplified string
	}{

		{
			original: `
					S1 * S2
					S2 -a (S3)
			`,
			simplified: "S1 -a (S3)",
		},
		{
			// seq like FSM
			original: `
					S1 *  S2
					S2 -a S3
					S3 *  S2
					S3 *  (S4)
			`,
			simplified: `
					S1 -a (S3)
					(S3) -a (S3)
			`,
		},
		{
			// optional transition FSM
			original: `
					S1 -a  S2
					S2 * (S3)
					S1 *  (S3)
			`,
			simplified: `
					(S1) -a (S2)
			`,
		},
	}

	for _, cas := range cases {
		t.Logf("FSM.simplify:\noriginal: %s\nexpected: %s", cas.original, cas.simplified)

		original := fsmtest.NewFsm(cas.original, testMatchers)

		original.Prepare()

		simplified := fsmtest.NewFsm(cas.simplified, testMatchers)

		require.Equal(t, fsmtest.FsmStr(simplified), fsmtest.FsmStr(original))
	}
}

func TestSort(t *testing.T) {
	s := fsmtest.NewFsm(`
		S1 * S2
		S1 * S3
		S1 ARG S2
		S1 -- S2
		S1 -abc S2
		S1 -a S3
	`, testMatchers)

	s.Prepare()

	require.Equal(t, []string{"-a", "-abc", "ARG", "--"}, fsmtest.TransitionStrs(s.Transitions))
}
