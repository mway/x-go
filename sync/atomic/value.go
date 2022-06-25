package atomic

import (
	"sync/atomic"
)

// Value is a strongly-typed atomic value. It is otherwise identical to the
// standard library's atomic.Value.
type Value[T any] struct {
	value atomic.Value
}

// NewValue creates a new Value that can store values of type T, initializing
// the Value with the given initial value.
func NewValue[T any](initial T) *Value[T] {
	v := &Value[T]{}
	v.value.Store(initial)
	return v
}

// CompareAndSwap performs an atomic compare and swap using oldval and newval.
// The return value indicates whether a swap took place.
func (v *Value[T]) CompareAndSwap(oldval T, newval T) bool {
	return v.value.CompareAndSwap(oldval, newval)
}

// Load loads the currently held T value.
func (v *Value[T]) Load() T {
	return v.value.Load().(T)
}

// Store stores the given T value.
func (v *Value[T]) Store(val T) {
	v.value.Store(val)
}

// Swap swaps the currently held T value with the given value, returning the
// previous T.
func (v *Value[T]) Swap(val T) T {
	return v.value.Swap(val).(T)
}
