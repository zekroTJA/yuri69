package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
)

const (
	endpointAuth     = "https://id.twitch.tv/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=user:read:email"
	endpointOauth    = "https://id.twitch.tv/oauth2/token"
	endpointValidate = "https://id.twitch.tv/oauth2/validate"
)

type TwitchOAuthConfig struct {
	ClientID     string
	ClientSecret string
}

// OnErrorFunc is the function to be used to handle errors during
// authentication.
type OnErrorFunc func(ctx *routing.Context, status int, msg string) error

// OnSuccessFuc is the func to be used to handle the successful
// authentication.
type OnSuccessFuc func(ctx *routing.Context, identity auth.Claims) error

type oAuthTokenResponse struct {
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
}

type UserIdentity struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
	Username string `json:"login"`
}

type TwitchOAuth struct {
	clientID     string
	clientSecret string
	redirectURI  string

	onError   OnErrorFunc
	onSuccess OnSuccessFuc
}

func NewTwitchOAuth(
	clientID string,
	clientSecret string,
	redirectURI string,
	onError OnErrorFunc,
	onSuccess OnSuccessFuc,
) *TwitchOAuth {
	return &TwitchOAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,

		onError:   onError,
		onSuccess: onSuccess,
	}
}

func (t *TwitchOAuth) HandlerInit(ctx *routing.Context) error {
	uri := fmt.Sprintf(endpointAuth, t.clientID, url.QueryEscape(t.redirectURI))
	ctx.Response.Header().Set("Location", uri)
	return ctx.WriteWithStatus(nil, http.StatusTemporaryRedirect)
}

func (t *TwitchOAuth) HandlerCallback(ctx *routing.Context) error {
	errCode := ctx.Query("error")
	if errCode != "" {
		return t.onError(ctx, http.StatusUnauthorized, errCode)
	}

	code := ctx.Query("code")

	// 1. Request getting bearer token by app auth code

	data := map[string][]string{
		"client_id":     {t.clientID},
		"client_secret": {t.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {t.redirectURI},
	}

	values := url.Values(data)
	req, err := http.NewRequest("POST", endpointOauth, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return t.onError(ctx, http.StatusInternalServerError, "failed executing request: "+err.Error())
	}

	resAuthBody := new(oAuthTokenResponse)
	err = parseJSONBody(res.Body, resAuthBody)
	if err != nil {
		return t.onError(ctx, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
	}

	if resAuthBody.Error != "" || resAuthBody.AccessToken == "" {
		return t.onError(ctx, http.StatusUnauthorized, "")
	}

	// 2. Request getting user ID

	req, err = http.NewRequest("GET", endpointValidate, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resAuthBody.AccessToken))

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return t.onError(ctx, http.StatusInternalServerError, "failed executing request: "+err.Error())
	}

	if res.StatusCode >= 300 {
		return t.onError(ctx, http.StatusUnauthorized, "")
	}

	var resValidate UserIdentity
	err = parseJSONBody(res.Body, &resValidate)
	if err != nil {
		return t.onError(ctx, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
	}

	if resValidate.UserID == "" {
		return t.onError(ctx, http.StatusUnauthorized, "")
	}

	return t.onSuccess(ctx, auth.Claims{
		UserID:   resValidate.UserID,
		Username: resValidate.Username,
		Scopes:   []string{string(auth.AuthOriginTwitch)},
	})
}

func parseJSONBody(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}
