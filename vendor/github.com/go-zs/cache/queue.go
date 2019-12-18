package cache

import (
	"container/list"
	"sync"
)

var lock sync.Mutex

type queue struct {
	*list.List
}

func NewQueue() *queue {
	return &queue{list.New()}
}

func (q *queue) PushFront(v interface{}) {
	defer lock.Unlock()
	lock.Lock()
	q.List.PushFront(v)
}

func (q *queue) Remove(e *list.Element) {
	defer lock.Unlock()
	lock.Lock()
	q.List.Remove(e)
}

func (q *queue) MoveToFront(e *list.Element) {
	defer lock.Unlock()
	lock.Lock()
	q.List.MoveToFront(e)
}
