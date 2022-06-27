package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/rlhandler"
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

	m := rlhandler.New(burst, reset)

	return func(ctx *routing.Context) error {
		id, err := identityGetter(ctx)
		if err != nil {
			return err
		}

		rl := m.Get(id)
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
