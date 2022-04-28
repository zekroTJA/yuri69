package models

import "time"

type AuthLoginResponse struct {
	AccessToken string    `json:"access_token"`
	Expires     time.Time `json:"expires"`
}
