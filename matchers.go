package cli

import (
	"fmt"
	"strings"
)

type upMatcher interface {
	match(args []string, c *parseContext) (bool, []string)
}

type upShortcut bool

func (u upShortcut) match(args []string, c *parseContext) (bool, []string) {
	return true, args
}

func (u upShortcut) String() string {
	return "*"
}

type upOptsEnd bool

func (u upOptsEnd) match(args []string, c *parseContext) (bool, []string) {
	c.rejectOptions = true
	return true, args
}

func (u upOptsEnd) String() string {
	return "--"
}

const (
	shortcut = upShortcut(true)
	optsEnd  = upOptsEnd(true)
)

func (arg *arg) match(args []string, c *parseContext) (bool, []string) {
	if len(args) == 0 {
		return false, args
	}
	if !c.rejectOptions && strings.HasPrefix(args[0], "-") && args[0] != "-" {
		return false, args
	}
	c.args[arg] = append(c.args[arg], args[0])
	return true, args[1:]
}

func (o *opt) match(args []string, c *parseContext) (bool, []string) {
	if len(args) == 0 || c.rejectOptions {
		return false, args
	}
	arg := args[0]
	switch {
	case strings.HasPrefix(arg, "--"):
		return o.matchLongOpt(args, c)
	case strings.HasPrefix(arg, "-"):
		return o.matchShortOpt(args, c)
	default:
		return false, args
	}
}

func (o *opt) matchLongOpt(args []string, c *parseContext) (bool, []string) {
	kv := strings.Split(args[0], "=")
	name := kv[0]

	for _, oname := range o.names {
		if name == oname {
			if len(kv) == 2 {
				c.opts[o] = append(c.opts[o], kv[1])
				return true, args[1:]
			}
			if o.isBool() {
				c.opts[o] = append(c.opts[o], "true")
				return true, args[1:]
			}
			if len(args) < 2 {
				return false, args
			}
			val := args[1]
			if strings.HasPrefix(val, "-") {
				return false, args
			}
			c.opts[o] = append(c.opts[o], val)
			return true, args[2:]
		}
	}
	return false, args
}

func (o *opt) matchShortOpt(args []string, c *parseContext) (bool, []string) {
	arg := args[0]
	if len(arg) < 2 {
		return false, args
	}
	name := arg[0:2]
	for _, oname := range o.names {
		if name == oname {
			switch {
			case o.isBool():

				if len(arg) == 2 {
					c.opts[o] = append(c.opts[o], "true")

					return true, args[1:]
				}
				rem := arg[2:]
				if strings.HasPrefix(rem, "=") {
					c.opts[o] = append(c.opts[o], rem[1:])

					return true, args[1:]
				}

				c.opts[o] = append(c.opts[o], "true")
				nargs := make([]string, len(args))
				nargs[0] = "-" + rem
				for i := 1; i < len(args); i++ {
					nargs[i] = args[i]
				}
				return true, nargs
			default:
				if len(arg) > 2 {
					val := arg[2:]
					if strings.HasPrefix(val, "=") {
						val = arg[3:]
					}
					c.opts[o] = append(c.opts[o], val)
					return true, args[1:]
				}
				if len(args) < 2 {
					return false, args
				}
				val := args[1]
				if strings.HasPrefix(val, "-") {
					return false, args
				}
				c.opts[o] = append(c.opts[o], val)
				return true, args[2:]
			}
		}
	}
	return false, args
}

type optsMatcher []*opt

func (om optsMatcher) try(args []string, c *parseContext) (bool, []string) {
	if len(args) == 0 || c.rejectOptions {
		return false, args
	}
	for _, o := range om {
		if ok, nargs := o.match(args, c); ok {
			return ok, nargs
		}
	}
	return false, args
}

func (om optsMatcher) match(args []string, c *parseContext) (bool, []string) {
	ok, nargs := om.try(args, c)
	if !ok {
		return false, args
	}

	for {
		ok, nnargs := om.try(nargs, c)
		if !ok {
			return true, nargs
		}
		nargs = nnargs
	}
}

func (om optsMatcher) String() string {
	return fmt.Sprintf("Opts(%v)", []*opt(om))

}
