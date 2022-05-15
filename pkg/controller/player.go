package controller

import (
	"math/rand"

	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/util"
)

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
