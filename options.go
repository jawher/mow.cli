package cli

import (
	"fmt"
	"reflect"
	"strings"
)

type opt struct {
	name      string
	desc      string
	envVar    string
	names     []string
	value     reflect.Value
	hideValue bool
}

type BoolOpt struct {
	BoolParam

	// A space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
	// The one letter names will then be called with a single dash (short option), the others with two (long options).
	Name string
	// The option description as will be shown in help messages
	Desc string
	// A space separated list of environment variables names to be used to initialize this option
	EnvVar string
	// The option's inital value
	Value bool
	// A boolean to display or not the current value of the option in the help message
	HideValue bool
}

type StringOpt struct {
	StringParam

	// A space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
	// The one letter names will then be called with a single dash (short option), the others with two (long options).
	Name string
	// The option description as will be shown in help messages
	Desc string
	// A space separated list of environment variables names to be used to initialize this option
	EnvVar string
	// The option's inital value
	Value string
	// A boolean to display or not the current value of the option in the help message
	HideValue bool
}

type IntOpt struct {
	IntParam

	// A space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
	// The one letter names will then be called with a single dash (short option), the others with two (long options).
	Name string
	// The option description as will be shown in help messages
	Desc string
	// A space separated list of environment variables names to be used to initialize this option
	EnvVar string
	// The option's inital value
	Value int
	// A boolean to display or not the current value of the option in the help message
	HideValue bool
}

type StringsOpt struct {
	StringsParam

	// A space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
	// The one letter names will then be called with a single dash (short option), the others with two (long options).
	Name string
	// The option description as will be shown in help messages
	Desc string
	// A space separated list of environment variables names to be used to initialize this option.
	// The env variable should contain a comma separated list of values
	EnvVar string
	// The option's inital value
	Value []string
	// A boolean to display or not the current value of the option in the help message
	HideValue bool
}

type IntsOpt struct {
	IntsParam

	// A space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
	// The one letter names will then be called with a single dash (short option), the others with two (long options).
	Name string
	// The option description as will be shown in help messages
	Desc string
	// A space separated list of environment variables names to be used to initialize this option.
	// The env variable should contain a comma separated list of values
	EnvVar string
	// The option's inital value
	Value []int
	// A boolean to display or not the current value of the option in the help message
	HideValue bool
}

/*
BoolOpt defines a boolean option on the command c named `name`, with an initial value of `value` and a description of `desc` which will be used in help messages.

The name is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).


The result should be stored in a variable (a pointer to a bool) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) BoolOpt(name string, value bool, desc string) *bool {
	return c.mkOpt(opt{name: name, desc: desc}, value).(*bool)
}

/*
StringOpt defines a string option on the command c named `name`, with an initial value of `value` and a description of `desc` which will be used in help messages.

The name is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).


The result should be stored in a variable (a pointer to a string) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) StringOpt(name string, value string, desc string) *string {
	return c.mkOpt(opt{name: name, desc: desc}, value).(*string)
}

/*
IntOpt defines an int option on the command c named `name`, with an initial value of `value` and a description of `desc` which will be used in help messages.

The name is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).


The result should be stored in a variable (a pointer to an int) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) IntOpt(name string, value int, desc string) *int {
	return c.mkOpt(opt{name: name, desc: desc}, value).(*int)
}

/*
StringsOpt defines a string slice option on the command c named `name`, with an initial value of `value` and a description of `desc` which will be used in help messages.

The name is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).


The result should be stored in a variable (a pointer to a string slice) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) StringsOpt(name string, value []string, desc string) *[]string {
	return c.mkOpt(opt{name: name, desc: desc}, value).(*[]string)
}

/*
IntsOpt defines an int slice option on the command c named `name`, with an initial value of `value` and a description of `desc` which will be used in help messages.

The name is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).


The result should be stored in a variable (a pointer to an int slice) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) IntsOpt(name string, value []int, desc string) *[]int {
	return c.mkOpt(opt{name: name, desc: desc}, value).(*[]int)
}

func (o *opt) isBool() bool {
	return o.value.Elem().Kind() == reflect.Bool
}

func (o *opt) match(args []string, c parseContext) (bool, int) {
	if len(args) == 0 || c.rejectOptions {
		return false, 0
	}
	for _, name := range o.names {
		if args[0] == name {
			if len(args) < 2 {
				return false, 0
			}
			val := args[1]
			c.opts[o] = append(c.opts[o], val)
			return true, 2
		}
	}
	return false, 0
}

func (o *opt) String() string {
	return fmt.Sprintf("Opt(%v)", o.names)
}

func (o *opt) get() interface{} {
	return o.value.Elem().Interface()
}
func (o *opt) set(s string) error {
	return vset(o.value, s)
}

func (c *Cmd) mkOpt(opt opt, defaultValue interface{}) interface{} {
	value := reflect.ValueOf(defaultValue)
	res := reflect.New(value.Type())

	vinit(res, opt.envVar, defaultValue)

	namesSl := strings.Split(opt.name, " ")
	for i, name := range namesSl {
		prefix := "-"
		if len(name) > 1 {
			prefix = "--"
		}
		namesSl[i] = prefix + name
	}

	opt.names = namesSl
	opt.value = res

	c.options = append(c.options, &opt)
	for _, name := range namesSl {
		c.optionsIdx[name] = &opt
	}

	return res.Interface()
}

type optsMatcher []*opt

func (om optsMatcher) try(visited map[*opt]bool, args []string, c parseContext) (bool, int) {
	if len(args) == 0 || c.rejectOptions {
		return false, 0
	}
	for _, o := range om {
		if v, found := visited[o]; !found || !v {
			if ok, cons := o.match(args, c); ok {
				visited[o] = true
				return ok, cons
			}
		}
	}
	return false, 0
}

func (om optsMatcher) match(args []string, c parseContext) (bool, int) {
	visited := map[*opt]bool{}
	ok, cons := om.try(visited, args, c)
	if !ok {
		return false, 0
	}
	consTot := cons
	for {
		ok, cons := om.try(visited, args[consTot:], c)
		if !ok {
			return true, consTot
		}
		consTot += cons
	}
	return true, consTot
}

func (om optsMatcher) String() string {
	return fmt.Sprintf("Opts(%v)", []*opt(om))

}
