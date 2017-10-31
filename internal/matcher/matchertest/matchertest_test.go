package matchertest

import (
	"fmt"
	"testing"
)

func TestMe(t *testing.T) {
	m := NewOptions("-abc")
	fmt.Printf("%s\n", m)
}
