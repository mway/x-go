package require_test

import (
	"errors"
	"fmt"
	"testing"

	"go.mway.dev/x/testing/require"
)

// TODO(mway): use an interface for testing.T and mock to test negative cases

var (
	errA = errors.New("lower")
	errB = fmt.Errorf("middle: %w", errA)
	errC = fmt.Errorf("upper: %w", errB)
)

func TestEqualErrorChains(t *testing.T) {
	cases := []struct {
		expect error
		actual error
	}{
		{
			expect: errC,
			actual: errC,
		},
	}

	for _, tt := range cases {
		t.Run(tt.expect.Error(), func(t *testing.T) {
			require.EqualErrorChains(t, tt.expect, tt.actual)
		})
	}
}

func TestContainsErrorChain(t *testing.T) {
	cases := []struct {
		expect error
		actual error
	}{
		{
			expect: errC,
			actual: wrap(wrap(errC, "extra"), "layers"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.expect.Error(), func(t *testing.T) {
			require.ContainsErrorChain(t, tt.expect, tt.actual)
		})
	}
}

func wrap(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	}

	return fmt.Errorf("%s: %w", msg, err)
}
