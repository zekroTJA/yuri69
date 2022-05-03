package lavalink

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/waterlink/v2"
	"github.com/lukasl-dev/waterlink/v2/track/query"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/discord"
)

type Lavalink struct {
	dc     *discord.Discord
	client *waterlink.Client
	conn   *waterlink.Connection
}

func New(c LavalinkConfig, dc *discord.Discord) (*Lavalink, error) {
	var t Lavalink
	var err error

	t.dc = dc

	creds := waterlink.Credentials{
		Authorization: c.Password,
		UserID:        snowflake.MustParse(t.dc.Session().State.User.ID),
		ResumeKey:     "yuri69session",
	}
	opts := waterlink.ConnectionOptions{
		HandleEventError: func(err error) {
			logrus.WithError(err).Error("Lavalink error")
		},
	}

	t.client, err = waterlink.NewClient(fmt.Sprintf("http://%s", c.Address), creds)
	if err != nil {
		return nil, err
	}

	t.conn, err = waterlink.Open(fmt.Sprintf("ws://%s", c.Address), creds, opts)
	if err != nil {
		return nil, err
	}

	t.dc.Session().AddHandler(t.handleVoiceServerUpdate)

	return &t, nil
}

func (t *Lavalink) Close() error {
	return t.conn.Close()
}

func (t *Lavalink) Play(guildID, ident string) error {
	tracks, err := t.client.LoadTracks(query.Of(ident))
	if err != nil {
		return err
	}

	logrus.
		WithField("type", tracks.LoadType).
		WithField("n", len(tracks.Tracks)).
		Debug("Tracks loaded")

	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	if len(tracks.Tracks) == 0 {
		return errors.New("no tracks have been loaded")
	}

	return t.conn.Guild(sf).PlayTrack(tracks.Tracks[0])
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
