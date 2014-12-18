package cli

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVConv(t *testing.T) {
	goodCases := []struct {
		s     string
		to    reflect.Type
		value interface{}
	}{
		{"test", reflect.TypeOf(""), "test"},
		{"true", reflect.TypeOf(true), true},
		{"false", reflect.TypeOf(true), false},
		{"42", reflect.TypeOf(42), 42},
		{"a,b", reflect.TypeOf([]string{}), []string{"a", "b"}},
		{"7,42", reflect.TypeOf([]int{}), []int{7, 42}},
	}

	for _, cas := range goodCases {
		v, err := vconv(cas.s, cas.to)
		require.Nil(t, err)
		require.Equal(t, cas.value, v.Interface())
	}

	badCases := []struct {
		s  string
		to reflect.Type
	}{
		{"42", reflect.TypeOf(true)},
		{"xx", reflect.TypeOf(true)},
		{"aa", reflect.TypeOf(42)},
		{"a,b", reflect.TypeOf([]int{})},
		{"7,42c", reflect.TypeOf([]int{})},
	}

	for _, cas := range badCases {
		_, err := vconv(cas.s, cas.to)
		require.NotNil(t, err)
	}
}
