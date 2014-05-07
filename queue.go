package yak

import (
	"container/list"
	"sync"
)

// Implement a simple FIFO queue

type Queue struct {
	*sync.Mutex
	*list.List
}

func NewQueue() *Queue {
	return &Queue{&sync.Mutex{}, list.New()}
}

func (q *Queue) Enqueue(v interface{}) {
	q.Lock()
	q.PushBack(v)
	q.Unlock()
}

func (q *Queue) Dequeue() interface{} {
	q.Lock()
	defer q.Unlock()
	e := q.Front()
	if e == nil {
		return nil
	}
	q.Remove(e)
	return e.Value
}

func (q *Queue) Empty() bool {
	q.Lock()
	defer q.Unlock()
	return q.Len() == 0
}
