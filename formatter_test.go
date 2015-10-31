package cli

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatters(t *testing.T) {
	cases := []struct {
		input    interface{}
		expected string
	}{
		{true, "true"},
		{false, "false"},

		{"", `""`},
		{"val", `"val"`},

		{42, "42"},

		{[]string{}, `[]`},
		{[]string{"a"}, `["a"]`},
		{[]string{"a", "b"}, `["a", "b"]`},

		{[]int{}, "[]"},
		{[]int{1}, "[1]"},
		{[]int{1, 2}, "[1, 2]"},
	}

	for _, cas := range cases {
		f := formatterFor(reflect.TypeOf(cas.input))
		require.Equal(t, cas.expected, f(cas.input), "formatting error for value %v (%T)", cas.input, cas.input)
	}
}
