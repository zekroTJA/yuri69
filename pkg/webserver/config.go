package webserver

import "github.com/zekrotja/yuri69/pkg/webserver/auth"

type DiscordOAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type WebserverConfig struct {
	BindAddress   string
	PublicAddress string
	DiscordOAuth  DiscordOAuthConfig
	Auth          auth.AuthConfig
}
