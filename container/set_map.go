package container

import (
	"sync"
)

type SetMap[Key comparable, Value any] struct {
	DataMap sync.Map
}

func NewSetMap[K comparable, V any]() *SetMap[K, V] {
	return &SetMap[K, V]{
		DataMap: sync.Map{},
	}
}

func (this *SetMap[Key, Value]) ContainsKey(key Key) bool {
	_, ok := this.DataMap.Load(key)
	return ok
}

func (this *SetMap[Key, Value]) Put(key Key, val Value) {
	var list *SafeSet[Value]
	obj, ok := this.DataMap.Load(key)
	if !ok {
		list = NewSafeSet[Value]()
	} else {
		list = obj.(*SafeSet[Value])
	}
	list.Add(val)
	this.DataMap.Store(key, list)
}

func (this *SetMap[Key, Value]) PutValues(key Key, values ...Value) {
	var list *SafeSet[Value]
	obj, ok := this.DataMap.Load(key)
	if !ok {
		list = NewSafeSet[Value]()
	} else {
		list = obj.(*SafeSet[Value])
	}

	for _, val := range values {
		list.Add(val)
	}
	this.DataMap.Store(key, list)
}

func (this *SetMap[Key, Value]) Get(key Key) []Value {
	result := make([]Value, 0)
	obj, ok := this.DataMap.Load(key)
	if !ok {
		return result
	}

	list := obj.(*SafeSet[Value])
	return list.ToSlice()
}

func (this *SetMap[Key, Value]) Len() int {
	rowCount := 0
	this.DataMap.Range(func(key, val any) bool {
		rowCount++
		return true
	})
	return rowCount
}

func (this *SetMap[Key, Value]) Range(fn func(key Key, val []Value)) {
	this.DataMap.Range(func(k, val any) bool {
		newKey := k.(Key)
		newVal := val.(*SafeSet[Value])
		fn(newKey, newVal.ToSlice())
		return true
	})
}

func (this *SetMap[Key, Value]) Delete(key Key) {
	this.DataMap.Delete(key)
}
