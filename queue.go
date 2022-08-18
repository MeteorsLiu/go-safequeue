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

func New[T any]() *Queue[T] {
	head := &elem[T]{}
	return &Queue[T]{head, head, 0}
}

// Enqueue
func (q *Queue[T]) Push(value T) bool {
	// pointer to the next element
	var next *unsafe.Pointer
	var p *elem[T]
	// we must read the pointer address carefully
	// so read it atomically
	tail := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)))
	new := unsafe.Pointer(&elem[T]{value: value})
	done := false
	// Repeat until CAS is true
	for !done {
		// p = tail
		p = (*elem[T])(tail)
		// pointer to the pointer to the next element, which avoids golang panic when the next pointer is nil
		next = (*unsafe.Pointer)(unsafe.Pointer(&p.next))
		done = atomic.CompareAndSwapPointer(next, nil, new)
		// Avoid panic
		if p != nil && !done {
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(tail), unsafe.Pointer(p), *next)
		}
	}

	if atomic.CompareAndSwapPointer((*unsafe.Pointer)(tail), unsafe.Pointer(p), new) {
		atomic.AddInt32(&q.len, 1)
		return true
	}
	// if CAS fail, some goroutine may enquenue
	// so do a retry-loop
	if p != nil {
		oldp := p
		for p.next != nil {
			p = p.next
		}

		for !atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&p.next)), nil, new) {
		}

		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(tail), unsafe.Pointer(oldp), new) {
			atomic.AddInt32(&q.len, 1)
			return true
		}
	}
	return false
}

// Dequeue
func (q *Queue[T]) Pop() (ret T, ok bool) {
	var p unsafe.Pointer
	for {
		p = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)))
		if (*elem[T])(p).next == nil {
			return
		}
		if atomic.CompareAndSwapPointer(&p, p, unsafe.Pointer((*elem[T])(p).next)) {
			break
		}
	}
	ret = (T)((*elem[T])(p).next.value)
	ok = true
	atomic.AddInt32(&q.len, -1)
	return
}
