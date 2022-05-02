package container

import (
	"container/list"
	"sync"
)

type SafeList[T any] struct {
	sync.RWMutex
	L *list.List
}

func NewSafeList[T any]() *SafeList[T] {
	result := &SafeList[T]{}
	result.L = list.New()
	return result
}

func (this *SafeList[T]) PushFront(v T) *list.Element {
	this.Lock()
	e := this.L.PushFront(v)
	this.Unlock()
	return e
}

func (this *SafeList[T]) PushBack(v T) *list.Element {
	this.Lock()
	e := this.L.PushBack(v)
	this.Unlock()
	return e
}

func (this *SafeList[T]) PushFrontBatch(vs []T) {
	this.Lock()
	for _, item := range vs {
		this.L.PushFront(item)
	}
	this.Unlock()
}

func (this *SafeList[T]) PopBack() T {
	this.Lock()

	if elem := this.L.Back(); elem != nil {
		item := this.L.Remove(elem)
		this.Unlock()
		return item
	}

	this.Unlock()
	return nil
}

func (this *SafeList[T]) PopBackBy(max int) []T {
	this.Lock()

	count := this.len()
	if count == 0 {
		this.Unlock()
		return []T{}
	}

	if count > max {
		count = max
	}

	items := make([]T, 0, count)
	for i := 0; i < count; i++ {
		item := this.L.Remove(this.L.Back())
		items = append(items, item)
	}

	this.Unlock()
	return items
}

func (this *SafeList[T]) PopBackAll() []T {
	this.Lock()

	count := this.len()
	if count == 0 {
		this.Unlock()
		return []T{}
	}

	items := make([]T, 0, count)
	for i := 0; i < count; i++ {
		item := this.L.Remove(this.L.Back())
		items = append(items, item)
	}

	this.Unlock()
	return items
}

func (this *SafeList[T]) Remove(e *list.Element) T {
	this.Lock()
	defer this.Unlock()
	return this.L.Remove(e)
}

func (this *SafeList) RemoveAll() {
	this.Lock()
	this.L = list.New()
	this.Unlock()
}

func (this *SafeList[T]) FrontAll() []T {
	this.RLock()
	defer this.RUnlock()

	count := this.len()
	if count == 0 {
		return []T{}
	}

	items := make([]T, 0, count)
	for e := this.L.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
	return items
}

func (this *SafeList[T]) BackAll() []T {
	this.RLock()
	defer this.RUnlock()

	count := this.len()
	if count == 0 {
		return []T{}
	}

	items := make([]T, 0, count)
	for e := this.L.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value)
	}
	return items
}

func (this *SafeList[T]) Front() T {
	this.RLock()

	if f := this.L.Front(); f != nil {
		this.RUnlock()
		return f.Value
	}

	this.RUnlock()
	return nil
}

func (this *SafeList) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.len()
}

func (this *SafeList) len() int {
	return this.L.Len()
}
