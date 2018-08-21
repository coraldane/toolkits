package container

import (
	"sync"
)

type DeepMap struct {
	table *sync.Map
}

func NewDeepMap() *DeepMap {
	inst := &DeepMap{}
	inst.table = &sync.Map{}
	return inst
}

func (this *DeepMap) Put(field, key, val interface{}) {
	children := this.GetChildren(field)
	if nil == children {
		children = &sync.Map{}
	}
	children.Store(key, val)
	this.table.Store(field, children)
}

func (this *DeepMap) Get(field, key interface{}) interface{} {
	children := this.GetChildren(field)
	if nil != children {
		retVal, ok := children.Load(key)
		if ok {
			return retVal
		}
	}
	return nil
}

func (this *DeepMap) ForEach(fn func(key, val interface{}) bool) {
	this.table.Range(fn)
}

func (this *DeepMap) GetChildren(field interface{}) *sync.Map {
	child, ok := this.table.Load(field)
	if ok {
		retVal, _ := child.(*sync.Map)
		return retVal
	}
	return nil
}

func (this *DeepMap) Remove(field, key interface{}) {
	children := this.GetChildren(field)
	if nil != children {
		children.Delete(key)
	}
}

func (this *DeepMap) RemoveChildren(field interface{}) {
	this.table.Delete(field)
}

func (this *DeepMap) Clear() {
	this.table.Range(func(key, val interface{}) bool {
		this.table.Delete(key)
		return true
	})
}
