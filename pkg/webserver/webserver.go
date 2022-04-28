package webserver

import (
	"fmt"
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/zekrotja/yuri69/pkg/discordoauth"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
)

type Webserver struct {
	bindAddress string
	router      *routing.Router
	server      *http.Server
	authHandler *auth.AuthHandler
}

func New(cfg WebserverConfig) (*Webserver, error) {
	var (
		t   Webserver
		err error
	)

	t.bindAddress = cfg.BindAddress
	t.router = routing.New()
	t.router.Use(
		content.TypeNegotiator(content.JSON),
	)

	t.authHandler, err = auth.New(cfg.Auth)
	if err != nil {
		return nil, err
	}

	oauth := discordoauth.NewDiscordOAuth(
		cfg.DiscordOAuth.ClientID,
		cfg.DiscordOAuth.ClientSecret,
		fmt.Sprintf("%s/api/v1/auth/oauthcallback", cfg.PublicAddress),
		t.onAuthError,
		t.authHandler.HandleLogin)

	gApi := t.router.Group("/api/v1")

	gAuth := gApi.Group("/auth")
	gAuth.Get("/login", oauth.HandlerInit)
	gAuth.Get("/oauthcallback", oauth.HandlerCallback)
	gAuth.Get("/refresh", t.authHandler.HandleRefresh)

	return &t, nil
}

func (t *Webserver) onAuthError(ctx *routing.Context, status int, msg string) error {
	return ctx.WriteWithStatus(msg, status)
}

func (t *Webserver) ListenAndServeBlocking() error {
	t.server = &http.Server{
		Addr:    t.bindAddress,
		Handler: t.router,
	}

	return t.server.ListenAndServe()
}
