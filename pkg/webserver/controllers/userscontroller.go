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
