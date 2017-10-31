package matcher

import (
	"testing"
	"time"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/jawher/mow.cli/internal/values"
	"github.com/stretchr/testify/require"
)

func TestOptsMatcher(t *testing.T) {
	opts := options{
		options: []*container.Container{
			{Names: []string{"-f", "--force"}, Value: values.NewBool(new(bool), false)},
			{Names: []string{"-g", "--green"}, Value: values.NewString(new(string), "")},
		},
		index: map[string]*container.Container{},
	}

	for _, o := range opts.options {
		for _, n := range o.Names {
			opts.index[n] = o
		}
	}

	cases := []struct {
		args  []string
		nargs []string
		val   [][]string
	}{
		{[]string{"-f", "x"}, []string{"x"}, [][]string{{"true"}, nil}},
		{[]string{"-f=false", "y"}, []string{"y"}, [][]string{{"false"}, nil}},
		{[]string{"--force", "x"}, []string{"x"}, [][]string{{"true"}, nil}},
		{[]string{"--force=false", "y"}, []string{"y"}, [][]string{{"false"}, nil}},

		{[]string{"-g", "x"}, []string{}, [][]string{nil, {"x"}}},
		{[]string{"-g=x", "y"}, []string{"y"}, [][]string{nil, {"x"}}},
		{[]string{"-gx", "y"}, []string{"y"}, [][]string{nil, {"x"}}},
		{[]string{"--green", "x"}, []string{}, [][]string{nil, {"x"}}},
		{[]string{"--green=x", "y"}, []string{"y"}, [][]string{nil, {"x"}}},

		{[]string{"-f", "-g", "x", "y"}, []string{"y"}, [][]string{{"true"}, {"x"}}},
		{[]string{"-g", "x", "-f", "y"}, []string{"y"}, [][]string{{"true"}, {"x"}}},
		{[]string{"-fg", "x", "y"}, []string{"y"}, [][]string{{"true"}, {"x"}}},
		{[]string{"-fgxxx", "y"}, []string{"y"}, [][]string{{"true"}, {"xxx"}}},
	}

	for _, cas := range cases {
		t.Logf("testing with args %#v", cas.args)
		pc := New()
		ok, nargs := opts.Match(cas.args, &pc)
		require.True(t, ok, "opts should match")
		require.Equal(t, cas.nargs, nargs, "opts should consume the option name")
		for i, opt := range opts.options {
			require.Equal(t, cas.val[i], pc.Opts[opt], "the option value for %v should be stored", opt)
		}

		pc = New()
		pc.RejectOptions = true
		nok, _ := opts.Match(cas.args, &pc)
		require.False(t, nok, "opts shouldn't match when rejectOptions flag is set")
	}
}

// Issue 55
func TestOptsMatcherInfiniteLoop(t *testing.T) {
	opts := options{
		options: []*container.Container{
			{Names: []string{"-g"}, Value: values.NewString(new(string), ""), ValueSetFromEnv: true},
		},
		index: map[string]*container.Container{},
	}

	for _, o := range opts.options {
		for _, n := range o.Names {
			opts.index[n] = o
		}
	}

	done := make(chan struct{}, 1)
	pc := New()
	go func() {
		opts.Match([]string{"-x"}, &pc)
		done <- struct{}{}
	}()

	select {
	case <-done:
		// nop, everything is good
	case <-time.After(5 * time.Second):
		t.Fatalf("Timed out after 5 seconds. Infinite loop in optsMatcher.")
	}

}
