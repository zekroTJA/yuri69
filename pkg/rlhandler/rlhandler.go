package rlhandler

import (
	"sync"
	"time"

	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap"
)

type RatelimitHandler struct {
	pool sync.Pool
	tm   *timedmap.TimedMap[string, *ratelimit.Limiter]
}

func New(burst int, reset time.Duration) RatelimitHandler {
	var t RatelimitHandler

	t.pool = sync.Pool{
		New: func() any {
			return ratelimit.NewLimiter(reset, burst)
		},
	}
	t.tm = timedmap.New[string, *ratelimit.Limiter](5 * time.Minute)

	return t
}

func (t RatelimitHandler) Get(key string) *ratelimit.Limiter {
	rl := t.tm.GetValue(key)
	if rl == nil {
		rl = t.pool.Get().(*ratelimit.Limiter)
		t.tm.Set(key, rl, time.Duration(rl.Burst())*rl.Limit())
	} else {
		t.tm.SetExpires(key, time.Duration(rl.Burst())*rl.Limit())
	}
	return rl
}

func (t RatelimitHandler) Update(burst int, reset time.Duration) {
	for _, rl := range t.tm.Snapshot() {
		rl.SetBurst(burst)
		rl.SetLimit(reset)
	}
}
