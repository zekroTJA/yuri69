package auth

import "github.com/zekrotja/yuri69/pkg/util"

type AuthOrigin string

const (
	AuthOriginDiscord = AuthOrigin("origin:discord")
	AuthOriginTwitch  = AuthOrigin("origin:twitch")
)

func (t Claims) IsAuthOrigin(origin AuthOrigin) bool {
	return util.Contains(t.Scopes, string(origin))
}
