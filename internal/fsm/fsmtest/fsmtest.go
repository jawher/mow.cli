package fsmtest

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jawher/mow.cli/internal/fsm"
	"github.com/jawher/mow.cli/internal/matcher"
)

type NopeMatcher struct{}

func (NopeMatcher) Match(args []string, c *matcher.ParseContext) (bool, []string) {
	return false, args
}

func (NopeMatcher) Priority() int {
	return 666
}

type YepMatcher struct{}

func (YepMatcher) Match(args []string, c *matcher.ParseContext) (bool, []string) {
	return true, args
}

func (YepMatcher) Priority() int {
	return 666
}

type TestMatcher struct {
	MatchFunc    func(args []string, c *matcher.ParseContext) (bool, []string)
	TestPriority int
}

func (t TestMatcher) Match(args []string, c *matcher.ParseContext) (bool, []string) {
	return t.MatchFunc(args, c)
}

func (t TestMatcher) Priority() int {
	return t.TestPriority
}

func NewFsm(spec string, matchers map[string]matcher.Matcher) *fsm.State {
	states := map[string]*fsm.State{}
	lines := strings.FieldsFunc(spec, func(r rune) bool { return r == '\n' })

	var res *fsm.State = nil

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			panic(fmt.Sprintf("Invalid line %q: syntax: START TR END", line))
		}
		sn, tn, en := parts[0], parts[1], parts[2]
		sn, sterm := stateNameTerm(sn)
		en, eterm := stateNameTerm(en)

		s, ok := states[sn]
		if !ok {
			s = fsm.NewState()
			states[sn] = s
		}
		s.Terminal = s.Terminal || sterm

		if res == nil {
			res = s
		}

		e, ok := states[en]
		if !ok {
			e = fsm.NewState()
			states[en] = e
		}
		e.Terminal = e.Terminal || eterm

		t, ok := matchers[tn]
		if !ok {
			panic(fmt.Sprintf("Unknown matcher %q in line %q", tn, line))
		}

		s.T(t, e)
	}
	return res
}

func TransitionStrs(trs fsm.StateTransitions) []string {
	var res []string
	for _, tr := range trs {
		fmt.Printf("%#v\n", tr.Matcher)
		res = append(res, fmt.Sprintf("%v", tr.Matcher))
	}
	return res
}

func stateNameTerm(name string) (string, bool) {
	if strings.HasPrefix(name, "(") {
		if strings.HasSuffix(name, ")") {
			name = name[1 : len(name)-1]
			return name, true
		}
		panic(fmt.Sprintf("Invalid state name %q", name))
	}
	return name, false
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

func stateName(s *fsm.State, sn *stateNames) string {
	id := sn.id(s)

	if !s.Terminal {
		return fmt.Sprintf("S%d", id)
	}
	return fmt.Sprintf("(S%d)", id)
}

func FsmStr(s *fsm.State) string {
	lines := fsmStrVis(s, mkStateNames(), map[*fsm.State]struct{}{})
	sort.Sort(lines)
	return strings.Join(lines, "\n")
}

func fsmStrVis(s *fsm.State, sn *stateNames, visited map[*fsm.State]struct{}) fsmStrings {
	if _, ok := visited[s]; ok {
		return nil
	}
	visited[s] = struct{}{}

	res := fsmStrings{}
	for _, tr := range s.Transitions {
		res = append(res, fmt.Sprintf("%s %v %s", stateName(s, sn), tr.Matcher, stateName(tr.Next, sn)))
		res = append(res, fsmStrVis(tr.Next, sn, visited)...)
	}

	return res
}

type fsmStrings []string

func (t fsmStrings) Len() int      { return len(t) }
func (t fsmStrings) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t fsmStrings) Less(i, j int) bool {
	a := strings.TrimFunc(t[i], isParen)
	b := strings.TrimFunc(t[j], isParen)
	return strings.Compare(a, b) < 0
}

func isParen(r rune) bool {
	return r == '(' || r == ')'
}
