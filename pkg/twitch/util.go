package twitch

import "sync"

type lockMap[TKey comparable, TVal any] struct {
	mtx sync.RWMutex
	m   map[TKey]TVal
}

func newLockMap[TKey comparable, TVal any]() *lockMap[TKey, TVal] {
	return &lockMap[TKey, TVal]{
		m: make(map[TKey]TVal),
	}
}

func (t *lockMap[TKey, TVal]) Set(key TKey, val TVal) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.m[key] = val
}

func (t *lockMap[TKey, TVal]) Get(key TKey) (TVal, bool) {
	t.mtx.RLock()
	defer t.mtx.RUnlock()

	v, ok := t.m[key]
	return v, ok
}

func (t *lockMap[TKey, TVal]) SetIfUnset(key TKey, val TVal) bool {
	_, ok := t.Get(key)
	if !ok {
		t.Set(key, val)
		return true
	}
	return false
}

func (t *lockMap[TKey, TVal]) Delete(key TKey) {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	delete(t.m, key)
}
