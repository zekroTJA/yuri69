package lavalink

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/waterlink/v2"
	"github.com/lukasl-dev/waterlink/v2/track"
	"github.com/lukasl-dev/waterlink/v2/track/query"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/discord"
)

type Lavalink struct {
	dc     *discord.Discord
	client *waterlink.Client
	conn   *waterlink.Connection

	address         string
	creds           waterlink.Credentials
	opts            waterlink.ConnectionOptions
	reconnectionTry int
}

func New(c LavalinkConfig, dc *discord.Discord, eventHandler func(any)) (*Lavalink, error) {
	var t Lavalink
	var err error

	t.address = c.Address
	t.dc = dc

	t.creds = waterlink.Credentials{
		Authorization: c.Password,
		UserID:        snowflake.MustParse(t.dc.Session().State.User.ID),
		ResumeKey:     "yuri69session",
	}
	t.opts = waterlink.ConnectionOptions{
		HandleEventError: t.handleErrors,
	}
	if eventHandler != nil {
		t.opts.EventHandler = waterlink.EventHandlerFunc(eventHandler)
	}

	t.client, err = waterlink.NewClient(fmt.Sprintf("http://%s", c.Address), t.creds)
	if err != nil {
		return nil, err
	}

	if err = t.Connect(); err != nil {
		return nil, err
	}

	t.dc.Session().AddHandler(t.handleVoiceServerUpdate)

	return &t, nil
}

func (t *Lavalink) Connect() error {
	if t.conn != nil && !t.conn.Closed() {
		return errors.New("connection already established")
	}
	var err error
	t.conn, err = waterlink.Open(fmt.Sprintf("ws://%s", t.address), t.creds, t.opts)
	return err
}

func (t *Lavalink) Close() error {
	return t.conn.Close()
}

func (t *Lavalink) Play(guildID, ident string) (track.Track, error) {
	tracks, err := t.client.LoadTracks(query.Of(ident))
	if err != nil {
		return track.Track{}, err
	}

	logrus.
		WithField("type", tracks.LoadType).
		WithField("n", len(tracks.Tracks)).
		Debug("Tracks loaded")

	if len(tracks.Tracks) == 0 {
		return track.Track{}, errors.New("no tracks have been loaded")
	}

	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return track.Track{}, err
	}

	track := tracks.Tracks[0]
	return track, t.conn.Guild(sf).PlayTrack(track)
}

func (t *Lavalink) Destroy(guildID string) error {
	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return t.conn.Guild(sf).Destroy()
}

func (t *Lavalink) Stop(guildID string) error {
	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return t.conn.Guild(sf).Stop()
}

func (t *Lavalink) SetVolume(guildID string, volume uint16) error {
	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return t.conn.Guild(sf).UpdateVolume(volume)
}

func (t *Lavalink) DecodeTrackId(uid string) (*track.Info, error) {
	return t.client.DecodeTrack(uid)
}

func (t *Lavalink) handleVoiceServerUpdate(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	logrus.
		WithField("guild", e.GuildID).
		WithField("sessionID", s.State.SessionID).
		Debugf("Update voice server: %+v", e)

	g := t.conn.Guild(snowflake.MustParse(e.GuildID))
	err := g.UpdateVoice(s.State.SessionID, e.Token, e.Endpoint)
	if err != nil {
		logrus.
			WithError(err).
			WithField("guild", e.GuildID).
			WithField("sessionID", s.State.SessionID).
			Error("Voice server update failed")
	}
}

func (t *Lavalink) handleErrors(err error) {
	if strings.Contains(err.Error(), "waterlink: connection: websocket: close") {
		logrus.WithError(err).Error("Lavalink connection closed")
		t.tryReconnecting()
		return
	}

	logrus.WithError(err).Error("Lavalink error")
}

func (t *Lavalink) tryReconnecting() {
	timeout := time.Duration(t.reconnectionTry)*500*time.Millisecond + time.Duration(rand.Intn(900)+100)*time.Millisecond

	if timeout > 30*time.Second {
		timeout = 30 * time.Second
	}

	logrus.
		WithField("try", t.reconnectionTry).
		WithField("timeout", timeout).
		Warn("Lavalink: Trying to reconnect ...")
	time.Sleep(timeout)

	t.reconnectionTry++
	err := t.Connect()
	if err != nil {
		logrus.WithError(err)
		t.tryReconnecting()
		return
	}

	t.reconnectionTry = 0
	logrus.Info("Lavalink connection re-established")
}
