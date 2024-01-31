package container

import (
	"github.com/coraldane/toolkits/concurrent"
)

type SafeSet[Value any] struct {
	M concurrent.Map
}

func NewSafeSet[Value any]() *SafeSet[Value] {
	return &SafeSet[Value]{
		M: concurrent.Map{},
	}
}

func (this *SafeSet[Value]) Add(val Value) {
	this.M.Store(val, true)
}

func (this *SafeSet[Value]) Remove(val Value) {
	this.M.Delete(val)
}

func (this *SafeSet[Value]) Clear() {
	this.M.Range(func(key, val any) bool {
		this.M.Delete(key)
		return true
	})
}

func (this *SafeSet[Value]) Contains(val Value) bool {
	_, ok := this.M.Load(val)
	return ok
}

func (this *SafeSet[Value]) Size() int {
	return this.M.Length()
}

func (this *SafeSet[Value]) ToSlice() []Value {
	result := make([]Value, 0)
	this.M.Range(func(key, val any) bool {
		keyObj := key.(Value)
		result = append(result, keyObj)
		return true
	})
	return result
}
