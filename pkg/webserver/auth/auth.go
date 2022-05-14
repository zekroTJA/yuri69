package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"github.com/zekrotja/jwt"
	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/debug"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/models"
)

const (
	issuer           = "yuri69_backend"
	otaTokenLifetime = 30 * time.Second
)

type AuthConfig struct {
	RefreshTokenKey      string
	RefreshTokenLifetime time.Duration
	AccessTokenKey       string
	AccessTokenLifetime  time.Duration
}

type AuthHandler struct {
	publicAddress string

	accessTokenHandler  JWTHandler
	refreshTokenHandler JWTHandler
	otaTokenHandler     JWTHandler
}

func New(config AuthConfig, publicAddress string) (*AuthHandler, error) {
	var (
		t   AuthHandler
		err error
	)

	t.publicAddress = publicAddress

	if config.AccessTokenLifetime == 0 {
		return nil, errors.New("AccessTokenLifetime must be larger than 0")
	}
	if config.RefreshTokenLifetime == 0 {
		return nil, errors.New("RefreshTokenLifetime must be larger than 0")
	}

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

	t.otaTokenHandler, err = NewJWTHandler(
		"", issuer, otaTokenLifetime)
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

	ctx.Response.Header().Set("Location", "/")
	ctx.Response.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func (t AuthHandler) HandleLogout(ctx *routing.Context) error {
	http.SetCookie(ctx.Response, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Domain:   ctx.Request.URL.Host,
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
		Secure:   !debug.Enabled(),
	})

	ctx.Response.Header().Set("Location", "/login")
	ctx.Response.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func (t AuthHandler) HandleRefresh(ctx *routing.Context) error {
	cookie, err := ctx.Request.Cookie("refreshToken")
	if err == http.ErrNoCookie {
		return errs.WrapUserError("no refresh token provided", http.StatusUnauthorized)
	}
	if err != nil {
		return err
	}

	claims, err := t.refreshTokenHandler.Verify(cookie.Value)
	if jwt.IsJWTError(err) {
		return errs.WrapUserError(err, http.StatusUnauthorized)
	}
	if err != nil {
		return err
	}

	return t.respondAccessToken(ctx, claims)
}

func (t AuthHandler) CheckAuthRaw(authToken string) (Claims, error) {
	return t.accessTokenHandler.Verify(authToken)
}

func (t AuthHandler) CheckAuth(ctx *routing.Context) error {
	authHeader := ctx.Request.Header.Get("authorization")
	if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		return errs.WrapUserError("invalid access token", http.StatusUnauthorized)
	}
	authToken := authHeader[len("bearer "):]

	claims, err := t.CheckAuthRaw(authToken)
	if jwt.IsJWTError(err) {
		return errs.WrapUserError("invalid access token", http.StatusUnauthorized)
	}
	if err != nil {
		return err
	}

	ctx.Set("claims", claims)
	ctx.Set("userid", claims.UserID)

	return nil
}

func (t AuthHandler) HandleGetOtaQR(ctx *routing.Context) error {
	userID, ok := ctx.Get("userid").(string)
	if !ok || userID == "" {
		return errs.WrapUserError("request is not authenticated")
	}

	var claims Claims
	claims.UserID = userID
	deadline := time.Now().Add(t.otaTokenHandler.Lifetime())
	otaToken, err := t.otaTokenHandler.Generate(claims)
	if err != nil {
		return err
	}

	authUrl := fmt.Sprintf("%s/auth/ota/login?token=%s", t.publicAddress, otaToken)

	qrData, err := qrcode.Encode(authUrl, qrcode.Medium, 256)
	if err != nil {
		return err
	}

	qrImageData := fmt.Sprintf("data:image/png;base64,%s",
		base64.StdEncoding.EncodeToString(qrData))

	return ctx.Write(models.OTAResponse{
		Deadline:   deadline,
		Token:      otaToken,
		QRCodeData: qrImageData,
	})
}

func (t AuthHandler) HandleOTALogin(ctx *routing.Context) error {
	otaToken := ctx.Query("token")
	if otaToken == "" {
		return errs.WrapUserError("no token specified", http.StatusUnauthorized)
	}

	claims, err := t.otaTokenHandler.Verify(otaToken)
	if jwt.IsJWTError(err) {
		return errs.WrapUserError("invalid OTA token", http.StatusUnauthorized)
	}
	if err != nil {
		return err
	}

	return t.HandleLogin(ctx, claims.UserID)
}

// --- Helpers ---

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
