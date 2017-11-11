package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShortcut(t *testing.T) {
	pc := &ParseContext{}
	args := []string{"a", "b"}
	ok, nargs := theShortcut.Match(args, pc)
	require.True(t, ok, "shortcut always matches")
	require.Equal(t, args, nargs, "shortcut doesn't touch the passed args")
	require.Equal(t, "*", theShortcut.String())
}
