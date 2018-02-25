package pool

import (
	"fmt"
)

type Entry struct {
	Key   interface{}
	Value interface{}
}

type LinkedHashMap struct {
	table  []*Entry
	keyMap map[interface{}]int
}

func NewLinkedHashMap() LinkedHashMap {
	retMap := LinkedHashMap{}
	retMap.keyMap = make(map[interface{}]int)
	return retMap
}

func (hashMap *LinkedHashMap) Put(key interface{}, value interface{}) {
	hash, exists := hashMap.getHashValue(key)
	if exists {
		e := hashMap.table[hash]
		e.Value = value
	} else {
		hashMap.table = append(hashMap.table, &Entry{key, value})
	}
}

func (hashMap *LinkedHashMap) Get(key interface{}) interface{} {
	hash, exists := hashMap.getHashValue(key)
	if exists {
		e := hashMap.table[hash]
		return e.Value
	}
	return nil
}

func (hashMap *LinkedHashMap) Remove(key interface{}) {
	hash, exists := hashMap.getHashValue(key)
	if exists {
		hashMap.table = append(hashMap.table[:hash], hashMap.table[hash+1:]...)
		delete(hashMap.keyMap, key)
	}
}

func (hashMap *LinkedHashMap) Clear() {
	for key, index := range hashMap.keyMap {
		delete(hashMap.keyMap, key)
		hashMap.table = append(hashMap.table[:index], hashMap.table[index+1:]...)
	}
}

func (hashMap *LinkedHashMap) ToSlice() []*Entry {
	return hashMap.table
}

func (hashMap *LinkedHashMap) ToMap() map[string]string {
	retMap := make(map[string]string)
	for _, e := range hashMap.table {
		if e != nil {
			retMap[fmt.Sprintf("%v", e.Key)] = fmt.Sprintf("%v", e.Value)
		}
	}
	return retMap
}

func (this *LinkedHashMap) getHashValue(key interface{}) (int, bool) {
	val, exists := this.keyMap[key]
	if !exists {
		val = len(this.keyMap) + 1
		this.keyMap[key] = val
	}
	return val, exists
}
