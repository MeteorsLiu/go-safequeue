package queue

import (
	"math/rand"
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

func generateSlice(num int) []int {
	s := make([]int, num)
	for i := 0; i < num; i++ {
		s = append(s, i)
	}
	return s
}

func generateQueue(num int) *Queue[int] {
	q := New[int]()
	for i := 0; i < num; i++ {
		q.Push(i)
	}
	return q
}

func BenchmarkGoMutexInsert(b *testing.B) {
	var lock sync.Mutex
	var wg sync.WaitGroup
	var s []int
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			lock.Lock()
			defer lock.Unlock()
			defer wg.Done()
			s = append(s, rand.Int())
		}()
	}
	wg.Wait()
}

func BenchmarkGoMutexRead(b *testing.B) {
	var lock sync.Mutex
	var wg sync.WaitGroup
	s := generateSlice(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			lock.Lock()
			defer lock.Unlock()
			defer wg.Done()
			s = s[1:]
		}()
	}
	wg.Wait()
}

func BenchmarkLockFreeInsert(b *testing.B) {
	var wg sync.WaitGroup
	q := New[int]()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			q.Push(rand.Int())
		}()
	}
	wg.Wait()
}

func BenchmarkLockFreeRead(b *testing.B) {
	var wg sync.WaitGroup
	s := generateQueue(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.Pop()
		}()
	}
	wg.Wait()
}
