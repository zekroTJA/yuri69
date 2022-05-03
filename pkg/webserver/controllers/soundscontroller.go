package controllers

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/util"
)

type soundsController struct {
	ct *controller.Controller
}

func NewSoundsController(r *routing.RouteGroup, ct *controller.Controller) {
	t := soundsController{ct: ct}
	r.Get("", t.handleSoundsList)
	r.Put("/upload", t.handleSoundsUpload)
	r.Post("/create", t.handleSoundsCreate)
	r.Post("/<id>", t.handleSoundsUpdate)
	r.Delete("/<id>", t.handleSoundsDelete)
	return
}

func (t *soundsController) handleSoundsList(ctx *routing.Context) error {
	sortOrder := ctx.Query("order")
	filterMust := util.SplitAndClean(ctx.Query("include"), ",")
	filterNot := util.SplitAndClean(ctx.Query("exclude"), ",")

	sounds, err := t.ct.ListSounds(sortOrder, filterMust, filterNot)
	if err != nil {
		return err
	}

	return ctx.Write(sounds)
}

func (t *soundsController) handleSoundsUpload(ctx *routing.Context) error {
	f, fh, err := ctx.Request.FormFile("file")
	if err != nil {
		return ctx.WriteWithStatus(err.Error(), http.StatusBadRequest)
	}

	ct := ctx.Query("type", fh.Header.Get("Content-Type"))
	if ct == "" {
		return errs.WrapUserError("no content type was specified")
	}

	id, deadline, err := t.ct.UploadSound(f, fh.Size, ct)
	if err != nil {
		return err
	}

	return ctx.Write(SoundUploadResponse{
		UploadId: id,
		Deadline: deadline,
	})
}

func (t *soundsController) handleSoundsCreate(ctx *routing.Context) error {
	var req CreateSoundRequest
	err := ctx.Read(&req)
	if err != nil {
		return ctx.WriteWithStatus(err.Error(), http.StatusBadRequest)
	}

	req.CreatorId, _ = ctx.Get("userid").(string)
	sound, err := t.ct.CreateSound(req)
	if err != nil {
		return err
	}

	return ctx.Write(sound)
}

func (t *soundsController) handleSoundsUpdate(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	id := ctx.Param("id")

	var req UpdateSoundRequest
	err := ctx.Read(&req)
	if err != nil {
		return ctx.WriteWithStatus(err.Error(), http.StatusBadRequest)
	}

	req.Uid = id
	newSound, err := t.ct.UpdateSound(req, userid)
	if err != nil {
		return err
	}

	return ctx.Write(newSound)
}

func (t *soundsController) handleSoundsDelete(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	id := ctx.Param("id")

	err := t.ct.RemoveSound(id, userid)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}
