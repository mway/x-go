package require

import (
	"testing"

	"go.mway.dev/x/testing/assert"
)

// EqualErrorChains asserts that an actual error chain is the same as an
// expected error chain. If not, it fails the test immediately.
func EqualErrorChains(t *testing.T, expect error, actual error, msgAndArgs ...any) {
	if !assert.EqualErrorChains(t, expect, actual, msgAndArgs...) {
		t.FailNow()
	}
}

// ContainsErrorChain asserts that an actual error chain contains all of the
// same underlying errors as the expected error chain, but allows for arbitrary
// chain order. If not, it fails the test immediately.
func ContainsErrorChain(t *testing.T, expect error, actual error, msgAndArgs ...any) {
	if !assert.ContainsErrorChain(t, expect, actual, msgAndArgs...) {
		t.FailNow()
	}
}
