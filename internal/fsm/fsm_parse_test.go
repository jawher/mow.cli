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

func TestParseErrorsWhenApplyReturnsFalse(t *testing.T) {
	s := fsm.NewState()
	err := s.Parse(nil)

	require.Error(t, err)
}

func TestParseNoErrorWhenApplyReturnsTrue(t *testing.T) {
	s := fsm.NewState()
	s.Terminal = true

	err := s.Parse(nil)

	require.NoError(t, err)
}

func TestParseFillsContainers(t *testing.T) {
	var (
		boolSetByUser    = false
		boolVar          = false
		stringsSetByUser = false
		stringsVar       []string
		boolCon          = &container.Container{
			Value:           values.NewBool(&boolVar, false),
			ValueSetFromEnv: true,
			ValueSetByUser:  &boolSetByUser,
		}
		stringsCon = &container.Container{
			Value:           values.NewStrings(&stringsVar, []string{"original", "value"}),
			ValueSetFromEnv: true,
			ValueSetByUser:  &stringsSetByUser,
		}
	)
	matchers := map[string]matcher.Matcher{
		"^": fsmtest.TestMatcher{
			TestPriority: 2,
			MatchFunc: func(args []string, c *matcher.ParseContext) (bool, []string) {
				c.Opts[boolCon] = []string{"true"}
				c.Args[stringsCon] = []string{"new", "value"}
				return true, nil
			},
		},
	}
	s := fsmtest.NewFsm(`
		S1 ^ (S2)
	`, matchers)

	s.Prepare()

	err := s.Parse([]string{"something"})

	require.NoError(t, err)

	require.False(t, boolCon.ValueSetFromEnv)
	require.True(t, boolSetByUser)
	require.True(t, boolVar)

	require.False(t, stringsCon.ValueSetFromEnv)
	require.True(t, stringsSetByUser)
	require.Equal(t, stringsVar, []string{"new", "value"})
}
