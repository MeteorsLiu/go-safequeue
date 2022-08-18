package queue

import "testing"

func TestQueue(t *testing.T) {
	tq := []int{1, 3, 5, 7, 9}
	q := New[int]()
	for _, v := range tq {
		t.Log(q.Push(v))
	}
	for {
		v, ok := q.Pop()
		if !ok {
			t.Log("End")
			return
		}
		t.Log(v)
	}
}
