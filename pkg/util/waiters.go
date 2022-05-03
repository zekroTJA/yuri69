package util

import (
	"sync"

	"github.com/zekrotja/yuri69/pkg/generic"
)

type Waiters[TKey any] struct {
	m generic.SyncMap[TKey, *sync.Cond]
}

func (t *Waiters[TKey]) Create(key TKey) *sync.Cond {
	cond := sync.NewCond(&sync.Mutex{})
	t.m.Store(key, cond)
	return cond
}

func (t *Waiters[TKey]) Get(key TKey) (*sync.Cond, bool) {
	return t.m.Load(key)
}

func (t *Waiters[TKey]) CreateAndWait(key TKey) {
	cond := t.Create(key)
	cond.L.Lock()
	defer cond.L.Unlock()
	cond.Wait()
}

func (t *Waiters[TKey]) BroadcastAndRemove(key TKey) bool {
	cond, ok := t.Get(key)
	if !ok {
		return false
	}

	cond.L.Lock()
	defer cond.L.Unlock()
	cond.Broadcast()
	return true
}
