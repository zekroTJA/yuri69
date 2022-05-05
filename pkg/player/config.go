package player

import (
	"time"

	"github.com/zekrotja/yuri69/pkg/lavalink"
)

type PlayerConfig struct {
	Hostname        string
	FastTriggerTime time.Duration

	Lavalink lavalink.LavalinkConfig
}
