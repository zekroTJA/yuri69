package discord

import "github.com/bwmarrin/discordgo"

type Discord struct {
	session *discordgo.Session
}

func New(c DiscordConfig) (*Discord, error) {
	var t Discord
	var err error

	t.session, err = discordgo.New("Bot " + c.Token)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Discord) Session() *discordgo.Session {
	return t.session
}

func (t *Discord) Open() error {
	return t.session.Open()
}

func (t *Discord) Close() error {
	return t.session.Close()
}
