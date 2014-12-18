package cli

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStringOpt(t *testing.T) {
	cmd := &Cmd{optionsIdx: map[string]*option{}}
	a := cmd.StringOpt("a", "test", "", nil)
	require.Equal(t, "test", *a)

	os.Setenv("B", "")
	b := cmd.StringOpt("b", "test", "", &OptExtra{EnvVar: "B"})
	require.Equal(t, "test", *b)

	os.Setenv("B", "mow")
	b = cmd.StringOpt("b", "test", "", &OptExtra{EnvVar: "B"})
	require.Equal(t, "mow", *b)

	os.Setenv("B", "")
	os.Setenv("C", "cli")
	os.Setenv("D", "mow")
	b = cmd.StringOpt("b", "test", "", &OptExtra{EnvVar: "B C D"})
	require.Equal(t, "cli", *b)
}

func TestBoolOpt(t *testing.T) {
	cmd := &Cmd{optionsIdx: map[string]*option{}}
	a := cmd.BoolOpt("a", true, "", nil)
	require.True(t, *a)

	os.Setenv("B", "")
	b := cmd.BoolOpt("b", false, "", &OptExtra{EnvVar: "B"})
	require.False(t, *b)

	trueValues := []string{"1", "true", "TRUE"}
	for _, tv := range trueValues {
		os.Setenv("B", tv)
		b = cmd.BoolOpt("b", false, "", &OptExtra{EnvVar: "B"})
		require.True(t, *b, "env=%s", tv)
	}

	falseValues := []string{"0", "false", "FALSE", "xyz"}
	for _, tv := range falseValues {
		os.Setenv("B", tv)
		b = cmd.BoolOpt("b", false, "", &OptExtra{EnvVar: "B"})
		require.False(t, *b, "env=%s", tv)
	}

	os.Setenv("B", "")
	os.Setenv("C", "false")
	os.Setenv("D", "true")
	b = cmd.BoolOpt("b", true, "", &OptExtra{EnvVar: "B C D"})
	require.False(t, *b)
}

func TestIntOpt(t *testing.T) {
	cmd := &Cmd{optionsIdx: map[string]*option{}}
	a := cmd.IntOpt("a", -1, "", nil)
	require.Equal(t, -1, *a)

	os.Setenv("B", "")
	b := cmd.IntOpt("b", -1, "", &OptExtra{EnvVar: "B"})
	require.Equal(t, -1, *b)

	goodValues := []int{1, 0, 33}
	for _, tv := range goodValues {
		os.Setenv("B", strconv.Itoa(tv))
		b := cmd.IntOpt("b", -1, "", &OptExtra{EnvVar: "B"})
		require.Equal(t, tv, *b, "env=%s", tv)
	}

	badValues := []string{"", "b", "q1", "_"}
	for _, tv := range badValues {
		os.Setenv("B", tv)
		b := cmd.IntOpt("b", -1, "", &OptExtra{EnvVar: "B"})
		require.Equal(t, -1, *b, "env=%s", tv)
	}

	os.Setenv("B", "")
	os.Setenv("C", "42")
	os.Setenv("D", "666")
	b = cmd.IntOpt("b", -1, "", &OptExtra{EnvVar: "B C D"})
	require.Equal(t, 42, *b)
}

func TestStringsOpt(t *testing.T) {
	cmd := &Cmd{optionsIdx: map[string]*option{}}
	v := []string{"test"}
	a := cmd.StringsOpt("a", v, "", nil)
	require.Equal(t, v, *a)

	os.Setenv("B", "")
	b := cmd.StringsOpt("b", v, "", &OptExtra{EnvVar: "B"})
	require.Equal(t, v, *b)

	os.Setenv("B", "mow")
	b = cmd.StringsOpt("b", nil, "", &OptExtra{EnvVar: "B"})
	require.Equal(t, []string{"mow"}, *b)

	os.Setenv("B", "mow, cli")
	b = cmd.StringsOpt("b", nil, "", &OptExtra{EnvVar: "B"})
	require.Equal(t, []string{"mow", "cli"}, *b)

	os.Setenv("B", "")
	os.Setenv("C", "test")
	os.Setenv("D", "xxx")
	b = cmd.StringsOpt("b", nil, "", &OptExtra{EnvVar: "B C D"})
	require.Equal(t, v, *b)
}

func TestIntsOpt(t *testing.T) {
	cmd := &Cmd{optionsIdx: map[string]*option{}}
	vi := []int{42}
	a := cmd.IntsOpt("a", vi, "", nil)
	require.Equal(t, vi, *a)

	os.Setenv("B", "")
	b := cmd.IntsOpt("b", vi, "", &OptExtra{EnvVar: "B"})
	require.Equal(t, vi, *b)

	os.Setenv("B", "666")
	b = cmd.IntsOpt("b", nil, "", &OptExtra{EnvVar: "B"})
	require.Equal(t, []int{666}, *b)

	os.Setenv("B", "1, 2 , 3")
	b = cmd.IntsOpt("b", nil, "", &OptExtra{EnvVar: "B"})
	require.Equal(t, []int{1, 2, 3}, *b)

	os.Setenv("B", "")
	os.Setenv("C", "abc")
	os.Setenv("D", "1, abc")
	os.Setenv("E", "42")
	os.Setenv("F", "666")
	b = cmd.IntsOpt("b", nil, "", &OptExtra{EnvVar: "B C D E F"})
	require.Equal(t, vi, *b)
}
