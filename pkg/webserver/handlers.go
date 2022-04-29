package webserver

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/util"
)

func (t *Webserver) handleSoundsList(ctx *routing.Context) error {
	sortOrder := ctx.Query("order")
	filterMust := util.SplitAndClean(ctx.Query("include"), ",")
	filterNot := util.SplitAndClean(ctx.Query("exclude"), ",")

	sounds, err := t.ct.ListSounds(sortOrder, filterMust, filterNot)
	if err != nil {
		return err
	}

	return ctx.Write(sounds)
}

func (t *Webserver) handleSoundsUpload(ctx *routing.Context) error {
	f, fh, err := ctx.Request.FormFile("file")
	if err != nil {
		return ctx.WriteWithStatus(err.Error(), http.StatusBadRequest)
	}

	ct := fh.Header.Get("Content-Type")
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

func (t *Webserver) handleSoundsCreate(ctx *routing.Context) error {
	var req CreateSoundRequest
	err := ctx.Read(&req)
	if err != nil {
		return ctx.WriteWithStatus(err.Error(), http.StatusBadRequest)
	}

	sound, err := t.ct.CreateSound(req)
	if err != nil {
		return err
	}

	return ctx.Write(sound)
}
