package auth

import (
	"errors"
	"net/http"
	"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/zekrotja/jwt"
	"github.com/zekrotja/yuri69/pkg/debug"
	"github.com/zekrotja/yuri69/pkg/models"
)

const issuer = "yuri69_backend"

type AuthConfig struct {
	RefreshTokenKey      string
	RefreshTokenLifetime time.Duration
	AccessTokenKey       string
	AccessTokenLifetime  time.Duration
}

type AuthHandler struct {
	accessTokenHandler  JWTHandler
	refreshTokenHandler JWTHandler
}

func New(config AuthConfig) (*AuthHandler, error) {
	var t AuthHandler

	if config.AccessTokenLifetime == 0 {
		return nil, errors.New("AccessTokenLifetime must be larger than 0")
	}
	if config.RefreshTokenLifetime == 0 {
		return nil, errors.New("RefreshTokenLifetime must be larger than 0")
	}

	var err error
	t.accessTokenHandler, err = NewJWTHandler(
		config.AccessTokenKey, issuer, config.AccessTokenLifetime)
	if err != nil {
		return nil, err
	}
	t.refreshTokenHandler, err = NewJWTHandler(
		config.RefreshTokenKey, issuer, config.RefreshTokenLifetime)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t AuthHandler) HandleLogin(ctx *routing.Context, userID string) error {
	var claims Claims
	claims.UserID = userID
	refreshToken, err := t.refreshTokenHandler.Generate(claims)
	if err != nil {
		return err
	}

	http.SetCookie(ctx.Response, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Domain:   ctx.Request.URL.Host,
		Path:     "/",
		MaxAge:   int(t.accessTokenHandler.Lifetime().Seconds()),
		HttpOnly: true,
		Secure:   !debug.Enabled(),
	})

	return t.respondAccessToken(ctx, claims)
}

func (t AuthHandler) HandleRefresh(ctx *routing.Context) error {
	cookie, err := ctx.Request.Cookie("refreshToken")
	if err == http.ErrNoCookie {
		return ctx.WriteWithStatus(nil, http.StatusUnauthorized)
	}
	if err != nil {
		return err
	}

	claims, err := t.refreshTokenHandler.Verify(cookie.Value)
	if jwt.IsJWTError(err) {
		return ctx.WriteWithStatus(nil, http.StatusUnauthorized)
	}
	if err != nil {
		return err
	}

	return t.respondAccessToken(ctx, claims)
}

func (t AuthHandler) respondAccessToken(ctx *routing.Context, claims Claims) error {
	accessTokenExpires := time.Now().Add(t.accessTokenHandler.Lifetime())
	accessToken, err := t.accessTokenHandler.Generate(claims)
	if err != nil {
		return err
	}

	return ctx.Write(models.AuthLoginResponse{
		AccessToken: accessToken,
		Expires:     accessTokenExpires,
	})
}
