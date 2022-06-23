package controller

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/kkdai/youtube/v2"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/database/dberrors"
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

	buff, err := io.ReadAll(r)
	if err != nil {
		return "", d, err
	}

	m, err := mimetype.DetectReader(bytes.NewReader(buff))
	if err != nil {
		return "", d, err
	}
	ext = mapExt(m.Extension())
	mimeType = m.String()

	id := xid.New().String()
	err = t.st.PutObject(static.BucketTemp, id, bytes.NewReader(buff), size, mimeType)
	if err != nil {
		return "", d, err
	}

	const lifetime = 5 * time.Minute
	t.pendingCrations.Set(id, ext[1:], lifetime, func(v string) {
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

	s, err := t.db.GetSound(req.Uid)
	if s.Uid == req.Uid {
		return Sound{}, errs.WrapUserError("sound with specified ID already exists")
	}
	if err != nil && err != dberrors.ErrNotFound {
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

	err = t.createSound(req.Sound, &buf, int64(buf.Len()))
	if err != nil {
		return Sound{}, err
	}

	err = t.resizeHistoryBuffer()
	return req.Sound, err
}

func (t *Controller) GetSound(uid string) (Sound, error) {
	sound, err := t.db.GetSound(uid)
	if err != nil {
		return Sound{}, err
	}

	user, err := t.dg.GetUser(sound.Creator.ID)
	if err != nil {
		return Sound{}, err
	}

	sound.Creator = UserSlimFromUser(*user)

	return sound, err
}

func (t *Controller) GetSoundReader(uid string) (io.ReadCloser, int64, error) {
	_, err := t.GetSound(uid)
	if err != nil {
		return nil, 0, err
	}

	return t.st.GetObject(static.BucketSounds, uid)
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

	if oldSound.Creator.ID != userID {
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
	newSound.Creator.ID = oldSound.Creator.ID
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

	if sound.Creator.ID != userID {
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

func (t *Controller) GetSoundFromYoutube(req CreateSoundRequest) (Sound, error) {
	if req.YouTube.URL == "" {
		return Sound{}, errs.WrapUserError("YouTube URL is empty")
	}
	if req.YouTube.EndTimeSeconds > 0 && req.YouTube.StartTimeSeconds > req.YouTube.EndTimeSeconds {
		return Sound{}, errs.WrapUserError("'end_time_seconds' must be larger than 'start_time_seconds'")
	}

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

	s, err := t.db.GetSound(req.Uid)
	if s.Uid == req.Uid {
		return Sound{}, errs.WrapUserError("sound with specified ID already exists")
	}
	if err != nil && err != dberrors.ErrNotFound {
		return Sound{}, err
	}

	client := youtube.Client{}
	video, err := client.GetVideo(req.YouTube.URL)
	if err != nil {
		return Sound{}, err
	}

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return Sound{}, errs.WrapUserError("the provided video does not have any audio streams")
	}
	formats.Sort()
	format := &formats[0]
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return Sound{}, err
	}

	var args []string
	if req.Normalize {
		args = append(args, "-af", "loudnorm=I=-16:TP=-0.3:LRA=11")
	}

	if req.YouTube.StartTimeSeconds > 0 || req.YouTube.EndTimeSeconds > 0 {
		args = append(args, "-ss", fmt.Sprintf("%.4f", req.YouTube.StartTimeSeconds))
	}
	if req.YouTube.EndTimeSeconds > 0 {
		args = append(args, "-t", fmt.Sprintf("%.4f",
			req.YouTube.EndTimeSeconds-req.YouTube.StartTimeSeconds))
	}

	mtyp := mimetype.Lookup(strings.SplitN(format.MimeType, ";", 2)[0])
	if len(formats) == 0 {
		return Sound{}, errs.WrapUserError(
			fmt.Sprintf("could not match any mime type to the extracted stream (%s)", format.MimeType))
	}
	var buf bytes.Buffer
	err = t.ffmpeg(stream, mtyp.Extension()[1:], &buf, "ogg", args...)
	if err != nil {
		return Sound{}, err
	}

	err = t.st.PutObject(static.BucketSounds, req.Uid, &buf, int64(buf.Len()), static.SoundsMime)
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

func (t *Controller) DownloadAllSounds() (rc io.ReadCloser, err error) {
	defer func() {
		if err != nil {
			rc.Close()
		}
	}()

	sounds, err := t.db.GetSounds()
	if err != nil {
		return nil, err
	}

	metaData, err := json.MarshalIndent(sounds, "", "  ")
	if err != nil {
		return nil, err
	}

	f, err := os.CreateTemp(".", "soundspkg-")
	if err != nil {
		return nil, err
	}

	rc = util.WrapReadCloser(f, func(err error) error {
		return os.Remove(f.Name())
	})

	gzipWriter := gzip.NewWriter(f)
	tarWriter := tar.NewWriter(gzipWriter)

	err = tarWriter.WriteHeader(&tar.Header{
		Name: "meta.json",
		Size: int64(len(metaData)),
		Mode: 0644,
	})
	if err != nil {
		return nil, err
	}
	_, err = tarWriter.Write(metaData)
	if err != nil {
		return nil, err
	}

	fileExt := static.SoundsMimeType.Extension()

	for _, sound := range sounds {
		sRc, size, err := t.st.GetObject(static.BucketSounds, sound.Uid)
		if err != nil {
			return nil, err
		}
		err = tarWriter.WriteHeader(&tar.Header{
			Name: path.Join("sounds", sound.Uid+fileExt),
			Size: size,
			Mode: 0644,
		})
		if err != nil {
			sRc.Close()
			return nil, err
		}
		_, err = io.CopyN(tarWriter, sRc, size)
		if err != nil {
			sRc.Close()
			return nil, err
		}
		sRc.Close()
	}

	err = tarWriter.Close()
	if err != nil {
		return nil, err
	}

	err = gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	return rc, nil
}

func (t *Controller) ImportSounds(userID string, f io.ReadCloser, mimeType string) (ImportResult, error) {
	ok, err := t.isAdmin(userID)
	if err != nil {
		return ImportResult{}, err
	}
	if !ok {
		return ImportResult{}, errs.WrapUserError(
			"you need admin privileges import sounds")
	}

	if !strings.HasPrefix(mimeType, "application/tar+gzip") &&
		!strings.HasPrefix(mimeType, "application/gzip") {
		return ImportResult{}, errs.WrapUserError(
			"currently, only archives of type 'application/tar+gzip' are supported")
	}

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return ImportResult{}, err
	}

	tarr := tar.NewReader(gzr)
	res := ImportResult{}
	var metaMap map[string]Sound
	var inBuff, outBuff bytes.Buffer
	copyBuf := make([]byte, 16*1024)

	for {
		header, err := tarr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return ImportResult{}, err
		}

		if header.Name == "meta.json" {
			var meta []Sound
			err = json.NewDecoder(tarr).Decode(&meta)
			if err != nil {
				return ImportResult{}, errs.WrapUserError(err)
			}
			metaMap = make(map[string]Sound)
			for _, sound := range meta {
				metaMap[sound.Uid] = sound
			}
			continue
		}

		if path.Dir(header.Name) == "sounds" {
			if metaMap == nil {
				return ImportResult{}, errs.WrapUserError("no metadata found")
			}

			uid := util.CleanBase(header.Name)
			meta, ok := metaMap[uid]
			if !ok {
				continue
			}

			_, err = t.db.GetSound(uid)
			if err == nil {
				res.Failed = append(res.Failed, SoundImportError{
					Uid:   uid,
					Error: "sound with uid already exists",
				})
				continue
			}

			inBuff.Reset()
			_, err := io.CopyBuffer(&inBuff, tarr, copyBuf)
			if err != nil {
				res.Failed = append(res.Failed, SoundImportError{
					Uid:   uid,
					Error: err.Error(),
				})
				continue
			}

			typ := mimetype.Detect(inBuff.Bytes())
			if typ == nil {
				res.Failed = append(res.Failed, SoundImportError{
					Uid:   uid,
					Error: "could not detect mime type from file stream",
				})
				continue
			}

			outBuff.Reset()
			err = t.ffmpeg(&inBuff, mapExt(typ.Extension())[1:], &outBuff, "ogg")
			if err != nil {
				res.Failed = append(res.Failed, SoundImportError{
					Uid:   uid,
					Error: err.Error(),
				})
				continue
			}

			err = t.createSound(meta, &outBuff, int64(outBuff.Len()))
			if err != nil {
				res.Failed = append(res.Failed, SoundImportError{
					Uid:   uid,
					Error: err.Error(),
				})
				continue
			}

			res.Successful = append(res.Successful, uid)
		}
	}

	return res, nil
}

// --- helpers ---

func (t *Controller) createSound(sound Sound, r io.Reader, size int64) (err error) {
	err = t.st.PutObject(static.BucketSounds, sound.Uid, r, size, static.SoundsMime)
	if err != nil {
		return err
	}

	sound.Created = time.Now()
	err = t.db.PutSound(sound)
	if err != nil {
		stErr := t.st.DeleteObject(static.BucketSounds, sound.Uid)
		if stErr != nil {
			logrus.
				WithError(stErr).
				WithField("id", sound.Uid).Error("Failed removing temp uploaded sound")
		}
		return err
	}

	t.Publish(ControllerEvent{
		IsBroadcast: true,
		Event: Event[any]{
			Type:    EventSoundCreated,
			Origin:  EventSenderController,
			Payload: sound,
		},
	})

	return nil
}
