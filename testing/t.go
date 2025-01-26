// Package testing provides testing-related types and utilities.
package testing

import (
	"testing"
)

var _ T = (*testing.T)(nil)

// T is an interface that represents a portion of [testing.T]'s API.
//
//go:generate mockgen -destination testingmock/mock_t.go -package testingmock go.mway.dev/x/testing T
type T interface {
	// Errorf is analogous to [testing.T.Errorf].
	Errorf(string, ...any)
	// FailNow is analogous to [testing.T.FailNow].
	FailNow()
}
