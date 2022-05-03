package controller

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/player"
	"github.com/zekrotja/yuri69/pkg/static"
	"github.com/zekrotja/yuri69/pkg/storage"
)

type Controller struct {
	db database.IDatabase
	st storage.IStorage
	pl *player.Player
	dg *discord.Discord

	ffmpegExec string

	pendingCrations *timedmap.TimedMap[string, string]
}

func New(
	db database.IDatabase,
	st storage.IStorage,
	pl *player.Player,
	dg *discord.Discord,
) (*Controller, error) {

	var (
		t   Controller
		err error
	)

	t.db = db
	t.st = st
	t.pl = pl
	t.dg = dg

	t.pendingCrations = timedmap.New[string, string](5 * time.Minute)

	t.ffmpegExec, err = exec.LookPath("ffmpeg")
	if errors.Is(err, exec.ErrNotFound) {
		return nil, errors.New("ffmpeg executable was not found")
	}

	return &t, nil
}

func (t *Controller) Close() error {
	for k := range t.pendingCrations.Snapshot() {
		err := t.st.DeleteObject(static.BucketTemp, k)
		if err != nil {
			logrus.WithError(err).WithField("id", k).Error("Failed removing temp uploaded sound")
		}
	}
	return nil
}

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
	if req.Uid == "" {
		return Sound{}, errs.WrapUserError("uid must be specified")
	}

	_, err := t.db.GetSound(req.Uid)
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

	return req.Sound, nil
}

func (t *Controller) ListSounds(
	order string,
	flagsMust []string,
	flagsNot []string,
) ([]Sound, error) {
	if order == "" {
		order = string(database.SortOrderCreated)
	}
	sounds, err := t.db.GetSounds(database.SortOrder(order), flagsMust, flagsNot)
	if err == database.ErrNotFound {
		sounds = []Sound{}
	} else if err != nil {
		return nil, err
	}

	return sounds, nil
}

func (t *Controller) UpdateSound(newSound UpdateSoundRequest, userID string) (Sound, error) {
	oldSound, err := t.db.GetSound(newSound.Uid)
	if err != nil {
		return Sound{}, err
	}

	// if oldSound.CreatorId != userID {
	// 	return errs.WrapUserError(
	// 		"you need to be either the creator of the sound or an admin to update it",
	// 		http.StatusForbidden)
	// }

	newSound.Created = oldSound.Created
	newSound.CreatorId = oldSound.CreatorId
	newSound.Uid = oldSound.Uid

	err = t.db.PutSound(newSound.Sound)
	if err != nil {
		return Sound{}, err
	}

	return newSound.Sound, nil
}

func (t *Controller) RemoveSound(id, userID string) error {
	sound, err := t.db.GetSound(id)
	if err != nil {
		return err
	}

	if sound.CreatorId != userID {
		return errs.WrapUserError(
			"you need to be either the creator of the sound or an admin to delete it",
			http.StatusForbidden)
	}

	err = t.db.RemoveSound(id)
	if err != nil {
		return err
	}

	err = t.st.DeleteObject(static.BucketSounds, id)
	return err
}

func (t *Controller) JoinChannel(userID string) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	return t.pl.Init(vs.GuildID, vs.ChannelID)
}

func (t *Controller) LeaveChannel(guildID, userID string) error {
	if guildID == "" {
		vs, ok := t.dg.FindUserVS(userID)
		if !ok {
			return errs.WrapUserError("you need to be in a voice channel to perform this action")
		}
		guildID = vs.GuildID
	}

	return t.pl.Destroy(guildID)
}

func (t *Controller) Play(userID, ident string) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	identLower := strings.ToLower(ident)
	if strings.HasPrefix(identLower, "https://youtu.be/") ||
		strings.HasPrefix(identLower, "https://youtube.com/watch?=v") ||
		strings.HasPrefix(identLower, "https://www.youtube.com/watch?=v") {

		return t.pl.Play(vs.GuildID, vs.ChannelID, ident)
	}

	volume, err := t.db.GetGuildVolume(vs.GuildID)
	if err == database.ErrNotFound {
		err = nil
		volume = 50
	}
	if err != nil {
		return err
	}

	if err = t.pl.Init(vs.GuildID, vs.ChannelID); err != nil {
		return err
	}

	if err = t.pl.SetVolume(vs.GuildID, uint16(volume)); err != nil {
		return err
	}

	return t.pl.PlaySound(vs.GuildID, vs.ChannelID, ident)
}

func (t *Controller) Stop(userID string) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	return t.pl.Stop(vs.GuildID)
}

func (t *Controller) SetVolume(userID string, volume int) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	if err := t.db.SetGuildVolume(vs.GuildID, volume); err != nil {
		return err
	}

	return t.pl.SetVolume(vs.GuildID, uint16(volume))
}

// --- Helpers ---

func (t *Controller) ffmpeg(
	in io.Reader,
	inTyp string,
	out io.Writer,
	outTyp string,
	args ...string,
) error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "-f", inTyp, "-i", "pipe:")
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, "-f", outTyp, "pipe:")

	var bufStdErr bytes.Buffer
	cmd := exec.Command(t.ffmpegExec, cmdArgs...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = &bufStdErr
	err := cmd.Run()

	if cmd.ProcessState != nil && cmd.ProcessState.ExitCode() != 0 {
		err = errors.New(bufStdErr.String())
	}

	return err
}
