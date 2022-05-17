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
	publicAddress        string
	checkApiTokenHandler func(string) (string, error)

	accessTokenHandler  JWTHandler
	refreshTokenHandler JWTHandler
	otaTokenHandler     JWTHandler
}

func New(
	config AuthConfig,
	publicAddress string,
	checkApiTokenHandler func(string) (string, error),
) (*AuthHandler, error) {
	var (
		t   AuthHandler
		err error
	)

	t.publicAddress = publicAddress
	t.checkApiTokenHandler = checkApiTokenHandler

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
		MaxAge:   int(t.refreshTokenHandler.Lifetime().Seconds()),
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
	var claims Claims

	if ok, refreshToken, err := getCookie(ctx, "refreshToken"); err != nil {
		return err
	} else if ok {
		claims, err = t.refreshTokenHandler.Verify(refreshToken)
		if jwt.IsJWTError(err) {
			return errs.WrapUserError(err, http.StatusUnauthorized)
		}
		if err != nil {
			return err
		}
	} else if ok, token := getAuthorizationToken(ctx, "basic"); ok {
		userid, err := t.checkApiTokenHandler(token)
		if err != nil {
			return errs.WrapUserError("invalid bearer token", http.StatusUnauthorized)
		}
		claims.UserID = userid
	} else {
		return errs.WrapUserError("no refresh or bearer token provided", http.StatusUnauthorized)
	}

	return t.respondAccessToken(ctx, claims)
}

func (t AuthHandler) CheckAuthRaw(authToken string) (Claims, error) {
	return t.accessTokenHandler.Verify(authToken)
}

func (t AuthHandler) CheckAuth(ctx *routing.Context) error {
	var (
		claims Claims
		err    error
	)

	if ok, token := getAuthorizationToken(ctx, "basic"); ok {
		userid, err := t.checkApiTokenHandler(token)
		if err != nil {
			return errs.WrapUserError("invalid basic token", http.StatusUnauthorized)
		}
		claims.UserID = userid
	} else if ok, token := getAuthorizationToken(ctx, "bearer"); ok {
		claims, err = t.CheckAuthRaw(token)
		if jwt.IsJWTError(err) {
			return errs.WrapUserError("invalid access token", http.StatusUnauthorized)
		}
		if err != nil {
			return err
		}
	} else {
		return errs.WrapUserError("invalid access token", http.StatusUnauthorized)
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

	authUrl := fmt.Sprintf("%s/api/v1/auth/ota/login?token=%s", t.publicAddress, otaToken)

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
	accessTokenDeadline := time.Now().Add(t.accessTokenHandler.Lifetime())
	accessToken, err := t.accessTokenHandler.Generate(claims)
	if err != nil {
		return err
	}

	return ctx.Write(models.AuthLoginResponse{
		AccessToken: accessToken,
		Deadline:    accessTokenDeadline,
	})
}

func getAuthorizationToken(ctx *routing.Context, typ string) (bool, string) {
	authHeader := ctx.Request.Header.Get("authorization")
	typ = typ + " "

	if !strings.HasPrefix(strings.ToLower(authHeader), typ) {
		return false, ""
	}

	token := authHeader[len(typ):]

	return token != "", token
}

func getCookie(ctx *routing.Context, name string) (bool, string, error) {
	cookie, err := ctx.Request.Cookie(name)
	if err == http.ErrNoCookie {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return cookie.Value != "", cookie.Value, nil
}
