package disgolink

import (
	"context"
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/zekrotja/yuri69/pkg/discord"
)

type Disgolink struct {
	dc     *discord.Discord
	client disgolink.Client
	node   disgolink.Node
}

func New(c *DisgolinkConfig, dc *discord.Discord, listener disgolink.EventListener) (t *Disgolink, err error) {
	t = new(Disgolink)
	t.dc = dc

	t.client = disgolink.New(snowflake.MustParse(t.dc.Session().State.User.ID), disgolink.WithListeners(listener))

	t.node, err = t.client.AddNode(context.Background(), disgolink.NodeConfig{
		Name:      "main",
		Address:   c.Address,
		Password:  c.Password,
		Secure:    false,
		SessionID: "yuri69session",
	})
	if err != nil {
		return nil, err
	}

	dc.Session().AddHandler(t.handleVoiceStateUpdate)
	dc.Session().AddHandler(t.handleVoiceServerUpdate)

	return t, nil
}

func (t *Disgolink) Close() error {
	t.client.Close()
	return nil
}

func (t *Disgolink) Play(guildID string, ident string) (*lavalink.Track, error) {
	guildIdSnowflake, err := snowflake.Parse(guildID)
	if err != nil {
		return nil, err
	}

	type Result struct {
		track *lavalink.Track
		err   error
	}

	cRes := make(chan Result, 1)

	t.client.BestNode().LoadTracksHandler(context.Background(), ident, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			player := t.client.Player(guildIdSnowflake)
			err := player.Update(context.Background(), lavalink.WithTrack(track))
			if err != nil {
				cRes <- Result{err: err}
			} else {
				cRes <- Result{track: &track}
			}
		},
		func(playlist lavalink.Playlist) {
			cRes <- Result{err: errors.New("playlist not supported")}
		},
		func(tracks []lavalink.Track) {
			cRes <- Result{err: errors.New("multiple tracks not supported")}
		},
		func() {
			cRes <- Result{err: errors.New("no track found")}
		},
		func(err error) {
			cRes <- Result{err: err}
		}))

	r := <-cRes
	return r.track, r.err
}

func (t *Disgolink) Destroy(guildID string) error {
	guildIdSnowflake, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return t.client.Player(guildIdSnowflake).Destroy(context.Background())
}

func (t *Disgolink) Stop(guildID string) error {
	guildIdSnowflake, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return t.client.Player(guildIdSnowflake).Update(context.Background(), lavalink.WithNullTrack())
}

func (t *Disgolink) SetVolume(guildID string, volume int) error {
	guildIdSnowflake, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return t.client.Player(guildIdSnowflake).Update(context.Background(), lavalink.WithVolume(volume))
}

func (t *Disgolink) DecodeTrack(uid string) (*lavalink.Track, error) {
	return t.client.BestNode().DecodeTrack(context.Background(), uid)
}

func (t *Disgolink) handleVoiceServerUpdate(_ *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	guildIdSnowflake, err := snowflake.Parse(e.GuildID)
	if err != nil {
		// TODO: err handling
		return
	}
	t.client.OnVoiceServerUpdate(context.Background(), guildIdSnowflake, e.Token, e.Endpoint)
}

func (t *Disgolink) handleVoiceStateUpdate(_ *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	guildIdSnowflake, err := snowflake.Parse(e.GuildID)
	if err != nil {
		// TODO: err handling
		return
	}

	var channelIdSnowflake *snowflake.ID
	if e.ChannelID != "" {
		channelIdSnowflakeVal, err := snowflake.Parse(e.ChannelID)
		if err != nil {
			// TODO: err handling
			return
		}
		channelIdSnowflake = &channelIdSnowflakeVal
	}

	t.client.OnVoiceStateUpdate(context.Background(), guildIdSnowflake, channelIdSnowflake, e.SessionID)
}
