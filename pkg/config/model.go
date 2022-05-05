package config

import (
	"time"

	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/lavalink"
	"github.com/zekrotja/yuri69/pkg/player"
	"github.com/zekrotja/yuri69/pkg/storage"
	"github.com/zekrotja/yuri69/pkg/webserver"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
)

var DefaultConfig = Config{
	Database: database.DatabaseConfig{
		Type: "nuts",
		Nuts: database.NutsConfig{
			Location: "data/db",
		},
	},
	Storage: storage.StorageConfig{
		Type: "file",
		File: storage.FileConfig{
			BasePath: "data/st",
		},
	},
	Webserver: webserver.WebserverConfig{
		BindAddress: "0.0.0.0:80",
		Auth: auth.AuthConfig{
			RefreshTokenLifetime: 90 * 24 * time.Hour,
			AccessTokenLifetime:  10 * time.Minute,
		},
	},
	Player: player.PlayerConfig{
		FastTriggerTime: 300 * time.Millisecond,
	},
}

type Config struct {
	Database  database.DatabaseConfig
	Storage   storage.StorageConfig
	Webserver webserver.WebserverConfig
	Discord   discord.DiscordConfig
	Lavalink  lavalink.LavalinkConfig
	Player    player.PlayerConfig
}
