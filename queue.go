package queue

import (
	"sync/atomic"
	"unsafe"
)

type elem[T any] struct {
	value  T
	next   *elem[T]
	refcnt uint8
}

type Queue[T any] struct {
	tail *elem[T]
	head *elem[T]
	len  int32
}

func New[T any](value T) *Queue[T] {
	head := &elem[T]{}
	return &Queue[T]{head, head, 0}
}

// Enqueue
func (q *Queue[T]) Push(value T) bool {
	// we must read the pointer address carefully
	// so read it atomically
	var next *unsafe.Pointer
	tail := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(q.tail)))
	if atomic.LoadInt32(&q.len) == 0 {
		next = (*unsafe.Pointer)(tail)
	} else {
		next = (*unsafe.Pointer)(unsafe.Pointer((*elem[T])(tail).next))
	}
	new := unsafe.Pointer(&elem[T]{value: value})
	for !atomic.CompareAndSwapPointer(next, nil, new) {
	}
	if atomic.CompareAndSwapPointer(next, *next, new) {
		atomic.AddInt32(&q.len, 1)
		return true
	}
	// if CAS fail, some goroutine may enquenue
	// so do Retry-loop
	for (*elem[T])(tail).next != nil {
		tail = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer((*elem[T])(tail).next)))
	}
	next = (*unsafe.Pointer)(unsafe.Pointer((*elem[T])(tail).next))
	for !atomic.CompareAndSwapPointer(next, nil, new) {
	}

	if atomic.CompareAndSwapPointer(next, *next, new) {
		atomic.AddInt32(&q.len, 1)
		return true
	}
	return false
}

// Dequeue
func (q *Queue[T]) Pop() (ret T, ok bool) {
	var p unsafe.Pointer
	for {
		p = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(q.head)))
		if (*elem[T])(p).next == nil {
			return
		}
		if !atomic.CompareAndSwapPointer(&p, p, unsafe.Pointer((*elem[T])(p).next)) {
			break
		}
	}
	ret = (T)((*elem[T])(p).next.value)
	ok = true
	atomic.AddInt32(&q.len, -1)
	return
}
