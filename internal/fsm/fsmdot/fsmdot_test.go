package fsmdot

import (
	"testing"

	"strings"

	"github.com/jawher/mow.cli/internal/fsm"
	"github.com/jawher/mow.cli/internal/fsm/fsmtest"
	"github.com/stretchr/testify/require"
)

func TestDot(t *testing.T) {
	s1 := fsm.NewState()
	s2 := fsm.NewState()
	s3 := fsm.NewState()
	s3.Terminal = true

	s1.T(fsmtest.YepMatcher{}, s2)
	s1.T(fsmtest.YepMatcher{}, s3)
	s2.T(fsmtest.NopeMatcher{}, s3)

	str := Dot(s1)

	str = strings.TrimSpace(str)

	expected := strings.TrimSpace(`
digraph G {
	rankdir=LR
	S1
	S1 -> S2 [label="<yep>"]
	S2
	S2 -> S3 [label="<nope>"]
	S3 [peripheries=2]
	S1 -> S3 [label="<yep>"]
}`)
	require.Equal(t, expected, str)
}
