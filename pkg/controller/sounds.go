package controller

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"sort"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/static"
	"github.com/zekrotja/yuri69/pkg/util"
)

func (t *Controller) UploadSound(
	r io.Reader,
	size int64,
	mimeType string,
) (string, time.Time, error) {
	var ext string
	var d time.Time

	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil {
		return "", d, err
	}
	if len(exts) != 0 {
		ext = exts[0][1:]
	} else {
		split := strings.Split(mimeType, "/")
		if len(split) != 2 {
			return "", d, errs.WrapUserError("the given mime type is not detectable")
		}
		ext = strings.ToLower(split[1])
	}

	id := xid.New().String()
	err = t.st.PutObject(static.BucketTemp, id, r, size, mimeType)
	if err != nil {
		return "", d, err
	}

	const lifetime = 5 * time.Minute
	t.pendingCrations.Set(id, ext, lifetime, func(v string) {
		t.st.DeleteObject(static.BucketTemp, id)
	})
	d = time.Now().Add(lifetime)
	return id, d, nil
}

func (t *Controller) CreateSound(req CreateSoundRequest) (Sound, error) {
	req.Sanitize()

	err := req.Check()
	if err != nil {
		return Sound{}, err
	}

	req.Uid = strings.ToLower(req.Uid)
	if util.Contains(reservedUids, req.Uid) {
		return Sound{}, errs.WrapUserError(
			fmt.Sprintf("UID '%s' is reserved and can not be used", req.Uid))
	}

	_, err = t.db.GetSound(req.Uid)
	if err == nil {
		return Sound{}, errs.WrapUserError("sound with specified ID already exists")
	}
	if err != nil && err != database.ErrNotFound {
		return Sound{}, err
	}

	typ := t.pendingCrations.GetValue(req.UploadId)
	if typ == "" {
		return Sound{}, errs.WrapUserError("no sound was uploaded or has been expired")
	}

	r, _, err := t.st.GetObject(static.BucketTemp, req.UploadId)
	if err != nil {
		return Sound{}, err
	}
	defer func() {
		r.Close()
		t.st.DeleteObject(static.BucketTemp, req.UploadId)
		t.pendingCrations.Remove(req.UploadId)
	}()

	var args []string
	if req.Normalize {
		args = append(args, "-af", "loudnorm=I=-16:TP=-0.3:LRA=11")
	}

	var buf bytes.Buffer
	err = t.ffmpeg(r, typ, &buf, "ogg", args...)
	if err != nil {
		return Sound{}, err
	}

	err = t.st.PutObject(static.BucketSounds, req.Uid, &buf, int64(buf.Len()), "audio/ogg")
	if err != nil {
		return Sound{}, err
	}

	req.Sound.Created = time.Now()
	err = t.db.PutSound(req.Sound)
	if err != nil {
		stErr := t.st.DeleteObject(static.BucketSounds, req.Uid)
		if stErr != nil {
			logrus.
				WithError(stErr).
				WithField("id", req.Uid).Error("Failed removing temp uploaded sound")
		}
		return Sound{}, err
	}

	t.Publish(ControllerEvent{
		IsBroadcast: true,
		Event: Event[any]{
			Type:    EventSoundCreated,
			Origin:  EventSenderController,
			Payload: req.Sound,
		},
	})

	err = t.resizeHistoryBuffer()
	return req.Sound, err
}

func (t *Controller) GetSound(uid string) (Sound, error) {
	sound, err := t.db.GetSound(uid)
	return sound, err
}

func (t *Controller) ListSounds(
	order string,
	tagsMust []string,
	tagsNot []string,
) ([]Sound, error) {
	sounds, err := t.listSoundsFiltered(tagsMust, tagsNot)
	if err != nil {
		return nil, err
	}

	if order == "" {
		order = string(SortOrderCreated)
	}

	var less func(i, j int) bool

	switch SortOrder(strings.ToLower(order)) {
	case SortOrderName:
		less = func(i, j int) bool {
			return sounds[i].String() < sounds[j].String()
		}
	case SortOrderCreated:
		less = func(i, j int) bool {
			return sounds[i].Created.After(sounds[j].Created)
		}
	default:
		return nil, errs.WrapUserError("invalid sort order")
	}

	sort.Slice(sounds, less)

	return sounds, nil
}

func (t *Controller) UpdateSound(newSound UpdateSoundRequest, userID string) (Sound, error) {
	oldSound, err := t.db.GetSound(newSound.Uid)
	if err != nil {
		return Sound{}, err
	}

	if oldSound.CreatorId != userID {
		ok, err := t.isAdmin(userID)
		if err != nil {
			return Sound{}, err
		}
		if !ok {
			return Sound{}, errs.WrapUserError(
				"you need admin privileges to edit a sound created by another user")
		}
	}

	newSound.Created = oldSound.Created
	newSound.CreatorId = oldSound.CreatorId
	newSound.Uid = oldSound.Uid

	err = t.db.PutSound(newSound.Sound)
	if err != nil {
		return Sound{}, err
	}

	t.Publish(ControllerEvent{
		IsBroadcast: true,
		Event: Event[any]{
			Type:    EventSoundUpdated,
			Origin:  EventSenderController,
			Payload: newSound.Sound,
		},
	})

	return newSound.Sound, nil
}

func (t *Controller) RemoveSound(id, userID string) error {
	sound, err := t.db.GetSound(id)
	if err != nil {
		return err
	}

	if sound.CreatorId != userID {
		ok, err := t.isAdmin(userID)
		if err != nil {
			return err
		}
		if !ok {
			return errs.WrapUserError(
				"you need admin privileges to remove a sound created by another user")
		}
	}

	err = t.db.RemoveSound(id)
	if err != nil {
		return err
	}

	err = t.st.DeleteObject(static.BucketSounds, id)
	if err != nil {
		return err
	}

	t.Publish(ControllerEvent{
		IsBroadcast: true,
		Event: Event[any]{
			Type:    EventSoundDeleted,
			Origin:  EventSenderController,
			Payload: sound,
		},
	})

	err = t.resizeHistoryBuffer()
	return err
}