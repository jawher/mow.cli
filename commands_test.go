package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	cases := []struct {
		opts   func(*Cmd)
		input  []string
		output []string
	}{
		{func(cmd *Cmd) {
		}, []string{}, []string{}},
		{func(cmd *Cmd) {
			cmd.BoolOpt("a", true, "", nil)
		}, []string{"ab"}, []string{"ab"}},
		{func(cmd *Cmd) {
			cmd.BoolOpt("a", true, "", nil)
		}, []string{"-a"}, []string{"-a"}},
		{func(cmd *Cmd) {
			cmd.BoolOpt("a", true, "", nil)
			cmd.BoolOpt("b", true, "", nil)
		}, []string{"-ab"}, []string{"-a", "-b"}},
		{func(cmd *Cmd) {
			cmd.StringOpt("s", "", "", nil)
		}, []string{"-shello"}, []string{"-s", "hello"}},
		{func(cmd *Cmd) {
			cmd.BoolOpt("a", true, "", nil)
			cmd.StringOpt("b", "", "", nil)
		}, []string{"-ab", "test"}, []string{"-a", "-b", "test"}},
	}
	for _, cas := range cases {
		cmd := &Cmd{optionsIdx: map[string]*option{}}
		cas.opts(cmd)
		nz, cons, err := cmd.normalize(cas.input)
		require.Nil(t, err, "Unexpected error %v", err)
		t.Logf("%v %v", nz, cons)
		require.Equal(t, cas.output, nz)

	}
}
