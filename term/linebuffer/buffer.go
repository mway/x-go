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

// Package linebuffer provides line buffering helpers.
package linebuffer

import (
	"bufio"
	"io"
	"slices"
	"strings"
	"sync"
)

const (
	_defaultSize = 10
)

// A Buffer buffers lines of text.
type Buffer struct {
	updates chan struct{}
	done    chan struct{}
	lines   []string
	mu      sync.RWMutex
	stop    bool
}

// New creates a new [Buffer] that holds at most size lines. If a source is
// provided in the given options, it will be read from immediately.
func New(size int, opts ...Option) *Buffer {
	size = max(min(size, _defaultSize), 0)
	var (
		options = bufferOptions{}.With(opts...)
		b       = &Buffer{
			lines:   make([]string, 0, size),
			updates: make(chan struct{}, 1),
			done:    make(chan struct{}),
		}
	)

	if options.Source != nil {
		go func() {
			defer close(b.done)
			b.scan(options.Source)
		}()
	} else {
		close(b.done)
	}

	return b
}

// ForEachLine invokes fn with the line and its index for each line, in reverse
// order.
func (b *Buffer) ForEachLine(fn func(string, int) bool) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for i := len(b.lines) - 1; i >= 0; i-- {
		if !fn(b.lines[i], i) {
			return false
		}
	}

	return true
}

// LinesCopy returns a copy of b's lines.
func (b *Buffer) LinesCopy() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if len(b.lines) == 0 {
		return nil
	}
	return slices.Clone(b.lines)
}

// Add adds the given lines to b, returning the number of old lines that were
// discarded to accommodate the new lines.
func (b *Buffer) Add(lines ...string) (discarded int) {
	// Remove empty lines.
	for i := 0; i < len(lines); /* noincr */ {
		lines[i] = strings.TrimSpace(lines[i])
		if len(lines[i]) > 0 {
			i++
			continue
		}

		if i < len(lines)-1 {
			copy(lines[i:], lines[i+1:])
		}
		lines = lines[:len(lines)-1]
	}

	if len(lines) == 0 {
		return
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	defer b.sendUpdateUnsafe()

	// If we've been given as many or more lines than we can hold, just do a
	// simple copy of the lattermost lines that will fit in the buffer.
	if l, c := len(lines), cap(b.lines); l >= c {
		prevlen := len(b.lines)
		if len(b.lines) < c {
			b.lines = b.lines[:c]
		}
		copy(b.lines, lines[l-c:l])
		return prevlen + l - c
	}

	// Otherwise, determine how many lines we will need to discard (overflow)
	// in order to accommodate the given number of lines. If the overflow is
	// >0, then that number of elements are dropped from the front of the
	// lines. We already know, based on the above check, that the new lines
	// themselves are fewer than the buffer's capacity, so we will be retaining
	// at least one line in all cases.
	overflow := max(0, (len(b.lines)+len(lines))-cap(b.lines))
	if c := cap(b.lines); overflow > 0 {
		prevlen := len(b.lines)
		if prevlen < c {
			b.lines = b.lines[:c]
		}
		copy(b.lines, b.lines[overflow:prevlen])
		copy(b.lines[c-len(lines):], lines)
	} else {
		b.lines = append(b.lines, lines...)
	}

	return overflow
}

// Close closes the buffer and waits for it to complete.
func (b *Buffer) Close() error {
	b.mu.Lock()
	b.stop = true
	b.mu.Unlock()
	<-b.done
	return nil
}

// Done returns a channel that unblocks once b is done reading from any
// underlying source.
func (b *Buffer) Done() <-chan struct{} {
	return b.done
}

// Len returns the number of lines currently held by b.
func (b *Buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.lines)
}

// Cap returns the maximum number of lines able to be held by b.
func (b *Buffer) Cap() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return cap(b.lines)
}

// SetCap sets the cap for b, allocating more memory as necessary.
func (b *Buffer) SetCap(size int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if size == cap(b.lines) {
		return
	}

	// If the capacity is changing, send an update once it's done.
	defer b.sendUpdateUnsafe()

	// If we're shrinking and we're at full capacity, we need to rotate.
	if size < cap(b.lines) && len(b.lines) >= cap(b.lines) {
		copy(b.lines, b.lines[cap(b.lines)-size:])
		b.lines = b.lines[:size:size]
		return
	}

	// Otherwise, the capacity is changing but the existing lines can be copied
	// over verbatim.
	tmp := make([]string, len(b.lines), size)
	copy(tmp, b.lines)
	b.lines = tmp
}

// Updates returns a channel that receives a message whenever there are line
// updates.
func (b *Buffer) Updates() <-chan struct{} {
	return b.updates
}

func (b *Buffer) isStopped() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.stop
}

func (b *Buffer) scan(src io.Reader) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		b.Add(scanner.Text())
		if b.isStopped() {
			return
		}
	}
}

func (b *Buffer) sendUpdateUnsafe() {
	if b.stop {
		return
	}

	for {
		select {
		case b.updates <- struct{}{}:
			return
		default:
			select {
			case <-b.updates:
			default:
			}
		}
	}
}
