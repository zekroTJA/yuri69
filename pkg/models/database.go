package models

import "time"

type Sound struct {
	Uid         string    `json:"uid"`
	DisplayName string    `json:"display_name"`
	Created     time.Time `json:"created_date"`
	CreatorId   string    `json:"creator_id"`
	Tags        []string  `json:"tags"`
}

func (t Sound) String() string {
	if t.DisplayName != "" {
		return t.DisplayName
	}
	return t.Uid
}
