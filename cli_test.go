package cli

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/stretchr/testify/require"

	"testing"
)

func TestTheCpCase(t *testing.T) {
	app := App("cp", "")
	app.Spec = "SRC... DST"

	src := app.StringsArg("SRC", nil, "", nil)
	dst := app.StringArg("DST", "", "", nil)

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
	x := app.BoolOpt("x", false, "", nil)
	y := app.StringOpt("y", "", "", nil)
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

func forkTest(testName string, fork func(), test func(err error)) {
	if os.Getenv("MOW_DO_IT") == "1" {
		fork()
	} else {
		cmd := exec.Command(os.Args[0], "-test.run="+testName)
		cmd.Stderr = ioutil.Discard
		cmd.Stdout = ioutil.Discard

		cmd.Env = append(os.Environ(), "MOW_DO_IT=1")
		test(cmd.Run())
	}
}

func TestHelpShortcut(t *testing.T) {
	forkTest("TestHelpShortcut",
		func() {
			app := App("x", "")
			app.Spec = "Y"

			app.StringArg("Y", "", "", nil)
			app.Run([]string{"x", "y", "-h", "z"})
		},
		func(err error) {
			fmt.Printf("test fork err %v\n", err)
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				return
			}
			t.Fatalf("process ran with err %v, want exit status != 0", err)
		})
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
	app := App("say", "")
	app.StringOpt("f", "", "", nil)
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
	forkTest("TestHelpShortcut",
		func() {
			app := App("x", "")
			app.Spec = "Y"

			app.StringArg("Y", "", "", nil)
			app.Run([]string{"x", "y", "z"})
		},
		func(err error) {
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				return
			}
			t.Fatalf("process ran with err %v, want exit status != 0", err)
		})
}

func TestPanicOnError(t *testing.T) {
	app := App("say", "")
	app.StringOpt("f", "", "", nil)
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
