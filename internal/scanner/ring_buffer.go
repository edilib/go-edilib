package scanner

import (
	"fmt"
	"io"
)

type RuneRingBuffer struct {
	buffer   []rune
	size     int
	indexOut int
	indexIn  int
}

func NewRuneRingBuffer(capacity int) *RuneRingBuffer {
	return &RuneRingBuffer{buffer: make([]rune, capacity), size: 0, indexIn: 0, indexOut: 0}
}

func (r *RuneRingBuffer) IsEmpty() bool {
	return r.size == 0
}

func (r *RuneRingBuffer) Size() int {
	return r.size
}

func (r *RuneRingBuffer) Peek(index int) (rune, error) {
	if index >= r.size {
		return 0, io.EOF
	}

	return r.buffer[(r.indexOut+index)%len(r.buffer)], nil
}

func (r *RuneRingBuffer) Capacity() int {
	return len(r.buffer)
}

func (r *RuneRingBuffer) AddAll(bs []rune) error {
	if r.size > len(r.buffer)-len(bs)-1 {
		return fmt.Errorf("no space left")
	}

	for i := 0; i < len(bs); i++ {
		r.buffer[r.indexIn] = bs[i]
		r.indexIn = (r.indexIn + 1) % len(r.buffer)
	}

	return nil

}

func (r *RuneRingBuffer) Add(b rune) error {
	if r.size >= len(r.buffer)-1 {
		return fmt.Errorf("no space left")
	}

	r.buffer[r.indexIn] = b
	r.indexIn = (r.indexIn + 1) % len(r.buffer)

	r.size = r.size + 1

	return nil
}

func (r *RuneRingBuffer) Remove() (rune, error) {
	if r.size == 0 {
		return 0, io.EOF
	}

	b := r.buffer[r.indexOut]
	r.indexOut = (r.indexOut + 1) % len(r.buffer)
	r.size = r.size - 1

	return b, nil
}
