package controller

import (
	"github.com/zekrotja/yuri69/pkg/database"
)

func (t *Controller) GetFastTrigger(userID string) (string, error) {
	ident, err := t.db.GetUserFastTrigger(userID)
	if err == database.ErrNotFound {
		err = nil
	}
	return ident, err
}

func (t *Controller) SetFastTrigger(userID, ident string) error {
	return t.db.SetUserFastTrigger(userID, ident)
}

func (t *Controller) GetFavorites(userID string) ([]string, error) {
	favs, err := t.db.GetFavorites(userID)
	if err == database.ErrNotFound {
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
