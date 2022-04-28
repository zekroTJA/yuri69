package generic

import "sync"

type SyncMap[TKey any, TVal any] struct {
	sync.Map
}

func (t *SyncMap[TKey, TVal]) Load(key TKey) (value TVal, ok bool) {
	vi, ok := t.Map.Load(key)
	if !ok {
		return value, false
	}
	value, ok = vi.(TVal)
	return value, ok
}

func (t *SyncMap[TKey, TVal]) Store(key TKey, value TVal) {
	t.Map.Store(key, value)
}

func (t *SyncMap[TKey, TVal]) LoadOrStore(key TKey, value TVal) (actual TVal, loaded bool) {
	ai, loaded := t.Map.LoadOrStore(key, value)
	actual, _ = ai.(TVal)
	return actual, loaded
}

func (t *SyncMap[TKey, TVal]) LoadAndDelete(key TKey) (value TVal, loaded bool) {
	vi, loaded := t.Map.LoadAndDelete(key)
	if loaded {
		value, _ = vi.(TVal)
	}
	return value, loaded
}

func (t *SyncMap[TKey, TVal]) Range(f func(key TKey, value TVal) bool) {
	t.Map.Range(func(ki, vi any) bool {
		k, _ := ki.(TKey)
		v, _ := vi.(TVal)
		return f(k, v)
	})
}
