package sync_test

import (
	"math"
	"sync"
	"testing"

	xsync "go.mway.dev/x/sync"
)

func BenchmarkWaitGroupAdd(b *testing.B) {
	b.Run("sync", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			wg := &sync.WaitGroup{}
			for pb.Next() {
				wg.Add(1)
			}
		})
	})

	b.Run("xsync", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			wg := &xsync.WaitGroup{}
			for pb.Next() {
				wg.Add(1)
			}
		})
	})
}

func BenchmarkWaitGroupDone(b *testing.B) {
	var wg xsync.WaitGroup
	wg.Add(math.MaxInt32)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Done()
	}
}
