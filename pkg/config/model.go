package config

import (
	"time"

	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/lavalink"
	"github.com/zekrotja/yuri69/pkg/storage"
	"github.com/zekrotja/yuri69/pkg/webserver"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
)

var DefaultConfig = Config{
	Storage: storage.StorageConfig{
		Type: "file",
		File: storage.FileConfig{
			BasePath: "data",
		},
	},
	Webserver: webserver.WebserverConfig{
		BindAddress: "0.0.0.0:80",
		Auth: auth.AuthConfig{
			RefreshTokenLifetime: 90 * 24 * time.Hour,
			AccessTokenLifetime:  10 * time.Minute,
		},
	},
}

type Config struct {
	Storage   storage.StorageConfig
	Webserver webserver.WebserverConfig
	Discord   discord.DiscordConfig
	Lavalink  lavalink.LavalinkConfig
}
