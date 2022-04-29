package webserver

import (
	"fmt"
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/access"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/fault"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/discordoauth"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
)

type Webserver struct {
	bindAddress string
	router      *routing.Router
	server      *http.Server
	authHandler *auth.AuthHandler
	ct          *controller.Controller
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
		access.CustomLogger(func(req *http.Request, res *access.LogResponseWriter, elapsed float64) {
			clientIP := access.GetClientIP(req)
			logrus.WithFields(logrus.Fields{
				"client": clientIP,
				"took":   fmt.Sprintf("%.3fms", elapsed),
			}).Debugf("%s %s %s", req.Method, req.URL.String(), req.Proto)
		}),
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

	gApi.Use(t.authHandler.CheckAuth)

	gSounds := gApi.Group("/sounds")
	gSounds.Get("", t.handleSoundsList)
	gSounds.Delete("/<id>", t.handleSoundsDelete)
	gSounds.Put("/upload", t.handleSoundsUpload)
	gSounds.Post("/create", t.handleSoundsCreate)

	return &t, nil
}
func (t *Webserver) ListenAndServeBlocking() error {
	t.server = &http.Server{
		Addr:    t.bindAddress,
		Handler: t.router,
	}

	return t.server.ListenAndServe()
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
