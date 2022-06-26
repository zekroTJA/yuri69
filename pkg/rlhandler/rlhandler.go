package rlhandler

import (
	"sync"
	"time"

	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap"
)

type RatetimitHandler struct {
	pool sync.Pool
	tm   *timedmap.TimedMap[string, *ratelimit.Limiter]
}

func New(burst int, reset time.Duration) RatetimitHandler {
	var t RatetimitHandler

	t.pool = sync.Pool{
		New: func() any {
			return ratelimit.NewLimiter(reset, burst)
		},
	}
	t.tm = timedmap.New[string, *ratelimit.Limiter](5 * time.Minute)

	return t
}

func (t RatetimitHandler) Get(key string) *ratelimit.Limiter {
	rl := t.tm.GetValue(key)
	exp := time.Duration(rl.Burst()) * rl.Limit()
	if rl == nil {
		rl = t.pool.Get().(*ratelimit.Limiter)
		t.tm.Set(key, rl, exp)
	} else {
		t.tm.SetExpires(key, exp)
	}
	return rl
}
