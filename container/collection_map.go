package container

import (
	"sync"
)

type CollectionMap[Key comparable, Value any] struct {
	DataMap sync.Map
}

func NewCollectionMap[K comparable, V any]() *CollectionMap[K, V] {
	return &CollectionMap[K, V]{
		DataMap: sync.Map{},
	}
}

func (this *CollectionMap[Key, Value]) Put(key Key, val Value) {
	var list *SafeList[Value]
	obj, ok := this.DataMap.Load(key)
	if !ok {
		list = NewSafeList[Value]()
	} else {
		list = obj.(*SafeList[Value])
	}
	list.PushBack(val)
	this.DataMap.Store(key, list)
}

func (this *CollectionMap[Key, Value]) PutValues(key Key, values ...Value) {
	var list *SafeList[Value]
	obj, ok := this.DataMap.Load(key)
	if !ok {
		list = NewSafeList[Value]()
	} else {
		list = obj.(*SafeList[Value])
	}

	for _, val := range values {
		list.PushBack(val)
	}
	this.DataMap.Store(key, list)
}

func (this *CollectionMap[Key, Value]) Get(key Key) []Value {
	result := make([]Value, 0)
	obj, ok := this.DataMap.Load(key)
	if !ok {
		return result
	}

	list := obj.(*SafeList[Value])
	return list.BackAll()
}

func (this *CollectionMap[Key, Value]) GetBackBy(key Key, max int) (int, []Value) {
	result := make([]Value, 0)
	obj, ok := this.DataMap.Load(key)
	if !ok {
		return 0, result
	}

	list := obj.(*SafeList[Value])
	return list.BackBy(max)
}

func (this *CollectionMap[Key, Value]) Len() int {
	rowCount := 0
	this.DataMap.Range(func(key, val any) bool {
		rowCount++
		return true
	})
	return rowCount
}

func (this *CollectionMap[Key, Value]) Range(fn func(key Key, val []Value)) {
	this.DataMap.Range(func(k, val any) bool {
		newKey := k.(Key)
		newVal := val.(*SafeList[Value])
		fn(newKey, newVal.BackAll())
		return true
	})
}

func (this *CollectionMap[Key, Value]) Delete(key Key) {
	this.DataMap.Delete(key)
}
