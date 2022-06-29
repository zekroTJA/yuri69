package controller

import (
	"github.com/zekrotja/yuri69/pkg/cryptoutil"
	"github.com/zekrotja/yuri69/pkg/database/dberrors"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
)

func (t *Controller) GetFastTrigger(userID string) (string, error) {
	ident, err := t.db.GetUserFastTrigger(userID)
	if err == dberrors.ErrNotFound {
		err = nil
	}
	return ident, err
}

func (t *Controller) SetFastTrigger(userID, ident string) error {
	return t.db.SetUserFastTrigger(userID, ident)
}

func (t *Controller) GetFavorites(userID string) ([]string, error) {
	favs, err := t.db.GetFavorites(userID)
	if err == dberrors.ErrNotFound {
		err = nil
	}
	if favs == nil {
		favs = make([]string, 0)
	}
	return favs, err
}

func (t *Controller) AddFavorite(userID, ident string) error {
	return t.db.AddFavorite(userID, ident)
}

func (t *Controller) RemoveFavorite(userID, ident string) error {
	return t.db.RemoveFavorite(userID, ident)
}

func (t *Controller) GetApiKey(userID string) (string, error) {
	return t.db.GetApiKey(userID)
}

func (t *Controller) GenerateApiKey(userID string) (string, error) {
	err := t.db.RemoveApiKey(userID)
	if err != nil && err != dberrors.ErrNotFound {
		return "", err
	}

	token, err := cryptoutil.GetRandBase64Str(32)
	if err != nil {
		return "", err
	}

	err = t.db.SetApiKey(userID, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (t *Controller) RemoveApiKey(userID string) error {
	return t.db.RemoveApiKey(userID)
}

func (t *Controller) GetUserByApiKey(token string) (string, error) {
	return t.db.GetUserByApiKey(token)
}

func (t *Controller) GetTwitchState(userid string) (TwitchState, error) {
	setting, err := t.db.GetTwitchSettings(userid)
	if err != nil && err != dberrors.ErrNotFound {
		return TwitchState{}, err
	}

	var state TwitchState
	state.TwitchSettings = setting
	state.Connected = setting.TwitchUserName != "" && t.tw.Joined(setting.TwitchUserName)

	return state, nil
}

func (t *Controller) UpdateTwitchSettings(userid string, setting *TwitchSettings, join bool) error {
	curr, err := t.GetTwitchState(userid)
	if err != nil {
		return err
	}

	if setting == nil {
		setting = &curr.TwitchSettings
	} else if curr.Connected && setting.TwitchUserName != curr.TwitchUserName {
		return errs.WrapUserError("twitch user name can not be changed whilest connected")
	}

	setting.UserID = userid
	err = t.db.SetTwitchSettings(*setting)
	if err != nil {
		return err
	}

	if join {
		if setting.TwitchUserName == "" {
			return errs.WrapUserError("unable to join: no twitch user name specified")
		}
		if err = t.tw.Join(userid, *setting); err != nil {
			return err
		}
	} else {
		t.tw.Update(*setting)
	}

	return nil
}

func (t *Controller) LeaveTwitch(userid string) error {
	setting, err := t.db.GetTwitchSettings(userid)
	if err != nil && err != dberrors.ErrNotFound {
		return errs.WrapUserError("not connected")
	}

	return t.tw.Leave(setting.TwitchUserName)
}
