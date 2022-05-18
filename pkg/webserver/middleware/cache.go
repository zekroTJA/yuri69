package middleware

import (
	"fmt"
	"strings"
	"time"

	routing "github.com/zekrotja/ozzo-routing/v2"
)

func Cache(maxAge time.Duration, mustRevalidate bool, public bool) routing.Handler {
	var vars []string

	if public {
		vars = append(vars, "public")
	} else {
		vars = append(vars, "private")
	}

	vars = append(vars, fmt.Sprintf("max-age=%.0f", maxAge.Seconds()))

	if mustRevalidate {
		vars = append(vars, "must-revalidate")
	}

	headerVal := strings.Join(vars, ", ")

	return func(ctx *routing.Context) error {
		ctx.Response.Header().Set("Cache-Control", headerVal)
		return nil
	}
}
