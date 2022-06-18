package controller

import (
	"errors"
	"math/rand"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/generic"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/player"
	"github.com/zekrotja/yuri69/pkg/static"
	"github.com/zekrotja/yuri69/pkg/storage"
	"github.com/zekrotja/yuri69/pkg/util"
)

var (
	reservedUids = []string{"random", "upload", "create", "downloadall"}
)

type ControllerEvent struct {
	IsBroadcast bool
	Receivers   []string
	Event       Event[any]
}

type Controller struct {
	*util.EventBus[ControllerEvent]

	ownerID string
	db      database.IDatabase
	st      storage.IStorage
	pl      *player.Player
	dg      *discord.Discord

	ffmpegExec string

	pendingCrations *timedmap.TimedMap[string, string]
	history         *generic.RingQueue[string]
}

func New(
	db database.IDatabase,
	st storage.IStorage,
	pl *player.Player,
	dg *discord.Discord,
	ownerID string,
) (*Controller, error) {

	var (
		t   Controller
		err error
	)

	rand.Seed(time.Now().UnixNano())

	t.EventBus = util.NewEventBus[ControllerEvent]()

	t.ownerID = ownerID
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

func (t *Controller) GetCurrentState(userID string) (EventStatePayload, error) {
	var (
		res EventStatePayload
		err error
	)

	res.IsAdmin = t.CheckAdmin(userID) == nil

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
