package cli

import (
	"fmt"
	"reflect"
	"strings"
)

type arg struct {
	name  string
	desc  string
	value reflect.Value
}

/*
Extra configuration for a command argument
*/
type ArgExtra struct {
	// A list of space separated environment variables names to be used to initialize the argument
	EnvVar string
}

func (arg *arg) match(args []string, c parseContext) (bool, int) {
	if len(args) == 0 {
		return false, 0
	}
	if strings.HasPrefix(args[0], "-") {
		return false, 0
	}
	c.args[arg] = append(c.args[arg], args[0])
	return true, 1
}

func (a *arg) String() string {
	return fmt.Sprintf("ARG(%s)", a.name)
}

func (a *arg) get() interface{} {
	return a.value.Elem().Interface()
}

func (a *arg) set(s string) error {
	return vset(a.value, s)
}

func (c *Cmd) mkArg(name string, defaultValue interface{}, desc string, extra *ArgExtra) interface{} {
	value := reflect.ValueOf(defaultValue)
	res := reflect.New(value.Type())

	envVars := ""
	if extra != nil {
		envVars = extra.EnvVar
	}
	vinit(res, envVars, defaultValue)

	arg := &arg{
		name,
		desc,
		res,
	}
	c.args = append(c.args, arg)
	c.argsIdx[name] = arg

	return res.Interface()
}

/*
BoolArg defines a boolean argument on the command c named `name`, with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: git clone REPO

	Arguments:
 	  REPO=""      $desc

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to a bool) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) BoolArg(name string, value bool, desc string, extra *ArgExtra) *bool {
	return c.mkArg(name, value, desc, extra).(*bool)
}

/*
StringArg defines a string argument on the command c named `name`, with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: git clone REPO

	Arguments:
 	  REPO=""      $desc

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to a string) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) StringArg(name string, value string, desc string, extra *ArgExtra) *string {
	return c.mkArg(name, value, desc, extra).(*string)
}

/*
StringsArg defines a string slice argument on the command c named `name`, with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: cp SRC...

	Arguments:
 	  SRC=[]      $desc

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to a string slice) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) StringsArg(name string, value []string, desc string, extra *ArgExtra) *[]string {
	return c.mkArg(name, value, desc, extra).(*[]string)
}

/*
IntArg defines an int argument on the command c named `name`, with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: foo COUNT

	Arguments:
 	  COUNT=0      $desc

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to an int) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) IntArg(name string, value int, desc string, extra *ArgExtra) *int {
	return c.mkArg(name, value, desc, extra).(*int)
}

/*
IntsArg defines an int slice argument on the command c named `name`, with an initial value of `value` and a description of `desc`
which will be used in help messages, e.g.:

	Usage: sort NUMBERS...

	Arguments:
 	  NUMBERS=[]      $desc

When needed, extra can be used to pass more options.

The result should be stored in a variable (a pointer to an int slice) which will be populated when the app is run and the call arguments get parsed
*/
func (c *Cmd) IntsArg(name string, value []int, desc string, extra *ArgExtra) *[]int {
	return c.mkArg(name, value, desc, extra).(*[]int)
}
