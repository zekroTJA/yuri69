package webserver

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/ozzo-routing/v2/content"
	"github.com/zekrotja/ozzo-routing/v2/cors"
	"github.com/zekrotja/ozzo-routing/v2/fault"
	"github.com/zekrotja/ozzo-routing/v2/file"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/debug"
	"github.com/zekrotja/yuri69/pkg/discordoauth"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
	"github.com/zekrotja/yuri69/pkg/webserver/controllers"
	"github.com/zekrotja/yuri69/pkg/webserver/ws"
)

var (
	//go:embed _webdist/*
	embeddedFS embed.FS
)

type Webserver struct {
	bindAddress string
	router      *routing.Router
	server      *http.Server
	authHandler *auth.AuthHandler
	ct          *controller.Controller
	hub         *ws.Hub
}

func New(cfg WebserverConfig, ct *controller.Controller) (*Webserver, error) {
	var (
		t   Webserver
		err error
	)

	t.ct = ct

	t.bindAddress = cfg.BindAddress
	t.router = routing.New()
	t.router.Use(
		content.TypeNegotiator(content.JSON),
		fault.Recovery(logrus.Debugf, t.errorHandler),
	)

	if debug.Enabled() {
		const corsOrigin = "http://localhost:3000"
		logrus.Warnf("CORS enabled for address %s", corsOrigin)
		t.router.Use(cors.Handler(cors.Options{
			AllowOrigins:     corsOrigin,
			AllowCredentials: true,
			AllowMethods:     "*",
			AllowHeaders:     "*",
		}))
	}

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

	t.hub = ws.NewHub(t.authHandler, t.ct)

	t.registerRoutes(oauth)
	t.hookFS()

	t.ct.SubscribeFunc(t.eventHandler)

	return &t, nil
}

func (t *Webserver) ListenAndServeBlocking() error {
	t.server = &http.Server{
		Addr:    t.bindAddress,
		Handler: t.router,
	}

	return t.server.ListenAndServe()
}

func (t *Webserver) registerRoutes(oauth *discordoauth.DiscordOAuth) {
	t.router.Any("/ws", t.hub.Upgrade)
	t.router.Get("/invite", func(ctx *routing.Context) error {
		url := fmt.Sprintf(
			"https://discord.com/api/oauth2/authorize?client_id=%s&permissions=36702208&scope=applications.commands%%20bot",
			t.ct.GetSelfUser().ID)
		ctx.Response.Header().Set("Location", url)
		ctx.Response.WriteHeader(http.StatusTemporaryRedirect)
		return nil
	})

	gApi := t.router.Group("/api/v1")

	gAuth := gApi.Group("/auth")
	gAuth.Get("/login", oauth.HandlerInit)
	gAuth.Get("/oauthcallback", oauth.HandlerCallback)
	gAuth.Get("/refresh", t.authHandler.HandleRefresh)

	gApi.Use(t.authHandler.CheckAuth)

	controllers.NewSoundsController(gApi.Group("/sounds"), t.ct)
	controllers.NewPlayerController(gApi.Group("/players"), t.ct)
	controllers.NewUsersController(gApi.Group("/users"), t.ct)
	controllers.NewGuildsController(gApi.Group("/guilds"), t.ct)
}

func (t *Webserver) hookFS() error {
	fsys, _ := fs.Sub(embeddedFS, "_webdist")
	_, err := fs.Stat(fsys, "index.html")
	if err != nil {
		logrus.WithError(err).Debug("Using external static web files")
		fsys = os.DirFS("web/dist")
	} else {
		logrus.Debug("Use embedded static web files")
	}

	t.router.Get("/*", file.Server(file.PathMap{
		"/":       "",
		"/assets": "/assets",
	}, file.ServerOptions{
		FS:           fsys,
		IndexFile:    "index.html",
		CatchAllFile: "index.html",
	}))

	return nil
}

func (t *Webserver) eventHandler(e controller.ControllerEvent) {
	var err error
	if e.IsBroadcast {
		err = t.hub.Broadcast(e.Event)
	} else {
		err = t.hub.BroadcastScoped(e.Event, e.Receivers)
	}
	if err != nil {
		logrus.WithError(err).Errorf("Broadcasting event failed: %+v", e)
	}
}

func (t *Webserver) onAuthError(ctx *routing.Context, status int, msg string) error {
	return ctx.WriteWithStatus(msg, status)
}

func (t *Webserver) errorHandler(ctx *routing.Context, err error) error {
	httpError, ok := errs.As[routing.HTTPError](err)
	if ok {
		ctx.Response.WriteHeader(httpError.StatusCode())
		return httpError
	}

	if err == database.ErrNotFound {
		err = errs.StatusError{
			Status:  http.StatusNotFound,
			Message: "Not Found",
		}
	}

	statusErr, ok := errs.As[errs.StatusError](err)
	if ok {
		ctx.Response.WriteHeader(statusErr.Status)
		return statusErr
	}

	return t.errorHandler(ctx, errs.StatusError{
		Status:  http.StatusInternalServerError,
		Message: err.Error(),
	})
}
