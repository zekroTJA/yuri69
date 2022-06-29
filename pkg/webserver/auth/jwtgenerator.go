package auth

import (
	"errors"
	"time"

	"github.com/zekrotja/jwt"
	"github.com/zekrotja/yuri69/pkg/cryptoutil"
)

type Claims struct {
	jwt.PublicClaims

	UserID   string
	Username string
	Scopes   []string
}

type JWTHandler struct {
	handler  jwt.Handler[Claims]
	lifetime time.Duration
	issuer   string
}

func NewJWTHandler(key string, issuer string, lifetime time.Duration) (JWTHandler, error) {
	var t JWTHandler

	t.lifetime = lifetime

	var bKey []byte
	var err error

	if key != "" {
		bKey = []byte(key)
	} else {
		bKey, err = cryptoutil.GetRandByteArray(64)
		if err != nil {
			return t, err
		}
	}

	t.handler = jwt.NewHandler[Claims](jwt.NewHmacSha256(bKey))

	return t, nil
}

func (t JWTHandler) Lifetime() time.Duration {
	return t.lifetime
}

func (t JWTHandler) Generate(claims Claims) (string, error) {
	claims.Iss = t.issuer
	claims.SetIat()
	claims.SetExpDuration(t.lifetime)
	claims.SetNbfTime(time.Now())

	return t.handler.EncodeAndSign(claims)
}

func (t JWTHandler) Verify(token string) (Claims, error) {
	claims, err := t.handler.DecodeAndValidate(token)
	if err != nil {
		return Claims{}, err
	}

	if claims.Iss != t.issuer {
		return Claims{}, errors.New("invalid issuer")
	}

	return claims, nil
}
