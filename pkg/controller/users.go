package controller

import (
	"github.com/zekrotja/yuri69/pkg/cryptoutil"
	"github.com/zekrotja/yuri69/pkg/database/dberrors"
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
