package cli

import (
	"flag"
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
			optionsIdx:    map[string]*option{},
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
	return cli.parse(args[1:])
}
