package container

import (
	"container/list"
	"sync"
)

type Deque struct {
	sync.RWMutex
	L *list.List
}

func NewDeque() *Deque {
	result := &Deque{}
	result.L = list.New()
	return result
}

func (this *Deque) PushFront(v any) *list.Element {
	this.Lock()
	e := this.L.PushFront(v)
	this.Unlock()
	return e
}

func (this *Deque) PushBack(v any) *list.Element {
	this.Lock()
	e := this.L.PushBack(v)
	this.Unlock()
	return e
}

func (this *Deque) PushFrontBatch(vs []any) {
	this.Lock()
	for _, item := range vs {
		this.L.PushFront(item)
	}
	this.Unlock()
}

func (this *Deque) PopBack() any {
	this.Lock()

	if elem := this.L.Back(); elem != nil {
		item := this.L.Remove(elem)
		this.Unlock()
		return item.(any)
	}

	this.Unlock()

	var res any
	return res
}

func (this *Deque) PopBackBy(max int) []any {
	this.Lock()

	count := this.len()
	if count == 0 {
		this.Unlock()
		return []any{}
	}

	if count > max {
		count = max
	}

	items := make([]any, 0, count)
	for i := 0; i < count; i++ {
		item := this.L.Remove(this.L.Back())
		items = append(items, item.(any))
	}

	this.Unlock()
	return items
}

func (this *Deque) PopBackAll() []any {
	this.Lock()

	count := this.len()
	if count == 0 {
		this.Unlock()
		return []any{}
	}

	items := make([]any, 0, count)
	for i := 0; i < count; i++ {
		item := this.L.Remove(this.L.Back())
		items = append(items, item.(any))
	}

	this.Unlock()
	return items
}

func (this *Deque) Remove(e *list.Element) any {
	this.Lock()
	defer this.Unlock()
	return this.L.Remove(e).(any)
}

func (this *Deque) RemoveAll() {
	this.Lock()
	this.L = list.New()
	this.Unlock()
}

func (this *Deque) FrontAll() []any {
	this.RLock()
	defer this.RUnlock()

	count := this.len()
	if count == 0 {
		return []any{}
	}

	items := make([]any, 0, count)
	for e := this.L.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(any))
	}
	return items
}

func (this *Deque) BackAll() []any {
	this.RLock()
	defer this.RUnlock()

	count := this.len()
	if count == 0 {
		return []any{}
	}

	items := make([]any, 0, count)
	for e := this.L.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value.(any))
	}
	return items
}

func (this *Deque) Front() any {
	this.RLock()

	if f := this.L.Front(); f != nil {
		this.RUnlock()
		return f.Value.(any)
	}

	this.RUnlock()
	var res any
	return res
}

func (this *Deque) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.len()
}

func (this *Deque) len() int {
	return this.L.Len()
}
