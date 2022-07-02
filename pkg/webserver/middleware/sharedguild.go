package middleware

import (
	"net/http"

	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
)

func SharedGuild(ct *controller.Controller) func(*routing.Context) error {
	return func(ctx *routing.Context) error {
		claims, ok := ctx.Get("claims").(auth.Claims)
		if !ok || claims.UserID == "" {
			return errs.WrapUserError("no identity claims found", http.StatusUnauthorized)
		}

		if !claims.IsAuthOrigin(auth.AuthOriginDiscord) {
			return errs.WrapUserError("invalid claims", http.StatusForbidden)
		}

		ok, err := ct.HasSharedGuild(claims.UserID)
		if err != nil {
			return err
		}
		if !ok {
			return errs.WrapUserError(
				"you need to share a guild with Yuri to access this resource", http.StatusForbidden)
		}

		return nil
	}
}
