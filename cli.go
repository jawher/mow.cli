package cli

import (
	"flag"
	"os"
)

/*
Cli represents the structure of a CLI app. It should be constructed using the App() function
*/
type Cli struct {
	*Cmd
}

/*
App creates a new and empty CLI app configured with the passed name and description.

name and description will be used to construct the help message for the app:

	Usage: $name [OPTIONS] COMMAND [arg...]

	$desc

*/
func App(name, desc string) *Cli {
	return &Cli{
		&Cmd{
			name:          name,
			desc:          desc,
			optionsIdx:    map[string]*opt{},
			argsIdx:       map[string]*arg{},
			ErrorHandling: flag.ExitOnError,
		},
	}
}

/*
Run uses the app configuration (specs, commands, ...) to parse the args slice
and to execute the matching command.

In case of an incorrect usage, and depending on the configured ErrorHandling policy,
it may return an error, panic or exit
*/
func (cli *Cli) Run(args []string) error {
	if err := cli.doInit(); err != nil {
		panic(err)
	}
	inFlow := &step{desc: "RootIn"}
	outFlow := &step{desc: "RootOut"}
	return cli.parse(args[1:], inFlow, inFlow, outFlow)
}

/*
Exit causes the app the exit with the specified exit code while giving the After interceptors a chance to run.
This should be used instead of os.Exit.
 */
func Exit(code int) {
	panic(exit(code))
}

type exit int

var exiter = func(code int) {
	os.Exit(code)
}
