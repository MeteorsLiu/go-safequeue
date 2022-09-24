package queue

import (
	"sync/atomic"
)

type elem[T any] struct {
	value  T
	next   atomic.Pointer[elem[T]]
	refcnt uint8
}

type Queue[T any] struct {
	tail atomic.Pointer[elem[T]]
	head atomic.Pointer[elem[T]]
	len  int32
}

func New[T any]() *Queue[T] {
	head := atomic.Pointer[elem[T]]{}
	head.Store(&elem[T]{})
	return &Queue[T]{head, head, 0}
}

// Enqueue
func (q *Queue[T]) Push(value T) {
	newq := &elem[T]{
		value: value,
	}
	succ := false
	var p atomic.Pointer[elem[T]]
	for !succ {
		p = q.tail
		if !p.Load().next.CompareAndSwap(nil, newq) {
			q.tail.CompareAndSwap(p.Load(), p.Load().next.Load())
		} else {
			succ = true
		}
	}
	if ok := q.tail.CompareAndSwap(p.Load(), newq); ok {
		_ = atomic.AddInt32(&q.len, 1)
	}
}

// Dequeue
func (q *Queue[T]) Pop() (ret T, ok bool) {
	succ := false
	var p atomic.Pointer[elem[T]]
	for !succ {
		p = q.head
		if p.Load().next.Load() == nil {
			return
		}
		succ = q.head.CompareAndSwap(p.Load(), p.Load().next.Load())
	}
	ret = p.Load().next.Load().value
	ok = true
	_ = atomic.AddInt32(&q.len, -1)
	return
}

func (q *Queue[T]) Len() int32 {
	return atomic.LoadInt32(&q.len)
}
