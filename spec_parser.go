package cli

import "fmt"

func uParse(c *Cmd) (*state, error) {
	tokens, err := uTokenize(c.Spec)
	if err != nil {
		return nil, err
	}

	p := &uParser{cmd: c, tokens: tokens}
	return p.parse()
}

type uParser struct {
	cmd    *Cmd
	tokens []*uToken

	tkpos int

	matchedToken *uToken

	rejectOptions bool
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

	type pair struct {
		s, e *state
	}
	comps := []pair{}

	if required {
		s, e := p.choice()
		comps = append(comps, pair{s, e})
	}
	for p.canAtom() {
		s, e := p.choice()
		comps = append(comps, pair{s, e})
	}
	if len(comps) == 0 {
		return start, end
	}

	appendComp := func(p pair) {
		for _, tr := range p.s.transitions {
			end.t(tr.matcher, tr.next)
		}
		end = p.e
	}

	copyPair := func(p pair) pair {
		s, e := fsmCopy(p.s, p.e)
		return pair{s, e}
	}

	var p0 *pair
	for _, p1 := range comps {
		if !p1.s.onlyOpts() {
			if p0 != nil {
				appendComp(*p0)
				p0 = nil
			}
			appendComp(p1)
			continue
		}
		if p0 == nil {
			x := p1
			p0 = &x
			continue
		}
		//p is onlyOpts, and there is a previous p0 which is also onlyOpts
		// generate order-less combination
		np := &pair{
			s: newState(p.cmd),
			e: newState(p.cmd),
		}

		p0copy := copyPair(*p0)
		p1copy := copyPair(p1)

		np.s.t(shortcut, p0.s)
		p0.e.t(shortcut, p1copy.s)

		np.s.t(shortcut, p1.s)
		p1.e.t(shortcut, p0copy.s)

		p1copy.e.t(shortcut, np.e)
		p0copy.e.t(shortcut, np.e)

		p0 = np
	}
	if p0 != nil {
		appendComp(*p0)
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
			p.back()
			panic(fmt.Sprintf("Undeclared arg %s", name))
		}
		end = start.t(arg, newState(p.cmd))
	case p.found(utOptions):
		if p.rejectOptions {
			p.back()
			panic("No options after --")
		}
		end = newState(p.cmd)
		start.t(optsMatcher(p.cmd.options), end)
	case p.found(utShortOpt):
		if p.rejectOptions {
			p.back()
			panic("No options after --")
		}
		name := p.matchedToken.val
		opt, declared := p.cmd.optionsIdx[name]
		if !declared {
			p.back()
			panic(fmt.Sprintf("Undeclared option %s", name))
		}
		end = start.t(opt, newState(p.cmd))
		p.found(utOptValue)
	case p.found(utLongOpt):
		if p.rejectOptions {
			p.back()
			panic("No options after --")
		}
		name := p.matchedToken.val
		opt, declared := p.cmd.optionsIdx[name]
		if !declared {
			p.back()
			panic(fmt.Sprintf("Undeclared option %s", name))
		}
		end = start.t(opt, newState(p.cmd))
		p.found(utOptValue)
	case p.found(utOptSeq):
		if p.rejectOptions {
			p.back()
			panic("No options after --")
		}
		end = newState(p.cmd)
		sq := p.matchedToken.val
		opts := []*opt{}
		for i := range sq {
			sn := sq[i : i+1]
			opt, declared := p.cmd.optionsIdx["-"+sn]
			if !declared {
				p.back()
				panic(fmt.Sprintf("Undeclared option %s", sn))
			}
			opts = append(opts, opt)
		}
		start.t(optsMatcher(opts), end)
	case p.found(utOpenPar):
		start, end = p.seq(true)
		p.expect(utClosePar)
	case p.found(utOpenSq):
		start, end = p.seq(true)
		start.t(shortcut, end)
		p.expect(utCloseSq)
	case p.found(utDoubleDash):
		p.rejectOptions = true
		end = start.t(optsEnd, newState(p.cmd))
		return start, end
	default:
		panic("Unexpected input: was expecting a command or a positional argument or an option")
	}
	if p.found(utRep) {
		end.t(shortcut, start)
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
	case p.is(utOptSeq):
		return true
	case p.is(utOpenPar):
		return true
	case p.is(utOpenSq):
		return true
	case p.is(utDoubleDash):
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

func (p *uParser) back() {
	p.tkpos--
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
