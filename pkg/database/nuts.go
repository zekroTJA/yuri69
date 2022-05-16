package database

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/xujiajun/nutsdb"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/util"
)

const (
	bucketSounds = "sounds"
	bucketGuilds = "guilds"
	bucketUsers  = "users"
	bucketStats  = "stats"
	bucketAdmins = "admins"
	bucketTokens = "tokens"
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
	return nuts_setValue(t, bucketSounds, nuts_key(sound.Uid), sound)
}

func (t *Nuts) RemoveSound(uid string) error {
	return t.remove(bucketSounds, []byte(uid))
}

func (t *Nuts) GetSounds() ([]Sound, error) {
	return nuts_listValues[Sound](t, bucketSounds, nil, nil)
}

func (t *Nuts) GetSound(uid string) (Sound, error) {
	return nuts_getValue[Sound](t, bucketSounds, nuts_key(uid))
}

func (t *Nuts) GetGuildVolume(guildID string) (int, error) {
	return nuts_getValue[int](t, bucketGuilds, nuts_key(guildID, "volume"))
}

func (t *Nuts) SetGuildVolume(guildID string, volume int) error {
	return nuts_setValue(t, bucketGuilds, nuts_key(guildID, "volume"), volume)
}

func (t *Nuts) GetUserFastTrigger(userID string) (string, error) {
	return nuts_getValue[string](t, bucketUsers, nuts_key(userID, "fasttrigger"))
}

func (t *Nuts) SetUserFastTrigger(userID, ident string) error {
	return nuts_setValue(t, bucketUsers, nuts_key(userID, "fasttrigger"), ident)
}

func (t *Nuts) GetGuildFilters(guildID string) (GuildFilters, error) {
	return nuts_getValue[GuildFilters](t, bucketGuilds, nuts_key(guildID, "filters"))
}

func (t *Nuts) SetGuildFilters(guildID string, f GuildFilters) error {
	return nuts_setValue(t, bucketGuilds, nuts_key(guildID, "filters"), f)
}

func (t *Nuts) PutPlaybackLog(e PlaybackLogEntry) error {
	return nuts_setValue(t, bucketStats, nuts_key(e.Id), e)
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
		log, err := nuts_unmarshal[PlaybackLogEntry](e.Value)
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

func (t *Nuts) GetPlaybackLogSize() (int, error) {
	var n int
	err := t.db.View(func(tx *nutsdb.Tx) error {
		entries, err := tx.GetAll(bucketStats)
		if err != nil {
			return t.wrapErr(err)
		}
		n = len(entries)
		return nil
	})
	return n, err
}

func (t *Nuts) GetPlaybackStats(guildID, userID string) ([]PlaybackStats, error) {
	logs, err := t.GetPlaybackLog(guildID, "", userID, 0, 0)
	if err != nil && err != ErrNotFound {
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

func (t *Nuts) GetAdmins() ([]string, error) {
	return nuts_listValues[string](t, bucketAdmins, nil, nil)
}

func (t *Nuts) AddAdmin(userID string) error {
	return nuts_setValue(t, bucketAdmins, nuts_key(userID), userID)
}

func (t *Nuts) RemoveAdmin(userID string) error {
	return t.remove(bucketAdmins, nuts_key(userID))
}

func (t *Nuts) IsAdmin(userID string) (bool, error) {
	v, err := nuts_getValue[string](t, bucketAdmins, []byte(userID))
	if err != nil && err != ErrNotFound {
		return false, err
	}
	return v == userID, nil
}

func (t *Nuts) GetFavorites(userID string) ([]string, error) {
	return nuts_getValue[[]string](t, bucketUsers, nuts_key(userID, "favs"))
}

func (t *Nuts) AddFavorite(userID, ident string) error {
	favs, err := t.GetFavorites(userID)
	if err != nil && err != ErrNotFound {
		return err
	}
	favs = util.AppendIfNotContains(favs, ident)
	return nuts_setValue(t, bucketUsers, nuts_key(userID, "favs"), favs)
}

func (t *Nuts) RemoveFavorite(userID, ident string) error {
	favs, err := t.GetFavorites(userID)
	if err != nil && err != ErrNotFound {
		return err
	}
	favs = util.Remove(favs, ident)
	return nuts_setValue(t, bucketUsers, nuts_key(userID, "favs"), favs)
}

func (t *Nuts) GetApiKey(userID string) (string, error) {
	return nuts_getValue[string](t, bucketUsers, nuts_key(userID, "apitoken"))
}

func (t *Nuts) GetUserByApiKey(token string) (string, error) {
	return nuts_getValue[string](t, bucketTokens, nuts_key(token))
}

func (t *Nuts) SetApiKey(userID, token string) error {
	err := nuts_setValue(t, bucketUsers, nuts_key(userID, "apitoken"), token)
	if err != nil {
		return err
	}
	return nuts_setValue(t, bucketTokens, nuts_key(token), userID)
}

func (t *Nuts) RemoveApiKey(userID string) error {
	token, err := t.GetApiKey(userID)
	if err != nil {
		return err
	}

	err = t.remove(bucketTokens, nuts_key(token))
	if err != nil {
		return err
	}

	return t.remove(bucketUsers, nuts_key(userID, "apitoken"))
}

// --- Internal ---

func nuts_getValue[TVal any](t *Nuts, bucket string, key []byte) (TVal, error) {
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

	v, err := nuts_unmarshal[TVal](e.Value)
	return v, err
}

func nuts_setValue[TVal any](t *Nuts, bucket string, key []byte, val TVal) error {
	data, err := nuts_marshal(val)
	if err != nil {
		return err
	}
	return t.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(bucket, key, data, 0)
	})
}

func nuts_listValues[TVal any](
	t *Nuts,
	bucket string,
	entryFilter func(*nutsdb.Entry) bool,
	valueFilter func(TVal) bool,
) ([]TVal, error) {
	var entries nutsdb.Entries
	err := t.db.View(func(tx *nutsdb.Tx) error {
		var err error
		entries, err = tx.GetAll(bucket)
		return t.wrapErr(err)
	})
	if err != nil {
		return nil, err
	}

	if entryFilter == nil {
		entryFilter = func(e *nutsdb.Entry) bool { return true }
	}
	if valueFilter == nil {
		valueFilter = func(t TVal) bool { return true }
	}

	vals := make([]TVal, 0, len(entries))
	for _, e := range entries {
		if !entryFilter(e) {
			continue
		}
		v, err := nuts_unmarshal[TVal](e.Value)
		if err != nil {
			return nil, err
		}
		if !valueFilter(v) {
			continue
		}
		vals = append(vals, v)
	}

	return vals, nil
}

func (t *Nuts) remove(bucket string, key []byte) error {
	return t.db.Update(func(tx *nutsdb.Tx) error {
		err := tx.Delete(bucket, key)
		return t.wrapErr(err)
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

func nuts_marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func nuts_unmarshal[T any](data []byte) (v T, err error) {
	err = json.Unmarshal(data, &v)
	return v, err
}

func nuts_key(elements ...string) []byte {
	return []byte(strings.Join(elements, keySeparator))
}
