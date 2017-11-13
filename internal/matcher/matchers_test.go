package matcher

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/require"
)

func TestIsShortcut(t *testing.T) {
	var (
		shortcut = NewShortcut()
		optsEnd  = NewOptsEnd()
		opt      = NewOpt(nil, nil)
		arg      = NewArg(nil)
		options  = NewOptions(nil, nil)
	)

	require.True(t, IsShortcut(shortcut))
	require.False(t, IsShortcut(optsEnd))
	require.False(t, IsShortcut(opt))
	require.False(t, IsShortcut(arg))
	require.False(t, IsShortcut(options))
}

func TestPriority(t *testing.T) {
	var (
		shortcut = NewShortcut()
		optsEnd  = NewOptsEnd()
		opt      = NewOpt(nil, nil)
		arg      = NewArg(nil)
		options  = NewOptions(nil, nil)
	)

	cases := []struct {
		a, b Matcher
	}{
		// shorcut always comes last
		{opt, shortcut},
		{arg, shortcut},
		{options, shortcut},
		{optsEnd, shortcut},

		// then comes optsEnd
		{opt, optsEnd},
		{arg, optsEnd},
		{options, optsEnd},

		// opt comes before the rest
		{opt, options},
		{opt, arg},

		// options comes before arg
		{options, arg},
	}

	for _, cas := range cases {
		t.Run(fmt.Sprintf("%T < %T", cas.a, cas.b), func(t *testing.T) {
			ap := cas.a.Priority()
			bp := cas.b.Priority()
			require.Truef(t, ap < bp, "%#v (priority=%d) should be < then %#v (priority=%d)", cas.a, ap, cas.b, bp)
		})
	}
}

func TestMatchersComparable(t *testing.T) {
	// this stupid-looking test is here for a reason:
	// It ensures that all exposed matchers from this package are comparable types
	matchers := []Matcher{
		NewShortcut(),
		NewOpt(nil, nil),
		NewOptions(nil, nil),
		NewArg(nil),
		NewOptsEnd(),
	}

	for _, m1 := range matchers {
		for _, m2 := range matchers {
			_ = m1 == m2
		}
	}
}
