package database

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/xujiajun/nutsdb"
	. "github.com/zekrotja/yuri69/pkg/models"
)

const (
	bucketSounds = "sounds"
	bucketGuilds = "guilds"
	bucketUsers  = "users"
	bucketStats  = "stats"

	keySeparator = ":"
)

type NutsConfig struct {
	Location string
}

type Nuts struct {
	db *nutsdb.DB
}

var _ (IDatabase) = (*Nuts)(nil)

func NewNuts(c NutsConfig) (*Nuts, error) {
	var (
		t   Nuts
		err error
	)

	opts := nutsdb.DefaultOptions
	opts.Dir = c.Location
	t.db, err = nutsdb.Open(opts)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Nuts) Close() error {
	return t.db.Close()
}

func (t *Nuts) PutSound(sound Sound) error {
	data, err := marshal(sound)
	if err != nil {
		return err
	}

	return t.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(bucketSounds, []byte(sound.Uid), data, 0)
	})
}

func (t *Nuts) RemoveSound(uid string) error {
	return t.db.Update(func(tx *nutsdb.Tx) error {
		err := tx.Delete(bucketSounds, []byte(uid))
		return t.wrapErr(err)
	})
}

func (t *Nuts) GetSounds() ([]Sound, error) {
	var entries nutsdb.Entries
	err := t.db.View(func(tx *nutsdb.Tx) error {
		var err error
		entries, err = tx.GetAll(bucketSounds)
		return t.wrapErr(err)
	})
	if err != nil {
		return nil, err
	}

	sounds := make([]Sound, 0, len(entries))
	for _, e := range entries {
		sound, err := unmarshal[Sound](e.Value)
		if err != nil {
			return nil, err
		}
		sounds = append(sounds, sound)
	}

	return sounds, nil
}

func (t *Nuts) GetSound(uid string) (Sound, error) {
	return getValue[Sound](t, bucketSounds, key(uid))
}

func (t *Nuts) GetGuildVolume(guildID string) (int, error) {
	return getValue[int](t, bucketGuilds, key(guildID, "volume"))
}

func (t *Nuts) SetGuildVolume(guildID string, volume int) error {
	return setValue(t, bucketGuilds, key(guildID, "volume"), volume)
}

func (t *Nuts) GetUserFastTrigger(userID string) (string, error) {
	return getValue[string](t, bucketUsers, key(userID, "fasttrigger"))
}

func (t *Nuts) SetUserFastTrigger(userID, ident string) error {
	return setValue(t, bucketUsers, key(userID, "fasttrigger"), ident)
}

func (t *Nuts) GetGuildFilters(guildID string) (GuildFilters, error) {
	return getValue[GuildFilters](t, bucketGuilds, key(guildID, "filters"))
}

func (t *Nuts) SetGuildFilters(guildID string, f GuildFilters) error {
	return setValue(t, bucketGuilds, key(guildID, "filters"), f)
}

func (t *Nuts) PutPlaybackLog(e PlaybackLogEntry) error {
	return setValue(t, bucketStats, key(e.Id), e)
}

func (t *Nuts) GetPlaybackLog(
	guildID, ident, userID string,
	limit, offset int,
) ([]PlaybackLogEntry, error) {

	var entries nutsdb.Entries
	err := t.db.View(func(tx *nutsdb.Tx) error {
		var err error
		entries, err = tx.GetAll(bucketStats)
		return t.wrapErr(err)
	})
	if err != nil {
		return nil, err
	}

	if offset >= len(entries) {
		return []PlaybackLogEntry{}, nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Meta.Timestamp > entries[j].Meta.Timestamp
	})

	size := limit
	if limit == 0 {
		size = len(entries)
	}
	logs := make([]PlaybackLogEntry, 0, size)
	n := 0
	for i, e := range entries {
		if i < offset {
			continue
		}
		if limit != 0 && n == limit {
			break
		}
		log, err := unmarshal[PlaybackLogEntry](e.Value)
		if err != nil {
			return nil, err
		}
		if guildID != "" && guildID != log.GuildID ||
			ident != "" && ident != log.Ident ||
			userID != "" && userID != log.UserID {
			continue
		}
		logs = append(logs, log)
		n++
	}

	return logs, nil
}

// --- Internal ---

func getValue[TVal any](t *Nuts, bucket string, key []byte) (TVal, error) {
	var (
		def TVal
		e   *nutsdb.Entry
	)

	err := t.db.View(func(tx *nutsdb.Tx) error {
		var err error
		e, err = tx.Get(bucket, key)
		return t.wrapErr(err)
	})
	if err != nil {
		return def, err
	}

	v, err := unmarshal[TVal](e.Value)
	return v, err
}

func setValue[TVal any](t *Nuts, bucket string, key []byte, val TVal) error {
	data, err := marshal(val)
	if err != nil {
		return err
	}
	return t.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(bucket, key, data, 0)
	})
}

func (t *Nuts) wrapErr(err error) error {
	if err == nil {
		return nil
	}
	if err == nutsdb.ErrKeyNotFound ||
		err == nutsdb.ErrNotFoundKey ||
		err == nutsdb.ErrBucketNotFound ||
		strings.HasPrefix(err.Error(), "bucket not found:") ||
		err == nutsdb.ErrBucketEmpty {
		return ErrNotFound
	}
	return err
}

func marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func unmarshal[T any](data []byte) (v T, err error) {
	err = json.Unmarshal(data, &v)
	return v, err
}

func key(elements ...string) []byte {
	return []byte(strings.Join(elements, keySeparator))
}
