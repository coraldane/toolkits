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

func (this *SafeList[T]) PopBack() T {
	this.Lock()

	if elem := this.L.Back(); elem != nil {
		item := this.L.Remove(elem)
		this.Unlock()
		return item.(T)
	}

	this.Unlock()

	var res T
	return res
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
		items = append(items, item.(T))
	}

	this.Unlock()
	return items
}

func (this *SafeList[T]) Remove(e *list.Element) T {
	this.Lock()
	defer this.Unlock()
	return this.L.Remove(e).(T)
}

func (this *SafeList[T]) RemoveAll() {
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
		items = append(items, e.Value.(T))
	}
	return items
}

func (this *SafeList[T]) FrontBy(max int, popItem bool) (int, []T) {
	this.Lock()

	count := this.len()
	if count == 0 {
		this.Unlock()
		return 0, []T{}
	}

	if count > max {
		count = max
	}

	items := make([]T, 0, count)
	index := 0
	for e := this.L.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(T))
		index++
		if popItem {
			this.L.Remove(e)
		}
		if index >= count {
			break
		}
	}

	this.Unlock()
	return count, items
}

func (this *SafeList[T]) BackBy(max int, popItem bool) (int, []T) {
	this.Lock()

	count := this.len()
	if count == 0 {
		this.Unlock()
		return 0, []T{}
	}

	if count > max {
		count = max
	}

	items := make([]T, 0, count)
	index := 0
	for e := this.L.Back(); e != nil; e = e.Prev() {
		items = append(items, e.Value.(T))
		if popItem {
			this.L.Remove(e)
		}
		index++
		if index >= count {
			break
		}
	}

	// reverse items at last
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}

	this.Unlock()
	return count, items
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
		items = append(items, e.Value.(T))
	}
	return items
}

func (this *SafeList[T]) Front() T {
	this.RLock()

	if f := this.L.Front(); f != nil {
		this.RUnlock()
		return f.Value.(T)
	}

	this.RUnlock()
	var res T
	return res
}

func (this *SafeList[T]) Len() int {
	this.RLock()
	defer this.RUnlock()
	return this.len()
}

func (this *SafeList[T]) len() int {
	return this.L.Len()
}
