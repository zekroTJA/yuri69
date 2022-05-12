package controllers

import (
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/player"
	"github.com/zekrotja/yuri69/pkg/util"
)

type playerController struct {
	ct *controller.Controller
}

func NewPlayerController(r *routing.RouteGroup, ct *controller.Controller) {
	t := playerController{ct: ct}
	r.Post("/join", t.handleJoin)
	r.Post("/leave", t.handleLeave)
	r.Post("/play/random", t.handlePlayRandom)
	r.Post("/play/external", t.handlePlayExternal)
	r.Post("/play/<ident>", t.handlePlay)
	r.Post("/stop", t.handleStop)
	r.Post("/volume", t.handleSetVolume)

	return
}

func (t *playerController) handleJoin(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	err := t.ct.JoinChannel(userid)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *playerController) handleLeave(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	err := t.ct.LeaveChannel("", userid)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *playerController) handlePlayRandom(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	filterMust := util.SplitAndClean(ctx.Query("include"), ",")
	filterNot := util.SplitAndClean(ctx.Query("exclude"), ",")

	err := t.ct.PlayRandom(userid, filterMust, filterNot)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *playerController) handlePlayExternal(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	ident := ctx.Query("url")

	err := t.ct.Play(userid, ident)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *playerController) handlePlay(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	ident := ctx.Param("ident")

	err := t.ct.Play(userid, ident)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *playerController) handleStop(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	err := t.ct.Stop(userid)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}

func (t *playerController) handleSetVolume(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)

	var req SetVolumeRequest
	if err := ctx.Read(&req); err != nil {
		return errs.WrapUserError(err)
	}

	if req.Volume < 1 || req.Volume > 200 {
		return errs.WrapUserError("Volume must be a value in range [1, 200]")
	}

	err := t.ct.SetVolume(userid, req.Volume)
	if err != nil && err != player.ErrNoGuildPlayer {
		return err
	}

	return ctx.Write(StatusOK)
}
