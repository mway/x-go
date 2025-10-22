// Copyright (c) 2024 Matt Way
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

// Package ptr provides helpers for working with pointers.
package ptr

// To returns a pointer to the value x.
func To[T any](x T) *T {
	return &x
}

// Load returns the dereferenced value of x. If x is nil, the zero value of T
// will be returned.
func Load[T any](x *T) (value T) {
	return LoadOr(x, value)
}

// LoadOr returns the dereferenced value of x. If x is nil, fallback will be
// returned.
func LoadOr[T any](x *T, fallback T) T {
	if x == nil {
		return fallback
	}
	return *x
}

// LoadOrElse returns the dereferenced value of x. If x is nil, the value
// returned by calling fn will be used.
func LoadOrElse[T any](x *T, fn func() T) T {
	if x == nil {
		return fn()
	}
	return *x
}

// A Pointer is a thin wrapper around a pointer to T.
type Pointer[T any] struct {
	ptr *T
}

// New creates a new [Pointer[T, P]], allocating and storing the given value.
func New[T any](value T) (p Pointer[T]) {
	p.Store(value)
	return p
}

// NewPtr creates a new [Pointer[T, P]], storing the given ptr.
func NewPtr[T any](ptr *T) (p Pointer[T]) {
	p.StorePtr(ptr)
	return p
}

// NewPtrCopy creates a new [Pointer[T, P]], storing a shallow copy of the
// given ptr.
func NewPtrCopy[T any](ptr *T) (p Pointer[T]) {
	p.StorePtrCopy(ptr)
	return p
}

// Call is a convenience method that calls fn with the (p.Load(), p.Held()).
func (p Pointer[T]) Call(fn func(T, bool)) {
	fn(p.Load(), p.Held())
}

// CallPtr is a convenience method that calls fn with (p.Raw(), p.Held()).
func (p Pointer[T]) CallPtr(fn func(*T, bool)) {
	fn(p.Raw(), p.Held())
}

// Clear resets p. Afterwords, p.Held() will return false.
func (p *Pointer[T]) Clear() {
	if p.ptr != nil {
		p.ptr = nil
	}
}

// Held indicates whether p holds a value. It is sugar for p.ptr != nil.
func (p Pointer[T]) Held() bool {
	return p.ptr != nil
}

// Load returns the dereferenced value of p. If p does not hold a value, the
// zero value of T will be returned instead.
func (p Pointer[T]) Load() T {
	return Load(p.ptr)
}

// LoadOr returns the dereferenced value of p. If p does not hold a value,
// fallback will be returned instead.
func (p Pointer[T]) LoadOr(fallback T) T {
	return LoadOr(p.ptr, fallback)
}

// LoadOrElse returns the dereferenced value of p. If p does not hold a value,
// the value returned by calling fn will be used.
func (p Pointer[T]) LoadOrElse(fn func() T) T {
	return LoadOrElse(p.ptr, fn)
}

// MaybeCall calls fn with the result of p.Load() if p currently holds a value.
func (p *Pointer[T]) MaybeCall(fn func(T)) bool {
	if p.Held() {
		fn(p.Load())
		return true
	}
	return false
}

// MaybeCallPtr calls fn with the result of p.Raw() if p currently holds a
// value.
func (p *Pointer[T]) MaybeCallPtr(fn func(*T)) bool {
	if p.Held() {
		fn(p.Raw())
		return true
	}
	return false
}

// MaybeStore calls p.Store(value) if p does not currently hold a value.
func (p *Pointer[T]) MaybeStore(value T) bool {
	if !p.Held() {
		p.Store(value)
		return true
	}
	return false
}

// MaybeStorePtr calls p.StorePtr(ptr) if p does not currently hold a value.
func (p *Pointer[T]) MaybeStorePtr(ptr *T) bool {
	if !p.Held() {
		p.StorePtr(ptr)
		return true
	}
	return false
}

// MaybeStorePtrCopy calls p.StorePtrCopy(ptr) if p does not currently hold a
// value.
func (p *Pointer[T]) MaybeStorePtrCopy(ptr *T) bool {
	if !p.Held() {
		p.StorePtrCopy(ptr)
		return true
	}
	return false
}

// Move moves p to dst. p.Held() will return false afterwards.
func (p *Pointer[T]) Move(dst *Pointer[T]) bool {
	if p == nil || !p.Held() || dst == nil || p.ptr == dst.ptr {
		return false
	}

	dst.ptr = p.ptr
	p.Clear()
	return true
}

// Raw returns p's underlying pointer.
func (p Pointer[T]) Raw() *T {
	return p.ptr
}

// Store stores the given value in p, allocating if necessary.
func (p *Pointer[T]) Store(value T) {
	if p.ptr == nil {
		p.ptr = &value
	} else {
		*p.ptr = value
	}
}

// StorePtr stores the given ptr in p.
func (p *Pointer[T]) StorePtr(ptr *T) {
	p.ptr = ptr
}

// StorePtrCopy stores a shallow copy of ptr in p.
func (p *Pointer[T]) StorePtrCopy(ptr *T) {
	tmp := *ptr
	p.StorePtr(&tmp)
}

// swap
