package database

import (
	"strings"

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

	GetAdmins() ([]string, error)
	AddAdmin(userID string) error
	RemoveAdmin(userID string) error
	IsAdmin(userID string) (bool, error)
}

type DatabaseConfig struct {
	Type string
	Nuts NutsConfig
}

func New(c DatabaseConfig) (IDatabase, error) {
	switch strings.ToLower(c.Type) {
	case "nuts", "local", "file":
		return NewNuts(c.Nuts)
	default:
		return nil, ErrUnsupportedProviderType
	}
}
