package fsmdot

import (
	"fmt"
	"strings"

	"github.com/jawher/mow.cli/internal/fsm"
)

// Dot generates a graphviz dot representation of an FSM
func Dot(s *fsm.State) string {
	trs := dot(s, mkStateNames(), map[*fsm.State]struct{}{})
	return fmt.Sprintf("digraph G {\n\trankdir=LR\n%s\n}\n", strings.Join(trs, "\n"))
}

func dot(s *fsm.State, sn *stateNames, visited map[*fsm.State]struct{}) []string {
	var res []string
	if _, ok := visited[s]; ok {
		return res
	}
	id := sn.id(s)
	visited[s] = struct{}{}

	attrs := ""
	if s.Terminal {
		attrs = " [peripheries=2]"
	}
	res = append(res, fmt.Sprintf("\tS%d%s", id, attrs))

	for _, tr := range s.Transitions {
		res = append(res, fmt.Sprintf("\tS%d -> S%d [label=\"%v\"]", id, sn.id(tr.Next), tr.Matcher))
		res = append(res, dot(tr.Next, sn, visited)...)
	}

	return res
}

func mkStateNames() *stateNames {
	return &stateNames{
		counter: 1,
		ids:     map[*fsm.State]int{},
	}
}

type stateNames struct {
	counter int
	ids     map[*fsm.State]int
}

func (sn *stateNames) id(s *fsm.State) int {
	res := sn.ids[s]
	if res != 0 {
		return res
	}
	res = sn.counter
	sn.ids[s] = res
	sn.counter++
	return res
}
