package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jawher/mow.cli/internal/container"
	"github.com/jawher/mow.cli/internal/flow"
)

/*
Cli represents the structure of a CLI app. It should be constructed using the App() function
*/
type Cli struct {
	*Cmd
	version *cliVersion
	exiter  func(code int) // REVIEW: might be desirable to have other options than just callback; the most common use I can imagine, other than `os.Exit`, is to simply capture the value again, for which a callback is both overkill and high friction.
	stdOut  io.Writer      // REVIEW: I brought this along for the ride, because it was there at package scope before, but... it's never used, afaict?
	stdErr  io.Writer
}

type cliVersion struct {
	version string
	option  *container.Container
}

/*
App creates a new and empty CLI app configured with the passed name and description.

name and description will be used to construct the help message for the app:

	Usage: $name [OPTIONS] COMMAND [arg...]

	$desc

*/
func App(name, desc string) *Cli {
	cli := &Cli{
		exiter: func(code int) { os.Exit(code) },
		stdOut: os.Stdout,
		stdErr: os.Stderr,
	}
	cli.Cmd = &Cmd{
		cli:           cli,
		name:          name,
		desc:          desc,
		optionsIdx:    map[string]*container.Container{},
		argsIdx:       map[string]*container.Container{},
		ErrorHandling: flag.ExitOnError,
	}
	return cli
}

/*
Version sets the version string of the CLI app together with the options that can be used to trigger
printing the version string via the CLI.

	Usage: appName --$name
	$version

*/
func (cli *Cli) Version(name, version string) {
	cli.Bool(BoolOpt{
		Name:      name,
		Value:     false,
		Desc:      "Show the version and exit",
		HideValue: true,
	})
	names := mkOptStrs(name)
	option := cli.optionsIdx[names[0]]
	cli.version = &cliVersion{version, option}
}

/*
SetStdout sets the CLI's concept of what is "standard out".

If SetStdout is not called, the default behavior is to use os.Stdout.

This is currently unused.
*/
func (cli *Cli) SetStdout(wr io.Writer) {
	cli.stdOut = wr
}

/*
SetStderr sets the CLI's concept of what is "standard error".

If SetStderr is not called, the default behavior is to use os.Stderr.

Information about parse errors is written to this stream,
as well as usage info or version info, if those are requested.
*/
func (cli *Cli) SetStderr(wr io.Writer) {
	cli.stdErr = wr
}

/*
SetExiter sets a callback to define what happens when an exit should happen with an exit code.

If SetExiter is not called, the default behavior is to call os.Exit
(which immediately halts the program).

Common uses of setting a custom exit function include gathering the code instead of halting the program
(which is often useful for writing tests of CLI behavior, for example).

SetExiter should not be used for cleanup hooks; use a Cmd.After callback for that.
*/
func (cli *Cli) SetExiter(exiter func(code int)) {
	cli.exiter = exiter
}

func (cli *Cli) parse(args []string, entry, inFlow, outFlow *flow.Step) error {
	// We overload Cmd.parse() and handle cases that only apply to the CLI command, like versioning
	// After that, we just call Cmd.parse() for the default behavior
	if cli.versionSetAndRequested(args) {
		cli.PrintVersion()
		cli.onError(errVersionRequested)
		return nil
	}
	return cli.Cmd.parse(args, entry, inFlow, outFlow)
}

func (cli *Cli) versionSetAndRequested(args []string) bool {
	return cli.version != nil && cli.isFirstItemAmong(args, cli.version.option.Names)
}

/*
PrintVersion prints the CLI app's version.
In most cases the library users won't need to call this method, unless
a more complex validation is needed.
*/
func (cli *Cli) PrintVersion() {
	fmt.Fprintln(cli.stdErr, cli.version.version)
}

/*
Run uses the app configuration (specs, commands, ...) to parse the args slice
and to execute the matching command.

In case of an incorrect usage, and depending on the configured ErrorHandling policy,
it may return an error, panic or exit
*/
func (cli *Cli) Run(args []string) error { // REVIEW: I would actually prefer this returned `(error, int)`... but, that's a breaking change.  We could also: use special error types for code; or, introduce a new function; or, do nothing, and require users to write a capturing thunk for an `exiter`.
	if err := cli.doInit(); err != nil {
		panic(err)
	}
	inFlow := &flow.Step{Desc: "RootIn", Exiter: cli.exiter}
	outFlow := &flow.Step{Desc: "RootOut", Exiter: cli.exiter}
	return cli.parse(args[1:], inFlow, inFlow, outFlow)
}

/*
ActionCommand is a convenience function to configure a command with an action.

cmd.ActionCommand(_, _, myFunc) is equivalent to cmd.Command(_, _, func(cmd *cli.Cmd) { cmd.Action = myFunc })
*/
func ActionCommand(action func()) CmdInitializer {
	return func(cmd *Cmd) {
		cmd.Action = action
	}
}

/*
Exit causes the app the exit with the specified exit code while giving the After interceptors a chance to run.
This should be used instead of os.Exit.

This function is implemented using a panic; nothing will occur after it is called
(except other deferred functions, and the After intercepters).
*/
func Exit(code int) {
	panic(flow.ExitCode(code))
}
