package discord

import (
	"github.com/bwmarrin/discordgo"
)

type Discord struct {
	session *discordgo.Session
}

func New(c DiscordConfig) (*Discord, error) {
	var t Discord
	var err error

	t.session, err = discordgo.New("Bot " + c.Token)
	t.session.State.TrackVoice = true
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Discord) Session() *discordgo.Session {
	return t.session
}

func (t *Discord) Open() error {
	cReady := make(chan struct{})

	t.session.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.Ready) {
		cReady <- struct{}{}
	})

	err := t.session.Open()
	if err != nil {
		return err
	}

	<-cReady

	return nil
}

func (t *Discord) Close() error {
	return t.session.Close()
}

func (t *Discord) FindUserVS(userID string) (discordgo.VoiceState, bool) {
	for _, g := range t.session.State.Guilds {
		for _, vs := range g.VoiceStates {
			for vs.UserID == userID {
				return *vs, true
			}
		}
	}
	return discordgo.VoiceState{}, false
}

func (t *Discord) UsersInGuildVoice(guildID string) ([]string, error) {
	g, err := t.session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(g.VoiceStates))
	for _, vs := range g.VoiceStates {
		if vs.UserID != t.session.State.User.ID {
			userIDs = append(userIDs, vs.UserID)
		}
	}

	return userIDs, nil
}

func (t *Discord) GetGuild(id string) (discordgo.Guild, error) {
	guild, err := t.session.State.Guild(id)
	if err == nil {
		return *guild, nil
	}

	guild, err = t.session.Guild(id)
	return *guild, err
}
