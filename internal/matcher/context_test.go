package matcher

import (
	"testing"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	c := New()

	require.NotNil(t, c.Args)
	require.Empty(t, c.Args)

	require.NotNil(t, c.Opts)
	require.Empty(t, c.Opts)

	require.NotNil(t, c.ExcludedOpts)
	require.Empty(t, c.ExcludedOpts)

	require.False(t, c.RejectOptions)
}

func TestMerge(t *testing.T) {
	var (
		a1 = &container.Container{}
		a2 = &container.Container{}
		a3 = &container.Container{}
		o1 = &container.Container{}
		o2 = &container.Container{}
		o3 = &container.Container{}
	)
	c1, c2 := New(), New()

	c1.Args[a1] = []string{"a1"}
	c1.Args[a3] = []string{"a3"}
	c1.Opts[o1] = []string{"o1"}
	c1.Opts[o3] = []string{"o3"}

	c2.Args[a1] = []string{"a1.2"}
	c2.Args[a2] = []string{"a2"}
	c2.Opts[o1] = []string{"o1.2"}
	c2.Opts[o2] = []string{"o2"}

	c1.Merge(c2)

	require.Equal(t, map[*container.Container][]string{
		a1: {"a1", "a1.2"},
		a2: {"a2"},
		a3: {"a3"},
	}, c1.Args)

	require.Equal(t, map[*container.Container][]string{
		o1: {"o1", "o1.2"},
		o2: {"o2"},
		o3: {"o3"},
	}, c1.Opts)
}
