package cli

import "fmt"

func uParse(c *Cmd) (*state, error) {
	tokens, err := uTokenize(c.Spec)
	if err != nil {
		return nil, err
	}

	p := &uParser{c, tokens, 0, nil}
	return p.parse()
}

type uParser struct {
	cmd    *Cmd
	tokens []*uToken

	tkpos int

	matchedToken *uToken
}

type upMatcher interface {
	match(args []string, c parseContext) (bool, int)
}

type upShortcut bool

func (u upShortcut) match(args []string, c parseContext) (bool, int) {
	return true, 0
}

func (u upShortcut) String() string {
	return "*"
}

const (
	shortcut = upShortcut(true)
)

type upExactly string

func (u upExactly) match(args []string, c parseContext) (bool, int) {
	if len(args) == 0 {
		return false, 0
	}
	if args[0] == string(u) {
		return true, 1
	}
	return false, 0
}

func (u upExactly) String() string {
	return "==" + string(u)
}

func (p *uParser) parse() (s *state, err error) {
	defer func() {
		if v := recover(); v != nil {
			pos := len(p.cmd.Spec)
			if !p.eof() {
				pos = p.token().pos
			}
			s = nil
			switch t, ok := v.(string); ok {
			case true:
				err = &parseError{p.cmd.Spec, t, pos}
			default:
				panic(v)
			}
		}
	}()
	err = nil
	var e *state
	s, e = p.seq(false)
	if !p.eof() {
		s = nil
		err = &parseError{p.cmd.Spec, "Unexpected input", p.token().pos}
		return
	}

	e.terminal = true
	s.simplify()

	return
}

func (p *uParser) seq(required bool) (*state, *state) {
	start := newState(p.cmd)
	end := start
	if required {
		start, end = p.choice()
	}
	for p.canAtom() {
		s, e := p.choice()
		for _, tr := range s.transitions {
			end.t(tr.matcher, tr.next)
		}
		end = e
	}

	return start, end
}

func (p *uParser) choice() (*state, *state) {
	start, end := newState(p.cmd), newState(p.cmd)

	add := func(s, e *state) {
		start.t(shortcut, s)
		e.t(shortcut, end)
	}

	add(p.atom())
	for p.found(utChoice) {
		add(p.atom())
	}
	return start, end
}

func (p *uParser) atom() (*state, *state) {
	start := newState(p.cmd)
	var end *state
	switch {
	case p.eof():
		panic("Unexpected end of input")
	case p.found(utPos):
		name := p.matchedToken.val
		arg, declared := p.cmd.argsIdx[name]
		if !declared {
			panic(fmt.Sprintf("Undeclared arg %s", name))
		}
		end = start.t(arg, newState(p.cmd))
	case p.found(utOptions):
		end = newState(p.cmd)
		for _, opt := range p.cmd.options {
			start.t(opt, end)
		}
		end.t(shortcut, start)
	case p.found(utShortOpt):
		name := p.matchedToken.val
		opt, declared := p.cmd.optionsIdx[name]
		if !declared {
			panic(fmt.Sprintf("Undeclared option %s", name))
		}
		end = start.t(opt, newState(p.cmd))

	case p.found(utLongOpt):
		name := p.matchedToken.val
		opt, declared := p.cmd.optionsIdx[name]
		if !declared {
			panic(fmt.Sprintf("Undeclared option %s", name))
		}
		end = start.t(opt, newState(p.cmd))
	case p.found(utOptSeq):
		end = newState(p.cmd)
		sq := p.matchedToken.val

		for i, _ := range sq {
			sn := sq[i : i+1]
			opt, declared := p.cmd.optionsIdx["-"+sn]
			if !declared {
				panic(fmt.Sprintf("Undeclared option %s", sn))
			}

			start.t(opt, end)
		}
	case p.found(utOpenPar):
		start, end = p.seq(true)
		p.expect(utClosePar)
	case p.found(utOpenSq):
		start, end = p.seq(true)
		start.t(shortcut, end)
		p.expect(utCloseSq)
	default:
		panic("Unexpected input: was expecting a command or a positional argument or an option")
	}
	if p.found(utRep) {
		start2, end2 := newState(p.cmd), newState(p.cmd)
		start2.t(shortcut, start)
		end.t(shortcut, end2)
		end2.t(shortcut, start2)
		start, end = start2, end2
	}
	return start, end
}

func (p *uParser) canAtom() bool {
	switch {
	case p.is(utPos):
		return true
	case p.is(utOptions):
		return true
	case p.is(utShortOpt):
		return true
	case p.is(utLongOpt):
		return true
	case p.is(utOpenPar):
		return true
	case p.is(utOpenSq):
		return true
	default:
		return false
	}
}

func (p *uParser) found(t uTokenType) bool {
	if p.is(t) {
		p.matchedToken = p.token()
		p.tkpos++
		return true
	}
	return false
}

func (p *uParser) is(t uTokenType) bool {
	if p.eof() {
		return false
	}
	return p.token().typ == t
}

func (p *uParser) expect(t uTokenType) {
	if !p.found(t) {
		panic(fmt.Sprintf("Was expecting %v", t))
	}
}

func (p *uParser) eof() bool {
	return p.tkpos >= len(p.tokens)
}

func (p *uParser) token() *uToken {
	if p.eof() {
		return nil
	}

	return p.tokens[p.tkpos]
}
