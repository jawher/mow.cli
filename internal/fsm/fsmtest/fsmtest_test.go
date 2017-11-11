package fsmtest

import (
	"testing"

	"github.com/jawher/mow.cli/internal/fsm"
	"github.com/jawher/mow.cli/internal/matcher"
	"github.com/stretchr/testify/require"
)

func TestNopeMatcher(t *testing.T) {
	m := NopeMatcher{}

	require.Equal(t, 666, m.Priority())

	pc := matcher.NewParseContext()
	args := []string{"a", "b"}
	ok, rem := m.Match(args, &pc)

	require.False(t, ok)
	require.Equal(t, args, rem)
}

func TestYepMatcher(t *testing.T) {
	m := YepMatcher{}

	require.Equal(t, 666, m.Priority())

	pc := matcher.NewParseContext()
	args := []string{"a", "b"}
	ok, rem := m.Match(args, &pc)

	require.True(t, ok)
	require.Equal(t, args, rem)
}

func TestTestMatcher(t *testing.T) {
	called := false
	args := []string{"a", "b"}

	m := TestMatcher{
		MatchFunc: func(targs []string, c *matcher.ParseContext) (bool, []string) {
			called = true
			require.Equal(t, args, targs)
			return true, targs
		},
		TestPriority: 7,
	}

	require.Equal(t, 7, m.Priority())
	pc := matcher.NewParseContext()
	ok, rem := m.Match(args, &pc)

	require.True(t, called)
	require.True(t, ok)
	require.Equal(t, args, rem)
}

func TestNewFsm(t *testing.T) {
	matchers := map[string]matcher.Matcher{
		"a": YepMatcher{},
		"b": NopeMatcher{},
	}

	s1 := NewFsm(`

		S1 a S2

		S2 b (S3)
	`, matchers)

	require.NotNil(t, s1)
	require.False(t, s1.Terminal)
	require.Equal(t, len(s1.Transitions), 1)

	ta := s1.Transitions[0]
	require.Equal(t, ta.Matcher, matchers["a"])

	s2 := ta.Next
	require.NotNil(t, s2)
	require.False(t, s2.Terminal)
	require.Equal(t, len(s2.Transitions), 1)

	tb := s2.Transitions[0]
	require.Equal(t, tb.Matcher, matchers["b"])

	s3 := tb.Next
	require.NotNil(t, s3)
	require.True(t, s3.Terminal)
	require.Equal(t, len(s3.Transitions), 0)
}

func TestTransitionStrs(t *testing.T) {
	trs := fsm.StateTransitions{
		&fsm.Transition{
			Matcher: NopeMatcher{},
		},
		&fsm.Transition{
			Matcher: YepMatcher{},
		},
	}
	actual := TransitionStrs(trs)

	require.Equal(t, []string{"<nope>", "<yep>"}, actual)
}

func TestFsmStr(t *testing.T) {
	s1 := fsm.NewState()
	s2 := fsm.NewState()
	s3 := fsm.NewState()
	s3.Terminal = true

	s1.T(YepMatcher{}, s2)
	s1.T(YepMatcher{}, s3)
	s2.T(NopeMatcher{}, s3)

	str := FsmStr(s1)

	require.Equal(t, "S1 <yep> (S3)\nS1 <yep> S2\nS2 <nope> (S3)", str)

}
