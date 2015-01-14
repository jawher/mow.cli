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

func badSpec(t *testing.T, spec string, init CmdInitializer) {
	cmd := &Cmd{
		name:       "test",
		optionsIdx: map[string]*opt{},
		argsIdx:    map[string]*arg{},
	}
	cmd.Spec = spec
	cmd.ErrorHandling = flag.ContinueOnError
	init(cmd)

	t.Logf("testing bad spec %s", spec)
	err := cmd.doInit()
	require.NotNil(t, err, "Bad spec %s should have failed to parse", spec)
	t.Logf("Bad spec %s did fail to parse with error: %v", spec, err)
}

func TestSpecBoolOpt(t *testing.T) {
	var f *bool
	init := func(c *Cmd) {
		f = c.BoolOpt("f force", false, "")
	}
	spec := "-f"

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)

	okCmd(t, spec, init, []string{"--force"})
	require.True(t, *f)

	okCmd(t, spec, init, []string{"-f=true"})
	require.True(t, *f)

	okCmd(t, spec, init, []string{"--force=true"})
	require.True(t, *f)

	okCmd(t, spec, init, []string{"--force=false"})
	require.False(t, *f)

	badCases := [][]string{
		{},
		{"-g"},
		{"-f", "-g"},
		{"-g", "-f"},
		{"-f", "true"},
		{"-f", "xxx"},
		{"xxx", "-f"},
	}
	for _, args := range badCases {
		failCmd(t, spec, init, args)
	}
}

func TestDefaultSpec(t *testing.T) {
	var (
		a *bool
		b *string
		c *string
	)

	type call struct {
		args []string
		a    bool
		b, c string
	}
	cases := []struct {
		init  func(cmd *Cmd)
		calls []call
	}{
		{
			func(cmd *Cmd) {
				a = cmd.BoolOpt("a", false, "")
				b = cmd.StringOpt("b", "", "")
				c = cmd.StringArg("C", "", "")
			},
			[]call{
				{[]string{"X"}, false, "", "X"},
				{[]string{"-a", "X"}, true, "", "X"},
				{[]string{"-b=Z", "X"}, false, "Z", "X"},
				{[]string{"-b=Z", "-a", "X"}, true, "Z", "X"},
				{[]string{"-a", "-b=Z", "X"}, true, "Z", "X"},
			},
		},
	}

	for _, cas := range cases {
		for _, cl := range cas.calls {
			okCmd(t, "", cas.init, cl.args)
			require.Equal(t, cl.a, *a)
			require.Equal(t, cl.b, *b)
			require.Equal(t, cl.c, *c)
		}
	}

}

func TestSpecOptFolding(t *testing.T) {
	var a, b, c *bool
	var d *string
	init := func(cmd *Cmd) {
		a = cmd.BoolOpt("a", false, "")
		b = cmd.BoolOpt("b", false, "")
		c = cmd.BoolOpt("c", false, "")

		d = cmd.StringOpt("d", "", "")
	}

	cases := []struct {
		spec    string
		args    []string
		a, b, c bool
		d       string
	}{
		{
			"[-abcd]", []string{},
			false, false, false,
			"",
		},
		{
			"[-abcd]", []string{"-ab"},
			true, true, false,
			"",
		},
		{
			"[-abcd]", []string{"-ba"},
			true, true, false,
			"",
		},

		{
			"[-abcd]", []string{"-ad", "TEST"},
			true, false, false,
			"TEST",
		},
		{
			"[-abcd]", []string{"-abd", "TEST"},
			true, true, false,
			"TEST",
		},
		{
			"[-abcd]", []string{"-abcd", "TEST"},
			true, true, true,
			"TEST",
		},
		{
			"[-abcd]", []string{"-bcd", "TEST"},
			false, true, true,
			"TEST",
		},
		{
			"[-abcd]", []string{"-cbd", "TEST"},
			false, true, true,
			"TEST",
		},
		{
			"[-abcd]", []string{"-ac"},
			true, false, true,
			"",
		},
		{
			"[-abcd]", []string{"-ca"},
			true, false, true,
			"",
		},
		{
			"[-abcd]", []string{"-cab"},
			true, true, true,
			"",
		},
	}

	for _, cas := range cases {
		okCmd(t, cas.spec, init, cas.args)
		require.Equal(t, cas.a, *a)
		require.Equal(t, cas.b, *b)
		require.Equal(t, cas.c, *c)
		require.Equal(t, cas.d, *d)
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

func TestSpecOptionsEndExplicit(t *testing.T) {
	var x *[]string
	init := func(c *Cmd) {
		x = c.StringsArg("X", nil, "")
	}
	spec := "-- X..."

	okCmd(t, spec, init, []string{"A"})
	require.Equal(t, []string{"A"}, *x)

	okCmd(t, spec, init, []string{"--", "A"})
	require.Equal(t, []string{"A"}, *x)

	okCmd(t, spec, init, []string{"--", "-x"})
	require.Equal(t, []string{"-x"}, *x)

	okCmd(t, spec, init, []string{"--", "A", "B"})
	require.Equal(t, []string{"A", "B"}, *x)
}

func TestSpecOptionsEndImplicit(t *testing.T) {
	var x *[]string
	var f *bool
	init := func(c *Cmd) {
		x = c.StringsArg("X", nil, "")
		f = c.BoolOpt("f", false, "")
	}
	spec := "-f|X..."

	okCmd(t, spec, init, []string{"A"})
	require.Equal(t, []string{"A"}, *x)

	okCmd(t, spec, init, []string{"--", "A"})
	require.Equal(t, []string{"A"}, *x)

	okCmd(t, spec, init, []string{"--", "-f"})
	require.Equal(t, []string{"-f"}, *x)

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
}

func TestSpecChoiceWithOptionsEndInLastPos(t *testing.T) {
	var x *[]string
	var f *bool
	init := func(c *Cmd) {
		x = c.StringsArg("X", nil, "")
		f = c.BoolOpt("f", false, "")
	}
	spec := "-f|(-- X...)"

	okCmd(t, spec, init, []string{"A", "B"})
	require.Equal(t, []string{"A", "B"}, *x)

	okCmd(t, spec, init, []string{"--", "-f", "B"})
	require.Equal(t, []string{"-f", "B"}, *x)

	okCmd(t, spec, init, []string{"--", "A", "B"})
	require.Equal(t, []string{"A", "B"}, *x)

	okCmd(t, spec, init, []string{"-f"})
	require.True(t, *f)
}

func TestSpecChoiceWithOptionsEnd(t *testing.T) {
	var x *[]string
	var f *bool
	init := func(c *Cmd) {
		x = c.StringsArg("X", nil, "")
		f = c.BoolOpt("f", false, "")
	}
	spec := "(-- X...)|-f"

	badSpec(t, spec, init)
}

func TestSpecOptionAfterOptionsEnd(t *testing.T) {
	var x *[]string
	var f *bool
	init := func(c *Cmd) {
		x = c.StringsArg("X", nil, "")
		f = c.BoolOpt("f", false, "")
	}

	spec := "-- X... -f"
	badSpec(t, spec, init)
}

func TestSpecOptionAfterOptionsEndInAChoice(t *testing.T) {
	var x *[]string
	var f, d *bool
	init := func(c *Cmd) {
		x = c.StringsArg("X", nil, "")
		f = c.BoolOpt("f", false, "")
		d = c.BoolOpt("d", false, "")
	}

	spec := "-f | (-- X...) -d"
	badSpec(t, spec, init)
}

func TestSpecOptionAfterOptionalOptionsEnd(t *testing.T) {
	var x *[]string
	var f, d *bool
	init := func(c *Cmd) {
		x = c.StringsArg("X", nil, "")
		f = c.BoolOpt("f", false, "")
		d = c.BoolOpt("d", false, "")
	}

	spec := "-f [-- X] -d"
	badSpec(t, spec, init)
}

func TestSpecOptionAfterOptionalOptionsEndInAChoice(t *testing.T) {
	var x *[]string
	var f, d *bool
	init := func(c *Cmd) {
		x = c.StringsArg("X", nil, "")
		f = c.BoolOpt("f", false, "")
		d = c.BoolOpt("d", false, "")
	}

	spec := "(-f | [-- X]) -d"
	badSpec(t, spec, init)
}

func TestSpecOptionAfterOptionsEndIsParsedAsArg(t *testing.T) {
	var cmd *string
	var args *[]string
	init := func(c *Cmd) {
		cmd = c.StringArg("CMD", "", "")
		args = c.StringsArg("ARG", nil, "")
	}

	spec := "-- CMD [ARG...]"
	okCmd(t, spec, init, []string{"--", "go", "test", "-v"})
}

func TestSpecSingleDash(t *testing.T) {
	var path *string
	var f *bool

	init := func(c *Cmd) {
		path = c.StringArg("PATH", "", "'-' can be used to read from stdin' ")
		f = c.BoolOpt("f", false, "")
	}

	spec := "[-f] PATH"

	okCmd(t, spec, init, []string{"TEST"})
	require.Equal(t, "TEST", *path)
	require.False(t, *f)

	okCmd(t, spec, init, []string{"-f", "TEST"})
	require.Equal(t, "TEST", *path)
	require.True(t, *f)

	okCmd(t, spec, init, []string{"-"})
	require.Equal(t, "-", *path)
	require.False(t, *f)

	okCmd(t, spec, init, []string{"-f", "-"})
	require.Equal(t, "-", *path)
	require.True(t, *f)

	okCmd(t, spec, init, []string{"--", "-"})
	require.Equal(t, "-", *path)
	require.False(t, *f)

	okCmd(t, spec, init, []string{"-f", "--", "-"})
	require.Equal(t, "-", *path)
	require.True(t, *f)

}
