package container

import (
	"gitee.com/coraldane/toolkits/concurrent"
)

type DeepSetMap[F, K comparable, V any] struct {
	table *concurrent.Map
}

func NewDeepSetMap[F, K comparable, V any]() *DeepSetMap[F, K, V] {
	inst := &DeepSetMap[F, K, V]{}
	inst.table = &concurrent.Map{}
	return inst
}

func (this *DeepSetMap[F, K, V]) ContainsKey(field F) bool {
	_, ok := this.table.Load(field)
	return ok
}

func (this *DeepSetMap[F, K, V]) Put(field F, key K, val V) {
	children := this.GetChildren(field)
	if nil == children {
		children = NewSetMap[K, V]()
	}
	children.Put(key, val)
	this.table.Store(field, children)
}

func (this *DeepSetMap[F, K, V]) Get(field F, key K) []V {
	children := this.GetChildren(field)
	if nil != children {
		return children.Get(key)
	}
	return nil
}

func (this *DeepSetMap[F, K, V]) Keys() []F {
	result := make([]F, 0)
	this.table.Range(func(key, val any) bool {
		result = append(result, key.(F))
		return true
	})
	return result
}

func (this *DeepSetMap[F, K, V]) Size() int {
	return this.table.Length()
}

func (this *DeepSetMap[F, K, V]) GetChildren(field F) *SetMap[K, V] {
	obj, ok := this.table.Load(field)
	var children *SetMap[K, V]
	if !ok {
		children = NewSetMap[K, V]()
	} else {
		children = obj.(*SetMap[K, V])
	}
	return children
}

func (this *DeepSetMap[F, K, V]) Remove(field F, key K) {
	children := this.GetChildren(field)
	if nil != children {
		children.Delete(key)
	}
}

func (this *DeepSetMap[F, K, V]) RemoveChildren(field F) {
	this.table.Delete(field)
}

func (this *DeepSetMap[F, K, V]) Clear() {
	this.table.Range(func(key, val interface{}) bool {
		this.table.Delete(key)
		return true
	})
}
