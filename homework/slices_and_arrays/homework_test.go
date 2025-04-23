package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Number interface {
	int | int8 | int16 | int32 | int64
}

type CircularQueue[T Number] struct {
	values []T
	head,
	tail,
	countElements int
}

func NewCircularQueue[T Number](size int) *CircularQueue[T] {
	return &CircularQueue[T]{
		values:        make([]T, size),
		head:          0,
		tail:          0,
		countElements: 0,
	}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}
	q.values[q.tail] = value
	q.countElements++
	q.tail = (q.tail + 1) % len(q.values)
	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}
	q.head = (q.head + 1) % len(q.values)
	q.countElements--
	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}
	return q.values[q.head]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}
	tail := (q.tail - 1 + len(q.values)) % len(q.values)
	return q.values[tail]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.countElements == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.countElements == len(q.values)
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	f := queue.Front()
	b := queue.Back()
	assert.Equal(t, 2, f)
	assert.Equal(t, 4, b)

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
