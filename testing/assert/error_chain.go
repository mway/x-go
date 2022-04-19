package assert

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// EqualErrorChains asserts that an actual error chain is the same as an
// expected error chain.
func EqualErrorChains(t *testing.T, expect error, actual error, msgAndArgs ...any) bool {
	for expect != nil && actual != nil {
		if !assert.ErrorIs(t, actual, expect, msgAndArgs...) {
			return false
		}

		expect = errors.Unwrap(expect)
		actual = errors.Unwrap(actual)
	}

	return assert.ErrorIs(t, actual, expect, msgAndArgs...)
}

// ContainsErrorChain asserts that an actual error chain contains all of the
// same underlying errors as the expected error chain, but allows for arbitrary
// chain order.
func ContainsErrorChain(t *testing.T, expect error, actual error, msgAndArgs ...any) bool {
	for expect != nil {
		if !assert.ErrorIs(t, actual, expect, msgAndArgs...) {
			return false
		}

		expect = errors.Unwrap(expect)
	}

	return true
}
