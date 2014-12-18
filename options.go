package cli

import (
	"fmt"
	"reflect"
	"strings"
)

type option struct {
	names []string
	desc  string
	value reflect.Value
}

type OptExtra struct {
	EnvVar string
}

func (o *option) isBool() bool {
	return o.value.Elem().Kind() == reflect.Bool
}

func (o *option) match(args []string, c parseContext) (bool, int) {
	if len(args) == 0 {
		return false, 0
	}
	for _, name := range o.names {
		if args[0] == name {
			if o.isBool() {
				c.opts[o] = append(c.opts[o], "true")
				return true, 1
			}
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

func (o *option) String() string {
	return fmt.Sprintf("Opt(%v)", o.names)
}

func (o *option) get() interface{} {
	return o.value.Elem().Interface()
}
func (o *option) set(s string) error {
	return vset(o.value, s)
}

func (c *Cmd) mkOpt(names string, defaultValue interface{}, desc string, extra *OptExtra) interface{} {
	value := reflect.ValueOf(defaultValue)
	res := reflect.New(value.Type())

	envVars := ""
	if extra != nil {
		envVars = extra.EnvVar
	}
	vinit(res, envVars, defaultValue)

	namesSl := strings.Split(names, " ")
	for i, name := range namesSl {
		prefix := "-"
		if len(name) > 1 {
			prefix = "--"
		}
		namesSl[i] = prefix + name
	}

	opt := &option{
		namesSl,
		desc,
		res,
	}
	c.options = append(c.options, opt)
	for _, name := range namesSl {
		c.optionsIdx[name] = opt
	}

	return res.Interface()
}

/*
BoolOpt defines a boolean option (flag) on the command c with the names `names` with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: git push [-f]

	Options:
 	  -f, --force=false      $desc

`names` is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to a bool) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) BoolOpt(names string, value bool, desc string, extra *OptExtra) *bool {
	return c.mkOpt(names, value, desc, extra).(*bool)
}

/*
StringOpt defines a string option on the command c with the names `names` with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: git clone [-o]

	Options:
 	  -o, --origin=""      $desc

`names` is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to a string) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) StringOpt(names string, value string, desc string, extra *OptExtra) *string {
	return c.mkOpt(names, value, desc, extra).(*string)
}

/*
StringsOpt defines a string slice option on the command c with the names `names` with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: docker run [-e...]

	Options:
 	  -e, --env=[]      $desc

`names` is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to a string slice) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) StringsOpt(names string, value []string, desc string, extra *OptExtra) *[]string {
	return c.mkOpt(names, value, desc, extra).(*[]string)
}

/*
IntOpt defines an int option on the command c with the names `names` with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: tail [-n]

	Options:
 	  -n, --number=0      $desc

`names` is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to an int) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) IntOpt(names string, value int, desc string, extra *OptExtra) *int {
	return c.mkOpt(names, value, desc, extra).(*int)
}

/*
IntsOpt defines an int slice option on the command c with the names `names` with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: bar [-n...]

	Options:
 	  -n, --number=[]      $desc

`names` is a space separated list of the option names *WITHOUT* the dashes, e.g. `f force` and *NOT* `-f --force`.
The one letter names will then be called with a single dash (short option), the others with two (long options).

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to an int slice) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) IntsOpt(names string, value []int, desc string, extra *OptExtra) *[]int {
	return c.mkOpt(names, value, desc, extra).(*[]int)
}
