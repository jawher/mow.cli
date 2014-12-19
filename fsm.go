package cli

import (
	"sort"
	"strings"

	"fmt"
)

type state struct {
	terminal    bool
	transitions transitions
	cmd         *Cmd
}

type transition struct {
	matcher upMatcher
	next    *state
}

type transitions []*transition

func (t transitions) Len() int      { return len(t) }
func (t transitions) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t transitions) Less(i, j int) bool {
	a, _ := t[i].matcher, t[j].matcher
	switch a.(type) {
	case upShortcut:
		return false
	case *arg:
		return false
	default:
		return true
	}

}

func newState(cmd *Cmd) *state {
	return &state{false, []*transition{}, cmd}
}

func (s *state) t(matcher upMatcher, next *state) *state {
	s.transitions = append(s.transitions, &transition{matcher, next})
	return next
}

func (s *state) replace(something, with *state) {
	in := incoming(s, something, map[*state]bool{})
	for _, tr := range in {
		tr.next = with
	}
	for _, tr := range something.transitions {
		with.t(tr.matcher, tr.next)
	}
}

func incoming(s, into *state, visited map[*state]bool) []*transition {
	res := []*transition{}
	if visited[s] {
		return res
	}
	visited[s] = true

	for _, tr := range s.transitions {
		if tr.next == into {
			res = append(res, tr)
		}
		res = append(res, incoming(tr.next, into, visited)...)
	}
	return res
}

func (s *state) dot() string {
	i := new(int)
	*i = 0
	trs := dot(s, i, map[*state]bool{}, map[*state]int{})
	return fmt.Sprintf("digraph G {\n\trankdir=LR\n%s\n}\n", strings.Join(trs, "\n"))
}

func id(s *state, counter *int, ids map[*state]int) int {
	id, ok := ids[s]
	if ok {
		return id
	}
	ids[s] = *counter
	*counter++
	return ids[s]
}
func dot(s *state, counter *int, visited map[*state]bool, ids map[*state]int) []string {
	res := []string{}
	if visited[s] {
		return res
	}
	visited[s] = true

	i := id(s, counter, ids)
	for _, tr := range s.transitions {
		res = append(res, fmt.Sprintf("\tS%d -> S%d [label=\"%v\"]", i, id(tr.next, counter, ids), tr.matcher))
		res = append(res, dot(tr.next, counter, visited, ids)...)
	}
	if s.terminal {
		res = append(res, fmt.Sprintf("S%d [peripheries=2]", i))
	}
	return res
}

type parseContext struct {
	args map[*arg][]string
	opts map[*opt][]string
}

func newParseContext() parseContext {
	return parseContext{map[*arg][]string{}, map[*opt][]string{}}
}

func (pc parseContext) merge(o parseContext) {
	for k, vs := range o.args {
		pc.args[k] = append(pc.args[k], vs...)
	}

	for k, vs := range o.opts {
		pc.opts[k] = append(pc.opts[k], vs...)
	}
}

func (s *state) parse(args []string) error {
	pc := newParseContext()
	if !s.apply(args, pc) {
		return fmt.Errorf("incorrect usage")
	}

	for opt, vs := range pc.opts {
		for _, v := range vs {
			if err := opt.set(v); err != nil {
				return err
			}
		}
	}

	for arg, vs := range pc.args {
		for _, v := range vs {
			if err := arg.set(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *state) apply(args []string, pc parseContext) bool {
	sort.Sort(s.transitions)

	type match struct {
		tr       *transition
		consumes int
		pc       parseContext
	}

	matches := []*match{}
	for _, tr := range s.transitions {
		fresh := newParseContext()
		if ok, cons := tr.matcher.match(args, fresh); ok {
			matches = append(matches, &match{tr, cons, fresh})
		}
	}

	for _, m := range matches {
		ok := m.tr.next.apply(args[m.consumes:], m.pc)
		if ok {
			pc.merge(m.pc)
			return true
		}
	}
	return s.terminal && len(args) == 0
}
