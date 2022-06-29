package controllers

import (
	"net/http"

	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
)

type twitchController struct {
	ct *controller.Controller
}

func NewTwitchController(r *routing.RouteGroup, ct *controller.Controller) {
	t := twitchController{ct: ct}
	r.Get("state", t.getState)
	r.Get("sounds", t.getSounds)
	r.Get("play/random", t.play)
	r.Get("play/<id>", t.play)
	return
}

func (t *twitchController) getState(ctx *routing.Context) error {
	claims, _ := ctx.Get("claims").(auth.Claims)

	state, err := t.ct.TwitchState(claims.Username)
	if err != nil {
		return err
	}

	return ctx.Write(state)
}

func (t *twitchController) getSounds(ctx *routing.Context) error {
	claims, _ := ctx.Get("claims").(auth.Claims)
	order := ctx.Query("order")

	sounds, err := t.ct.TwitchListSounds(claims.Username, order)
	if err != nil {
		return err
	}

	return ctx.Write(sounds)
}

func (t *twitchController) play(ctx *routing.Context) error {
	claims, _ := ctx.Get("claims").(auth.Claims)
	ident := ctx.Param("id")

	ok, res, err := t.ct.TwitchPlay(claims.Username, ident)
	if err != nil {
		return err
	}

	var payload StatusWithReservation
	payload.Ratelimit = res

	if !ok {
		payload.StatusModel.Status = http.StatusTooManyRequests
		payload.StatusModel.Message = "you have been rate limited"
		return ctx.WriteWithStatus(payload, http.StatusTooManyRequests)
	}

	payload.StatusModel = StatusOK
	return ctx.Write(payload)
}
