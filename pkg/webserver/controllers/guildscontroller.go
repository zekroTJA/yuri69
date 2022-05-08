package controllers

import (
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
)

type guildsController struct {
	ct *controller.Controller
}

func NewGuildsController(r *routing.RouteGroup, ct *controller.Controller) {
	t := guildsController{ct: ct}
	r.Get("/filters", t.handleGetFilters)
	r.Post("/filters", t.handleSetFilters)
	return
}

func (t *guildsController) handleGetFilters(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	filters, err := t.ct.GetGuildFilter(userid)
	if err != nil {
		return err
	}

	return ctx.Write(filters)
}

func (t *guildsController) handleSetFilters(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	var req GuildFilters
	if err := ctx.Read(&req); err != nil {
		return errs.WrapUserError(err)
	}

	err := t.ct.SetGuildFilter(userid, req)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}
