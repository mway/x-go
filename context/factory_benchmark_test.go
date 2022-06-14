package context_test

import (
	"testing"
	"time"

	"go.mway.dev/chrono/clock"
	"go.mway.dev/x/context"
)

func BenchmarkDebouncedFactory_Get(b *testing.B) {
	bench := func(b *testing.B, factory *context.DebouncedFactory) {
		b.ReportAllocs()
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				factory.Get()
			}
		})
	}

	b.Run("fake clock", func(b *testing.B) {
		factory, _ := newFactory(b, time.Second)
		bench(b, factory)
	})

	b.Run("real clock", func(b *testing.B) {
		factory, _ := newFactory(
			b,
			time.Second,
			context.WithClock(clock.NewMonotonicClock()),
			context.WithDebounce(time.Second),
		)
		bench(b, factory)
	})
}
