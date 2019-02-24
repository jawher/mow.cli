package values

import (
	"testing"

	"flag"

	"fmt"

	"github.com/stretchr/testify/require"
)

func TestBoolParam(t *testing.T) {
	var into bool

	param := NewBool(&into, false)

	require.True(t, param.IsBoolFlag())

	cases := []struct {
		input  string
		err    bool
		result bool
		string string
	}{
		{"true", false, true, "true"},
		{"false", false, false, "false"},
		{"123", true, false, ""},
		{"", true, false, ""},
	}

	for _, cas := range cases {
		t.Logf("testing .Set() with %q", cas.input)

		err := param.Set(cas.input)

		if cas.err {
			require.Errorf(t, err, "value %q should have returned an error", cas.input)
			continue
		}

		require.Equal(t, cas.result, into)
		require.Equal(t, cas.string, param.String())
	}

	defCases := []struct {
		value     bool
		isDefault bool
	}{
		{value: false, isDefault: true},
		{value: true, isDefault: false},
	}

	for _, cas := range defCases {
		t.Run(fmt.Sprintf("%v", cas.value), func(t *testing.T) {

			t.Logf("testing .IsDefault() with %v", cas.value)

			v := NewBool(&into, cas.value)

			require.Equal(t, cas.isDefault, v.IsDefault())
		})
	}
}

func TestStringParam(t *testing.T) {
	var into string

	param := NewString(&into, "")

	cases := []struct {
		input  string
		string string
	}{
		{"a", `"a"`},
		{"", `""`},
	}

	for _, cas := range cases {
		t.Run(cas.input, func(t *testing.T) {

			t.Logf("testing with %q", cas.input)

			err := param.Set(cas.input)

			require.NoError(t, err)

			require.Equal(t, cas.input, into)
			require.Equal(t, cas.string, param.String())
		})
	}

	defCases := []struct {
		value     string
		isDefault bool
	}{
		{value: "", isDefault: true},
		{value: "a", isDefault: false},
	}

	for _, cas := range defCases {
		t.Run(cas.value, func(t *testing.T) {

			t.Logf("testing .IsDefault() with %v", cas.value)

			v := NewString(&into, cas.value)

			require.Equal(t, cas.isDefault, v.IsDefault())
		})
	}
}

func TestIntParam(t *testing.T) {
	var into int

	param := NewInt(&into, 0)

	cases := []struct {
		input  string
		err    bool
		result int
		string string
	}{
		{"12", false, 12, "12"},
		{"0", false, 0, "0"},
		{"01", false, 1, "1"},
		{"", true, 0, ""},
		{"abc", true, 0, ""},
	}

	for _, cas := range cases {
		t.Run(cas.input, func(t *testing.T) {
			t.Logf("testing with %q", cas.input)

			err := param.Set(cas.input)

			if cas.err {
				require.Errorf(t, err, "value %q should have returned an error", cas.input)
				return
			}

			require.Equal(t, cas.result, into)
			require.Equal(t, cas.string, param.String())
		})
	}

	var v flag.Value = NewInt(&into, 0)
	_, ok := v.(DefaultValued)

	require.False(t, ok, "*IntValue should not implement DefaultValued")
}

func TestFloat64Param(t *testing.T) {
	var into float64

	param := NewFloat64(&into, 0)

	cases := []struct {
		input  string
		err    bool
		result float64
		string string
	}{
		{"12", false, 12, "12"},
		{"0", false, 0, "0"},
		{"01", false, 1, "1"},
		{"3.14", false, 3.14, "3.14"},
		{"00.123456789", false, 0.123456789, "0.123456789"},
		{"", true, 0, ""},
		{"abc", true, 0, ""},
	}

	for _, cas := range cases {
		t.Run(cas.input, func(t *testing.T) {
			t.Logf("testing with %q", cas.input)

			err := param.Set(cas.input)

			if cas.err {
				require.Errorf(t, err, "value %q should have returned an error", cas.input)
				return
			}

			require.Equal(t, cas.result, into)
			require.Equal(t, cas.string, param.String())
		})
	}

	var v flag.Value = NewFloat64(&into, 0)
	_, ok := v.(DefaultValued)

	require.False(t, ok, "*IntValue should not implement DefaultValued")
}

func TestStringsParam(t *testing.T) {
	var into []string
	param := NewStrings(&into, nil)

	param.Set("a")
	param.Set("b")

	require.Equal(t, []string{"a", "b"}, into)
	require.Equal(t, `["a", "b"]`, param.String())

	param.Clear()

	require.Empty(t, into)
	require.Equal(t, `[]`, param.String())

	defCases := []struct {
		value     []string
		isDefault bool
	}{
		{value: nil, isDefault: true},
		{value: []string{}, isDefault: true},
		{value: []string{""}, isDefault: false},
		{value: []string{"a"}, isDefault: false},
		{value: []string{"a", "b"}, isDefault: false},
	}

	for _, cas := range defCases {
		t.Run(fmt.Sprintf("%#v", cas.value), func(t *testing.T) {
			t.Logf("testing .IsDefault() with %v", cas.value)

			v := NewStrings(&into, cas.value)

			require.Equal(t, cas.isDefault, v.IsDefault())
		})
	}
}

func TestIntsParam(t *testing.T) {
	var into []int
	param := NewInts(&into, nil)

	err := param.Set("1")
	require.NoError(t, err)

	err = param.Set("2")
	require.NoError(t, err)

	require.Equal(t, []int{1, 2}, into)

	require.Equal(t, `[1, 2]`, param.String())

	err = param.Set("c")
	require.Error(t, err)
	require.Equal(t, []int{1, 2}, into)

	param.Clear()

	require.Empty(t, into)
	require.Equal(t, `[]`, param.String())

	defCases := []struct {
		value     []int
		isDefault bool
	}{
		{value: nil, isDefault: true},
		{value: []int{}, isDefault: true},
		{value: []int{0}, isDefault: false},
		{value: []int{1}, isDefault: false},
		{value: []int{1, 2}, isDefault: false},
	}

	for _, cas := range defCases {
		t.Run(fmt.Sprintf("%#v", cas.value), func(t *testing.T) {
			t.Logf("testing .IsDefault() with %v", cas.value)

			v := NewInts(&into, cas.value)

			require.Equal(t, cas.isDefault, v.IsDefault())
		})
	}
}

func TestFloats64Param(t *testing.T) {
	var into []float64
	param := NewFloats64(&into, nil)

	err := param.Set("1.1")
	require.NoError(t, err)

	err = param.Set("2.2")
	require.NoError(t, err)

	require.Equal(t, []float64{1.1, 2.2}, into)

	require.Equal(t, `[1.1, 2.2]`, param.String())

	err = param.Set("c")
	require.Error(t, err)
	require.Equal(t, []float64{1.1, 2.2}, into)

	param.Clear()

	require.Empty(t, into)
	require.Equal(t, `[]`, param.String())

	defCases := []struct {
		value     []float64
		isDefault bool
	}{
		{value: nil, isDefault: true},
		{value: []float64{}, isDefault: true},
		{value: []float64{0}, isDefault: false},
		{value: []float64{1}, isDefault: false},
		{value: []float64{1, 2}, isDefault: false},
	}

	for _, cas := range defCases {
		t.Run(fmt.Sprintf("%#v", cas.value), func(t *testing.T) {
			t.Logf("testing .IsDefault() with %v", cas.value)

			v := NewFloats64(&into, cas.value)

			require.Equal(t, cas.isDefault, v.IsDefault())
		})
	}
}
