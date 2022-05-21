package controller

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/player"
	"github.com/zekrotja/yuri69/pkg/util"
)

var extMappings = map[string]string{
	".oga": ".ogg",
	".ogv": ".ogg",
	".ogx": ".ogg",
}

func mapExt(ext string) string {
	v := extMappings[ext]
	if v == "" {
		return ext
	}
	return v
}

func (t *Controller) isAdmin(userID string) (bool, error) {
	if userID == t.ownerID {
		return true, nil
	}

	return t.db.IsAdmin(userID)
}

func (t *Controller) ffmpeg(
	in io.Reader,
	inTyp string,
	out io.Writer,
	outTyp string,
	args ...string,
) error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, "-f", inTyp, "-i", "pipe:", "-map", "0:a:0")
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
	if err != nil && err != database.ErrNotFound {
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
	filters, err := t.db.GetGuildFilters(vs.GuildID)
	if err != nil && err != database.ErrNotFound {
		return err
	}

	if len(filters.Exclude) > 0 {
		sound, err := t.db.GetSound(ident)
		if err != nil {
			return err
		}

		if util.ContainsAny(filters.Exclude, sound.Tags) {
			return errs.WrapUserError("you are not allowed to paly excluded sounds")
		}
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
