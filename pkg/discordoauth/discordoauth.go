// package discordoauth provides http handlers
// to authenticate with via the Discord OAuth2
// endpoint.
package discordoauth

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
	endpointOauth = "https://discord.com/api/oauth2/token"
	endpointMe    = "https://discord.com/api/users/@me"
	endpointAuth  = "https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify"
)

// OnErrorFunc is the function to be used to handle errors during
// authentication.
type OnErrorFunc func(ctx *routing.Context, status int, msg string) error

// OnSuccessFuc is the func to be used to handle the successful
// authentication.
type OnSuccessFuc func(ctx *routing.Context, claims auth.Claims) error

// DiscordOAuth provides http handlers for
// authenticating a discord User by your Discord
// OAuth application.
type DiscordOAuth struct {
	clientID     string
	clientSecret string
	redirectURI  string

	onError   OnErrorFunc
	onSuccess OnSuccessFuc
}

type oAuthTokenResponse struct {
	Error       string `json:"error"`
	AccessToken string `json:"access_token"`
}

type getUserMeResponse struct {
	Error string `json:"error"`
	ID    string `json:"id"`
}

// NewDiscordOAuth returns a new instance of DiscordOAuth.
func NewDiscordOAuth(
	clientID string,
	clientSecret string,
	redirectURI string,
	onError OnErrorFunc,
	onSuccess OnSuccessFuc,
) *DiscordOAuth {
	if onError == nil {
		onError = func(*routing.Context, int, string) error { return nil }
	}
	if onSuccess == nil {
		onSuccess = func(*routing.Context, auth.Claims) error { return nil }
	}

	return &DiscordOAuth{
		clientID:     clientID,
		clientSecret: clientSecret,
		redirectURI:  redirectURI,

		onError:   onError,
		onSuccess: onSuccess,
	}
}

// HandlerInit returns a redirect response to the OAuth Apps
// authentication page.
func (d *DiscordOAuth) HandlerInit(ctx *routing.Context) error {
	uri := fmt.Sprintf(endpointAuth, d.clientID, url.QueryEscape(d.redirectURI))
	ctx.Response.Header().Set("Location", uri)
	return ctx.WriteWithStatus(nil, http.StatusTemporaryRedirect)
}

// HandlerCallback will be requested by discordapp.com on successful
// app authentication. This handler will check the validity of the passed
// authorization code by getting a bearer token and trying to get self
// user data by requesting them using the bearer token.
// If this fails, onError will be called. Else, onSuccess will be
// called passing the userID of the user authenticated.
func (d *DiscordOAuth) HandlerCallback(ctx *routing.Context) error {
	code := ctx.Query("code")

	// 1. Request getting bearer token by app auth code

	data := map[string][]string{
		"client_id":     {d.clientID},
		"client_secret": {d.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {d.redirectURI},
		"scope":         {"identify"},
	}

	values := url.Values(data)
	req, err := http.NewRequest("POST", endpointOauth, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		d.onError(ctx, http.StatusInternalServerError, "failed executing request: "+err.Error())
		return nil
	}

	if res.StatusCode >= 300 {
		d.onError(ctx, http.StatusUnauthorized, "")
		return nil
	}

	resAuthBody := new(oAuthTokenResponse)
	err = parseJSONBody(res.Body, resAuthBody)
	if err != nil {
		d.onError(ctx, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
		return nil
	}

	if resAuthBody.Error != "" || resAuthBody.AccessToken == "" {
		d.onError(ctx, http.StatusUnauthorized, "")
		return nil
	}

	// 2. Request getting user ID

	req, err = http.NewRequest("GET", endpointMe, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", resAuthBody.AccessToken))

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		d.onError(ctx, http.StatusInternalServerError, "failed executing request: "+err.Error())
		return nil
	}

	if res.StatusCode >= 300 {
		d.onError(ctx, http.StatusUnauthorized, "")
		return nil
	}

	resGetMe := new(getUserMeResponse)
	err = parseJSONBody(res.Body, resGetMe)
	if err != nil {
		d.onError(ctx, http.StatusInternalServerError, "failed parsing Discord API response: "+err.Error())
		return nil
	}

	if resGetMe.Error != "" || resGetMe.ID == "" {
		d.onError(ctx, http.StatusUnauthorized, "")
		return nil
	}

	d.onSuccess(ctx, auth.Claims{
		UserID: resGetMe.ID,
		Scopes: []string{"origin:discord"},
	})

	return nil
}

func parseJSONBody(body io.Reader, v interface{}) error {
	return json.NewDecoder(body).Decode(v)
}
