package matcher

import (
	"testing"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/stretchr/testify/require"
)

func TestArgMatcher(t *testing.T) {
	a := &container.Container{Name: "X"}
	argMatcher := arg{arg: a}

	require.Equal(t, "X", argMatcher.String())

	{
		pc := NewParseContext()
		args := []string(nil)
		ok, nargs := argMatcher.Match(args, &pc)
		require.False(t, ok, "arg should not match")
		require.Nil(t, nargs, "arg should not consume anything")
		require.Nil(t, pc.Args[a], "arg should not store anything")
	}
	{
		pc := NewParseContext()
		args := []string{"a", "b"}
		ok, nargs := argMatcher.Match(args, &pc)
		require.True(t, ok, "arg should match")
		require.Equal(t, []string{"b"}, nargs, "arg should consume the matched value")
		require.Equal(t, []string{"a"}, pc.Args[a], "arg should stored the matched value")
	}
	{
		pc := NewParseContext()
		ok, _ := argMatcher.Match([]string{"-v"}, &pc)
		require.False(t, ok, "arg should not match options")
	}
	{
		pc := NewParseContext()
		pc.RejectOptions = true
		ok, _ := argMatcher.Match([]string{"-v"}, &pc)
		require.True(t, ok, "arg should match options when the reject flag is set")
	}
}
