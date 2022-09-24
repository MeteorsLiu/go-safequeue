package queue

import (
	"sync/atomic"
)

type elem[T any] struct {
	value  T
	next   atomic.Pointer[elem[T]]
	refcnt uint32
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
	var p *elem[T]
	for !succ {
		p = q.tail.Load()
		if !p.next.CompareAndSwap(nil, newq) {
			q.tail.CompareAndSwap(p, p.next.Load())
		} else {
			succ = true
		}
	}
	if ok := q.tail.CompareAndSwap(p, newq); ok {
		_ = atomic.AddInt32(&q.len, 1)
	}
}

// Dequeue
func (q *Queue[T]) Pop() (ret T, ok bool) {
	succ := false
	var p *elem[T]
	for !succ {
		p = q.head.Load()
		if p.next.Load() == nil {
			return
		}
		succ = q.head.CompareAndSwap(p, p.next.Load())
	}
	ret = p.next.Load().value
	ok = true
	_ = atomic.AddInt32(&q.len, -1)
	return
}

func (q *Queue[T]) Len() int32 {
	return atomic.LoadInt32(&q.len)
}
