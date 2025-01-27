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

// Package stack provides stack-based types and helpers.
package stack

// A Stack is a LIFO queue that holds values of type T.
type Stack[T any] struct {
	data []T
}

// NewStack creates a new [Stack[T]] with the given initial capacity.
func NewStack[T any](size int) *Stack[T] {
	return &Stack[T]{
		data: make([]T, 0, size),
	}
}

// Push pushes x on top of the stack.
func (s *Stack[T]) Push(x T) {
	s.data = append(s.data, x)
}

// Top returns the value on top of the stack.
func (s *Stack[T]) Top() T {
	return s.data[len(s.data)-1]
}

// Pop pops the top value off of the stack and returns it.
func (s *Stack[T]) Pop() T {
	x := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return x
}

// Len returns the number of values held by the stack.
func (s *Stack[T]) Len() int {
	return len(s.data)
}

// Cap returns the stack's current capacity.
func (s *Stack[T]) Cap() int {
	return cap(s.data)
}

// MaybeTop returns the value on top of the stack if there is one. The boolean
// return indicates whether the T value is valid.
func (s *Stack[T]) MaybeTop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.Top(), true
}

// MaybePop pops the top value off of the stack and returns it, if there is
// one. The boolean return indicates whether the T is valid.
func (s *Stack[T]) MaybePop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.Pop(), true
}

// PeekEach calls fn for each value on the stack. If fn returns false,
// iteration will stop and the function will return immediately. If there are
// no values on the stack, fn will not be called.
func (s *Stack[T]) PeekEach(fn func(T) bool) {
	if len(s.data) == 0 {
		return
	}
	for i := len(s.data) - 1; i >= 0; i-- {
		if !fn(s.data[i]) {
			return
		}
	}
}

// PopEach pops the top value off the stack and passes it to fn while the stack
// is not empty. If fn returns false, no more values will be popped and the
// function will return immediately. If there are no values on the stack, fn
// will not be called.
func (s *Stack[T]) PopEach(fn func(T) bool) {
	for {
		if x, ok := s.MaybePop(); !ok || !fn(x) {
			break
		}
	}
}
