package cli

import (
	"sort"
	"strings"

	"fmt"
)

type state struct {
	id          int
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

var _id = 0

func newState(cmd *Cmd) *state {
	_id++
	return &state{_id, false, []*transition{}, cmd}
}

func (s *state) t(matcher upMatcher, next *state) *state {
	s.transitions = append(s.transitions, &transition{matcher, next})
	return next
}

func (s *state) has(tr *transition) bool {
	for _, t := range s.transitions {
		if t.next == tr.next && t.matcher == tr.matcher {
			return true
		}
	}
	return false
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

func removeItemAt(idx int, arr transitions) transitions {
	res := make([]*transition, len(arr)-1)
	copy(res, arr[:idx])
	copy(res[idx:], arr[idx+1:])
	return res
}

func fsmCopy(s, e *state) (*state, *state) {
	cache := map[*state]*state{}
	s.copy(cache)
	return cache[s], cache[e]
}

func (s *state) copy(cache map[*state]*state) *state {
	c := newState(s.cmd)
	c.terminal = s.terminal
	cache[s] = c
	for _, tr := range s.transitions {
		dest := tr.next
		if destCopy, found := cache[dest]; found {
			c.transitions = append(c.transitions, &transition{tr.matcher, destCopy})
			continue
		}
		c.transitions = append(c.transitions, &transition{tr.matcher, dest.copy(cache)})
	}
	return c
}

func (s *state) simplify() {
	simplify(s, s, map[*state]bool{})
}

func simplify(start, s *state, visited map[*state]bool) {
	if visited[s] {
		return
	}
	visited[s] = true
	for _, tr := range s.transitions {
		simplify(start, tr.next, visited)
	}
	for s.simplifySelf(start) {
	}
}

func (s *state) simplifySelf(start *state) bool {
	for idx, tr := range s.transitions {
		if _, ok := tr.matcher.(upShortcut); ok {
			next := tr.next
			s.transitions = removeItemAt(idx, s.transitions)
			for _, tr := range next.transitions {
				if !s.has(tr) {
					s.transitions = append(s.transitions, tr)
				}
			}
			if next.terminal {
				s.terminal = true
			}
			return true
		}
	}
	return false
}

func (s *state) onlyOpts() bool {
	return onlyOpts(s, map[*state]bool{})
}

func onlyOpts(s *state, visited map[*state]bool) bool {
	if visited[s] {
		return true
	}
	visited[s] = true

	for _, tr := range s.transitions {
		switch tr.matcher.(type) {
		case *arg:
			return false
		}
		if !onlyOpts(tr.next, visited) {
			return false
		}
	}

	return true
}

func (s *state) dot() string {
	trs := dot(s, map[*state]bool{})
	return fmt.Sprintf("digraph G {\n\trankdir=LR\n%s\n}\n", strings.Join(trs, "\n"))
}

func dot(s *state, visited map[*state]bool) []string {
	res := []string{}
	if visited[s] {
		return res
	}
	visited[s] = true

	for _, tr := range s.transitions {
		res = append(res, fmt.Sprintf("\tS%d -> S%d [label=\"%v\"]", s.id, tr.next.id, tr.matcher))
		res = append(res, dot(tr.next, visited)...)
	}
	if s.terminal {
		res = append(res, fmt.Sprintf("\tS%d [peripheries=2]", s.id))
	}
	return res
}

type parseContext struct {
	args          map[*arg][]string
	opts          map[*opt][]string
	rejectOptions bool
}

func newParseContext() parseContext {
	return parseContext{map[*arg][]string{}, map[*opt][]string{}, false}
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
	if s.terminal && len(args) == 0 {
		return true
	}
	sort.Sort(s.transitions)

	if len(args) > 0 && args[0] == "--" && !pc.rejectOptions {
		pc.rejectOptions = true
		args = args[1:]
	}

	type match struct {
		tr       *transition
		consumes int
		pc       parseContext
	}

	matches := []*match{}
	for _, tr := range s.transitions {
		fresh := newParseContext()
		fresh.rejectOptions = pc.rejectOptions
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
	return false
}
