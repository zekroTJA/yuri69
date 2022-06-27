package controller

import (
	"github.com/zekrotja/yuri69/pkg/database/dberrors"
	"github.com/zekrotja/yuri69/pkg/errs"
	. "github.com/zekrotja/yuri69/pkg/models"
)

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
