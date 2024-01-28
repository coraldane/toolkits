package container

import (
	"gitee.com/coraldane/toolkits/concurrent"
)

// DeepMap /**
/**
Data structure
{ F: { K : V } }
*/
type DeepMap[F, K comparable, V any] struct {
	table *concurrent.Map
}

func NewDeepMap[F, K comparable, V any]() *DeepMap[F, K, V] {
	inst := &DeepMap[F, K, V]{}
	inst.table = &concurrent.Map{}
	return inst
}

func (this *DeepMap[F, K, V]) ContainsKey(key F) bool {
	_, ok := this.table.Load(key)
	return ok
}

func (this *DeepMap[F, K, V]) Put(field F, key K, val V) {
	children := this.GetChildren(field)
	if nil == children {
		children = NewLinkedMap[K, V]()
	}
	children.Put(key, val)
	this.table.Store(field, children)
}

func (this *DeepMap[F, K, V]) Get(field F, key K) (any, bool) {
	children := this.GetChildren(field)
	if nil != children {
		return children.Get(key)
	}
	return nil, false
}

func (this *DeepMap[F, K, V]) Keys() []F {
	result := make([]F, 0)
	this.table.Range(func(key, val any) bool {
		result = append(result, key.(F))
		return true
	})
	return result
}

func (this *DeepMap[F, K, V]) Size() int {
	return this.table.Length()
}

func (this *DeepMap[F, K, V]) GetChildren(field F) *LinkedMap[K, V] {
	obj, ok := this.table.Load(field)
	var children *LinkedMap[K, V]
	if !ok {
		children = NewLinkedMap[K, V]()
	} else {
		children = obj.(*LinkedMap[K, V])
	}
	return children
}

func (this *DeepMap[F, K, V]) Remove(field F, key K) {
	children := this.GetChildren(field)
	if nil != children {
		children.Delete(key)
	}
}

func (this *DeepMap[F, K, V]) RemoveChildren(field F) {
	this.table.Delete(field)
}

func (this *DeepMap[F, K, V]) Clear() {
	this.table.Range(func(key, val interface{}) bool {
		this.table.Delete(key)
		return true
	})
}
