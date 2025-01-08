package iterator

import "github.com/wxy365/basal/fn"

type mapIterator[K comparable, V any] struct {
	tgt     map[K]V
	pipe    chan MapEntry[K, V]
	cur     MapEntry[K, V]
	next    MapEntry[K, V]
	hasNext bool
}

func (m *mapIterator[K, V]) HasNext() bool {
	return m.hasNext
}

func (m *mapIterator[K, V]) Next() MapEntry[K, V] {
	m.cur = m.next
	m.next, m.hasNext = <-m.pipe
	return m.cur
}

func (m *mapIterator[K, V]) Remove() {
	delete(m.tgt, m.cur.GetKey())
}

func (m *mapIterator[K, V]) ForEach(consumer fn.Consumer[MapEntry[K, V]]) {
	for m.HasNext() {
		consumer(m.Next())
	}
}

func (m *mapIterator[K, V]) Close() {
	close(m.pipe)
}

func NewMapEntry[K comparable, V any](key K, val V) MapEntry[K, V] {
	return &defaultMapEntry[K, V]{
		key:   key,
		value: val,
	}
}

type MapEntry[K comparable, V any] interface {
	GetKey() K
	GetValue() V
	Set(k K, v V)
}

type defaultMapEntry[K comparable, V any] struct {
	key   K
	value V
	owner map[K]V
}

func (m *defaultMapEntry[K, V]) GetKey() K {
	return m.key
}

func (m *defaultMapEntry[K, V]) GetValue() V {
	return m.value
}

func (m *defaultMapEntry[K, V]) Set(k K, v V) {
	if k == m.key {
		m.owner[k] = v
		m.value = v
	} else {
		delete(m.owner, m.key)
		m.owner[k] = v
		m.key = k
		m.value = v
	}
}
