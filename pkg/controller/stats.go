package controller

import (
	"github.com/zekrotja/yuri69/pkg/database/dberrors"
	. "github.com/zekrotja/yuri69/pkg/models"
)

func (t *Controller) GetPlaybackLog(
	guildID, ident, userID string,
	limit, offset int,
) ([]PlaybackLogEntry, error) {
	logs, err := t.db.GetPlaybackLog(guildID, ident, userID, limit, offset)
	if err == dberrors.ErrNotFound {
		err = nil
	}

	return logs, err
}

func (t *Controller) GetPlaybackStats(
	guildID, userID string,
) ([]PlaybackStats, error) {
	return t.db.GetPlaybackStats(guildID, userID)
}

func (t *Controller) GetState() (StateStats, error) {
	var state StateStats

	sounds, err := t.db.GetSounds()
	if err != nil && err != dberrors.ErrNotFound {
		return StateStats{}, err
	}
	state.NSoudns = len(sounds)

	state.NPlays, err = t.db.GetPlaybackLogSize()
	if err != nil && err != dberrors.ErrNotFound {
		return StateStats{}, err
	}

	return state, nil
}
