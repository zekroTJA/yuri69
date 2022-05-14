package models

import (
	"time"

	"github.com/zekrotja/yuri69/pkg/util"
)

const (
	EventSoundCreated       = "soundcreated"
	EventSoundUpdated       = "soundupdated"
	EventSoundDeleted       = "sounddeleted"
	EventVolumeUpdated      = "volumeupdated"
	EventGuildFilterUpdated = "guildfilterupdated"

	EventSenderController = "controller"
	EventSenderPlayer     = "player"
)

type EventType struct {
	Type string `json:"type"`
}

type Event[T any] struct {
	Type    string `json:"type"`
	Origin  string `json:"origin,omitempty"`
	Payload T      `json:"payload,omitempty"`
}

type EventAuthPromptPayload struct {
	Deadline  time.Time `json:"deadline"`
	TokenType string    `json:"token_type"`
}

type EventAuthRequest struct {
	Token string `json:"token"`
}

type GuildInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

type EventVoiceJoinPayload struct {
	Volume  int          `json:"volume,omitempty"`
	Filters GuildFilters `json:"filters,omitempty"`
	Guild   GuildInfo    `json:"guild,omitempty"`
}

type EventStatePayload struct {
	EventVoiceJoinPayload

	Connected bool `json:"connected"`
	Joined    bool `json:"joined"`
	IsAdmin   bool `json:"is_admin"`
}

type EventErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WrapErrorEvent(err error, code ...int) Event[any] {
	return Event[any]{
		Type: "error",
		Payload: EventErrorPayload{
			Code:    util.Opt(code, 500),
			Message: err.Error(),
		},
	}
}
