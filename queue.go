package yak

import (
	"container/list"
)

// Implement a simple FIFO queue

type Queue struct {
	*list.List
}

func NewQueue() *Queue {
	return &Queue{list.New()}
}

func (q *Queue) Enqueue(v interface{}) {
	q.PushBack(v)
}

func (q *Queue) Dequeue() interface{} {
	e := q.Front()
	if e == nil {
		return nil
	}
	q.Remove(e)
	return e.Value
}

func (q *Queue) Empty() bool {
    return q.Len() == 0
}
