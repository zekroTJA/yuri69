package player

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/lukasl-dev/waterlink/v2/event"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/generic"
	"github.com/zekrotja/yuri69/pkg/lavalink"
	"github.com/zekrotja/yuri69/pkg/static"
	"github.com/zekrotja/yuri69/pkg/storage"
	"github.com/zekrotja/yuri69/pkg/util"
)

type Player struct {
	*util.EventBus[Event]

	dc *discord.Discord
	st storage.IStorage
	ll *lavalink.Lavalink

	hostname string

	vcs     generic.SyncMap[string, voiceConnection]
	waiters util.Waiters[string]

	router *routing.Router
	server *http.Server

	autoLeaveTimer *time.Timer
	trackCache     *timedmap.TimedMap[string, string]
}

type voiceConnection struct {
	GuildID   string
	ChannelID string
}

func NewPlayer(
	c PlayerConfig,
	dc *discord.Discord,
	st storage.IStorage,
) (*Player, error) {

	var (
		t   Player
		err error
	)

	t.hostname = c.Hostname
	if t.hostname == "" {
		var err error
		t.hostname, err = os.Hostname()
		if err != nil {
			return nil, err
		}
	}

	t.EventBus = util.NewEventBus[Event](100)
	t.trackCache = timedmap.New[string, string](5 * time.Minute)

	t.dc = dc
	t.st = st

	t.ll, err = lavalink.New(c.Lavalink, dc, t.handleEvent)
	if err != nil {
		return nil, err
	}

	t.router = routing.New()
	t.router.Get("/file/<id>", t.handleGetFile)

	t.dc.Session().AddHandler(t.handleVoiceUpdate)

	t.SubscribeFunc(func(e Event) {
		entry := logrus.WithFields(logrus.Fields{
			"type":  e.Type,
			"ident": e.Ident,
			"guild": e.GuildID,
		})

		if e.Err != nil {
			entry.WithError(err).Error("Lavalink Event")
		} else {
			entry.Debug("Lavalink Event")
		}
	})

	return &t, nil
}

func (t *Player) ListenAndServeBlocking() error {
	t.server = &http.Server{
		Addr:    "0.0.0.0:6969",
		Handler: t.router,
	}
	return t.server.ListenAndServe()
}

func (t *Player) Init(guildID, channelID string) error {
	vc, ok := t.vcs.Load(guildID)
	if !ok || vc.ChannelID != channelID {
		err := t.dc.Session().ChannelVoiceJoinManual(guildID, channelID, false, true)
		if err != nil {
			return err
		}
		t.waiters.CreateAndWait(guildID)
	}
	return nil
}

func (t *Player) PlaySound(guildID, channelID, ident string) error {
	return t.Play(guildID, channelID,
		fmt.Sprintf("http://%s:6969/file/%s", t.hostname, ident), ident)
}

func (t *Player) Play(guildID, channelID, url, ident string) error {
	track, err := t.ll.Play(guildID, url)
	if track.ID != "" {
		if ident == "" {
			ident = url
		}
		t.trackCache.Set(track.Info.URI, ident,
			time.Duration(track.Info.Length)*time.Microsecond+30*time.Second)
	}
	return err
}

func (t *Player) Destroy(guildID string) error {
	_, ok := t.vcs.Load(guildID)
	if !ok {
		return ErrNoGuildPlayer
	}

	return t.dc.Session().ChannelVoiceJoinManual(guildID, "", false, true)
}

func (t *Player) Stop(guildID string) error {
	_, ok := t.vcs.Load(guildID)
	if !ok {
		return ErrNoGuildPlayer
	}

	return t.ll.Stop(guildID)
}

func (t *Player) SetVolume(guildID string, volume uint16) error {
	_, ok := t.vcs.Load(guildID)
	if !ok {
		return ErrNoGuildPlayer
	}

	return t.ll.SetVolume(guildID, volume)
}

func (t *Player) Close() error {
	return t.ll.Close()
}

// --- Internal stuff ---

func (t *Player) handleGetFile(ctx *routing.Context) error {
	id := ctx.Param("id")

	r, _, err := t.st.GetObject(static.BucketSounds, id)
	if err != nil {
		return ctx.WriteWithStatus("", http.StatusBadRequest)
	}
	defer r.Close()

	_, err = io.Copy(ctx.Response, r)
	if err != nil {
		logrus.WithError(err).Error("Writing file to response body failed")
		return ctx.WriteWithStatus("", http.StatusInternalServerError)
	}

	ctx.Response.WriteHeader(http.StatusOK)
	return nil
}

func (t *Player) handleVoiceUpdate(_ *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	if e.BeforeUpdate == nil && e.ChannelID != "" {
		t.onVoiceJoin(e)
	} else if e.BeforeUpdate != nil && e.ChannelID != "" && e.BeforeUpdate.ChannelID != e.ChannelID {
		t.onVoiceMove(e)
	} else if e.BeforeUpdate != nil && e.ChannelID == "" {
		t.onVoiceLeave(e)
	}
}

func (t *Player) initVs(guildID, channelID string) {
	t.vcs.Store(guildID, voiceConnection{
		GuildID:   guildID,
		ChannelID: channelID,
	})
	t.waiters.BroadcastAndRemove(guildID)
}

func (t *Player) onVoiceJoin(e *discordgo.VoiceStateUpdate) {
	if e.UserID == t.dc.Session().State.User.ID {
		t.initVs(e.GuildID, e.ChannelID)
		logrus.
			WithField("guildID", e.GuildID).
			WithField("chanID", e.ChannelID).
			Debug("Voice state created")
	} else {
		t.cancelAutoLeave(e.VoiceState)
	}
}

func (t *Player) onVoiceMove(e *discordgo.VoiceStateUpdate) {
	if e.UserID == t.dc.Session().State.User.ID {
		t.initVs(e.GuildID, e.ChannelID)
		logrus.
			WithField("guildID", e.GuildID).
			WithField("chanID", e.ChannelID).
			Debug("Voice state updated")
		t.autoLeave(e.VoiceState)
		t.cancelAutoLeave(e.VoiceState)
	} else {
		t.autoLeave(e.BeforeUpdate)
		t.cancelAutoLeave(e.VoiceState)
	}
}

func (t *Player) onVoiceLeave(e *discordgo.VoiceStateUpdate) {
	if e.UserID == t.dc.Session().State.User.ID {
		t.vcs.Delete(e.GuildID)
		t.ll.Destroy(e.GuildID)
		logrus.
			WithField("guildID", e.GuildID).
			WithField("chanID", e.ChannelID).
			Debug("Voice state removed")
	} else {
		t.autoLeave(e.BeforeUpdate)
	}
}

func (t *Player) autoLeave(e *discordgo.VoiceState) {
	if t.autoLeaveTimer == nil {
		vc, ok := t.vcs.Load(e.GuildID)
		if !ok || vc.ChannelID != e.ChannelID {
			return
		}
		vConns, err := t.getChannelVoiceConnections(e.GuildID, e.ChannelID)
		if err != nil {
			return
		}
		if vConns == 0 {
			logrus.WithField("guildID", e.GuildID).Debug("Trigger autoleave timer")
			t.autoLeaveTimer = time.AfterFunc(5*time.Second, func() {
				t.dc.Session().ChannelVoiceJoinManual(e.GuildID, "", false, true)
			})
		}
	}
}

func (t *Player) cancelAutoLeave(e *discordgo.VoiceState) {
	if t.autoLeaveTimer != nil {
		vc, ok := t.vcs.Load(e.GuildID)
		if !ok || vc.ChannelID != e.ChannelID {
			return
		}
		vConns, err := t.getChannelVoiceConnections(e.GuildID, e.ChannelID)
		if err != nil {
			return
		}
		if vConns != 0 {
			logrus.WithField("guildID", e.GuildID).Debug("Clear autoleave timer")
			t.autoLeaveTimer.Stop()
			t.autoLeaveTimer = nil
		}
	}
}

func (t *Player) getChannelVoiceConnections(guildID, channelID string) (int, error) {
	guild, err := t.dc.Session().State.Guild(guildID)
	if err != nil {
		return 0, err
	}

	selfID := t.dc.Session().State.User.ID
	n := 0
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == channelID && vs.UserID != selfID {
			n++
		}
	}

	return n, nil
}

func (t *Player) getIdentFromCache(trackID string) string {
	info, err := t.ll.DecodeTrackId(trackID)
	if err != nil {
		logrus.WithError(err).
			WithField("trackID", trackID).
			Error("Failed decoding track id")
		return ""
	}

	return t.trackCache.GetValue(info.URI)
}

func (t *Player) handleEvent(e any) {
	switch et := e.(type) {
	// case event.PlayerUpdate:
	// case event.Stats:

	case event.WebSocketClosed:
		logrus.Warn("Lavalink websokcet connection closed. Trying reconnect ...")
		// TODO: Try reconnect

	case event.TrackException:
		ident := t.getIdentFromCache(et.TrackID)
		if ident == "" {
			return
		}
		t.Publish(Event{
			Type:    EventPlayException,
			Ident:   ident,
			GuildID: et.GuildID.String(),
			Err:     errors.New(et.Error),
		})
	case event.TrackStuck:
		ident := t.getIdentFromCache(et.TrackID)
		if ident == "" {
			return
		}
		t.Publish(Event{
			Type:    EventPlayStuck,
			Ident:   ident,
			GuildID: et.GuildID.String(),
		})
	case event.TrackStart:
		ident := t.getIdentFromCache(et.TrackID)
		if ident == "" {
			return
		}
		t.Publish(Event{
			Type:    EventPlayStart,
			Ident:   ident,
			GuildID: et.GuildID.String(),
		})
	case event.TrackEnd:
		ident := t.getIdentFromCache(et.TrackID)
		if ident == "" {
			return
		}
		t.Publish(Event{
			Type:    EventPlayEnd,
			Ident:   ident,
			GuildID: et.GuildID.String(),
		})
	}
}
