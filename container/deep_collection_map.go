package container

import (
	"sync"
)

type DeepCollectionMap[F, K comparable, V any] struct {
	table *sync.Map
}

func NewDeepCollectionMap[F, K comparable, V any]() *DeepCollectionMap[F, K, V] {
	inst := &DeepCollectionMap[F, K, V]{}
	inst.table = &sync.Map{}
	return inst
}

func (this *DeepCollectionMap[F, K, V]) Put(field F, key K, val V) {
	children := this.GetChildren(field)
	if nil == children {
		children = NewCollectionMap[K, V]()
	}
	children.Put(key, val)
	this.table.Store(field, children)
}

func (this *DeepCollectionMap[F, K, V]) Get(field F, key K) []V {
	children := this.GetChildren(field)
	if nil != children {
		return children.Get(key)
	}
	return nil
}

func (this *DeepCollectionMap[F, K, V]) Keys() []F {
	result := make([]F, 0)
	this.table.Range(func(key, val any) bool {
		result = append(result, key.(F))
		return true
	})
	return result
}

func (this *DeepCollectionMap[F, K, V]) Len() int {
	rowCount := 0
	this.table.Range(func(key, val any) bool {
		rowCount++
		return true
	})
	return rowCount
}

func (this *DeepCollectionMap[F, K, V]) GetChildren(field F) *CollectionMap[K, V] {
	obj, ok := this.table.Load(field)
	var children *CollectionMap[K, V]
	if !ok {
		children = NewCollectionMap[K, V]()
	} else {
		children = obj.(*CollectionMap[K, V])
	}
	return children
}

func (this *DeepCollectionMap[F, K, V]) Remove(field F, key K) {
	children := this.GetChildren(field)
	if nil != children {
		children.Delete(key)
	}
}

func (this *DeepCollectionMap[F, K, V]) RemoveChildren(field F) {
	this.table.Delete(field)
}

func (this *DeepCollectionMap[F, K, V]) Clear() {
	this.table.Range(func(key, val interface{}) bool {
		this.table.Delete(key)
		return true
	})
}
