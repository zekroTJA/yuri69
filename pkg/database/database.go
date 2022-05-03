package database

import (
	"strings"

	. "github.com/zekrotja/yuri69/pkg/models"
)

type SortOrder string

const (
	SortOrderName    = SortOrder("name")
	SortOrderCreated = SortOrder("created")
)

type IDatabase interface {
	Close() error

	PutSound(sound Sound) error
	RemoveSound(uid string) error
	GetSounds(sortOrder SortOrder, tagsMust, tagsNot []string) ([]Sound, error)
	GetSound(uid string) (Sound, error)

	SetGuildVolume(guildID string, volume int) error
	GetGuildVolume(guildID string) (int, error)
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
