// Package gomock provides supplemental gomock types and helpers.
package gomock

import (
	"fmt"
	"strings"

	"go.uber.org/mock/gomock"

	"go.mway.dev/x/container/set"
)

// MatchFunc is a [gomock.Matcher] that evaluates itself against a parameter.
type MatchFunc[T any] func(have T) bool

// Matches indicates whether x matches f. Matches returns false if either x is
// not a T, or if f(x) returns false.
func (f MatchFunc[T]) Matches(x any) bool {
	if t, ok := x.(T); ok {
		return f(t)
	}
	return false
}

// String returns the string representation of f.
func (f MatchFunc[T]) String() string {
	return fmt.Sprintf("%T", f)
}

// Match is a convenience alias for MatchFunc[T](...).
func Match[T any](check func(T) bool) gomock.Matcher {
	return MatchFunc[T](check)
}

// MatchSubstring returns a [gomock.Matcher] that asserts that a given
// parameter is a string and contains the given substring.
func MatchSubstring(want string) gomock.Matcher {
	return Match(func(have string) bool {
		return strings.Contains(have, want)
	})
}

// MatchPrefix returns a [gomock.Matcher] that asserts that a given parameter
// is a string and contains the given prefix.
func MatchPrefix(want string) gomock.Matcher {
	return Match(func(have string) bool {
		return strings.HasPrefix(have, want)
	})
}

// MatchSuffix returns a [gomock.Matcher] that asserts that a given parameter
// is a string and contains the given suffix.
func MatchSuffix(want string) gomock.Matcher {
	return Match(func(have string) bool {
		return strings.HasSuffix(have, want)
	})
}

// MatchSliceContains returns a [gomock.Matcher] that asserts that a given
// parameter is []T and contains the given values.
func MatchSliceContains[T comparable](want ...T) gomock.Matcher {
	if len(want) == 0 {
		return Match(func([]T) bool {
			return true
		})
	}

	return Match(func(have []T) bool {
		give := set.New(have...)
		return give.ContainsAll(want...)
	})
}

// MatchMapContains returns a [gomock.Matcher] that asserts that a given
// parameter is map[K]V and contains the given key/value pairs.
func MatchMapContains[K comparable, V comparable](
	want map[K]V,
) gomock.Matcher {
	return Match(func(have map[K]V) bool {
		for wantK, wantV := range want {
			if haveV, ok := have[wantK]; !ok || haveV != wantV {
				return false
			}
		}
		return true
	})
}

// MatchMapContainsKeys returns a [gomock.Matcher] that asserts that a given
// parameter is map[K]V and contains the given keys.
func MatchMapContainsKeys[K comparable, V any](want ...K) gomock.Matcher {
	if len(want) == 0 {
		return Match(func(map[K]V) bool {
			return true
		})
	}

	return Match(func(have map[K]V) bool {
		for _, wantK := range want {
			if _, ok := have[wantK]; !ok {
				return false
			}
		}
		return true
	})
}
