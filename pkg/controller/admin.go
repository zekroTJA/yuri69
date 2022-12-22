package controller

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/yuri69/pkg/database/dberrors"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
)

func (t *Controller) GetAdmins(executorID string) ([]User, error) {
	if err := t.CheckAdmin(executorID); err != nil {
		return nil, err
	}

	adminIDs, err := t.db.GetAdmins()
	if err != nil && err != dberrors.ErrNotFound {
		return nil, err
	}

	admins := make([]User, 0, len(adminIDs)+1)

	owner, err := t.dg.Session().User(t.ownerID)
	if err != nil {
		return nil, err
	}
	uOwner := UserFromUser(*owner)
	uOwner.IsOwner = true
	admins = append(admins, uOwner)

	for _, id := range adminIDs {
		u, err := t.dg.Session().User(id)
		if err != nil {
			u = &discordgo.User{
				ID: id,
			}
		}
		admins = append(admins, UserFromUser(*u))
	}

	return admins, nil
}

func (t *Controller) SetAdmin(executorID, userID string) (User, error) {
	if err := t.CheckAdmin(executorID); err != nil {
		return User{}, err
	}

	user, err := t.dg.Session().User(userID)
	if err != nil {
		return User{}, errs.WrapUserError("user with this ID does not exist", http.StatusBadRequest)
	}

	if userID != t.ownerID {
		err = t.db.AddAdmin(userID)
		if err != nil {
			return User{}, err
		}
	}

	uUser := UserFromUser(*user)
	return uUser, nil
}

func (t *Controller) RemoveAdmin(executorID, userID string) error {
	if err := t.CheckAdmin(executorID); err != nil {
		return err
	}

	if userID == t.ownerID {
		return errs.WrapUserError("owner can not be removed from admins")
	}

	return t.db.RemoveAdmin(userID)
}

func (t *Controller) CheckAdmin(userID string) error {
	if userID == t.ownerID {
		return nil
	}

	ok, err := t.db.IsAdmin(userID)
	if err != nil {
		return err
	}

	if !ok {
		return errs.WrapUserError("admin privileges required", http.StatusForbidden)
	}

	return nil
}

func (t *Controller) GetGuilds(userID string) ([]*discordgo.Guild, error) {
	if err := t.CheckAdmin(userID); err != nil {
		return nil, err
	}

	return t.dg.Session().State.Guilds, nil
}

func (t *Controller) RemoveGuild(userID, guildID string) error {
	if err := t.CheckAdmin(userID); err != nil {
		return err
	}

	return t.dg.Session().GuildLeave(guildID)
}
