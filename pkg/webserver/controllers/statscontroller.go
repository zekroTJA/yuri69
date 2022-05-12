package controllers

import (
	"sort"

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

	sort.Slice(log, func(i, j int) bool {
		return log[i].Timestamp.After(log[j].Timestamp)
	})

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
