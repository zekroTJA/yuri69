package util

import (
	"sync"

	"github.com/rs/xid"
)

type EventBus[T any] struct {
	buffSize int

	m    sync.Mutex
	subs map[string]chan T
}

func NewEventBus[T any](buffSize ...int) *EventBus[T] {
	return &EventBus[T]{
		buffSize: Opt(buffSize, 1000),
		subs:     map[string]chan T{},
	}
}

func (t *EventBus[T]) Publish(v T) {
	for _, s := range t.subs {
		s <- v
	}
}

func (t *EventBus[T]) Subscribe() (chan T, func()) {
	id := xid.New().String()

	ch := make(chan T, t.buffSize)

	t.m.Lock()
	defer t.m.Unlock()

	t.subs[id] = ch

	unsub := func() {
		delete(t.subs, id)
		close(ch)
	}

	return ch, unsub
}

func (t *EventBus[T]) SubscribeFunc(f func(T)) func() {
	ch, unsub := t.Subscribe()

	go func() {
		for v := range ch {
			f(v)
		}
	}()

	return unsub
}
