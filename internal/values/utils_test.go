package values

import (
	"reflect"
	"testing"

	"flag"

	"os"

	"github.com/stretchr/testify/require"
)

func TestIsBool(t *testing.T) {
	require.True(t, IsBool(NewBool(new(bool), false)))

	require.False(t, IsBool(NewString(new(string), "")))
	require.False(t, IsBool(NewInt(new(int), 0)))
	require.False(t, IsBool(NewStrings(new([]string), nil)))
	require.False(t, IsBool(NewInts(new([]int), nil)))
}

func TestSetFromEnv(t *testing.T) {
	cases := []struct {
		desc     string
		envVars  string
		setup    func() (flag.Value, interface{})
		expected bool
		val      interface{}
	}{
		{
			desc:    "No env var",
			envVars: "",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "DoDo")
				var s string
				return NewString(&s, "default"), &s
			},
			expected: false,
			val:      "default",
		},
		{
			desc:    "One set env var",
			envVars: "A",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "DoDo")
				var s string
				return NewString(&s, "default"), &s
			},
			expected: true,
			val:      "DoDo",
		},
		{
			desc:    "One unset env var",
			envVars: "A",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "")
				var s string
				return NewString(&s, "default"), &s
			},
			expected: false,
			val:      "default",
		},
		{
			desc:    "Two env var, both set",
			envVars: "A B",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "Aval")
				os.Setenv("B", "Bval")
				var s string
				return NewString(&s, "default"), &s
			},
			expected: true,
			val:      "Aval",
		},
		{
			desc:    "Two env var, first unset, second set",
			envVars: "A B",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "")
				os.Setenv("B", "Bval")
				var s string
				return NewString(&s, "default"), &s
			},
			expected: true,
			val:      "Bval",
		},
		{
			desc:    "One set env var, invalid value",
			envVars: "A",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "XXX")
				var i int
				return NewInt(&i, 12), &i
			},
			expected: false,
			val:      12,
		},
		{
			desc:    "Two env var, both set, first invalid",
			envVars: "A B",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "XXX")
				os.Setenv("B", "16")
				var i int
				return NewInt(&i, 0), &i
			},
			expected: true,
			val:      16,
		},
		{
			desc:    "Three env vars, all set, first and last invalid",
			envVars: "A B C",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "XXX")
				os.Setenv("B", "16")
				os.Setenv("C", "YYY")
				var i int
				return NewInt(&i, 0), &i
			},
			expected: true,
			val:      16,
		},
		{
			desc:    "Three env vars, all set, first invalid",
			envVars: "A B C",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "XXX")
				os.Setenv("B", "16")
				os.Setenv("C", "32")
				var i int
				return NewInt(&i, 0), &i
			},
			expected: true,
			val:      16,
		},
		{
			desc:    "Multi value, one unset env var",
			envVars: "A",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "")
				var i []int
				return NewInts(&i, []int{1, 2}), &i
			},
			expected: false,
			val:      []int{1, 2},
		},
		{
			desc:    "Multi value, one set env var",
			envVars: "A",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "7")
				var i []int
				return NewInts(&i, []int{1, 2}), &i
			},
			expected: true,
			val:      []int{7},
		},
		{
			desc:    "Multi value, one set env var (2)",
			envVars: "A",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "7, 8, 9")
				var i []int
				return NewInts(&i, []int{1, 2}), &i
			},
			expected: true,
			val:      []int{7, 8, 9},
		},
		{
			desc:    "Multi value, 2 env vars, first unset",
			envVars: "A B",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "")
				os.Setenv("B", "7, 8, 9")
				var i []int
				return NewInts(&i, []int{1, 2}), &i
			},
			expected: true,
			val:      []int{7, 8, 9},
		},
		{
			desc:    "Multi value, 2 env vars, first invalid",
			envVars: "A B",
			setup: func() (flag.Value, interface{}) {
				os.Setenv("A", "10, 11, b")
				os.Setenv("B", "7, 8, 9")
				var i []int
				return NewInts(&i, []int{1, 2}), &i
			},
			expected: true,
			val:      []int{7, 8, 9},
		},
	}

	for _, cas := range cases {
		t.Run(cas.desc, func(t *testing.T) {
			t.Logf("Case: %s", cas.desc)

			val, into := cas.setup()

			actual := SetFromEnv(val, cas.envVars)

			require.Equal(t, cas.expected, actual)

			typ := reflect.TypeOf(into)
			if typ.Kind() != reflect.Ptr {
				t.Fatalf("config func did not return a pointer")
			}
			actualValue := reflect.ValueOf(into).Elem().Interface()

			require.Equal(t, cas.val, actualValue)
		})
	}
}
