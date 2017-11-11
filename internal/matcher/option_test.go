package matcher

import (
	"testing"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/jawher/mow.cli/internal/values"
	"github.com/stretchr/testify/require"
)

func TestBoolOptMatcher(t *testing.T) {
	forceOpt := &container.Container{Names: []string{"-f", "--force"}, Value: values.NewBool(new(bool), false)}

	optMatcher := opt{
		theOne: forceOpt,
		index: map[string]*container.Container{
			"-f":      forceOpt,
			"--force": forceOpt,
			"-g":      {Names: []string{"-g"}, Value: values.NewBool(new(bool), false)},
			"-x":      {Names: []string{"-x"}, Value: values.NewBool(new(bool), false)},
			"-y":      {Names: []string{"-y"}, Value: values.NewBool(new(bool), false)},
		},
	}

	require.Equal(t, "-f", optMatcher.String())

	cases := []struct {
		args  []string
		nargs []string
		val   []string
	}{
		{[]string{"-f", "x"}, []string{"x"}, []string{"true"}},
		{[]string{"-f=true", "x"}, []string{"x"}, []string{"true"}},
		{[]string{"-f=false", "x"}, []string{"x"}, []string{"false"}},
		{[]string{"--force", "x"}, []string{"x"}, []string{"true"}},
		{[]string{"--force=true", "x"}, []string{"x"}, []string{"true"}},
		{[]string{"--force=false", "x"}, []string{"x"}, []string{"false"}},
		{[]string{"-fgxy", "x"}, []string{"-gxy", "x"}, []string{"true"}},
		{[]string{"-gfxy", "x"}, []string{"-gxy", "x"}, []string{"true"}},
		{[]string{"-gxfy", "x"}, []string{"-gxy", "x"}, []string{"true"}},
		{[]string{"-gxyf", "x"}, []string{"-gxy", "x"}, []string{"true"}},
	}
	for _, cas := range cases {
		t.Logf("Testing case: %#v", cas)
		pc := New()
		ok, nargs := optMatcher.Match(cas.args, &pc)
		require.True(t, ok, "opt should match")
		require.Equal(t, cas.nargs, nargs, "opt should consume the option name")
		require.Equal(t, cas.val, pc.Opts[forceOpt], "true should stored as the option's value")

		pc = New()
		pc.RejectOptions = true
		nok, _ := optMatcher.Match(cas.args, &pc)
		require.False(t, nok, "opt shouldn't match when rejectOptions flag is set")
	}
}

func TestOptMatcher(t *testing.T) {
	names := []string{"-f", "--force"}
	opts := []*container.Container{
		{Names: names, Value: values.NewString(new(string), "")},
		{Names: names, Value: values.NewInt(new(int), 0)},
		{Names: names, Value: values.NewStrings(new([]string), nil)},
		{Names: names, Value: values.NewInts(new([]int), nil)},
	}

	cases := []struct {
		args  []string
		nargs []string
		val   []string
	}{
		{[]string{"-f", "x"}, []string{}, []string{"x"}},
		{[]string{"-f=x", "y"}, []string{"y"}, []string{"x"}},
		{[]string{"-fx", "y"}, []string{"y"}, []string{"x"}},
		{[]string{"-afx", "y"}, []string{"-a", "y"}, []string{"x"}},
		{[]string{"-af", "x", "y"}, []string{"-a", "y"}, []string{"x"}},
		{[]string{"--force", "x"}, []string{}, []string{"x"}},
		{[]string{"--force=x", "y"}, []string{"y"}, []string{"x"}},
	}

	for _, cas := range cases {
		for _, forceOpt := range opts {
			t.Logf("Testing with args %#v and optValue %#v", cas.args, forceOpt.Value)
			optMatcher := opt{
				theOne: forceOpt,
				index: map[string]*container.Container{
					"-f":      forceOpt,
					"--force": forceOpt,
					"-a":      {Names: []string{"-a"}, Value: values.NewBool(new(bool), false)},
				},
			}

			require.Equal(t, "-f", optMatcher.String())

			pc := New()
			ok, nargs := optMatcher.Match(cas.args, &pc)
			require.True(t, ok, "opt %#v should match args %v, %v", forceOpt, cas.args, values.IsBool(forceOpt.Value))
			require.Equal(t, cas.nargs, nargs, "opt should consume the option name")
			require.Equal(t, cas.val, pc.Opts[forceOpt], "true should stored as the option's value")

			pc = New()
			pc.RejectOptions = true
			nok, _ := optMatcher.Match(cas.args, &pc)
			require.False(t, nok, "opt shouldn't match when rejectOptions flag is set")
		}
	}
}

func TestOptNegatives(t *testing.T) {
	names := []string{"-f", "--force"}
	opts := []*container.Container{
		{Names: names, Value: values.NewString(new(string), "")},
		{Names: names, Value: values.NewInt(new(int), 0)},
		{Names: names, Value: values.NewStrings(new([]string), nil)},
		{Names: names, Value: values.NewInts(new([]int), nil)},
	}

	cases := []struct {
		args []string
	}{
		{[]string{"-"}},
		{[]string{"-", "x"}},
		{[]string{"--", "y"}},
		{[]string{"-c"}},
		{[]string{"--qui"}},
		{[]string{"-b"}},
		{[]string{"-b", "-z"}},
		{[]string{"f", "-z"}},
		{[]string{"-f="}},
		{[]string{"--force="}},
		{[]string{"-b=", "-z"}},
	}

	for _, cas := range cases {
		for _, forceOpt := range opts {
			t.Logf("Testing args %#v with optValue: %#v", cas.args, forceOpt.Value)
			optMatcher := opt{
				theOne: forceOpt,
				index: map[string]*container.Container{
					"-f":      forceOpt,
					"--force": forceOpt,
					"-a":      {Names: []string{"-a"}, Value: values.NewBool(new(bool), false)},
					"-b":      {Names: []string{"-a"}, Value: values.NewString(new(string), "")},
				},
			}

			pc := New()
			ok, nargs := optMatcher.Match(cas.args, &pc)
			require.False(t, ok, "opt %#v should not match args %v, %v", forceOpt, cas.args, values.IsBool(forceOpt.Value))
			require.Equal(t, cas.args, nargs, "opt should not have consumed anything")
		}
	}
}
