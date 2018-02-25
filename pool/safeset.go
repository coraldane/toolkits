package pool

import (
	"sync"
)

type SafeSet struct {
	sync.RWMutex
	M map[interface{}]bool
}

func NewSafeSet() *SafeSet {
	return &SafeSet{
		M: make(map[interface{}]bool),
	}
}

func (this *SafeSet) Add(val interface{}) {
	this.Lock()
	this.M[val] = true
	this.Unlock()
}

func (this *SafeSet) Remove(val interface{}) {
	this.Lock()
	delete(this.M, val)
	this.Unlock()
}

func (this *SafeSet) Clear() {
	this.Lock()
	this.M = make(map[interface{}]bool)
	this.Unlock()
}

func (this *SafeSet) Contains(val interface{}) bool {
	this.RLock()
	_, exists := this.M[val]
	this.RUnlock()
	return exists
}

func (this *SafeSet) Size() int {
	this.RLock()
	len := len(this.M)
	this.RUnlock()
	return len
}

func (this *SafeSet) ToSlice() []interface{} {
	this.RLock()
	defer this.RUnlock()

	count := len(this.M)
	if count == 0 {
		return []interface{}{}
	}

	r := []interface{}{}
	for val := range this.M {
		r = append(r, val)
	}

	return r
}
