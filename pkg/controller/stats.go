package controller

import (
	"sort"

	"github.com/zekrotja/yuri69/pkg/database"
	. "github.com/zekrotja/yuri69/pkg/models"
)

func (t *Controller) GetPlaybackLog(
	guildID, ident, userID string,
	limit, offset int,
) ([]PlaybackLogEntry, error) {
	logs, err := t.db.GetPlaybackLog(guildID, ident, userID, limit, offset)
	if err == database.ErrNotFound {
		err = nil
	}

	return logs, err
}

func (t *Controller) GetPlaybackStats(
	guildID, userID string,
) ([]PlaybackStats, error) {
	logs, err := t.db.GetPlaybackLog(guildID, "", userID, 0, 0)
	if err != nil && err != database.ErrNotFound {
		return nil, err
	}

	statsMap := make(map[string]int)
	for _, log := range logs {
		statsMap[log.Ident]++
	}

	counts := make([]PlaybackStats, 0, len(statsMap))
	for ident, count := range statsMap {
		counts = append(counts, PlaybackStats{
			Ident: ident,
			Count: count,
		})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	return counts, nil
}

func (t *Controller) GetState() (StateStats, error) {
	var state StateStats

	sounds, err := t.db.GetSounds()
	if err != nil && err != database.ErrNotFound {
		return StateStats{}, err
	}
	state.NSoudns = len(sounds)

	state.NPlays, err = t.db.GetPlaybackLogSize()
	if err != nil && err != database.ErrNotFound {
		return StateStats{}, err
	}

	return state, nil
}
