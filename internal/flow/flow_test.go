package flow

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestStepCallsDo(t *testing.T) {
	called := false
	step := &Step{
		Do: func() {
			called = true
		},
	}

	step.Run(nil)

	require.True(t, called, "Step's do wasn't called")
}

func TestStepCallsSuccessAfterDo(t *testing.T) {
	calls := 0
	step := &Step{
		Do: func() {
			require.Equal(t, 0, calls, "Do should be called first")
			calls++
		},
		Success: &Step{
			Do: func() {
				require.Equal(t, 1, calls, "Success should be called second")
				calls++
			},
		},
		Error: &Step{
			Do: func() {
				t.Fatalf("Error should not have been called")
			},
		},
	}

	step.Run(nil)

	require.Equal(t, 2, calls, "Both do and success should be called")
}

func TestStepCallsErrorIfDoPanics(t *testing.T) {
	defer func() { recover() }()
	calls := 0
	step := &Step{
		Do: func() {
			require.Equal(t, 0, calls, "Do should be called first")
			calls++
			panic(42)
		},
		Success: &Step{
			Do: func() {
				t.Fatalf("Success should not have been called")
			},
		},
		Error: &Step{
			Do: func() {
				require.Equal(t, 1, calls, "Error should be called second")
				calls++
			},
		},
	}

	step.Run(nil)

	require.Equal(t, 2, calls, "Both do and error should be called")
}

func TestStepCallsOsExitIfAskedTo(t *testing.T) {
	exitCode := -1
	step := &Step{Exiter: func(x int) {
		exitCode = x
	}}

	step.Run(ExitCode(42))

	require.Equal(t, exitCode, 42, "should have called exit with 42")
}

func TestStepRethrowsPanic(t *testing.T) {
	defer func() {
		require.Equal(t, 42, recover(), "should panicked with the same value")
	}()

	step := &Step{}

	step.Run(42)

	t.Fatalf("Should have panicked")
}

func TestStepShouldNopIfNoSuccessNorPanic(t *testing.T) {

	step := &Step{Exiter: func(x int) {
		t.Fatalf("Should not have called exit")
	}}

	step.Run(nil)
}
