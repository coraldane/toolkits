package container

import (
	"sync"
	"sync/atomic"
)

type SafeSet[Value any] struct {
	M sync.Map
}

func NewSafeSet[Value any]() *SafeSet[Value] {
	return &SafeSet[Value]{
		M: sync.Map{},
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
	var rowCount int32
	this.M.Range(func(key, val any) bool {
		atomic.AddInt32(&rowCount, 1)
		return true
	})
	return int(rowCount)
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
