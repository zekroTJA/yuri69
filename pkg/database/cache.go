package database

import (
	"strings"
	"sync"
)

const (
	cacheKeySeparator = ":"
)

type DatabaseCache struct {
	IDatabase

	cache sync.Map
}

var _ IDatabase = (*DatabaseCache)(nil)

func WrapCache(db IDatabase, err error) (IDatabase, error) {
	if err != nil {
		return nil, err
	}

	var t DatabaseCache
	t.IDatabase = db

	return &t, nil
}

func (t *DatabaseCache) SetGuildVolume(guildID string, volume int) error {
	t.cache.Store(ckey("guilds", guildID, "volume"), volume)
	return t.IDatabase.SetGuildVolume(guildID, volume)
}

func (t *DatabaseCache) GetGuildVolume(guildID string) (int, error) {
	var err error

	vi, _ := t.cache.Load(ckey("guilds", guildID, "volume"))
	v, ok := vi.(int)
	if !ok {
		v, err = t.GetGuildVolume(guildID)
	}

	return v, err
}

// --- Felpers ---

func ckey(elements ...string) string {
	return strings.Join(elements, cacheKeySeparator)
}
