package sync_test

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mway.dev/x/sync"
)

func TestWaitGroupAddInc(t *testing.T) {
	var (
		target = math.MaxInt16
		add    sync.WaitGroup
		inc    sync.WaitGroup
	)

	add.Add(target)

	for i := 0; i < target; i++ {
		inc.Inc()
	}

	require.Equal(t, target, add.Len())
	require.Equal(t, add.Len(), inc.Len())
}

func TestWaitGroupDone(t *testing.T) {
	var (
		target = math.MaxInt16
		wg     sync.WaitGroup
	)

	require.NotPanics(t, func() {
		wg.Add(target)
		for i := 0; i < target; i++ {
			wg.Done()
			require.Equal(t, target-i-1, wg.Len())
		}
		require.Equal(t, 0, wg.Len())
	})
}

func TestWaitGroupWait(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()

		for wg.Len() > 0 {
			<-ticker.C
			wg.Done()
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		require.Equal(t, 0, wg.Len())
	case <-time.After(time.Second):
		require.FailNow(t, "WaitGroup did not unblock")
	}
}

func TestWaitGroupLen(t *testing.T) {
	var wg sync.WaitGroup
	require.Equal(t, 0, wg.Len())

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		require.Equal(t, i+1, wg.Len())
	}

	for wg.Len() > 0 {
		prev := wg.Len()
		wg.Done()
		require.Equal(t, prev-1, wg.Len())
	}

	require.NotPanics(t, func() {
		wg.Add(math.MaxInt32)
	})

	require.Panics(t, func() {
		wg.Add(1)
	})
}
