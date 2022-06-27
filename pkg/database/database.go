package database

import (
	"strings"

	"github.com/zekrotja/yuri69/pkg/database/dberrors"
	"github.com/zekrotja/yuri69/pkg/database/nuts"
	"github.com/zekrotja/yuri69/pkg/database/postgres"
	. "github.com/zekrotja/yuri69/pkg/models"
)

type IDatabase interface {
	Close() error

	PutSound(sound Sound) error
	RemoveSound(uid string) error
	GetSounds() ([]Sound, error)
	GetSound(uid string) (Sound, error)

	GetGuildVolume(guildID string) (int, error)
	SetGuildVolume(guildID string, volume int) error

	GetUserFastTrigger(userID string) (string, error)
	SetUserFastTrigger(userID, ident string) error

	GetGuildFilters(guildID string) (GuildFilters, error)
	SetGuildFilters(guildID string, f GuildFilters) error

	PutPlaybackLog(e PlaybackLogEntry) error
	GetPlaybackLog(guildID, ident, userID string, limit, offset int) ([]PlaybackLogEntry, error)
	GetPlaybackLogSize() (int, error)
	GetPlaybackStats(guildID, userID string) ([]PlaybackStats, error)

	GetAdmins() ([]string, error)
	AddAdmin(userID string) error
	RemoveAdmin(userID string) error
	IsAdmin(userID string) (bool, error)

	GetFavorites(userID string) ([]string, error)
	AddFavorite(userID, ident string) error
	RemoveFavorite(userID, ident string) error

	GetApiKey(userID string) (string, error)
	GetUserByApiKey(token string) (string, error)
	SetApiKey(userID, token string) error
	RemoveApiKey(userID string) error

	SetTwitchSettings(s TwitchSettings) error
	GetTwitchSettings(userid string) (TwitchSettings, error)
}

type IMigrate interface {
	Migrate() error
}

type DatabaseConfig struct {
	Type     string
	Nuts     nuts.NutsConfig
	Postgres postgres.PostgresConfig
}

func New(c DatabaseConfig) (IDatabase, error) {
	var (
		db  IDatabase
		err error
	)

	switch strings.ToLower(c.Type) {
	case "nuts", "local", "file":
		db, err = nuts.NewNuts(c.Nuts)
	case "postgres", "pg", "postgresql":
		db, err = postgres.NewPostgres(c.Postgres)
	default:
		err = dberrors.ErrUnsupportedProviderType
	}

	if err != nil {
		return db, err
	}

	if mg, ok := db.(IMigrate); ok {
		err = mg.Migrate()
	}

	return db, err
}
