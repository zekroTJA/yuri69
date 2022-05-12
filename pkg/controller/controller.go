package controller

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/generic"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/player"
	"github.com/zekrotja/yuri69/pkg/static"
	"github.com/zekrotja/yuri69/pkg/storage"
	"github.com/zekrotja/yuri69/pkg/util"
)

var (
	reservedUids = []string{"random", "upload", "create"}
)

type ControllerEvent struct {
	IsBroadcast bool
	Receivers   []string
	Event       Event[any]
}

type Controller struct {
	*util.EventBus[ControllerEvent]

	db database.IDatabase
	st storage.IStorage
	pl *player.Player
	dg *discord.Discord

	ffmpegExec string

	pendingCrations *timedmap.TimedMap[string, string]
	history         *generic.RingQueue[string]
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

	rand.Seed(time.Now().UnixNano())

	t.EventBus = util.NewEventBus[ControllerEvent]()

	t.db = db
	t.st = st
	t.pl = pl
	t.dg = dg

	t.pendingCrations = timedmap.New[string, string](5 * time.Minute)

	t.history = generic.NewRingQueue[string](1)
	if err = t.resizeHistoryBuffer(); err != nil {
		return nil, err
	}

	t.ffmpegExec, err = exec.LookPath("ffmpeg")
	if errors.Is(err, exec.ErrNotFound) {
		return nil, errors.New("ffmpeg executable was not found")
	}

	t.pl.SubscribeFunc(t.playerEventHandler)

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
		return errs.WrapUserError(
			"you need to be either the creator of the sound or an admin to delete it",
			http.StatusForbidden)
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

	return t.play(vs, ident)
}

func (t *Controller) PlayRandom(userID string, tagsMust []string, tagsNot []string) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	var (
		guildFilters GuildFilters
		err          error
	)
	if len(tagsMust) == 0 || len(tagsNot) == 0 {
		guildFilters, err = t.db.GetGuildFilters(vs.GuildID)
		if err != nil && err != database.ErrNotFound {
			return err
		}
	}

	if len(tagsMust) == 0 {
		tagsMust = guildFilters.Include
	}
	if len(tagsNot) == 0 {
		tagsNot = guildFilters.Exclude
	}

	sounds, err := t.listSoundsFiltered(tagsMust, tagsNot)
	if err != nil {
		return err
	}

	if len(sounds) == 0 {
		return nil
	}

	historySnapshot := t.history.Snapshot()

	var sound Sound
	// If no sound could be picked which is not in the history
	// in 10 tries, play the last randomly picked sound anyway
	// to prevent long wait times.
	for i := 0; i < 10; i++ {
		rng := rand.Intn(len(sounds))
		sound = sounds[rng]
		if !util.Contains(historySnapshot, sound.Uid) {
			break
		}
	}

	if err = t.play(vs, sound.Uid); err != nil {
		return nil
	}

	t.history.Enqueue(sound.Uid)
	return nil
}

func (t *Controller) Stop(userID string) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	return t.pl.Stop(vs.GuildID)
}

func (t *Controller) GetVolume(userID string) (int, error) {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return 0, errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	v, err := t.db.GetGuildVolume(vs.GuildID)
	if err != nil && err != database.ErrNotFound {
		return 0, err
	}

	return v, nil
}

func (t *Controller) SetVolume(userID string, volume int) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	if err := t.db.SetGuildVolume(vs.GuildID, volume); err != nil {
		return err
	}

	err := t.pl.SetVolume(vs.GuildID, uint16(volume))
	if err != nil {
		return err
	}

	return t.publishToGuildUsers(vs.GuildID, Event[any]{
		Type:   EventVolumeUpdated,
		Origin: EventSenderController,
		Payload: SetVolumeRequest{
			Volume: volume,
		},
	})
}

func (t *Controller) GetFastTrigger(userID string) (string, error) {
	ident, err := t.db.GetUserFastTrigger(userID)
	if err == database.ErrNotFound {
		err = nil
	}
	return ident, err
}

func (t *Controller) SetFastTrigger(userID, ident string) error {
	return t.db.SetUserFastTrigger(userID, ident)
}

func (t *Controller) GetGuildFilter(userID string) (GuildFilters, error) {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return GuildFilters{},
			errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	f, err := t.db.GetGuildFilters(vs.GuildID)
	if err == database.ErrNotFound {
		err = nil
	}

	return f, err
}

func (t *Controller) SetGuildFilter(userID string, f GuildFilters) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	f.Sanitize()
	err := f.Check()
	if err != nil {
		return err
	}

	err = t.db.SetGuildFilters(vs.GuildID, f)
	if err != nil {
		return err
	}

	return t.publishToGuildUsers(vs.GuildID, Event[any]{
		Type:    EventGuildFilterUpdated,
		Origin:  EventSenderController,
		Payload: f,
	})
}

func (t *Controller) GetCurrentState(userID string) (EventStatePayload, error) {
	var (
		res EventStatePayload
		err error
	)

	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return res, nil
	}

	res.Connected = true
	res.Joined = t.pl.HasPlayer(vs.GuildID)

	res.EventVoiceJoinPayload, err = t.getVoiceJoinPayload(vs.GuildID)
	if err != nil {
		return EventStatePayload{}, err
	}

	return res, nil
}

func (t *Controller) GetSelfUser() discordgo.User {
	return *t.dg.Session().State.User
}

func (t *Controller) GetPlaybackLog(
	guildID, ident, userID string,
	limit, offset int,
) ([]PlaybackLogEntry, error) {
	logs, err := t.db.GetPlaybackLog(guildID, ident, userID, limit, offset)
	if err == database.ErrNotFound {
		err = nil
	}

	return logs, err
}

func (t *Controller) GetPlaybackStats(
	guildID, userID string,
) ([]PlaybackStats, error) {
	logs, err := t.db.GetPlaybackLog(guildID, "", userID, 0, 0)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}

	statsMap := make(map[string]int)
	for _, log := range logs {
		statsMap[log.Ident]++
	}

	counts := make([]PlaybackStats, 0, len(statsMap))
	for ident, count := range statsMap {
		counts = append(counts, PlaybackStats{
			Ident: ident,
			Count: count,
		})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	return counts, nil
}

func (t *Controller) GetState() (StateStats, error) {
	var state StateStats

	sounds, err := t.db.GetSounds()
	if err != nil && err != database.ErrNotFound {
		return StateStats{}, err
	}
	state.NSoudns = len(sounds)

	state.NPlays, err = t.db.GetPlaybackLogSize()
	if err != nil && err != database.ErrNotFound {
		return StateStats{}, err
	}

	return state, nil
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

func (t *Controller) listSoundsFiltered(tagsMust []string, tagsNot []string) ([]Sound, error) {
	sounds, err := t.db.GetSounds()
	if err == database.ErrNotFound {
		sounds = []Sound{}
	} else if err != nil {
		return nil, err
	}

	if len(tagsMust) != 0 || len(tagsNot) != 0 {
		newSounds := make([]Sound, 0, len(sounds))
		for _, sound := range sounds {
			if !util.ContainsAll(sound.Tags, tagsMust) || util.ContainsAny(sound.Tags, tagsNot) {
				continue
			}
			newSounds = append(newSounds, sound)
		}
		sounds = newSounds
	}

	return sounds, nil
}

func (t *Controller) resizeHistoryBuffer() error {
	sounds, err := t.db.GetSounds()
	if err != nil {
		return err
	}

	nSounds := len(sounds)

	nBuff := nSounds / 5
	if nBuff > 10 {
		nBuff = 10
	}

	t.history.Resize(nBuff)
	logrus.WithField("size", nBuff).Debug("Resized history buffer")

	return nil
}

func (t *Controller) play(vs discordgo.VoiceState, ident string) error {
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

	fmt.Println(ident)
	identLower := strings.ToLower(ident)
	if strings.HasPrefix(identLower, "https://") {
		err = t.pl.Play(vs.GuildID, vs.ChannelID, ident, ident)
	} else {
		err = t.pl.PlaySound(vs.GuildID, vs.ChannelID, ident)
	}

	if err != nil {
		return err
	}

	return t.db.PutPlaybackLog(PlaybackLogEntry{
		Id:        xid.New().String(),
		Ident:     ident,
		GuildID:   vs.GuildID,
		UserID:    vs.UserID,
		Timestamp: time.Now(),
	})
}

func (t *Controller) execFastTrigger(guildID, userID string) {
	ident, err := t.db.GetUserFastTrigger(userID)
	if err != nil && err != database.ErrNotFound {
		logrus.WithError(err).WithFields(logrus.Fields{
			"guildid": guildID,
			"userid":  userID,
		}).Error("Getting fast trigger setting failed")
		return
	}
	if ident == "" {
		return
	}

	if strings.ToLower(ident) == "random" {
		err = t.PlayRandom(userID, nil, nil)
	} else {
		err = t.Play(userID, ident)
	}
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"guildid": guildID,
			"userid":  userID,
			"ident":   ident,
		}).Error("Playing fast trigegr sound failed")
	}
}

func (t *Controller) publishToGuildUsers(guildID string, e Event[any]) error {
	users, err := t.dg.UsersInGuildVoice(guildID)
	if err != nil {
		return err
	}
	t.Publish(ControllerEvent{
		Receivers: users,
		Event:     e,
	})
	return nil
}

func (t *Controller) playerEventHandler(e player.Event) {
	switch e.Type {
	case player.EventFastTrigger:
		t.execFastTrigger(e.GuildID, e.UserID)

	case player.EventVoiceJoin:
		ep, err := t.getVoiceJoinPayload(e.GuildID)
		if err != nil {
			return
		}
		t.publishToGuildUsers(e.GuildID, Event[any]{
			Type:    string(e.Type),
			Origin:  EventSenderPlayer,
			Payload: ep,
		})

	case player.EventVoiceInit:
		ep, err := t.getVoiceJoinPayload(e.GuildID)
		if err != nil {
			return
		}
		t.Publish(ControllerEvent{
			Receivers: []string{e.UserID},
			Event: Event[any]{
				Type:    string(e.Type),
				Origin:  EventSenderPlayer,
				Payload: ep,
			},
		})
	case player.EventVoiceDeinit:
		t.Publish(ControllerEvent{
			Receivers: []string{e.UserID},
			Event: Event[any]{
				Type:   string(e.Type),
				Origin: EventSenderPlayer,
			},
		})

	default:
		t.publishToGuildUsers(e.GuildID, Event[any]{
			Type:    string(e.Type),
			Origin:  EventSenderPlayer,
			Payload: e,
		})
	}
}

func (t *Controller) getVoiceJoinPayload(guildID string) (EventVoiceJoinPayload, error) {
	var (
		e   EventVoiceJoinPayload
		err error
	)

	guild, err := t.dg.GetGuild(guildID)
	if err != nil {
		return EventVoiceJoinPayload{}, err
	}
	e.Guild.ID = guild.ID
	e.Guild.Name = guild.Name
	e.Guild.IconUrl = guild.IconURL()

	e.Volume, err = t.db.GetGuildVolume(guildID)
	if err == database.ErrNotFound {
		err = nil
		e.Volume = 50
	}
	if err != nil {
		return EventVoiceJoinPayload{}, err
	}
	e.Filters, err = t.db.GetGuildFilters(guildID)
	if err != nil && err != database.ErrNotFound {
		return EventVoiceJoinPayload{}, err
	}

	return e, nil
}
