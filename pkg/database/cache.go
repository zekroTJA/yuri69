package database

import (
	"strings"
	"sync"

	. "github.com/zekrotja/yuri69/pkg/models"
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

func (t *DatabaseCache) GetSounds() ([]Sound, error) {
	var err error
	key := ckey("sounds")

	vi, _ := t.cache.Load(key)
	v, ok := vi.([]Sound)
	if !ok {
		v, err = t.IDatabase.GetSounds()
		if err != nil {
			return nil, err
		}
		t.cache.Store(key, v)
	}

	r := make([]Sound, len(v))
	copy(r, v)

	return r, nil
}

func (t *DatabaseCache) PutSound(sound Sound) error {
	t.cache.Delete("sounds")
	return t.IDatabase.PutSound(sound)
}

func (t *DatabaseCache) RemoveSound(uid string) error {
	t.cache.Delete("sound")
	return t.IDatabase.RemoveSound(uid)
}

func (t *DatabaseCache) GetGuildVolume(guildID string) (int, error) {
	var err error
	key := ckey("guilds", guildID, "volume")

	vi, _ := t.cache.Load(key)
	v, ok := vi.(int)
	if !ok {
		v, err = t.IDatabase.GetGuildVolume(guildID)
		if err != nil {
			return 0, err
		}
		t.cache.Store(key, v)
	}

	return v, nil
}

func (t *DatabaseCache) SetGuildVolume(guildID string, volume int) error {
	t.cache.Store(ckey("guilds", guildID, "volume"), volume)
	return t.IDatabase.SetGuildVolume(guildID, volume)
}

func (t *DatabaseCache) GetUserFastTrigger(userID string) (string, error) {
	var err error
	key := ckey("users", userID, "fasttrigger")

	vi, _ := t.cache.Load(key)
	v, ok := vi.(string)
	if !ok {
		v, err = t.IDatabase.GetUserFastTrigger(userID)
		if err != nil {
			return "", err
		}
		t.cache.Store(key, v)
	}

	return v, nil
}

func (t *DatabaseCache) SetUserFastTrigger(userID, ident string) error {
	t.cache.Store(ckey("users", userID, "fasttrigger"), ident)
	return t.IDatabase.SetUserFastTrigger(userID, ident)
}

func (t *DatabaseCache) GetGuildFilters(guildID string) (GuildFilters, error) {
	var err error
	key := ckey("guilds", guildID, "filters")

	vi, _ := t.cache.Load(key)
	v, ok := vi.(GuildFilters)
	if !ok {
		v, err = t.IDatabase.GetGuildFilters(guildID)
		if err != nil {
			return GuildFilters{}, err
		}
		t.cache.Store(key, v)
	}

	return v, nil
}

func (t *DatabaseCache) SetGuildFilters(guildID string, f GuildFilters) error {
	t.cache.Store(ckey("guilds", guildID, "filters"), f)
	return t.IDatabase.SetGuildFilters(guildID, f)
}

// --- Felpers ---

func ckey(elements ...string) string {
	return strings.Join(elements, cacheKeySeparator)
}
