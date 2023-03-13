package container

import (
	"bytes"
	"encoding/json"
	"sort"
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
	values     map[K]V
	escapeHTML bool
}

func NewLinkedMap[K comparable, V any]() *LinkedMap[K, V] {
	inst := &LinkedMap[K, V]{}
	inst.keys = []K{}
	inst.values = make(map[K]V)
	inst.escapeHTML = true
	return inst
}

func (o *LinkedMap[K, V]) SetEscapeHTML(on bool) {
	o.escapeHTML = on
}

func (o *LinkedMap[K, V]) Get(key K) (V, bool) {
	val, exists := o.values[key]
	return val, exists
}

func (o *LinkedMap[K, V]) Set(key K, value V) {
	_, exists := o.values[key]
	if !exists {
		o.keys = append(o.keys, key)
	}
	o.values[key] = value
}

func (o *LinkedMap[K, V]) Delete(key K) {
	// check key is in use
	_, ok := o.values[key]
	if !ok {
		return
	}
	// remove from keys
	for i, k := range o.keys {
		if k == key {
			o.keys = append(o.keys[:i], o.keys[i+1:]...)
			break
		}
	}
	// remove from values
	delete(o.values, key)
}

func (o *LinkedMap[K, V]) Keys() []K {
	return o.keys
}

// SortKeys Sort the map keys using your sort func
func (o *LinkedMap[K, V]) SortKeys(sortFunc func(keys []K)) {
	sortFunc(o.keys)
}

// Sort Sort the map using your sort func
func (o *LinkedMap[K, V]) Sort(lessFunc func(a *Pair[K, V], b *Pair[K, V]) bool) {
	pairs := make([]*Pair[K, V], len(o.keys))
	for i, key := range o.keys {
		pairs[i] = &Pair[K, V]{key, o.values[key]}
	}

	sort.Sort(ByPair[K, V]{pairs, lessFunc})

	for i, pair := range pairs {
		o.keys[i] = pair.key
	}
}

func (o LinkedMap[K, V]) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(o.escapeHTML)
	for i, k := range o.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		// add key
		if err := encoder.Encode(k); err != nil {
			return nil, err
		}
		buf.WriteByte(':')
		// add value
		if err := encoder.Encode(o.values[k]); err != nil {
			return nil, err
		}
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}
