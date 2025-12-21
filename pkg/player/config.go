package player

import (
	"time"

	"github.com/zekrotja/yuri69/pkg/disgolink"
)

type PlayerConfig struct {
	Hostname        string
	FastTriggerTime time.Duration

	Lavalink disgolink.DisgolinkConfig
}
