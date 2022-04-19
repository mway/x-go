package require_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.mway.dev/x/testing/internal/requiremock"
	xrequire "go.mway.dev/x/testing/require"
)

var (
	errA = errors.New("lower")
	errB = fmt.Errorf("middle: %w", errA)
	errC = fmt.Errorf("upper: %w", errB)
)

func TestEqualErrorChains(t *testing.T) {
	cases := []struct {
		expect     error
		actual     error
		expectFail bool
	}{
		{
			expect:     errA,
			actual:     errA,
			expectFail: false,
		},
		{
			expect:     errA,
			actual:     errB,
			expectFail: true,
		},
		{
			expect:     errA,
			actual:     errors.New("random"),
			expectFail: true,
		},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("%v vs %v", tt.expect, tt.actual)
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockT := requiremock.NewMockTestingT(ctrl)
			if tt.expectFail {
				mockT.EXPECT().FailNow()
				mockT.EXPECT().
					Errorf(gomock.Any(), gomock.Any()).
					Do(func(msg string, args ...any) {
						require.Len(t, args, 1)

						str, ok := args[0].(string)
						require.True(t, ok)
						require.Contains(t, str, "Target error should be in err chain")
					})
			}

			xrequire.EqualErrorChains(mockT, tt.expect, tt.actual)
		})
	}
}

func TestContainsErrorChain(t *testing.T) {
	cases := []struct {
		expect     error
		actual     error
		expectFail bool
	}{
		{
			expect:     errC,
			actual:     wrap(wrap(errC, "extra"), "layers"),
			expectFail: false,
		},
		{
			expect:     wrap(errC, "extra"),
			actual:     wrap(wrap(errC, "extra"), "layers"),
			expectFail: true,
		},
	}

	for _, tt := range cases {
		name := fmt.Sprintf("%v vs %v", tt.expect, tt.actual)
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockT := requiremock.NewMockTestingT(ctrl)
			if tt.expectFail {
				mockT.EXPECT().FailNow()
				mockT.EXPECT().
					Errorf(gomock.Any(), gomock.Any()).
					Do(func(msg string, args ...any) {
						require.Len(t, args, 1)

						str, ok := args[0].(string)
						require.True(t, ok)
						require.Contains(t, str, "Target error should be in err chain")
					})
			}

			xrequire.ContainsErrorChain(mockT, tt.expect, tt.actual)
		})
	}
}

func wrap(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	}

	return fmt.Errorf("%s: %w", msg, err)
}
