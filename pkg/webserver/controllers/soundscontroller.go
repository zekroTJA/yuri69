package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/static"
	"github.com/zekrotja/yuri69/pkg/util"
)

type soundsController struct {
	ct *controller.Controller
}

func NewSoundsController(r *routing.RouteGroup, ct *controller.Controller) {
	t := soundsController{ct: ct}
	r.Get("", t.handleList)
	r.Put("/upload", t.handleUpload)
	r.Post("/create", t.handleCreate)
	r.Get("/<id>", t.handleGet)
	r.Get("/<id>/download", t.handleGetDownload)
	r.Post("/<id>", t.handleUpdate)
	r.Delete("/<id>", t.handleDelete)
	return
}

func (t *soundsController) handleList(ctx *routing.Context) error {
	sortOrder := ctx.Query("order")
	filterMust := util.SplitAndClean(ctx.Query("include"), ",")
	filterNot := util.SplitAndClean(ctx.Query("exclude"), ",")

	sounds, err := t.ct.ListSounds(sortOrder, filterMust, filterNot)
	if err != nil {
		return err
	}

	return ctx.Write(sounds)
}

func (t *soundsController) handleGet(ctx *routing.Context) error {
	uid := ctx.Param("id")
	sound, err := t.ct.GetSound(uid)
	if err != nil {
		return err
	}
	return ctx.Write(sound)
}

func (t *soundsController) handleGetDownload(ctx *routing.Context) error {
	uid := ctx.Param("id")
	r, _, err := t.ct.GetSoundReader(uid)
	if err != nil {
		return err
	}

	ext := mimetype.Lookup(static.SoundsMime).Extension()
	ctx.Response.Header().Set("Content-Type", static.SoundsMime)
	ctx.Response.Header().Set("Content-Disposition",
		fmt.Sprintf("atatchment; filename=\"%s%s\"", uid, ext))
	ctx.Response.WriteHeader(http.StatusOK)

	_, err = io.Copy(ctx.Response, r)
	if err != nil {
		return err
	}

	return nil
}

func (t *soundsController) handleUpload(ctx *routing.Context) error {
	f, fh, err := ctx.Request.FormFile("file")
	if err != nil {
		return errs.WrapUserError(err)
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

func (t *soundsController) handleCreate(ctx *routing.Context) error {
	var req CreateSoundRequest
	err := ctx.Read(&req)
	if err != nil {
		return errs.WrapUserError(err)
	}

	req.CreatorId, _ = ctx.Get("userid").(string)

	var sound Sound
	if req.YouTube != nil {
		sound, err = t.ct.GetSoundFromYoutube(req)
	} else {
		sound, err = t.ct.CreateSound(req)
	}
	if err != nil {
		return err
	}

	return ctx.Write(sound)
}

func (t *soundsController) handleUpdate(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	id := ctx.Param("id")

	var req UpdateSoundRequest
	err := ctx.Read(&req)
	if err != nil {
		return errs.WrapUserError(err)
	}

	req.Uid = id
	newSound, err := t.ct.UpdateSound(req, userid)
	if err != nil {
		return err
	}

	return ctx.Write(newSound)
}

func (t *soundsController) handleDelete(ctx *routing.Context) error {
	userid, _ := ctx.Get("userid").(string)
	id := ctx.Param("id")

	err := t.ct.RemoveSound(id, userid)
	if err != nil {
		return err
	}

	return ctx.Write(StatusOK)
}
