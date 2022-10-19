package container

import (
	"reflect"
	"sync"
)

type SafeMap[Key comparable, Value any] struct {
	M sync.Map
}

func NewSafeMap[Key comparable, Value any]() *SafeMap[Key, Value] {
	return &SafeMap[Key, Value]{
		M: sync.Map{},
	}
}

func (this *SafeMap[Key, Value]) Put(key Key, val Value) {
	this.M.Store(key, val)
}

func (this *SafeMap[Key, Value]) Get(key Key) (any, bool) {
	val, ok := this.M.Load(key)
	return val, ok
}

func (this *SafeMap[Key, Value]) Delete(key Key) {
	this.M.Delete(key)
}

func (this *SafeMap[Key, Value]) Clear() {
	this.M.Range(func(key, val any) bool {
		this.M.Delete(key)
		return true
	})
}

func (this *SafeMap[Key, Value]) Keys() []Key {
	result := make([]Key, 0)
	this.M.Range(func(k, v any) bool {
		keyObj := k.(Key)
		result = append(result, keyObj)
		return true
	})
	return result
}

func (this *SafeMap[Key, Value]) Range(fn func(key Key, val Value) bool) {
	this.M.Range(func(k, v any) bool {
		keyObj := k.(Key)
		valObj := v.(Value)
		res := fn(keyObj, valObj)
		return res
	})
}

func (this *SafeMap[Key, Value]) ContainsKey(key Key) bool {
	var found bool
	this.M.Range(func(k, val any) bool {
		keyObj := k.(Key)
		if reflect.DeepEqual(keyObj, key) {
			found = true
			return false
		}
		return true
	})
	return found
}

func (this *SafeMap[Key, Value]) Size() int {
	rowCount := 0
	this.M.Range(func(key, val any) bool {
		rowCount++
		return true
	})
	return rowCount
}

func (this *SafeMap[Key, Value]) IsEmpty() bool {
	return this.Size() == 0
}
