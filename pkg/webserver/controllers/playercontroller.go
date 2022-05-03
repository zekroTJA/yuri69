package controllers

import (
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	. "github.com/zekrotja/yuri69/pkg/models"
)

type playerController struct {
	ct *controller.Controller
}

func NewPlayerController(r *routing.RouteGroup, ct *controller.Controller) {
	t := playerController{ct: ct}
	r.Post("/join", t.handleJoin)
	r.Post("/leave", t.handleLeave)
	r.Post("/play/<ident>", t.handlePlay)
	r.Post("/stop", t.handleStop)

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
