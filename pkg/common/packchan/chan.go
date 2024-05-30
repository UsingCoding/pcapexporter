package common

import "slices"

func NewPackChan[T any](bufSize int, c chan []T) *PackChan[T] {
	return &PackChan[T]{
		bufSize: bufSize,
		c:       c,
	}
}

type PackChan[T any] struct {
	buf []T
	c   chan []T

	bufSize int
}

func (c *PackChan[T]) Send(t T) {
	if c.buf == nil {
		c.buf = make([]T, 0, c.bufSize)
	}

	c.buf = append(c.buf, t)

	if len(c.buf) < c.bufSize {
		return
	}

	c.c <- slices.Clone(c.buf)

	c.buf = nil
}

func (c *PackChan[T]) Recv() <-chan []T {
	return c.c
}

func (c *PackChan[T]) Close() {
	close(c.c)
	c.buf = nil
}
