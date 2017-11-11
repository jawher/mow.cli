package flowdot

import (
	"fmt"
	"strings"

	"github.com/jawher/mow.cli/internal/flow"
)

/*
Dot generates a graphviz dot string representing the a flow
*/
func Dot(s *flow.Step) string {
	trs := flowDot(s, map[*flow.Step]bool{})
	return fmt.Sprintf("digraph G {\n\trankdir=LR\n%s\n}\n", strings.Join(trs, "\n"))
}

func flowDot(s *flow.Step, visited map[*flow.Step]bool) []string {
	var res []string
	if visited[s] {
		return res
	}
	visited[s] = true

	if s.Success != nil {
		res = append(res, fmt.Sprintf("\t\"%s\" -> \"%s\" [label=\"ok\"]", s.Desc, s.Success.Desc))
		res = append(res, flowDot(s.Success, visited)...)
	}
	if s.Error != nil {
		res = append(res, fmt.Sprintf("\t\"%s\" -> \"%s\" [label=\"ko\"]", s.Desc, s.Error.Desc))
		res = append(res, flowDot(s.Error, visited)...)
	}
	return res
}
