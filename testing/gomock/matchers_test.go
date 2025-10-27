package gomock_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/testing/gomock"
)

func TestMatchFunc(t *testing.T) {
	const have = "abc123"

	matches := gomock.Match(func(x string) bool {
		return x == have
	})
	require.True(t, matches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	notmatches := gomock.Match(func(x string) bool {
		return x != have
	})
	require.False(t, notmatches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())
}

func TestMatchSubstring(t *testing.T) {
	const have = "abc123"

	matches := gomock.MatchSubstring("abc")
	require.True(t, matches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	matches = gomock.MatchSubstring("bc12")
	require.True(t, matches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	matches = gomock.MatchSubstring("abc123")
	require.True(t, matches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	notmatches := gomock.MatchSubstring("xyz")
	require.False(t, notmatches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	notmatches = gomock.MatchSubstring("abc1234")
	require.False(t, notmatches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())
}

func TestMatchPrefix(t *testing.T) {
	const have = "abc123"

	matches := gomock.MatchPrefix("abc")
	require.True(t, matches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	matches = gomock.MatchPrefix("abc123")
	require.True(t, matches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	notmatches := gomock.MatchPrefix("xyz")
	require.False(t, notmatches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	notmatches = gomock.MatchPrefix("bc123")
	require.False(t, notmatches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())
}

func TestMatchSuffix(t *testing.T) {
	const have = "abc123"

	matches := gomock.MatchSuffix("123")
	require.True(t, matches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	matches = gomock.MatchSuffix("abc123")
	require.True(t, matches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	notmatches := gomock.MatchSuffix("xyz")
	require.False(t, notmatches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())

	notmatches = gomock.MatchSuffix("abc12")
	require.False(t, notmatches.Matches(have))
	require.Equal(t, "gomock.MatchFunc[string]", matches.String())
}

func TestMatchSliceContains(t *testing.T) {
	have := []int{1, 2, 3, 4, 5}

	matches := gomock.MatchSliceContains(4, 2)
	require.True(t, matches.Matches(have))
	matches = gomock.MatchSliceContains(2, 4)
	require.True(t, matches.Matches(have))
	matches = gomock.MatchSliceContains(2, 2)
	require.True(t, matches.Matches(have))
	matches = gomock.MatchSliceContains[int]()
	require.True(t, matches.Matches(have))

	notmatches := gomock.MatchSliceContains(5, 6)
	require.False(t, notmatches.Matches(have))
	notmatches = gomock.MatchSliceContains(6, 7)
	require.False(t, notmatches.Matches(have))
	notmatches = gomock.MatchSliceContains[string]()
	require.False(t, notmatches.Matches(have))
}

func TestMatchMapContains(t *testing.T) {
	have := map[int]int{
		1: 101,
		2: 102,
		3: 103,
		4: 104,
		5: 105,
	}

	matches := gomock.MatchMapContains[int, int](nil)
	require.True(t, matches.Matches(have))
	matches = gomock.MatchMapContains(map[int]int{})
	require.True(t, matches.Matches(have))
	matches = gomock.MatchMapContains(map[int]int{
		1: 101,
	})
	require.True(t, matches.Matches(have))
	matches = gomock.MatchMapContains(map[int]int{
		1: 101,
		3: 103,
		5: 105,
	})
	require.True(t, matches.Matches(have))

	notmatches := gomock.MatchMapContains(map[int]int{
		1: 102,
	})
	require.False(t, notmatches.Matches(have))
	notmatches = gomock.MatchMapContains(map[int]int{
		10: 1000,
	})
	require.False(t, notmatches.Matches(have))
	notmatches = gomock.MatchMapContains[string, string](nil)
	require.False(t, notmatches.Matches(have))
}

func TestMatchMapContainsKeys(t *testing.T) {
	have := map[int]struct{}{
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
	}

	matches := gomock.MatchMapContainsKeys[int, struct{}]()
	require.True(t, matches.Matches(have))

	matches = gomock.MatchMapContainsKeys[int, struct{}](1)
	require.True(t, matches.Matches(have))

	matches = gomock.MatchMapContainsKeys[int, struct{}](5, 4, 3, 2, 1)
	require.True(t, matches.Matches(have))

	notmatches := gomock.MatchMapContainsKeys[int, struct{}](5, 6)
	require.False(t, notmatches.Matches(have))

	notmatches = gomock.MatchMapContainsKeys[int, struct{}](6, 7)
	require.False(t, notmatches.Matches(have))

	notmatches = gomock.MatchMapContainsKeys[string, string]()
	require.False(t, notmatches.Matches(have))
}
