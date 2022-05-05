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
