package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/yuri69/pkg/util"
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

	t.session.AddHandlerOnce(func(s *discordgo.Session, e *discordgo.Ready) {
		cReady <- struct{}{}
	})

	err := t.session.Open()
	if err != nil {
		return err
	}

	<-cReady

	rotateStatus(t.session, 1*time.Minute)

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

func (t *Discord) GetUser(id string) (*discordgo.User, error) {
	return t.session.User(id)
}

func (t *Discord) HasSharedGuild(userID string) (bool, error) {
	for _, guild := range t.session.State.Guilds {
		member, err := t.getMember(guild.ID, userID)
		if err != nil && !util.IsErrCode(err, discordgo.ErrCodeUnknownMember) {
			return false, err
		}
		if member != nil {
			return true, nil
		}
	}

	return false, nil
}

func (t *Discord) getMember(guildID, userID string) (*discordgo.Member, error) {
	member, err := t.session.State.Member(guildID, userID)
	if err != nil && err != discordgo.ErrStateNotFound {
		return nil, err
	}

	member, err = t.session.GuildMember(guildID, userID)
	if err != nil {
		return nil, err
	}

	t.session.State.MemberAdd(member)
	return member, nil
}
