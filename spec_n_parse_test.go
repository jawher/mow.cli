package cli

import (
	"flag"

	"github.com/stretchr/testify/require"

	"testing"
)

func okCmd(t *testing.T, spec string, init CmdInitializer, args []string) {
	cmd := &Cmd{
		name:       "test",
		optionsIdx: map[string]*opt{},
		argsIdx:    map[string]*arg{},
	}
	cmd.Spec = spec
	cmd.ErrorHandling = flag.ContinueOnError
	init(cmd)

	err := cmd.doInit()
	require.Nil(t, err, "should parse")
	t.Logf("testing spec %s with args: %v", spec, args)
	err = cmd.parse(args)
	require.Nil(t, err, "cmd parse should't fail")
}

func failCmd(t *testing.T, spec string, init CmdInitializer, args []string) {
	cmd := &Cmd{
		name:       "test",
		optionsIdx: map[string]*opt{},
		argsIdx:    map[string]*arg{},
	}
	cmd.Spec = spec
	cmd.ErrorHandling = flag.ContinueOnError
	init(cmd)

	err := cmd.doInit()
	require.Nil(t, err, "should parse")
	t.Logf("testing spec %s with args: %v", spec, args)
	err = cmd.parse(args)
	require.NotNil(t, err, "cmd parse should have failed")
}

func TestSpecBoolOpt(t *testing.T) {
	var f *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
	}
	spec := "-f"
	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)

	badCases := [][]string{
		{},
		{"-g"},
		{"-f", "-g"},
		{"-g", "-f"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecStrOpt(t *testing.T) {
	var f *string
	init := func(c *Cmd) {
		f = c.StringOpt("f", "", "")
	}
	spec := "-f"

	cases := [][]string{
		{"-fValue"},
		{"-f", "Value"},
		{"-f=Value"},
	}
	for _, args := range cases {
		okCmd(t, spec, init, args)
		require.Equal(t, "Value", *f)
	}

	badCases := [][]string{
		{},
		{"-g"},
		{"-f", "-g"},
		{"-g", "-f"},
		{"-f", "xxx", "yyy"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecIntOpt(t *testing.T) {
	var f *int
	init := func(c *Cmd) {
		f = c.IntOpt("f", -1, "")
	}

	spec := "-f"
	cases := [][]string{
		{"-f42"},
		{"-f", "42"},
		{"-f=42"},
	}
	for _, args := range cases {
		okCmd(t, spec, init, args)
		require.Equal(t, 42, *f)
	}

	badCases := [][]string{
		{},
		{"-f", "x"},
		{"-g"},
		{"-f", "-g"},
		{"-g", "-f"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecStrsOpt(t *testing.T) {
	var f *[]string
	init := func(c *Cmd) {
		f = c.StringsOpt("f", nil, "")
	}
	spec := "-f..."
	cases := [][]string{
		{"-fA"},
		{"-f", "A"},
		{"-f=A"},
	}
	for _, args := range cases {
		okCmd(t, spec, init, args)
		require.Equal(t, []string{"A"}, *f)
	}

	cases = [][]string{
		{"-fA", "-f", "B"},
		{"-f", "A", "-f", "B"},
		{"-f=A", "-fB"},
	}
	for _, args := range cases {
		okCmd(t, spec, init, args)
		require.Equal(t, []string{"A", "B"}, *f)
	}

	badCases := [][]string{
		{},
		{"-g"},
		{"-f", "-g"},
		{"-g", "-f"},
		{"-f", "xxx", "yyy"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecIntsOpt(t *testing.T) {
	var f *[]int
	init := func(c *Cmd) {
		f = c.IntsOpt("f", nil, "")
	}
	spec := "-f..."
	cases := [][]string{
		{"-f1"},
		{"-f", "1"},
		{"-f=1"},
	}
	for _, args := range cases {
		okCmd(t, spec, init, args)
		require.Equal(t, []int{1}, *f)
	}

	cases = [][]string{
		{"-f1", "-f", "2"},
		{"-f", "1", "-f", "2"},
		{"-f=1", "-f2"},
	}
	for _, args := range cases {
		okCmd(t, spec, init, args)
		require.Equal(t, []int{1, 2}, *f)
	}

	badCases := [][]string{
		{},
		{"-f", "b"},
		{"-f", "3", "-f", "c"},
		{"-g"},
		{"-f", "-g"},
		{"-g", "-f"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecOptionalOpt(t *testing.T) {
	var f *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
	}
	spec := "[-f]"
	okCmd(t, "[-f]", init, []string{"-f"})
	require.True(t, *f)

	okCmd(t, spec, init, []string{})
	require.False(t, *f)

	badCases := [][]string{
		{"-g"},
		{"-f", "-g"},
		{"-g", "-f"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecArg(t *testing.T) {
	var s *string
	init := func(c *Cmd) {
		s = c.StringArg("ARG", "", "")
	}
	spec := "ARG"
	okCmd(t, spec, init, []string{"value"})
	require.Equal(t, "value", *s)

	badCases := [][]string{
		{},
		{"-g"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecOptionalArg(t *testing.T) {
	var s *string
	init := func(c *Cmd) {
		s = c.StringArg("ARG", "", "")
	}
	spec := "[ARG]"

	okCmd(t, spec, init, []string{"value"})
	require.Equal(t, "value", *s)

	okCmd(t, spec, init, []string{})
	require.Equal(t, "", *s)

	badCases := [][]string{
		{"-g"},
		{"-g", "-f"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}

}

func TestSpecOptionChoice(t *testing.T) {
	var f, g *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
		g = c.BoolOpt("g", false, "")
	}
	spec := "-f|-g"

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-g"})
	require.False(t, *f)
	require.True(t, *g)

	badCases := [][]string{
		{},
		{"-f", "-g"},
		{"-f", "-s"},
		{"-g", "-f"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecOptional2OptionChoice(t *testing.T) {
	var f, g *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
		g = c.BoolOpt("g", false, "")
	}
	spec := "[-f|-g]"

	okCmd(t, spec, init, []string{})
	require.False(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-g"})
	require.False(t, *f)
	require.True(t, *g)

	badCases := [][]string{
		{"-s"},
		{"-f", "-g"},
		{"-g", "-f"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecRepeatable2OptionChoice(t *testing.T) {
	var f, g *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
		g = c.BoolOpt("g", false, "")
	}
	spec := "(-f|-g)..."

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-g"})
	require.False(t, *f)
	require.True(t, *g)

	okCmd(t, spec, init, []string{"-f", "-g"})
	require.True(t, *f)
	require.True(t, *g)

	okCmd(t, spec, init, []string{"-g", "-f"})
	require.True(t, *f)
	require.True(t, *g)

	badCases := [][]string{
		{"-s"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecRepeatableOptional2OptionChoice(t *testing.T) {
	var f, g *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
		g = c.BoolOpt("g", false, "")
	}
	spec := "[-f|-g]..."

	okCmd(t, spec, init, []string{})
	require.False(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-g"})
	require.False(t, *f)
	require.True(t, *g)

	okCmd(t, spec, init, []string{"-f", "-g"})
	require.True(t, *f)
	require.True(t, *g)

	okCmd(t, spec, init, []string{"-g", "-f"})
	require.True(t, *f)
	require.True(t, *g)

	badCases := [][]string{
		{"-s"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestSpecOption3Choice(t *testing.T) {
	var f, g, h *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
		g = c.BoolOpt("g", false, "")
		h = c.BoolOpt("x", false, "")
	}
	spec := "-f|-g|-x"

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
	require.False(t, *g)
	require.False(t, *h)

	okCmd(t, spec, init, []string{"-g"})
	require.False(t, *f)
	require.True(t, *g)
	require.False(t, *h)

	okCmd(t, spec, init, []string{"-x"})
	require.False(t, *f)
	require.False(t, *g)
	require.True(t, *h)
}

func TestSpecOptionalOption3Choice(t *testing.T) {
	var f, g, h *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
		g = c.BoolOpt("g", false, "")
		h = c.BoolOpt("x", false, "")
	}
	spec := "[-f|-g|-x]"

	okCmd(t, spec, init, []string{})
	require.False(t, *f)
	require.False(t, *g)
	require.False(t, *h)

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
	require.False(t, *g)
	require.False(t, *h)

	okCmd(t, spec, init, []string{"-g"})
	require.False(t, *f)
	require.True(t, *g)
	require.False(t, *h)

	okCmd(t, spec, init, []string{"-x"})
	require.False(t, *f)
	require.False(t, *g)
	require.True(t, *h)
}

func TestSpecC1(t *testing.T) {
	var f, g *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
		g = c.BoolOpt("g", false, "")
	}
	spec := "-f|-g..."
	// spec = "[-f|-g...]"

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-g"})
	require.False(t, *f)
	require.True(t, *g)

	okCmd(t, spec, init, []string{"-g", "-g"})
	require.False(t, *f)
	require.True(t, *g)
}

func TestSpecC2(t *testing.T) {
	var f, g *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f", false, "")
		g = c.BoolOpt("g", false, "")
	}
	spec := "[-f|-g...]"

	okCmd(t, spec, init, []string{})
	require.False(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
	require.False(t, *g)

	okCmd(t, spec, init, []string{"-g"})
	require.False(t, *f)
	require.True(t, *g)

	okCmd(t, spec, init, []string{"-g", "-g"})
	require.False(t, *f)
	require.True(t, *g)
}

func TestSpecCpCase(t *testing.T) {
	var f, g *[]string
	init := func(c *Cmd) {
		f = c.StringsArg("SRC", nil, "")
		g = c.StringsArg("DST", nil, "")
	}
	spec := "SRC... DST"

	okCmd(t, spec, init, []string{"A", "B"})
	require.Equal(t, []string{"A"}, *f)
	require.Equal(t, []string{"B"}, *g)

	okCmd(t, spec, init, []string{"A", "B", "C"})
	require.Equal(t, []string{"A", "B"}, *f)
	require.Equal(t, []string{"C"}, *g)

	okCmd(t, spec, init, []string{"A", "B", "C", "D"})
	require.Equal(t, []string{"A", "B", "C"}, *f)
	require.Equal(t, []string{"D"}, *g)
}

func TestSpecC3(t *testing.T) {
	var f, g *[]string
	init := func(c *Cmd) {
		f = c.StringsArg("SRC", nil, "")
		g = c.StringsArg("DST", nil, "")
	}
	spec := "(SRC... DST) | SRC"

	okCmd(t, spec, init, []string{"A"})
	require.Equal(t, []string{"A"}, *f)
	require.Equal(t, 0, len(*g))

	okCmd(t, spec, init, []string{"A", "B"})
	require.Equal(t, []string{"A"}, *f)
	require.Equal(t, []string{"B"}, *g)

	okCmd(t, spec, init, []string{"A", "B", "C"})
	require.Equal(t, []string{"A", "B"}, *f)
	require.Equal(t, []string{"C"}, *g)

	okCmd(t, spec, init, []string{"A", "B", "C", "D"})
	require.Equal(t, []string{"A", "B", "C"}, *f)
	require.Equal(t, []string{"D"}, *g)
}

func TestSpecC5(t *testing.T) {
	var f, g *[]string
	var x *bool
	init := func(c *Cmd) {
		f = c.StringsArg("SRC", nil, "")
		g = c.StringsArg("DST", nil, "")
		x = c.BoolOpt("x", false, "")
	}
	spec := "(SRC... -x DST) | (SRC... DST)"

	okCmd(t, spec, init, []string{"A", "B"})
	require.Equal(t, []string{"A"}, *f)
	require.Equal(t, []string{"B"}, *g)
	require.False(t, *x)

	okCmd(t, spec, init, []string{"A", "B", "C"})
	require.Equal(t, []string{"A", "B"}, *f)
	require.Equal(t, []string{"C"}, *g)
	require.False(t, *x)

	okCmd(t, spec, init, []string{"A", "B", "-x", "C"})
	require.Equal(t, []string{"A", "B"}, *f)
	require.Equal(t, []string{"C"}, *g)
	require.True(t, *x)

}
