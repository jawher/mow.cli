package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWords(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"", []string{}},
		{" ", []string{}},
		{"  ", []string{}},
		{"a", []string{"a"}},
		{"  a  ", []string{"a"}},
		{"  a  b  ", []string{"a", "b"}},
	}

	for _, cas := range cases {
		t.Logf("Testing input %q", cas.input)
		actual := words(cas.input)
		require.Equal(t, cas.expected, actual)
	}
}
