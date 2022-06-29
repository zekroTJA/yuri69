package controllers

import (
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
)

type usersController struct {
	ct *controller.Controller
}

func NewUsersController(r *routing.RouteGroup, ct *controller.Controller) {
	t := usersController{ct: ct}
	r.Get("/settings/fasttrigger", t.handleGetFastTrigger)
	r.Post("/settings/fasttrigger", t.handleSetFastTrigger)
	r.Get("/settings/favorites", t.handleGetFavorites)
	r.Put("/settings/favorites/<ident>", t.handlePutFavorite)
	r.Delete("/settings/favorites/<ident>", t.handleDeleteFavorite)
	r.Get("/settings/apikey", t.handleGetApiKey)
	r.Post("/settings/apikey", t.handlePostApiKey)
	r.Delete("/settings/apikey", t.handleDeleteApiKey)
	r.Get("/settings/twitch/state", t.getTwitchState)
	r.Post("/settings/twitch/settings", t.postTwitchSettings)
	r.Post("/settings/twitch/join", t.postTwitchJoin)
	r.Post("/settings/twitch/leave", t.postTwitchLeave)
	return
}

func (t *usersController) handleGetFastTrigger(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	ident, err := t.ct.GetFastTrigger(userid)
	if err != nil {
		return err
	}

	return ctx.Write(FastTrigger{FastTrigger: ident})
}

func (t *usersController) handleSetFastTrigger(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	var req FastTrigger
	if err := ctx.Read(&req); err != nil {
		return errs.WrapUserError(err)
	}

	err := t.ct.SetFastTrigger(userid, req.FastTrigger)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *usersController) handleGetFavorites(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	favs, err := t.ct.GetFavorites(userid)
	if err != nil {
		return err
	}

	return ctx.Write(favs)
}

func (t *usersController) handlePutFavorite(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	ident := ctx.Param("ident")

	if ident == "" {
		return errs.WrapUserError("favorite must be specified")
	}

	err := t.ct.AddFavorite(userid, ident)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *usersController) handleDeleteFavorite(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	ident := ctx.Param("ident")

	if ident == "" {
		return errs.WrapUserError("favorite must be specified")
	}

	err := t.ct.RemoveFavorite(userid, ident)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *usersController) handleGetApiKey(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	token, err := t.ct.GetApiKey(userid)
	if err != nil {
		return err
	}

	return ctx.Write(ApiKey{ApiKey: token})
}

func (t *usersController) handlePostApiKey(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	token, err := t.ct.GenerateApiKey(userid)
	if err != nil {
		return err
	}

	return ctx.Write(ApiKey{ApiKey: token})
}

func (t *usersController) handleDeleteApiKey(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	err := t.ct.RemoveApiKey(userid)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *usersController) getTwitchState(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	state, err := t.ct.GetTwitchState(userid)
	if err != nil {
		return err
	}

	return ctx.Write(state)
}

func (t *usersController) postTwitchSettings(ctx *routing.Context) error {
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

func (t *usersController) postTwitchJoin(ctx *routing.Context) error {
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

func (t *usersController) postTwitchLeave(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	err := t.ct.LeaveTwitch(userid)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}
