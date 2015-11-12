package cli

import (
	"flag"

	"github.com/stretchr/testify/require"

	"testing"
)

func TestTheCpCase(t *testing.T) {
	app := App("cp", "")
	app.Spec = "SRC... DST"

	src := app.Strings(StringsArg{Name: "SRC", Value: nil, Desc: ""})
	dst := app.String(StringArg{Name: "DST", Value: "", Desc: ""})

	ex := false
	app.Action = func() {
		ex = true
	}
	app.Run([]string{"cp", "x", "y", "z"})

	require.Equal(t, []string{"x", "y"}, *src)
	require.Equal(t, "z", *dst)

	require.True(t, ex, "Exec wasn't called")
}

func TestImplicitSpec(t *testing.T) {
	app := App("test", "")
	x := app.Bool(BoolOpt{Name: "x", Value: false, Desc: ""})
	y := app.String(StringOpt{Name: "y", Value: "", Desc: ""})
	called := false
	app.Action = func() {
		called = true
	}
	app.ErrorHandling = flag.ContinueOnError

	err := app.Run([]string{"test", "-x", "-y", "hello"})

	require.Nil(t, err)
	require.True(t, *x)
	require.Equal(t, "hello", *y)

	require.True(t, called, "Exec wasn't called")
}

func TestHelpShortcut(t *testing.T) {
	defer suppressOutput()()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 2, &exitCalled)()

	app := App("x", "")
	app.Spec = "Y"

	app.String(StringArg{Name: "Y", Value: "", Desc: ""})

	actionCalled := false
	app.Action = func() {
		actionCalled = true
	}
	app.Run([]string{"x", "y", "-h", "z"})

	require.False(t, actionCalled, "action should not have been called")
	require.True(t, exitCalled, "exit should have been called")
}

func TestHelpMessage(t *testing.T) {
	var out, err string
	defer captureAndRestoreOutput(&out, &err)()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 2, &exitCalled)()

	app := App("app", "App Desc")
	app.Spec = "[-o] ARG"

	app.String(StringOpt{Name: "o opt", Value: "", Desc: "Option"})
	app.String(StringArg{Name: "ARG", Value: "", Desc: "Argument"})

	app.Action = func() {}
	app.Run([]string{"app", "-h"})

	help := `
Usage: app [-o] ARG

App Desc

Arguments:
  ARG=""       Argument

Options:
  -o, --opt=""   Option
`

	require.Equal(t, help, err)
}

func TestLongHelpMessage(t *testing.T) {
	var out, err string
	defer captureAndRestoreOutput(&out, &err)()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 2, &exitCalled)()

	app := App("app", "App Desc")
	app.LongDesc = "Longer App Desc"
	app.Spec = "[-o] ARG"

	app.String(StringOpt{Name: "o opt", Value: "", Desc: "Option"})
	app.String(StringArg{Name: "ARG", Value: "", Desc: "Argument"})

	app.Action = func() {}
	app.Run([]string{"app", "-h"})

	help := `
Usage: app [-o] ARG

Longer App Desc

Arguments:
  ARG=""       Argument

Options:
  -o, --opt=""   Option
`

	require.Equal(t, help, err)
}

func TestVersionShortcut(t *testing.T) {
	defer suppressOutput()()
	exitCalled := false
	defer exitShouldBeCalledWith(t, 0, &exitCalled)()

	app := App("cp", "")
	app.Version("v version", "cp 1.2.3")

	actionCalled := false
	app.Action = func() {
		actionCalled = true
	}

	app.Run([]string{"cp", "--version"})

	require.False(t, actionCalled, "action should not have been called")
	require.True(t, exitCalled, "exit should have been called")
}

func TestSubCommands(t *testing.T) {
	app := App("say", "")

	hi, bye := false, false

	app.Command("hi", "", func(cmd *Cmd) {
		cmd.Action = func() {
			hi = true
		}
	})

	app.Command("byte", "", func(cmd *Cmd) {
		cmd.Action = func() {
			bye = true
		}
	})

	app.Run([]string{"say", "hi"})
	require.True(t, hi, "hi should have been called")
	require.False(t, bye, "byte should NOT have been called")
}

func TestContinueOnError(t *testing.T) {
	defer exitShouldNotCalled(t)()
	defer suppressOutput()()

	app := App("say", "")
	app.String(StringOpt{Name: "f", Value: "", Desc: ""})
	app.Spec = "-f"
	app.ErrorHandling = flag.ContinueOnError
	called := false
	app.Action = func() {
		called = true
	}

	err := app.Run([]string{"say"})
	require.NotNil(t, err)
	require.False(t, called, "Exec should NOT have been called")
}

func TestExitOnError(t *testing.T) {
	defer suppressOutput()()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 2, &exitCalled)()

	app := App("x", "")
	app.Spec = "Y"

	app.String(StringArg{Name: "Y", Value: "", Desc: ""})
	app.Run([]string{"x", "y", "z"})
	require.True(t, exitCalled, "exit should have been called")
}

func TestPanicOnError(t *testing.T) {
	defer suppressOutput()()

	app := App("say", "")
	app.String(StringOpt{Name: "f", Value: "", Desc: ""})
	app.Spec = "-f"
	app.ErrorHandling = flag.PanicOnError
	called := false
	app.Action = func() {
		called = true
	}

	defer func() {
		if r := recover(); r != nil {
			require.False(t, called, "Exec should NOT have been called")
		} else {

		}
	}()
	app.Run([]string{"say"})
	t.Fatalf("wanted panic")
}
