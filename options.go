package cli

import (
	"fmt"
	"reflect"
	"strings"
)

// BoolOpt describes a boolean option
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

// StringOpt describes a string option
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

// IntOpt describes an int option
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

// StringsOpt describes a string slice option
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

// IntsOpt describes an int slice option
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

type opt struct {
	name          string
	desc          string
	envVar        string
	names         []string
	helpFormatter func(interface{}) string
	value         reflect.Value
	hideValue     bool
}

func (o *opt) isBool() bool {
	return o.value.Elem().Kind() == reflect.Bool
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

func mkOptStrs(optName string) []string {
	namesSl := strings.Split(optName, " ")
	for i, name := range namesSl {
		prefix := "-"
		if len(name) > 1 {
			prefix = "--"
		}
		namesSl[i] = prefix + name
	}
	return namesSl
}

func (c *Cmd) mkOpt(opt opt, defaultValue interface{}) interface{} {
	value := reflect.ValueOf(defaultValue)
	res := reflect.New(value.Type())

	opt.helpFormatter = formatterFor(value.Type())

	vinit(res, opt.envVar, defaultValue)

	opt.names = mkOptStrs(opt.name)
	opt.value = res

	c.options = append(c.options, &opt)
	for _, name := range opt.names {
		c.optionsIdx[name] = &opt
	}

	return res.Interface()
}
