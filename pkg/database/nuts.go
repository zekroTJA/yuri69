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

func (t *Nuts) GetSounds(sortOrder SortOrder, tagsMust, tagsNot []string) ([]Sound, error) {
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
		if !util.ContainsAll(sound.Tags, tagsMust) || util.ContainsAny(sound.Tags, tagsNot) {
			continue
		}
		sounds = append(sounds, sound)
	}

	var less func(i, j int) bool

	switch sortOrder {
	case SortOrderName:
		less = func(i, j int) bool {
			return sounds[i].String() < sounds[j].String()
		}
	case SortOrderCreated:
		less = func(i, j int) bool {
			return sounds[i].Created.Before(sounds[j].Created)
		}
	default:
		less = func(i, j int) bool { return false }
	}

	sort.Slice(sounds, less)
	return sounds, nil
}

func (t *Nuts) GetSound(uid string) (Sound, error) {
	var e *nutsdb.Entry
	err := t.db.View(func(tx *nutsdb.Tx) error {
		var err error
		e, err = tx.Get(bucketSounds, []byte(uid))
		return t.wrapErr(err)

	})
	if err != nil {
		return Sound{}, err
	}

	sound, err := unmarshal[Sound](e.Value)
	return sound, err
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
