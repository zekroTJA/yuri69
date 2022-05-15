package controller

import (
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
)

func (t *Controller) GetGuildFilter(userID string) (GuildFilters, error) {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return GuildFilters{},
			errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	f, err := t.db.GetGuildFilters(vs.GuildID)
	if err == database.ErrNotFound {
		err = nil
	}

	return f, err
}

func (t *Controller) SetGuildFilter(userID string, f GuildFilters) error {
	vs, ok := t.dg.FindUserVS(userID)
	if !ok {
		return errs.WrapUserError("you need to be in a voice channel to perform this action")
	}

	f.Sanitize()
	err := f.Check()
	if err != nil {
		return err
	}

	err = t.db.SetGuildFilters(vs.GuildID, f)
	if err != nil {
		return err
	}

	return t.publishToGuildUsers(vs.GuildID, Event[any]{
		Type:    EventGuildFilterUpdated,
		Origin:  EventSenderController,
		Payload: f,
	})
}
