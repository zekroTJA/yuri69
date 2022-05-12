package controllers

import (
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/util"
)

type statsController struct {
	ct *controller.Controller
}

func NewStatsController(r *routing.RouteGroup, ct *controller.Controller) {
	t := statsController{ct: ct}
	r.Get("/log", t.handleGetLog)
	r.Get("/count", t.handleGetCount)
	r.Get("/state", t.handleGetState)
	return
}

func (t *statsController) handleGetLog(ctx *routing.Context) error {
	guildid := ctx.Query("guildid")
	userid := ctx.Query("userid")
	ident := ctx.Query("ident")

	limit, err := util.QueryInt(ctx, "limit", 100)
	if err != nil {
		return errs.WrapUserError(err)
	}
	if limit < 0 {
		return errs.WrapUserError("limit must be larger than 0")
	}

	offset, err := util.QueryInt(ctx, "offset", 0)
	if err != nil {
		return errs.WrapUserError(err)
	}
	if offset < 0 {
		return errs.WrapUserError("offset must be larger than 0")
	}

	log, err := t.ct.GetPlaybackLog(guildid, ident, userid, limit, offset)
	if err != nil && err != database.ErrNotFound {
		return err
	}

	if log == nil {
		log = make([]models.PlaybackLogEntry, 0)
	}

	return ctx.Write(log)
}

func (t *statsController) handleGetCount(ctx *routing.Context) error {
	guildid := ctx.Query("guildid")
	userid := ctx.Query("userid")

	stats, err := t.ct.GetPlaybackStats(guildid, userid)
	if err != nil && err != database.ErrNotFound {
		return err
	}

	return ctx.Write(stats)
}

func (t *statsController) handleGetState(ctx *routing.Context) error {
	state, err := t.ct.GetState()
	if err != nil {
		return err
	}

	return ctx.Write(state)
}
