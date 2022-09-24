package queue

import (
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	tq := []int{1, 3, 5, 7, 9}
	q := New[int]()
	for _, v := range tq {
		q.Push(v)
		t.Logf("Push %d", v)
	}
	for {
		v, ok := q.Pop()
		t.Log(v, ok)
		if !ok {
			t.Log("return")
			return
		}
	}
}

func TestQueueParallel(t *testing.T) {
	tq := []int{1, 3, 5, 7, 9}
	q := New[int]()
	var wg sync.WaitGroup
	wg.Add(len(tq))
	id := make(chan int)
	for _, v := range tq {
		go func() {
			defer wg.Done()
			goID := <-id
			q.Push(goID)
			t.Logf("Push %d, Len: %d", goID, q.Len())
		}()
		id <- v
	}
	wg.Wait()
	for {
		v, ok := q.Pop()
		t.Log(v, ok)
		if !ok {
			t.Log("return")
			return
		}
	}
}
