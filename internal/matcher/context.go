package matcher

import "github.com/jawher/mow.cli/internal/container"

type ParseContext struct {
	Args          map[*container.Container][]string
	Opts          map[*container.Container][]string
	ExcludedOpts  map[*container.Container]struct{}
	RejectOptions bool
}

func New() ParseContext {
	return ParseContext{
		Args:          map[*container.Container][]string{},
		Opts:          map[*container.Container][]string{},
		ExcludedOpts:  map[*container.Container]struct{}{},
		RejectOptions: false,
	}
}

func (pc ParseContext) Merge(o ParseContext) {
	for k, vs := range o.Args {
		pc.Args[k] = append(pc.Args[k], vs...)
	}

	for k, vs := range o.Opts {
		pc.Opts[k] = append(pc.Opts[k], vs...)
	}
}
