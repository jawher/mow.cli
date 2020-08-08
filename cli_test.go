package cli

import (
	"bytes"
	"flag"

	"github.com/stretchr/testify/require"

	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"fmt"

	"github.com/jawher/mow.cli/internal/flow"
)

func TestTheCpCase(t *testing.T) {
	app := App("cp", "")
	app.Spec = "SRC... DST"

	src := app.Strings(StringsArg{Name: "SRC", Value: nil, Desc: ""})
	dst := app.String(StringArg{Name: "DST", Value: "", Desc: ""})

	ex := false
	app.Action = func() {
		ex = true

		require.Equal(t, []string{"x", "y"}, *src)
		require.Equal(t, "z", *dst)
	}

	require.NoError(t,
		app.Run([]string{"cp", "x", "y", "z"}))

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

	require.NoError(t, err)
	require.True(t, *x)
	require.Equal(t, "hello", *y)

	require.True(t, called, "Exec wasn't called")
}

func TestExit(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Logf("Panicked with %v", p)
			require.Equal(t, flow.ExitCode(666), p)
			return
		}
		t.Fatalf("Should have panicked")
	}()

	Exit(666)

	t.Fatalf("Should have panicked")
}

func TestInvalidSpec(t *testing.T) {
	app := App("test", "")
	app.Spec = "X"

	called := false
	app.Action = func() {
		called = true
	}

	defer func() {
		if p := recover(); p != nil {
			t.Logf("Panicked with %v", p)
			require.False(t, called, "action should not have been called")
			return
		}
		t.Fatalf("Should have panicked")
	}()

	require.NoError(t,
		app.Run([]string{"test", "-x", "-y", "hello"}))

	t.Fatalf("Should have panicked")
}

func TestDuplicateOptionName(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Logf("Panicked with %v", p)
			return
		}
		t.Fatalf("Should have panicked")
	}()

	app := App("test", "")
	app.BoolOpt("f force", false, "")
	app.StringOpt("f file", "", "")

	t.Fatalf("Should have panicked")
}

func TestDuplicateArgName(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Logf("Panicked with %v", p)
			return
		}
		t.Fatalf("Should have panicked")
	}()

	app := App("test", "")
	app.StringArg("ARG", "", "")
	app.StringArg("ARG", "", "")

	t.Fatalf("Should have panicked")
}

func TestInvalidArgName(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Logf("Panicked with %v", p)
			return
		}
		t.Fatalf("Should have panicked")
	}()

	app := App("test", "")
	app.StringArg("arg", "", "")

	t.Fatalf("Should have panicked")
}

func runAppAndCheckValue(t *testing.T, name string, callArgs []string, expected interface{}, appConfigurer func(app *Cli) interface{}) {
	t.Run(fmt.Sprintf("%s %+v", name, callArgs), func(t *testing.T) {
		t.Logf("Testing %+v", callArgs)

		app := App("app", "")
		app.ErrorHandling = flag.ContinueOnError
		opt := appConfigurer(app)

		ex := false
		app.Action = func() {
			ex = true
			require.Equal(t, expected, reflect.ValueOf(opt).Elem().Interface())
		}
		err := app.Run(callArgs)

		require.NoError(t, err)
		require.True(t, ex, "Exec wasn't called")
	})
}

func TestAppWithBoolOption(t *testing.T) {

	cases := []struct {
		args             []string
		expectedOptValue bool
	}{
		{[]string{"app"}, false},
		{[]string{"app", "-o"}, true},
		{[]string{"app", "-o=true"}, true},
		{[]string{"app", "-o=false"}, false},

		{[]string{"app", "--option"}, true},
		{[]string{"app", "--option=true"}, true},
		{[]string{"app", "--option=false"}, false},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.BoolOpt("o option", false, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.Bool(BoolOpt{
				Name: "o option",
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val bool
			app.BoolOptPtr(&val, "o option", false, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val bool
			app.BoolPtr(&val, BoolOpt{
				Name: "o option",
			})
			return &val
		})
	}
}

func TestAppWithStringOption(t *testing.T) {

	cases := []struct {
		args             []string
		expectedOptValue string
	}{
		{[]string{"app"}, "default"},
		{[]string{"app", "-o", "user"}, "user"},
		{[]string{"app", "-o=user"}, "user"},

		{[]string{"app", "--option", "user"}, "user"},
		{[]string{"app", "--option=user"}, "user"},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.StringOpt("o option", "default", "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.String(StringOpt{
				Name:  "o option",
				Value: "default",
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val string
			app.StringOptPtr(&val, "o option", "default", "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val string
			app.StringPtr(&val, StringOpt{
				Name:  "o option",
				Value: "default",
			})
			return &val
		})
	}
}

func TestAppWithIntOption(t *testing.T) {

	cases := []struct {
		args             []string
		expectedOptValue int
	}{
		{[]string{"app"}, 3},
		{[]string{"app", "-o", "16"}, 16},
		{[]string{"app", "-o=16"}, 16},

		{[]string{"app", "--option", "16"}, 16},
		{[]string{"app", "--option=16"}, 16},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.IntOpt("o option", 3, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.Int(IntOpt{
				Name:  "o option",
				Value: 3,
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val int
			app.IntOptPtr(&val, "o option", 3, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val int
			app.IntPtr(&val, IntOpt{
				Name:  "o option",
				Value: 3,
			})
			return &val
		})
	}
}

func TestAppWithFloat64Option(t *testing.T) {

	cases := []struct {
		args             []string
		expectedOptValue float64
	}{
		{[]string{"app"}, 3.14},
		{[]string{"app", "-o", "16.0001"}, 16.0001},
		{[]string{"app", "-o=16.0001"}, 16.0001},

		{[]string{"app", "--option", "16.0001"}, 16.0001},
		{[]string{"app", "--option=16.0001"}, 16.0001},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.Float64Opt("o option", 3.14, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.Float64(Float64Opt{
				Name:  "o option",
				Value: 3.14,
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val float64
			app.Float64OptPtr(&val, "o option", 3.14, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val float64
			app.Float64Ptr(&val, Float64Opt{
				Name:  "o option",
				Value: 3.14,
			})
			return &val
		})
	}
}

func TestAppWithStringsOption(t *testing.T) {

	cases := []struct {
		args             []string
		expectedOptValue []string
	}{
		{[]string{"app"}, []string{"a", "b"}},
		{[]string{"app", "-o", "x"}, []string{"x"}},
		{[]string{"app", "-o", "x", "-o=y"}, []string{"x", "y"}},

		{[]string{"app", "--option", "x"}, []string{"x"}},
		{[]string{"app", "--option", "x", "--option=y"}, []string{"x", "y"}},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.StringsOpt("o option", []string{"a", "b"}, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.Strings(StringsOpt{
				Name:  "o option",
				Value: []string{"a", "b"},
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val []string
			app.StringsOptPtr(&val, "o option", []string{"a", "b"}, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val []string
			app.StringsPtr(&val, StringsOpt{
				Name:  "o option",
				Value: []string{"a", "b"},
			})
			return &val
		})
	}
}

func TestAppWithIntsOption(t *testing.T) {

	cases := []struct {
		args             []string
		expectedOptValue []int
	}{
		{[]string{"app"}, []int{1, 2}},
		{[]string{"app", "-o", "10"}, []int{10}},
		{[]string{"app", "-o", "10", "-o=11"}, []int{10, 11}},

		{[]string{"app", "--option", "10"}, []int{10}},
		{[]string{"app", "--option", "10", "--option=11"}, []int{10, 11}},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.IntsOpt("o option", []int{1, 2}, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.Ints(IntsOpt{
				Name:  "o option",
				Value: []int{1, 2},
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val []int
			app.IntsOptPtr(&val, "o option", []int{1, 2}, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val []int
			app.IntsPtr(&val, IntsOpt{
				Name:  "o option",
				Value: []int{1, 2},
			})
			return &val
		})
	}
}

func TestAppWithFloats64Option(t *testing.T) {

	cases := []struct {
		args             []string
		expectedOptValue []float64
	}{
		{[]string{"app"}, []float64{1.1, 2.2}},
		{[]string{"app", "-o", "10.05"}, []float64{10.05}},
		{[]string{"app", "-o", "10.05", "-o=11.993"}, []float64{10.05, 11.993}},

		{[]string{"app", "--option", "10.05"}, []float64{10.05}},
		{[]string{"app", "--option", "10.05", "--option=11.993"}, []float64{10.05, 11.993}},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.Floats64Opt("o option", []float64{1.1, 2.2}, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			return app.Floats64(Floats64Opt{
				Name:  "o option",
				Value: []float64{1.1, 2.2},
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val []float64
			app.Floats64OptPtr(&val, "o option", []float64{1.1, 2.2}, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedOptValue, func(app *Cli) interface{} {
			var val []float64
			app.Floats64Ptr(&val, Floats64Opt{
				Name:  "o option",
				Value: []float64{1.1, 2.2},
			})
			return &val
		})
	}
}

func TestAppWithBoolArg(t *testing.T) {

	cases := []struct {
		args             []string
		expectedArgValue bool
	}{
		{[]string{"app"}, false},
		{[]string{"app", "true"}, true},
		{[]string{"app", "false"}, false},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			return app.BoolArg("ARG", false, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			return app.Bool(BoolArg{
				Name: "ARG",
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			var val bool
			app.BoolArgPtr(&val, "ARG", false, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			var val bool
			app.BoolPtr(&val, BoolArg{
				Name: "ARG",
			})
			return &val
		})
	}
}

func TestAppWithStringArg(t *testing.T) {

	cases := []struct {
		args             []string
		expectedArgValue string
	}{
		{[]string{"app"}, "default"},
		{[]string{"app", "user"}, "user"},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			return app.StringArg("ARG", "default", "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			return app.String(StringArg{
				Name:  "ARG",
				Value: "default",
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			var val string
			app.StringArgPtr(&val, "ARG", "default", "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			var val string
			app.StringPtr(&val, StringArg{
				Name:  "ARG",
				Value: "default",
			})
			return &val
		})
	}
}

func TestAppWithIntArg(t *testing.T) {

	cases := []struct {
		args             []string
		expectedArgValue int
	}{
		{[]string{"app"}, 3},
		{[]string{"app", "16"}, 16},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			return app.IntArg("ARG", 3, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			return app.Int(IntArg{
				Name:  "ARG",
				Value: 3,
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			var val int
			app.IntArgPtr(&val, "ARG", 3, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			var val int
			app.IntPtr(&val, IntArg{
				Name:  "ARG",
				Value: 3,
			})
			return &val
		})
	}
}

func TestAppWithFloat64Arg(t *testing.T) {

	cases := []struct {
		args             []string
		expectedArgValue float64
	}{
		{[]string{"app"}, 3.14},
		{[]string{"app", "16.123456789"}, 16.123456789},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			return app.Float64Arg("ARG", 3.14, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			return app.Float64(Float64Arg{
				Name:  "ARG",
				Value: 3.14,
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			var val float64
			app.Float64ArgPtr(&val, "ARG", 3.14, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG]"

			var val float64
			app.Float64Ptr(&val, Float64Arg{
				Name:  "ARG",
				Value: 3.14,
			})
			return &val
		})
	}
}

func TestAppWithStringsArg(t *testing.T) {

	cases := []struct {
		args             []string
		expectedArgValue []string
	}{
		{[]string{"app"}, []string{"a", "b"}},
		{[]string{"app", "x"}, []string{"x"}},
		{[]string{"app", "x", "y"}, []string{"x", "y"}},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			return app.StringsArg("ARG", []string{"a", "b"}, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			return app.Strings(StringsArg{
				Name:  "ARG",
				Value: []string{"a", "b"},
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			var val []string
			app.StringsArgPtr(&val, "ARG", []string{"a", "b"}, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			var val []string
			app.StringsPtr(&val, StringsArg{
				Name:  "ARG",
				Value: []string{"a", "b"},
			})
			return &val
		})
	}
}

func TestAppWithIntsArg(t *testing.T) {

	cases := []struct {
		args             []string
		expectedArgValue []int
	}{
		{[]string{"app"}, []int{1, 2}},
		{[]string{"app", "10"}, []int{10}},
		{[]string{"app", "10", "11"}, []int{10, 11}},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			return app.IntsArg("ARG", []int{1, 2}, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			return app.Ints(IntsArg{
				Name:  "ARG",
				Value: []int{1, 2},
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			var val []int
			app.IntsArgPtr(&val, "ARG", []int{1, 2}, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			var val []int
			app.IntsPtr(&val, IntsArg{
				Name:  "ARG",
				Value: []int{1, 2},
			})
			return &val
		})
	}
}

func TestAppWithFloats64Arg(t *testing.T) {

	cases := []struct {
		args             []string
		expectedArgValue []float64
	}{
		{[]string{"app"}, []float64{1.1, 2.2}},
		{[]string{"app", "10.123"}, []float64{10.123}},
		{[]string{"app", "10.123", "11.995"}, []float64{10.123, 11.995}},
	}

	for _, cas := range cases {
		runAppAndCheckValue(t, "short", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			return app.Floats64Arg("ARG", []float64{1.1, 2.2}, "")
		})
		runAppAndCheckValue(t, "struct", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			return app.Floats64(Floats64Arg{
				Name:  "ARG",
				Value: []float64{1.1, 2.2},
			})
		})
		runAppAndCheckValue(t, "short-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			var val []float64
			app.Floats64ArgPtr(&val, "ARG", []float64{1.1, 2.2}, "")
			return &val
		})
		runAppAndCheckValue(t, "struct-ptr", cas.args, cas.expectedArgValue, func(app *Cli) interface{} {
			app.Spec = "[ARG...]"

			var val []float64
			app.Floats64Ptr(&val, Floats64Arg{
				Name:  "ARG",
				Value: []float64{1.1, 2.2},
			})
			return &val
		})
	}
}

func testHelpAndVersionWithOptionsEnd(flag string, t *testing.T) {
	t.Logf("Testing help/version with --: flag=%q", flag)
	defer suppressOutput()()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 0, &exitCalled)()

	app := App("x", "")
	app.Version("v version", "1.0")
	app.Spec = "CMD"

	cmd := app.String(StringArg{Name: "CMD", Value: "", Desc: ""})

	actionCalled := false
	app.Action = func() {
		actionCalled = true
		require.Equal(t, flag, *cmd)
	}

	require.NoError(t,
		app.Run([]string{"x", "--", flag}))

	require.True(t, actionCalled, "action should have been called")
	require.False(t, exitCalled, "exit should not have been called")
}

func TestHelpAndVersionWithOptionsEnd(t *testing.T) {
	for _, flag := range []string{"-h", "--help", "-v", "--version"} {
		t.Run(flag, func(t *testing.T) {
			testHelpAndVersionWithOptionsEnd(flag, t)
		})
	}
}

var genGolden = flag.Bool("g", false, "Generate golden file(s)")

func TestHelpMessage(t *testing.T) {
	cases := []struct {
		name     string
		params   []string
		env      map[string]string
		exitCode int
	}{
		{name: "top-help", params: []string{"app", "-h"}},
		{name: "top-help-i-user", params: []string{"app", "-i=5"}, exitCode: 2},
		{name: "top-help-i-env", params: []string{"app"}, env: map[string]string{"INT1": "25"}, exitCode: 2},
		{name: "command1", params: []string{"app", "command1", "-h"}},
		{name: "command2", params: []string{"app", "command2", "-h"}},
		{name: "command3", params: []string{"app", "command3", "-h"}},
		{name: "command3-child1", params: []string{"app", "command3", "child1", "-h"}},
		{name: "command3-child2", params: []string{"app", "command3", "child2", "-h"}},
		{name: "command4", params: []string{"app", "command4", "-h"}},
	}
	for _, cas := range cases {
		cas := cas
		t.Run(cas.name, func(t *testing.T) {
			t.Logf("case: %+v", cas)
			var out, stdErr string
			defer captureAndRestoreOutput(&out, &stdErr)()
			defer setAndRestoreEnv(cas.env)()

			exitCalled := false
			defer exitShouldBeCalledWith(t, cas.exitCode, &exitCalled)()

			app := App("app", "App Desc")
			app.Spec = "[-bdsuikqs] [BOOL1 STR1 INT3...]"

			// Options
			app.Bool(BoolOpt{Name: "b bool1 u uuu", Value: false, EnvVar: "BOOL1", Desc: "Bool Option 1"})
			app.Bool(BoolOpt{Name: "bool2", Value: true, EnvVar: " ", Desc: "Bool Option 2"})
			app.Bool(BoolOpt{Name: "d", Value: true, EnvVar: "BOOL3", Desc: "Bool Option 3", HideValue: true})

			app.String(StringOpt{Name: "s str1", Value: "", EnvVar: "STR1", Desc: "String Option 1"})
			app.String(StringOpt{Name: "str2", Value: "a value", Desc: "String Option 2"})
			app.String(StringOpt{Name: "v", Value: "another value", EnvVar: "STR3", Desc: "String Option 3", HideValue: true})

			app.Int(IntOpt{Name: "i int1", Value: 0, EnvVar: "INT1 ALIAS_INT1"})
			app.Int(IntOpt{Name: "int2", Value: 1, EnvVar: "INT2", Desc: "Int Option 2"})
			app.Int(IntOpt{Name: "k", Value: 1, EnvVar: "INT3", Desc: "Int Option 3", HideValue: true})

			app.Strings(StringsOpt{Name: "x strs1", Value: nil, EnvVar: "STRS1", Desc: "Strings Option 1"})
			app.Strings(StringsOpt{Name: "strs2", Value: []string{"value1", "value2"}, EnvVar: "STRS2", Desc: "Strings Option 2"})
			app.Strings(StringsOpt{Name: "z", Value: []string{"another value"}, EnvVar: "STRS3", Desc: "Strings Option 3", HideValue: true})

			app.Ints(IntsOpt{Name: "q ints1", Value: nil, EnvVar: "INTS1", Desc: "Ints Option 1"})
			app.Ints(IntsOpt{Name: "ints2", Value: []int{1, 2, 3}, EnvVar: "INTS2", Desc: "Ints Option 2"})
			app.Ints(IntsOpt{Name: "j", Value: []int{1}, EnvVar: "INTS3", Desc: "Ints Option 3", HideValue: true})

			// Args
			app.Bool(BoolArg{Name: "BOOL1", Value: false, EnvVar: "BOOL1", Desc: "Bool Argument 1"})
			app.Bool(BoolArg{Name: "BOOL2", Value: true, Desc: "Bool Argument 2"})
			app.Bool(BoolArg{Name: "BOOL3", Value: true, EnvVar: "BOOL3", Desc: "Bool Argument 3", HideValue: true})

			app.String(StringArg{Name: "STR1", Value: "", EnvVar: "STR1", Desc: "String Argument 1"})
			app.String(StringArg{Name: "STR2", Value: "a value", EnvVar: "STR2", Desc: "String Argument 2"})
			app.String(StringArg{Name: "STR3", Value: "another value", EnvVar: "STR3", Desc: "String Argument 3", HideValue: true})

			app.Int(IntArg{Name: "INT1", Value: 0, EnvVar: "INT1", Desc: "Int Argument 1"})
			app.Int(IntArg{Name: "INT2", Value: 1, EnvVar: "INT2", Desc: "Int Argument 2"})
			app.Int(IntArg{Name: "INT3", Value: 1, EnvVar: "INT3", Desc: "Int Argument 3", HideValue: true})

			app.Strings(StringsArg{Name: "STRS1", Value: nil, EnvVar: "STRS1", Desc: "Strings Argument 1"})
			app.Strings(StringsArg{Name: "STRS2", Value: []string{"value1", "value2"}, EnvVar: "STRS2"})
			app.Strings(StringsArg{Name: "STRS3", Value: []string{"another value"}, EnvVar: "STRS3", Desc: "Strings Argument 3", HideValue: true})

			app.Ints(IntsArg{Name: "INTS1", Value: nil, EnvVar: "INTS1", Desc: "Ints Argument 1"})
			app.Ints(IntsArg{Name: "INTS2", Value: []int{1, 2, 3}, EnvVar: "INTS2", Desc: "Ints Argument 2"})
			app.Ints(IntsArg{Name: "INTS3", Value: []int{1}, EnvVar: "INTS3", Desc: "Ints Argument 3", HideValue: true})

			app.Command("command1", "command1 description", func(cmd *Cmd) {})
			app.Command("command2", "command2 description", func(cmd *Cmd) {})
			app.Command("command3", "command3 description", func(cmd *Cmd) {
				cmd.Command("child1", "child1 description", func(cmd *Cmd) {
					cmd.StringArg("ARG1", "", "arg1 desc")
				})
				cmd.Command("child2", "child2 description", func(cmd *Cmd) {
					cmd.Hidden = true
					cmd.StringOpt("o opt", "", "opt desc")
				})
			})
			app.Command("command4", "command4 description", func(cmd *Cmd) {
				cmd.Hidden = true
			})

			fmt.Printf("calling app with %+v\n", cas.params)
			require.NoError(t,
				app.Run(cas.params))

			filename := fmt.Sprintf("testdata/help-output-%s.txt", cas.name)

			if *genGolden {
				require.NoError(t,
					ioutil.WriteFile(filename, []byte(stdErr), 0644))
			}

			expected, e := ioutil.ReadFile(filename)
			require.NoError(t, e, "Failed to read the expected help output from %s", filename)

			require.Equal(t, string(expected), stdErr)
		})
	}
}

func TestLongHelpMessage(t *testing.T) {
	var out, err string
	defer captureAndRestoreOutput(&out, &err)()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 0, &exitCalled)()

	app := App("app", "App Desc")
	app.LongDesc = "Longer App Desc"
	app.Spec = "[-o] ARG"

	app.String(StringOpt{Name: "o opt", Value: "", Desc: "Option"})
	app.String(StringArg{Name: "ARG", Value: "", Desc: "Argument"})

	app.Action = func() {}
	require.NoError(t,
		app.Run([]string{"app", "-h"}))

	if *genGolden {
		require.NoError(t,
			ioutil.WriteFile("testdata/long-help-output.txt.golden", []byte(err), 0644))
	}

	expected, e := ioutil.ReadFile("testdata/long-help-output.txt")
	require.NoError(t, e, "Failed to read the expected help output from testdata/long-help-output.txt")

	require.Equal(t, expected, []byte(err))
}

func TestMultiLineDescInHelpMessage(t *testing.T) {
	var out, err string
	defer captureAndRestoreOutput(&out, &err)()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 0, &exitCalled)()

	app := App("app", "App Desc")
	app.LongDesc = "Longer App Desc"
	app.Spec = "[-o] ARG"

	app.String(StringOpt{Name: "o opt", Value: "default", EnvVar: "XXX_TEST", Desc: "Option\ndoes something\n  another line"})
	app.Bool(BoolOpt{Name: "f force", Value: false, EnvVar: "YYY_TEST", Desc: "Force\ndoes something\n  another line"})
	app.String(StringArg{Name: "ARG", Value: "", Desc: "Argument\nDescription\nMultiple\nLines"})

	app.Action = func() {}
	require.NoError(t,
		app.Run([]string{"app", "-h"}))

	if *genGolden {
		require.NoError(t,
			ioutil.WriteFile("testdata/multi-line-desc-help-output.txt.golden", []byte(err), 0644))
	}

	expected, e := ioutil.ReadFile("testdata/multi-line-desc-help-output.txt")
	require.NoError(t, e, "Failed to read the expected help output from testdata/long-help-output.txt")

	require.Equal(t, expected, []byte(err))
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

	require.NoError(t,
		app.Run([]string{"cp", "--version"}))

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

	require.NoError(t,
		app.Run([]string{"say", "hi"}))
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

func TestContinueOnErrorWithHelpAndVersion(t *testing.T) {
	defer exitShouldNotCalled(t)()
	defer suppressOutput()()

	app := App("say", "")
	app.Version("v", "1.0")
	app.String(StringOpt{Name: "f", Value: "", Desc: ""})
	app.Spec = "-f"
	app.ErrorHandling = flag.ContinueOnError
	called := false
	app.Action = func() {
		called = true
	}

	{
		err := app.Run([]string{"say", "-h"})
		require.Nil(t, err)
		require.False(t, called, "Exec should NOT have been called")
	}

	{
		err := app.Run([]string{"say", "-v"})
		require.Nil(t, err)
		require.False(t, called, "Exec should NOT have been called")
	}
}

func TestExitOnError(t *testing.T) {
	defer suppressOutput()()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 2, &exitCalled)()

	app := App("x", "")
	app.ErrorHandling = flag.ExitOnError
	app.Spec = "Y"

	app.String(StringArg{Name: "Y", Value: "", Desc: ""})

	require.Error(t,
		app.Run([]string{"x", "y", "z"}))
	require.True(t, exitCalled, "exit should have been called")
}

func TestExitOnErrorWithHelp(t *testing.T) {
	defer suppressOutput()()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 0, &exitCalled)()

	app := App("x", "")
	app.Spec = "Y"
	app.ErrorHandling = flag.ExitOnError

	app.String(StringArg{Name: "Y", Value: "", Desc: ""})

	require.NoError(t,
		app.Run([]string{"x", "-h"}))
	require.True(t, exitCalled, "exit should have been called")
}

func TestExitOnErrorWithVersion(t *testing.T) {
	defer suppressOutput()()

	exitCalled := false
	defer exitShouldBeCalledWith(t, 0, &exitCalled)()

	app := App("x", "")
	app.Version("v", "1.0")
	app.Spec = "Y"
	app.ErrorHandling = flag.ExitOnError

	app.String(StringArg{Name: "Y", Value: "", Desc: ""})

	require.NoError(t,
		app.Run([]string{"x", "-v"}))
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
		}
	}()
	require.NoError(t,
		app.Run([]string{"say"}))
	t.Fatalf("wanted panic")
}

func TestOptSetByUser(t *testing.T) {
	cases := []struct {
		desc     string
		config   func(*Cli, *bool)
		args     []string
		expected bool
	}{
		// OPTS
		// String
		{
			desc: "String Opt, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.String(StringOpt{Name: "f", Value: "a", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "String Opt, not set by user, env value",
			config: func(c *Cli, s *bool) {
				os.Setenv("MOW_VALUE", "value")
				c.String(StringOpt{Name: "f", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "String Opt, set by user",
			config: func(c *Cli, s *bool) {
				c.String(StringOpt{Name: "f", Value: "a", SetByUser: s})
			},
			args:     []string{"test", "-f=hello"},
			expected: true,
		},

		// Bool
		{
			desc: "Bool Opt, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Bool(BoolOpt{Name: "f", Value: true, SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Bool Opt, not set by user, env value",
			config: func(c *Cli, s *bool) {
				os.Setenv("MOW_VALUE", "true")
				c.Bool(BoolOpt{Name: "f", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Bool Opt, set by user",
			config: func(c *Cli, s *bool) {
				c.Bool(BoolOpt{Name: "f", SetByUser: s})
			},
			args:     []string{"test", "-f"},
			expected: true,
		},

		// Int
		{
			desc: "Int Opt, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Int(IntOpt{Name: "f", Value: 42, SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Int Opt, not set by user, env value",
			config: func(c *Cli, s *bool) {
				os.Setenv("MOW_VALUE", "33")
				c.Int(IntOpt{Name: "f", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Int Opt, set by user",
			config: func(c *Cli, s *bool) {
				c.Int(IntOpt{Name: "f", SetByUser: s})
			},
			args:     []string{"test", "-f=666"},
			expected: true,
		},

		// Ints
		{
			desc: "Ints Opt, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Ints(IntsOpt{Name: "f", Value: []int{42}, SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Ints Opt, not set by user, env value",
			config: func(c *Cli, s *bool) {
				os.Setenv("MOW_VALUE", "11,22,33")
				c.Ints(IntsOpt{Name: "f", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Ints Opt, set by user",
			config: func(c *Cli, s *bool) {
				c.Ints(IntsOpt{Name: "f", SetByUser: s})
			},
			args:     []string{"test", "-f=666"},
			expected: true,
		},

		// Strings
		{
			desc: "Strings Opt, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Strings(StringsOpt{Name: "f", Value: []string{"aaa"}, SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Strings Opt, not set by user, env value",
			config: func(c *Cli, s *bool) {
				os.Setenv("MOW_VALUE", "a,b,c")
				c.Strings(StringsOpt{Name: "f", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Strings Opt, set by user",
			config: func(c *Cli, s *bool) {
				c.Strings(StringsOpt{Name: "f", SetByUser: s})
			},
			args:     []string{"test", "-f=ccc"},
			expected: true,
		},
	}

	for _, cas := range cases {
		t.Run(cas.desc, func(t *testing.T) {
			t.Log(cas.desc)

			setByUser := false
			app := App("test", "")

			cas.config(app, &setByUser)

			called := false
			app.Action = func() {
				called = true
			}

			require.NoError(t,
				app.Run(cas.args))

			require.True(t, called, "action should have been called")
			require.Equal(t, cas.expected, setByUser)
		})
	}
}

func TestArgSetByUser(t *testing.T) {
	cases := []struct {
		desc     string
		config   func(*Cli, *bool)
		args     []string
		expected bool
	}{
		// ARGS
		// String
		{
			desc: "String Arg, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				c.String(StringArg{Name: "ARG", Value: "a", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "String Arg, not set by user, env value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "value")
				c.String(StringArg{Name: "ARG", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "String Arg, set by user",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				c.String(StringArg{Name: "ARG", Value: "a", SetByUser: s})
			},
			args:     []string{"test", "aaa"},
			expected: true,
		},

		// Bool
		{
			desc: "Bool Arg, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				c.Bool(BoolArg{Name: "ARG", Value: true, SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Bool Arg, not set by user, env value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "true")
				c.Bool(BoolArg{Name: "ARG", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Bool Arg, set by user",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				c.Bool(BoolArg{Name: "ARG", SetByUser: s})
			},
			args:     []string{"test", "true"},
			expected: true,
		},

		// Int
		{
			desc: "Int Arg, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				c.Int(IntArg{Name: "ARG", Value: 42, SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Int Arg, not set by user, env value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "33")
				c.Int(IntArg{Name: "ARG", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Int Arg, set by user",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG]"
				c.Int(IntArg{Name: "ARG", SetByUser: s})
			},
			args:     []string{"test", "666"},
			expected: true,
		},

		// Ints
		{
			desc: "Ints Arg, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG...]"
				c.Ints(IntsArg{Name: "ARG", Value: []int{42}, SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Ints Arg, not set by user, env value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG...]"
				os.Setenv("MOW_VALUE", "11,22,33")
				c.Ints(IntsArg{Name: "ARG", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Ints Arg, set by user",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG...]"
				c.Ints(IntsArg{Name: "ARG", SetByUser: s})
			},
			args:     []string{"test", "333", "666"},
			expected: true,
		},

		// Strings
		{
			desc: "Strings Arg, not set by user, default value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG...]"
				c.Strings(StringsArg{Name: "ARG", Value: []string{"aaa"}, SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Strings Arg, not set by user, env value",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG...]"
				os.Setenv("MOW_VALUE", "a,b,c")
				c.Strings(StringsArg{Name: "ARG", EnvVar: "MOW_VALUE", SetByUser: s})
			},
			args:     []string{"test"},
			expected: false,
		},
		{
			desc: "Strings Arg, set by user",
			config: func(c *Cli, s *bool) {
				c.Spec = "[ARG...]"
				c.Strings(StringsArg{Name: "ARG", SetByUser: s})
			},
			args:     []string{"test", "aaa", "ccc"},
			expected: true,
		},
	}

	for _, cas := range cases {
		t.Run(cas.desc, func(t *testing.T) {
			t.Log(cas.desc)

			setByUser := false
			app := App("test", "")

			cas.config(app, &setByUser)

			called := false
			app.Action = func() {
				called = true
			}

			require.NoError(t,
				app.Run(cas.args))

			require.True(t, called, "action should have been called")
			require.Equal(t, cas.expected, setByUser)
		})
	}

}

func TestOptSetByEnv(t *testing.T) {
	cases := []struct {
		desc     string
		config   func(*Cli) interface{}
		args     []string
		expected interface{}
	}{
		// OPTS
		// String
		{
			desc: "String Opt, empty env var",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "")
				return c.String(StringOpt{Name: "f", Value: "default", EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: "default",
		},
		{
			desc: "String Opt, env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "env")
				return c.String(StringOpt{Name: "f", Value: "default", EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: "env",
		},
		{
			desc: "String Opt, env set, set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "env")
				return c.String(StringOpt{Name: "f", Value: "default", EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "-f=user"},
			expected: "user",
		},

		// Bool
		{
			desc: "Bool Opt, empty env var",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "")
				return c.Bool(BoolOpt{Name: "f", Value: true, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: true,
		},
		{
			desc: "Bool Opt, env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "true")
				return c.Bool(BoolOpt{Name: "f", Value: false, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: true,
		},
		{
			desc: "Bool Opt, multi env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE1", "xxx")
				os.Setenv("MOW_VALUE2", "true")
				return c.Bool(BoolOpt{Name: "f", Value: false, EnvVar: "MOW_VALUE1 MOW_VALUE2"})
			},
			args:     []string{"test"},
			expected: true,
		},
		{
			desc: "Bool Opt, env set bad value, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "xxx")
				return c.Bool(BoolOpt{Name: "f", Value: true, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: true,
		},
		{
			desc: "Bool Opt, env set, set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "false")
				return c.Bool(BoolOpt{Name: "f", Value: false, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "-f"},
			expected: true,
		},

		// Int
		{
			desc: "Int Opt, empty env var",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "")
				return c.Int(IntOpt{Name: "f", Value: 42, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: 42,
		},
		{
			desc: "Int Opt, env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "42")
				return c.Int(IntOpt{Name: "f", Value: 17, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: 42,
		},
		{
			desc: "Int Opt, multi env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE1", "xxx")
				os.Setenv("MOW_VALUE2", "42")
				return c.Int(IntOpt{Name: "f", Value: 17, EnvVar: "MOW_VALUE1 MOW_VALUE2"})
			},
			args:     []string{"test"},
			expected: 42,
		},
		{
			desc: "Int Opt, env set bad value, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "xxx")
				return c.Int(IntOpt{Name: "f", Value: 42, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: 42,
		},
		{
			desc: "Int Opt, env set, set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "42")
				return c.Int(IntOpt{Name: "f", Value: 17, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "-f=72"},
			expected: 72,
		},

		// Strings
		{
			desc: "Strings Opt, empty env var",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "")
				return c.Strings(StringsOpt{Name: "f", Value: []string{"oh", "ai"}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []string{"oh", "ai"},
		},
		{
			desc: "Strings Opt, env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "do, re, mi")
				return c.Strings(StringsOpt{Name: "f", Value: []string{"oh", "ai"}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []string{"do", "re", "mi"},
		},
		{
			desc: "Strings Opt, multi env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE1", "")
				os.Setenv("MOW_VALUE2", "do, re, mi")
				return c.Strings(StringsOpt{Name: "f", Value: []string{"oh", "ai"}, EnvVar: "MOW_VALUE1 MOW_VALUE2"})
			},
			args:     []string{"test"},
			expected: []string{"do", "re", "mi"},
		},
		{
			desc: "Strings Opt, env set, set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "do, re")
				c.Spec = "-f..."
				return c.Strings(StringsOpt{Name: "f", Value: []string{"oh", "ai"}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "-f=mi", "-f=fa"},
			expected: []string{"mi", "fa"},
		},

		// Ints
		{
			desc: "Ints Opt, empty env var",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "")
				return c.Ints(IntsOpt{Name: "f", Value: []int{1, 2}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []int{1, 2},
		},
		{
			desc: "Ints Opt, env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "11, 13, 17")
				return c.Ints(IntsOpt{Name: "f", Value: []int{1, 2}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []int{11, 13, 17},
		},
		{
			desc: "Ints Opt, multi env set, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE1", "10, 20, xxx")
				os.Setenv("MOW_VALUE2", "11, 13, 17")
				return c.Ints(IntsOpt{Name: "f", Value: []int{1, 2}, EnvVar: "MOW_VALUE1 MOW_VALUE2"})
			},
			args:     []string{"test"},
			expected: []int{11, 13, 17},
		},
		{
			desc: "Ints Opt, env set bad value, not set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "xxx")
				return c.Ints(IntsOpt{Name: "f", Value: []int{1, 2}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []int(nil), //TODO: this is bad and you should feel bad
		},
		{
			desc: "Ints Opt, env set, set by user",
			config: func(c *Cli) interface{} {
				os.Setenv("MOW_VALUE", "42")
				c.Spec = "-f..."
				return c.Ints(IntsOpt{Name: "f", Value: []int{1, 2}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "-f=5", "-f=7"},
			expected: []int{5, 7},
		},
	}

	for _, cas := range cases {
		t.Run(cas.desc, func(t *testing.T) {
			t.Log(cas.desc)

			app := App("test", "")

			pointer := cas.config(app)

			called := false
			app.Action = func() {
				called = true
			}

			require.NoError(t,
				app.Run(cas.args))

			typ := reflect.TypeOf(pointer)
			if typ.Kind() != reflect.Ptr {
				t.Fatalf("config func did not return a pointer")
			}
			actualValue := reflect.ValueOf(pointer).Elem().Interface()

			require.True(t, called, "action should have been called")
			require.Equal(t, cas.expected, actualValue)
		})
	}
}

func TestArgSetByEnv(t *testing.T) {
	cases := []struct {
		desc     string
		config   func(*Cli) interface{}
		args     []string
		expected interface{}
	}{
		// ARGs
		// String
		{
			desc: "String Arg, empty env var",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "")
				return c.String(StringArg{Name: "ARG", Value: "default", EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: "default",
		},
		{
			desc: "String Arg, env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "env")
				return c.String(StringArg{Name: "ARG", Value: "default", EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: "env",
		},
		{
			desc: "String Arg, env set, set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "env")
				return c.String(StringArg{Name: "ARG", Value: "default", EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "user"},
			expected: "user",
		},

		// Bool
		{
			desc: "Bool Arg, empty env var",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "")
				return c.Bool(BoolArg{Name: "ARG", Value: true, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: true,
		},
		{
			desc: "Bool Arg, env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "true")
				return c.Bool(BoolArg{Name: "ARG", Value: false, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: true,
		},
		{
			desc: "Bool Arg, multi env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE1", "xxx")
				os.Setenv("MOW_VALUE2", "true")
				return c.Bool(BoolArg{Name: "ARG", Value: false, EnvVar: "MOW_VALUE1 MOW_VALUE2"})
			},
			args:     []string{"test"},
			expected: true,
		},
		{
			desc: "Bool Arg, env set bad value, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "xxx")
				return c.Bool(BoolArg{Name: "ARG", Value: true, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: true,
		},
		{
			desc: "Bool Arg, env set, set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "false")
				return c.Bool(BoolArg{Name: "ARG", Value: false, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "TRUE"},
			expected: true,
		},

		// Int
		{
			desc: "Int Arg, empty env var",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "")
				return c.Int(IntArg{Name: "ARG", Value: 42, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: 42,
		},
		{
			desc: "Int Arg, env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "42")
				return c.Int(IntArg{Name: "ARG", Value: 17, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: 42,
		},
		{
			desc: "Int Arg, multi env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE1", "xxx")
				os.Setenv("MOW_VALUE2", "42")
				return c.Int(IntArg{Name: "ARG", Value: 17, EnvVar: "MOW_VALUE1 MOW_VALUE2"})
			},
			args:     []string{"test"},
			expected: 42,
		},
		{
			desc: "Int Arg, env set bad value, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "xxx")
				return c.Int(IntArg{Name: "ARG", Value: 42, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: 42,
		},
		{
			desc: "Int Arg, env set, set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "42")
				return c.Int(IntArg{Name: "ARG", Value: 17, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "72"},
			expected: 72,
		},

		// Strings
		{
			desc: "Strings Arg, empty env var",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "")
				return c.Strings(StringsArg{Name: "ARG", Value: []string{"oh", "ai"}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []string{"oh", "ai"},
		},
		{
			desc: "Strings Arg, env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "do, re, mi")
				return c.Strings(StringsArg{Name: "ARG", Value: []string{"oh", "ai"}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []string{"do", "re", "mi"},
		},
		{
			desc: "Strings Arg, multi env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE1", "")
				os.Setenv("MOW_VALUE2", "do, re, mi")
				return c.Strings(StringsArg{Name: "ARG", Value: []string{"oh", "ai"}, EnvVar: "MOW_VALUE1 MOW_VALUE2"})
			},
			args:     []string{"test"},
			expected: []string{"do", "re", "mi"},
		},
		{
			desc: "Strings Arg, env set, set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "ARG..."
				os.Setenv("MOW_VALUE", "do, re")
				return c.Strings(StringsArg{Name: "ARG", Value: []string{"oh", "ai"}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "mi", "fa"},
			expected: []string{"mi", "fa"},
		},

		// Ints
		{
			desc: "Ints Arg, empty env var",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "")
				return c.Ints(IntsArg{Name: "ARG", Value: []int{1, 2}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []int{1, 2},
		},
		{
			desc: "Ints Arg, env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "11, 13, 17")
				return c.Ints(IntsArg{Name: "ARG", Value: []int{1, 2}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []int{11, 13, 17},
		},
		{
			desc: "Ints Arg, multi env set, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE1", "10, 20, xxx")
				os.Setenv("MOW_VALUE2", "11, 13, 17")
				return c.Ints(IntsArg{Name: "ARG", Value: []int{1, 2}, EnvVar: "MOW_VALUE1 MOW_VALUE2"})
			},
			args:     []string{"test"},
			expected: []int{11, 13, 17},
		},
		{
			desc: "Ints Arg, env set bad value, not set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "[ARG]"
				os.Setenv("MOW_VALUE", "xxx")
				return c.Ints(IntsArg{Name: "ARG", Value: []int{1, 2}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test"},
			expected: []int(nil), //TODO: this is bad and you should feel bad
		},
		{
			desc: "Ints Arg, env set, set by user",
			config: func(c *Cli) interface{} {
				c.Spec = "ARG..."
				os.Setenv("MOW_VALUE", "42")
				return c.Ints(IntsArg{Name: "ARG", Value: []int{1, 2}, EnvVar: "MOW_VALUE"})
			},
			args:     []string{"test", "5", "7"},
			expected: []int{5, 7},
		},
	}

	for _, cas := range cases {
		t.Run(cas.desc, func(t *testing.T) {
			t.Log(cas.desc)

			app := App("test", "")
			app.ErrorHandling = flag.ContinueOnError

			pointer := cas.config(app)

			called := false
			app.Action = func() {
				called = true
			}

			require.NoError(t,
				app.Run(cas.args))

			typ := reflect.TypeOf(pointer)
			if typ.Kind() != reflect.Ptr {
				t.Fatalf("config func did not return a pointer")
			}
			actualValue := reflect.ValueOf(pointer).Elem().Interface()

			require.True(t, called, "action should have been called")
			require.Equal(t, cas.expected, actualValue)
		})
	}
}

func TestCommandAction(t *testing.T) {

	called := false

	app := App("app", "")

	app.Command("a", "", ActionCommand(func() { called = true }))

	require.NoError(t,
		app.Run([]string{"app", "a"}))

	require.True(t, called, "commandAction should be called")

}

func TestCommandAliases(t *testing.T) {
	defer suppressOutput()()

	cases := []struct {
		args          []string
		errorExpected bool
	}{
		{
			args:          []string{"say", "hello"},
			errorExpected: false,
		},
		{
			args:          []string{"say", "hi"},
			errorExpected: false,
		},
		{
			args:          []string{"say", "hello hi"},
			errorExpected: true,
		},
		{
			args:          []string{"say", "hello", "hi"},
			errorExpected: true,
		},
	}

	for _, cas := range cases {
		t.Run(fmt.Sprintf("%+v", cas.args), func(t *testing.T) {
			app := App("say", "")
			app.ErrorHandling = flag.ContinueOnError

			called := false

			app.Command("hello hi", "", func(cmd *Cmd) {
				cmd.Action = func() {
					called = true
				}
			})

			err := app.Run(cas.args)

			if cas.errorExpected {
				require.Error(t, err, "Run() should have returned with an error")
				require.False(t, called, "action should not have been called")
			} else {
				require.NoError(t, err, "Run() should have returned without an error")
				require.True(t, called, "action should have been called")
			}
		})
	}
}

func TestSubcommandAliases(t *testing.T) {
	cases := []struct {
		args []string
	}{
		{
			args: []string{"app", "foo", "bar", "baz"},
		},
		{
			args: []string{"app", "foo", "bar", "z"},
		},
		{
			args: []string{"app", "foo", "b", "baz"},
		},
		{
			args: []string{"app", "f", "bar", "baz"},
		},
		{
			args: []string{"app", "f", "b", "baz"},
		},
		{
			args: []string{"app", "f", "b", "z"},
		},
		{
			args: []string{"app", "foo", "b", "z"},
		},
		{
			args: []string{"app", "f", "bar", "z"},
		},
	}

	for _, cas := range cases {
		t.Run(fmt.Sprintf("%+v", cas.args), func(t *testing.T) {
			app := App("app", "")
			app.ErrorHandling = flag.ContinueOnError

			called := false

			app.Command("foo f", "", func(cmd *Cmd) {
				cmd.Command("bar b", "", func(cmd *Cmd) {
					cmd.Command("baz z", "", func(cmd *Cmd) {
						cmd.Action = func() {
							called = true
						}
					})
				})
			})

			err := app.Run(cas.args)

			require.NoError(t, err, "Run() should have returned without an error")
			require.True(t, called, "action should have been called")
		})
	}
}

func TestBeforeAndAfterFlowOrder(t *testing.T) {
	counter := 0

	app := App("app", "")

	app.Before = callChecker(t, 0, &counter)
	app.Command("c", "", func(c *Cmd) {
		c.Before = callChecker(t, 1, &counter)
		c.Command("cc", "", func(cc *Cmd) {
			cc.Before = callChecker(t, 2, &counter)
			cc.Action = callChecker(t, 3, &counter)
			cc.After = callChecker(t, 4, &counter)
		})
		c.After = callChecker(t, 5, &counter)
	})
	app.After = callChecker(t, 6, &counter)

	require.NoError(t,
		app.Run([]string{"app", "c", "cc"}))
	require.Equal(t, 7, counter)
}

func TestBeforeAndAfterFlowOrderWhenOneBeforePanics(t *testing.T) {
	defer func() {
		r := recover()
		require.Equal(t, 42, r)
	}()

	counter := 0

	app := App("app", "")

	app.Before = callChecker(t, 0, &counter)
	app.Command("c", "", func(c *Cmd) {
		c.Before = callChecker(t, 1, &counter)
		c.Command("cc", "", func(cc *Cmd) {
			cc.Before = callCheckerAndPanic(t, 42, 2, &counter)
			cc.Action = func() {
				t.Fatalf("should not have been called")
			}
			cc.After = func() {
				t.Fatalf("should not have been called")
			}
		})
		c.After = callChecker(t, 3, &counter)
	})
	app.After = callChecker(t, 4, &counter)

	require.NoError(t,
		app.Run([]string{"app", "c", "cc"}))
	require.Equal(t, 5, counter)
}

func TestBeforeAndAfterFlowOrderWhenOneAfterPanics(t *testing.T) {
	defer func() {
		e := recover()
		require.Equal(t, 42, e)
	}()

	counter := 0

	app := App("app", "")

	app.Before = callChecker(t, 0, &counter)
	app.Command("c", "", func(c *Cmd) {
		c.Before = callChecker(t, 1, &counter)
		c.Command("cc", "", func(cc *Cmd) {
			cc.Before = callChecker(t, 2, &counter)
			cc.Action = callChecker(t, 3, &counter)
			cc.After = callCheckerAndPanic(t, 42, 4, &counter)
		})
		c.After = callChecker(t, 5, &counter)
	})
	app.After = callChecker(t, 6, &counter)

	require.NoError(t,
		app.Run([]string{"app", "c", "cc"}))
	require.Equal(t, 7, counter)
}

func TestBeforeAndAfterFlowOrderWhenMultipleAftersPanic(t *testing.T) {
	defer func() {
		e := recover()
		require.Equal(t, 666, e)
	}()

	counter := 0

	app := App("app", "")

	app.Before = callChecker(t, 0, &counter)
	app.Command("c", "", func(c *Cmd) {
		c.Before = callChecker(t, 1, &counter)
		c.Command("cc", "", func(cc *Cmd) {
			cc.Before = callChecker(t, 2, &counter)
			cc.Action = callChecker(t, 3, &counter)
			cc.After = callCheckerAndPanic(t, 42, 4, &counter)
		})
		c.After = callChecker(t, 5, &counter)
	})
	app.After = callCheckerAndPanic(t, 666, 6, &counter)

	require.NoError(t,
		app.Run([]string{"app", "c", "cc"}))
	require.Equal(t, 7, counter)
}

func exitShouldBeCalledWith(t *testing.T, wantedExitCode int, called *bool) func() {
	oldExiter := exiter
	exiter = func(code int) {
		require.Equal(t, wantedExitCode, code, "unwanted exit code")
		*called = true
	}
	return func() { exiter = oldExiter }
}

func exitShouldNotCalled(t *testing.T) func() {
	oldExiter := exiter
	exiter = func(code int) {
		t.Errorf("exit should not have been called")
	}
	return func() { exiter = oldExiter }
}

func suppressOutput() func() {
	return captureAndRestoreOutput(nil, nil)
}

func setAndRestoreEnv(env map[string]string) func() {
	backup := map[string]string{}
	for k, v := range env {
		backup[k] = os.Getenv(k)
		os.Setenv(k, v)
	}

	return func() {
		for k, v := range backup {
			os.Setenv(k, v)
		}
	}
}

func captureAndRestoreOutput(out, err *string) func() {
	oldStdOut := stdOut
	oldStdErr := stdErr

	if out == nil {
		stdOut = ioutil.Discard
	} else {
		stdOut = trapWriter(out)
	}
	if err == nil {
		stdErr = ioutil.Discard
	} else {
		stdErr = trapWriter(err)
	}

	return func() {
		stdOut = oldStdOut
		stdErr = oldStdErr
	}
}

func trapWriter(writeTo *string) *writerTrap {
	return &writerTrap{
		buffer:  bytes.NewBuffer(nil),
		writeTo: writeTo,
	}
}

type writerTrap struct {
	buffer  *bytes.Buffer
	writeTo *string
}

func (w *writerTrap) Write(p []byte) (n int, err error) {
	n, err = w.buffer.Write(p)
	if err == nil {
		*(w.writeTo) = w.buffer.String()
	}
	return
}

func callChecker(t *testing.T, wanted int, counter *int) func() {
	return func() {
		t.Logf("checker: wanted: %d, got %d", wanted, *counter)
		require.Equal(t, wanted, *counter)
		*counter++
	}
}

func callCheckerAndPanic(t *testing.T, panicValue interface{}, wanted int, counter *int) func() {
	return func() {
		t.Logf("checker: wanted: %d, got %d", wanted, *counter)
		require.Equal(t, wanted, *counter)
		*counter++
		panic(panicValue)
	}
}
