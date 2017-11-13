package fsm_test

import (
	"testing"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/jawher/mow.cli/internal/fsm"
	"github.com/jawher/mow.cli/internal/fsm/fsmtest"
	"github.com/jawher/mow.cli/internal/matcher"
	"github.com/jawher/mow.cli/internal/values"
	"github.com/stretchr/testify/require"
)

func TestApplyTerminalStateNoArgs(t *testing.T) {
	s := fsm.NewState()
	s.Terminal = true

	err := s.Parse(nil)

	require.NoError(t, err)
}

func TestApply(t *testing.T) {
	var (
		testArgs = []string{"1", "2", "3"}
		optStrs  []string
		optCon   = &container.Container{
			Value: values.NewStrings(&optStrs, nil),
		}
		argStrs []string
		argCon  = &container.Container{
			Value: values.NewStrings(&argStrs, nil),
		}
		calls []string
	)
	matchers := map[string]matcher.Matcher{
		"a": fsmtest.TestMatcher{
			TestPriority: 2,
			MatchFunc: func(args []string, c *matcher.ParseContext) (bool, []string) {
				require.Equal(t, testArgs, args)

				calls = append(calls, "a")

				c.Opts[optCon] = []string{"a.opt"}
				c.Args[argCon] = []string{"a.arg"}
				return true, args[1:]
			},
		},
		"b": fsmtest.TestMatcher{
			TestPriority: 1,
			MatchFunc: func(args []string, c *matcher.ParseContext) (bool, []string) {
				require.Equal(t, testArgs, args)

				calls = append(calls, "b")

				c.Opts[optCon] = []string{"b.opt"}
				c.Args[argCon] = []string{"b.arg"}
				return true, args[1:]
			},
		},
		"c": fsmtest.TestMatcher{
			TestPriority: 1,
			MatchFunc: func(args []string, c *matcher.ParseContext) (bool, []string) {
				require.Equal(t, testArgs[1:], args, "second stage matchers should be called with the rem args")

				calls = append(calls, "c")

				c.Opts[optCon] = []string{"c.opt"}
				c.Args[argCon] = []string{"c.arg"}
				return true, args[1:]
			},
		},
		"d": fsmtest.TestMatcher{
			TestPriority: 1,
			MatchFunc: func(args []string, c *matcher.ParseContext) (bool, []string) {
				require.Equal(t, testArgs[1:], args, "second stage matchers should be called with the rem args")

				calls = append(calls, "d")

				c.Opts[optCon] = []string{"d.opt"}
				c.Args[argCon] = []string{"d.arg"}
				return false, args[1:]
			},
		},
		"e": fsmtest.TestMatcher{
			TestPriority: 1,
			MatchFunc: func(args []string, c *matcher.ParseContext) (bool, []string) {
				require.Equal(t, testArgs[2:], args, "third stage matchers should be called with the rem args")

				calls = append(calls, "e")

				c.Opts[optCon] = []string{"e.opt"}
				c.Args[argCon] = []string{"e.arg"}
				return true, nil
			},
		},
	}
	s := fsmtest.NewFsm(`
		S1 a S2
		S1 b S3

		S2 c S4
		S3 d S4

		S4 e (S5)
	`, matchers)

	s.Prepare()

	err := s.Parse(testArgs)

	require.Equal(t, []string{"b", "a", "d", "c", "e"}, calls)

	require.NoError(t, err)

	require.Equal(t, []string{"a.opt", "c.opt", "e.opt"}, optStrs)
	require.Equal(t, []string{"a.arg", "c.arg", "e.arg"}, argStrs)
}

func TestApplyRejectOptions(t *testing.T) {
	var (
		testArgs = []string{"1", "--", "2"}
		calls    []string
	)
	matchers := map[string]matcher.Matcher{
		"a": fsmtest.TestMatcher{
			MatchFunc: func(args []string, c *matcher.ParseContext) (bool, []string) {
				require.Equal(t, testArgs, args)
				require.False(t, c.RejectOptions)

				calls = append(calls, "a")
				return true, args[1:]
			},
		},
		"b": fsmtest.TestMatcher{
			MatchFunc: func(args []string, c *matcher.ParseContext) (bool, []string) {
				require.Equal(t, []string{"2"}, args)
				require.True(t, c.RejectOptions)

				calls = append(calls, "b")
				return true, args[1:]
			},
		},
	}
	s := fsmtest.NewFsm(`
		S1 a S2
		S2 b S3
	`, matchers)

	s.Prepare()

	err := s.Parse(testArgs)

	require.NoError(t, err)

	require.Equal(t, []string{"a", "b"}, calls)

}
