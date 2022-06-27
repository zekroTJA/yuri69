package controllers

import (
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
)

type twitchController struct {
	ct *controller.Controller
}

func NewTwitchController(r *routing.RouteGroup, ct *controller.Controller) {
	t := twitchController{ct: ct}
	r.Get("/state", t.getState)
	r.Post("/settings", t.postSettings)
	r.Post("/join", t.postJoin)
	r.Post("/leave", t.postLeave)
	return
}

func (t *twitchController) getState(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	state, err := t.ct.GetTwitchState(userid)
	if err != nil {
		return err
	}

	return ctx.Write(state)
}

func (t *twitchController) postSettings(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	var settings TwitchSettings
	if err := ctx.Read(&settings); err != nil {
		return errs.WrapUserError(err)
	}

	err := t.ct.UpdateTwitchSettings(userid, &settings, false)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *twitchController) postJoin(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	var settings *TwitchSettings
	if ctx.Request.ContentLength > 0 {
		settings = new(TwitchSettings)
		if err := ctx.Read(&settings); err != nil {
			return errs.WrapUserError(err)
		}
	}

	err := t.ct.UpdateTwitchSettings(userid, settings, true)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *twitchController) postLeave(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	err := t.ct.LeaveTwitch(userid)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}
