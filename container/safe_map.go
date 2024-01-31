package container

import "github.com/coraldane/toolkits/concurrent"

type SafeMap[Key comparable, Value any] struct {
	M concurrent.Map
}

func NewSafeMap[Key comparable, Value any]() *SafeMap[Key, Value] {
	return &SafeMap[Key, Value]{
		M: concurrent.Map{},
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
	_, ok := this.M.Load(key)
	return ok
}

func (this *SafeMap[Key, Value]) Size() int {
	return this.M.Length()
}

func (this *SafeMap[Key, Value]) IsEmpty() bool {
	return this.Size() == 0
}
