package flow

type ExitCode int

type Step struct {
	Do      func()
	Success *Step
	Error   *Step
	Desc    string
	Exiter  func(code int)
}

func (s *Step) Run(p interface{}) {
	s.callDo(p)

	switch {
	case s.Success != nil:
		s.Success.Run(p)
	case p == nil:
		return
	default:
		if code, ok := p.(ExitCode); ok {
			if s.Exiter != nil {
				s.Exiter(int(code))
			}
			return
		}
		panic(p)
	}
}

func (s *Step) callDo(p interface{}) {
	if s.Do == nil {
		return
	}
	defer func() {
		if e := recover(); e != nil {
			if s.Error == nil {
				panic(p)
			}
			s.Error.Run(e)
		}
	}()
	s.Do()
}
