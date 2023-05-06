package container

import (
	"bytes"
	"encoding/json"
	"sort"
	"sync"
)

type Pair[K comparable, V any] struct {
	key   K
	value V
}

func (kv *Pair[K, V]) Key() K {
	return kv.key
}

func (kv *Pair[K, V]) Value() V {
	return kv.value
}

type ByPair[K comparable, V any] struct {
	Pairs    []*Pair[K, V]
	LessFunc func(a *Pair[K, V], j *Pair[K, V]) bool
}

func (a ByPair[K, V]) Len() int           { return len(a.Pairs) }
func (a ByPair[K, V]) Swap(i, j int)      { a.Pairs[i], a.Pairs[j] = a.Pairs[j], a.Pairs[i] }
func (a ByPair[K, V]) Less(i, j int) bool { return a.LessFunc(a.Pairs[i], a.Pairs[j]) }

type LinkedMap[K comparable, V any] struct {
	keys       []K
	values     sync.Map
	escapeHTML bool
}

func NewLinkedMap[K comparable, V any]() *LinkedMap[K, V] {
	inst := &LinkedMap[K, V]{}
	inst.keys = []K{}
	inst.values = sync.Map{}
	inst.escapeHTML = true
	return inst
}

func (this *LinkedMap[K, V]) SetEscapeHTML(on bool) {
	this.escapeHTML = on
}

func (this *LinkedMap[K, V]) Get(key K) (V, bool) {
	val, exists := this.values.Load(key)
	return val.(V), exists
}

func (this *LinkedMap[K, V]) Set(key K, value V) {
	this.values.Store(key, value)
}

func (this *LinkedMap[K, V]) Delete(key K) {
	// check key is in use
	this.values.LoadAndDelete(key)
	// remove from keys
	for i, k := range this.keys {
		if k == key {
			this.keys = append(this.keys[:i], this.keys[i+1:]...)
			break
		}
	}
}

func (this *LinkedMap[K, V]) Keys() []K {
	return this.keys
}

// SortKeys Sort the map keys using your sort func
func (this *LinkedMap[K, V]) SortKeys(sortFunc func(keys []K)) {
	sortFunc(this.keys)
}

// Sort Sort the map using your sort func
func (this *LinkedMap[K, V]) Sort(lessFunc func(a *Pair[K, V], b *Pair[K, V]) bool) {
	pairs := make([]*Pair[K, V], len(this.keys))
	for i, key := range this.keys {
		val, ok := this.values.Load(key)
		if ok {
			pairs[i] = &Pair[K, V]{key, val.(V)}
		}
	}

	sort.Sort(ByPair[K, V]{pairs, lessFunc})

	for i, pair := range pairs {
		this.keys[i] = pair.key
	}
}

func (this *LinkedMap[K, V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(this.escapeHTML)
	for i, k := range this.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		// add key
		if err := encoder.Encode(k); err != nil {
			return nil, err
		}
		buf.WriteByte(':')
		// add value
		val, _ := this.Get(k)
		if err := encoder.Encode(val); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
