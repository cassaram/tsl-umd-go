package tsl

import "sync"

type safeMap[K comparable, V any] struct {
	data map[K]V
	mute sync.Mutex
}

func makeSafeMap[K comparable, V any]() safeMap[K, V] {
	return safeMap[K, V]{
		data: make(map[K]V),
	}
}

func (d safeMap[K, V]) get(Key K) V {
	d.mute.Lock()
	defer d.mute.Unlock()

	return d.data[Key]
}

func (d safeMap[K, V]) set(Key K, Value V) {
	d.mute.Lock()
	defer d.mute.Unlock()

	d.data[Key] = Value
}
