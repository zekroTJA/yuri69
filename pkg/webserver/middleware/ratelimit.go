package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap"
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/errs"
)

type IdentityGetter func(ctx *routing.Context) (string, error)

func IdentityLookup(key string) IdentityGetter {
	return func(ctx *routing.Context) (string, error) {
		v, ok := ctx.Get(key).(string)
		if !ok || v == "" {
			return "", fmt.Errorf("could not get '%s' from ctx", key)
		}
		return v, nil
	}
}

func RateLimit(burst int, reset time.Duration, identityGetter IdentityGetter) routing.Handler {

	pool := sync.Pool{
		New: func() any {
			return ratelimit.NewLimiter(reset, burst)
		},
	}
	tm := timedmap.New[string, *ratelimit.Limiter](5 * time.Minute)

	return func(ctx *routing.Context) error {
		id, err := identityGetter(ctx)
		if err != nil {
			return err
		}

		rl := tm.GetValue(id)
		if rl == nil {
			rl = pool.Get().(*ratelimit.Limiter)
			tm.Set(id, rl, time.Duration(burst)*reset)
		} else {
			tm.SetExpires(id, time.Duration(burst)*reset)
		}

		ok, res := rl.Reserve()

		ctx.Response.Header().Set("X-RateLimit-Limit", strconv.Itoa(res.Burst))
		ctx.Response.Header().Set("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))
		ctx.Response.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(res.Reset.Unix()/1000)))

		if !ok {
			return errs.WrapUserError("too many requests", http.StatusTooManyRequests)
		}

		return nil
	}
}
