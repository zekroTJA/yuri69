package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/util"
)

var (
	uidRx = regexp.MustCompile(`^[a-z0-9_.-]{1,30}$`)
)

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

func (t Sound) Check() error {
	if t.Uid == "" {
		return errs.WrapUserError("uid must be specified")
	}

	if !uidRx.MatchString(t.Uid) {
		return errs.WrapUserError("malformed uid")
	}

	if util.HasDuplicates(t.Tags) {
		return errs.WrapUserError("'tags' has duplicate entries")
	}

	return nil
}

func (t *Sound) Sanitize() {
	util.ApplyToAll(t.Tags, strings.ToLower)
}

type GuildFilters struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

func (t GuildFilters) Check() error {
	if util.HasDuplicates(t.Include) {
		return errs.WrapUserError("'include' has duplicate elements")
	}
	if util.HasDuplicates(t.Exclude) {
		return errs.WrapUserError("'exclude' has duplicate elements")
	}
	return nil
}

func (t *GuildFilters) Sanitize() {
	util.ApplyToAll(t.Include, strings.ToLower)
	util.ApplyToAll(t.Exclude, strings.ToLower)
}

type PlaybackLogEntry struct {
	Id        string    `json:"id"`
	Ident     string    `json:"ident"`
	GuildID   string    `json:"guild_id"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}
