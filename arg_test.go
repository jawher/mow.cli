package cli

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringArg(t *testing.T) {
	cmd := &Cmd{argsIdx: map[string]*arg{}}
	a := cmd.StringArg("a", "test", "", nil)
	require.Equal(t, "test", *a)

	os.Setenv("B", "")
	b := cmd.StringArg("b", "test", "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, "test", *b)

	os.Setenv("B", "mow")
	b = cmd.StringArg("b", "test", "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, "mow", *b)

	os.Setenv("B", "")
	os.Setenv("C", "cli")
	os.Setenv("D", "mow")
	b = cmd.StringArg("b", "test", "", &ArgExtra{EnvVar: "B C D"})
	require.Equal(t, "cli", *b)
}

func TestBoolArg(t *testing.T) {
	cmd := &Cmd{argsIdx: map[string]*arg{}}
	a := cmd.BoolArg("a", true, "", nil)
	require.True(t, *a)

	os.Setenv("B", "")
	b := cmd.BoolArg("b", false, "", &ArgExtra{EnvVar: "B"})
	require.False(t, *b)

	trueValues := []string{"1", "true", "TRUE"}
	for _, tv := range trueValues {
		os.Setenv("B", tv)
		b = cmd.BoolArg("b", false, "", &ArgExtra{EnvVar: "B"})
		require.True(t, *b, "env=%s", tv)
	}

	falseValues := []string{"0", "false", "FALSE", "xyz"}
	for _, tv := range falseValues {
		os.Setenv("B", tv)
		b = cmd.BoolArg("b", false, "", &ArgExtra{EnvVar: "B"})
		require.False(t, *b, "env=%s", tv)
	}

	os.Setenv("B", "")
	os.Setenv("C", "false")
	os.Setenv("D", "true")
	b = cmd.BoolArg("b", true, "", &ArgExtra{EnvVar: "B C D"})
	require.False(t, *b)
}

func TestIntArg(t *testing.T) {
	cmd := &Cmd{argsIdx: map[string]*arg{}}
	a := cmd.IntArg("a", -1, "", nil)
	require.Equal(t, -1, *a)

	os.Setenv("B", "")
	b := cmd.IntArg("b", -1, "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, -1, *b)

	goodValues := []int{1, 0, 33}
	for _, tv := range goodValues {
		os.Setenv("B", strconv.Itoa(tv))
		b := cmd.IntArg("b", -1, "", &ArgExtra{EnvVar: "B"})
		require.Equal(t, tv, *b, "env=%s", tv)
	}

	badValues := []string{"", "b", "q1", "_"}
	for _, tv := range badValues {
		os.Setenv("B", tv)
		b := cmd.IntArg("b", -1, "", &ArgExtra{EnvVar: "B"})
		require.Equal(t, -1, *b, "env=%s", tv)
	}

	os.Setenv("B", "")
	os.Setenv("C", "42")
	os.Setenv("D", "666")
	b = cmd.IntArg("b", -1, "", &ArgExtra{EnvVar: "B C D"})
	require.Equal(t, 42, *b)
}

func TestStringsArg(t *testing.T) {
	cmd := &Cmd{argsIdx: map[string]*arg{}}
	v := []string{"test"}
	a := cmd.StringsArg("a", v, "", nil)
	require.Equal(t, v, *a)

	os.Setenv("B", "")
	b := cmd.StringsArg("b", v, "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, v, *b)

	os.Setenv("B", "mow")
	b = cmd.StringsArg("b", nil, "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, []string{"mow"}, *b)

	os.Setenv("B", "mow, cli")
	b = cmd.StringsArg("b", nil, "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, []string{"mow", "cli"}, *b)

	os.Setenv("B", "")
	os.Setenv("C", "test")
	os.Setenv("D", "xxx")
	b = cmd.StringsArg("b", nil, "", &ArgExtra{EnvVar: "B C D"})
	require.Equal(t, v, *b)
}

func TestIntsArg(t *testing.T) {
	cmd := &Cmd{argsIdx: map[string]*arg{}}
	vi := []int{42}
	a := cmd.IntsArg("a", vi, "", nil)
	require.Equal(t, vi, *a)

	os.Setenv("B", "")
	b := cmd.IntsArg("b", vi, "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, vi, *b)

	os.Setenv("B", "666")
	b = cmd.IntsArg("b", nil, "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, []int{666}, *b)

	os.Setenv("B", "1, 2 , 3")
	b = cmd.IntsArg("b", nil, "", &ArgExtra{EnvVar: "B"})
	require.Equal(t, []int{1, 2, 3}, *b)

	os.Setenv("B", "")
	os.Setenv("C", "abc")
	os.Setenv("D", "1, abc")
	os.Setenv("E", "42")
	os.Setenv("F", "666")
	b = cmd.IntsArg("b", nil, "", &ArgExtra{EnvVar: "B C D E F"})
	require.Equal(t, vi, *b)
}
