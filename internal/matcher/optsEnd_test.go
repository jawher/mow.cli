package matcher

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOptsEnd(t *testing.T) {
	require.Equal(t, "--", theOptsEnd.String())

	pc := &ParseContext{}
	args := []string{"a", "b"}

	ok, nargs := theOptsEnd.Match(args, pc)

	require.True(t, ok, "optsEnd always matches")
	require.Equal(t, args, nargs, "optsEnd doesn't touch the passed args")
	require.True(t, pc.RejectOptions, "optsEnd sets the rejectOptions flag")
}
