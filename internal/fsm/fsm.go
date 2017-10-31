package fsm

import (
	"sort"

	"fmt"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/jawher/mow.cli/internal/matcher"
	"github.com/jawher/mow.cli/internal/values"
)

type State struct {
	Terminal    bool
	Transitions StateTransitions
}

type Transition struct {
	Matcher matcher.Matcher
	Next    *State
}

type StateTransitions []*Transition

func (t StateTransitions) Len() int      { return len(t) }
func (t StateTransitions) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t StateTransitions) Less(i, j int) bool {
	a, b := t[i].Matcher, t[j].Matcher
	return a.Priority() < b.Priority()
}

func NewState() *State {
	return &State{Transitions: []*Transition{}}
}

func (s *State) T(matcher matcher.Matcher, next *State) *State {
	s.Transitions = append(s.Transitions, &Transition{Matcher: matcher, Next: next})
	return next
}

func (s *State) Prepare() {
	simplify(s, s, map[*State]bool{})
	sortTransitions(s, map[*State]bool{})
}

func sortTransitions(s *State, visited map[*State]bool) {
	if visited[s] {
		return
	}
	visited[s] = true

	sort.Sort(s.Transitions)

	for _, tr := range s.Transitions {
		sortTransitions(tr.Next, visited)
	}
}

func simplify(start, s *State, visited map[*State]bool) {
	if visited[s] {
		return
	}
	visited[s] = true
	for _, tr := range s.Transitions {
		simplify(start, tr.Next, visited)
	}
	for s.simplifySelf(start) {
	}
}

func (s *State) simplifySelf(start *State) bool {
	for idx, tr := range s.Transitions {
		if matcher.IsShortcut(tr.Matcher) {
			next := tr.Next
			s.Transitions = removeTransitionAt(idx, s.Transitions)
			for _, tr := range next.Transitions {
				if !s.has(tr) {
					s.Transitions = append(s.Transitions, tr)
				}
			}
			if next.Terminal {
				s.Terminal = true
			}
			return true
		}
	}
	return false
}

func removeTransitionAt(idx int, arr StateTransitions) StateTransitions {
	res := make([]*Transition, len(arr)-1)
	copy(res, arr[:idx])
	copy(res[idx:], arr[idx+1:])
	return res
}

func (s *State) has(tr *Transition) bool {
	for _, t := range s.Transitions {
		if t.Next == tr.Next && t.Matcher == tr.Matcher {
			return true
		}
	}
	return false
}

func (s *State) Parse(args []string) error {
	pc := matcher.New()
	ok := s.apply(args, pc)
	if !ok {
		return fmt.Errorf("incorrect usage")
	}

	if err := fillContainers(pc.Opts); err != nil {
		return err
	}
	if err := fillContainers(pc.Args); err != nil {
		return err
	}

	return nil
}

func fillContainers(containers map[*container.Container][]string) error {
	for con, vs := range containers {
		if multiValued, ok := con.Value.(values.MultiValued); ok {
			multiValued.Clear()
		}
		for _, v := range vs {
			if err := con.Value.Set(v); err != nil {
				return err
			}
		}

		con.ValueSetFromEnv = false
		if con.ValueSetByUser != nil {
			*con.ValueSetByUser = true
		}
	}
	return nil
}

func (s *State) apply(args []string, pc matcher.ParseContext) bool {
	if s.Terminal && len(args) == 0 {
		return true
	}

	if len(args) > 0 {
		arg := args[0]

		if !pc.RejectOptions && arg == "--" {
			pc.RejectOptions = true
			args = args[1:]
		}
	}

	type match struct {
		tr  *Transition
		rem []string
		pc  matcher.ParseContext
	}

	var matches []*match
	for _, tr := range s.Transitions {
		fresh := matcher.New()
		fresh.RejectOptions = pc.RejectOptions
		if ok, rem := tr.Matcher.Match(args, &fresh); ok {
			matches = append(matches, &match{tr, rem, fresh})
		}
	}

	for _, m := range matches {
		if ok := m.tr.Next.apply(m.rem, m.pc); ok {
			pc.Merge(m.pc)
			return true
		}
	}

	return false
}
