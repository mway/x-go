// Copyright (c) 2025 Matt Way
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE THE SOFTWARE.

package env_test

import (
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"go.mway.dev/x/env"
	"go.mway.dev/x/must"
)

func TestWith(t *testing.T) {
	var (
		wantKey   = fmt.Sprintf("%s_%d", t.Name(), time.Now().UnixNano())
		wantValue = "TEST VALUE: " + t.Name()
	)

	haveValue, haveOK := os.LookupEnv(wantKey)
	require.False(t, haveOK)
	require.Empty(t, haveValue)

	env.With(func() {
		haveValue, haveOK := os.LookupEnv(wantKey)
		require.True(t, haveOK)
		require.Equal(t, wantValue, haveValue)
	}, wantKey, wantValue)
}

func TestWithError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		var (
			wantKey   = fmt.Sprintf("%s_%d", t.Name(), time.Now().UnixNano())
			wantValue = "TEST VALUE: " + t.Name()
		)

		haveValue, haveOK := os.LookupEnv(wantKey)
		require.False(t, haveOK)
		require.Empty(t, haveValue)

		haveErr := env.WithError(func() error {
			haveValue, haveOK := os.LookupEnv(wantKey)
			require.True(t, haveOK)
			require.Equal(t, wantValue, haveValue)
			return nil
		}, wantKey, wantValue)
		require.NoError(t, haveErr)

		wantErr := fmt.Errorf("test error: %s", t.Name())
		haveErr = env.WithError(func() error {
			return wantErr
		}, wantKey, wantValue)
		require.ErrorIs(t, haveErr, wantErr)
	})

	t.Run("with error", func(t *testing.T) {
		wantErr := errors.New(t.Name())
		require.ErrorIs(t, env.WithError(func() error {
			return wantErr
		}, "doesn't", "matter"), wantErr)
	})
}

func TestWithout(t *testing.T) {
	wantKey := fmt.Sprintf("%s_%d", t.Name(), time.Now().UnixNano())

	v, err := env.NewVarWithValue(wantKey, "should be unused")
	require.NoError(t, err)
	defer must.NotError(v.Restore)

	env.Without(func() {
		haveValue, haveOK := os.LookupEnv(wantKey)
		require.False(t, haveOK)
		require.Empty(t, haveValue)
	}, wantKey)
}

func TestWithoutError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		wantKey := fmt.Sprintf("%s_%d", t.Name(), time.Now().UnixNano())

		v, err := env.NewVarWithValue(wantKey, "should be unused")
		require.NoError(t, err)
		defer must.NotError(v.Restore)

		require.NoError(t, env.WithoutError(func() error {
			haveValue, haveOK := os.LookupEnv(wantKey)
			require.False(t, haveOK)
			require.Empty(t, haveValue)
			return nil
		}, wantKey))
	})

	t.Run("with error", func(t *testing.T) {
		wantErr := errors.New(t.Name())
		require.ErrorIs(t, env.WithoutError(func() error {
			return wantErr
		}, "doesn't matter"), wantErr)
	})
}

func TestWithAll(t *testing.T) {
	env.WithAll(func() {
		name := "key1"
		for {
			x := os.Getenv(name)
			if len(x) == 0 {
				break
			}
			name = x
		}
		require.Equal(t, t.Name(), name)
	}, map[string]string{
		"key1": "key2",
		"key2": "key3",
		"key3": "key4",
		"key4": "key5",
		"key5": t.Name(),
	})
}

func TestWithAllError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		require.NoError(t, env.WithAllError(func() error {
			name := "key1"
			for {
				x := os.Getenv(name)
				if len(x) == 0 {
					break
				}
				name = x
			}
			require.Equal(t, t.Name(), name)
			return nil
		}, map[string]string{
			"key1": "key2",
			"key2": "key3",
			"key3": "key4",
			"key4": "key5",
			"key5": t.Name(),
		}))
	})

	t.Run("with error", func(t *testing.T) {
		wantErr := errors.New(t.Name())
		require.ErrorIs(t, env.WithAllError(func() error {
			return wantErr
		}, map[string]string{
			"var1": "value1",
		}), wantErr)
	})
}
