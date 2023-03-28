package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/internal/embedded"
	"github.com/zekrotja/yuri69/pkg/database/dberrors"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/util"
)

type PostgresConfig struct {
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type Postgres struct {
	db *sql.DB
}

func NewPostgres(c PostgresConfig) (*Postgres, error) {
	var (
		t   Postgres
		err error
	)

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		c.Host, c.Port, c.Database, c.Username, c.Password)
	t.db, err = sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = t.db.Ping()
	if err != nil {
		return nil, err
	}

	goose.SetBaseFS(embedded.Migrations)
	goose.SetDialect("postgres")
	goose.SetLogger(logrus.StandardLogger())
	err = goose.Up(t.db, "migrations/postgres")
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Postgres) Close() error {
	return t.db.Close()
}

func (t *Postgres) PutSound(sound Sound) error {
	oldSound, err := t.GetSound(sound.Uid)
	exists := oldSound.Uid == sound.Uid
	if err != nil && err != dberrors.ErrNotFound {
		return err
	}

	if exists {
		err = t.tx(func(tx *sql.Tx) error {
			_, err := tx.Exec(`
				UPDATE sounds
				SET "displayname" = $2
				WHERE "uid" = $1
			`, sound.Uid, sound.DisplayName)
			if err != nil {
				return err
			}

			addedTags, removedTags := util.Diff(oldSound.Tags, sound.Tags)
			for _, tag := range removedTags {
				_, err := tx.Exec(`DELETE FROM sounds_tags WHERE "sound" = $1 AND "tag" = $2`,
					sound.Uid, tag)
				if err != nil {
					return err
				}
			}
			for _, tag := range addedTags {
				_, err := tx.Exec(`INSERT INTO sounds_tags ("sound", "tag") VALUES ($1, $2)`,
					sound.Uid, tag)
				if err != nil {
					return err
				}
			}
			return nil
		})
		return err
	}

	err = t.tx(func(tx *sql.Tx) error {
		_, err := tx.Exec(`
			INSERT INTO sounds ("uid", "displayname", "created", "creatorid")
			VALUES ($1, $2, $3, $4)
		`, sound.Uid, sound.DisplayName, sound.Created, sound.Creator.ID)
		if err != nil {
			return err
		}

		for _, tag := range sound.Tags {
			_, err = tx.Exec(`
				INSERT INTO sounds_tags ("sound", "tag")
				VALUES ($1, $2)
			`, sound.Uid, tag)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (t *Postgres) RemoveSound(uid string) error {
	return pg_delete(t, "sounds", "uid", uid)
}

func (t *Postgres) GetSounds() ([]Sound, error) {
	rows, err := t.db.Query(`
		SELECT "uid", "displayname", "created", "creatorid", "tag"
		FROM sounds
		LEFT JOIN sounds_tags
		ON sounds."uid" = sounds_tags."sound"
	`)
	if err != nil {
		return nil, t.wrapErr(err)
	}

	soundsMap := make(map[string]*Sound)
	for rows.Next() {
		var s Sound
		var tag sql.NullString
		err = rows.Scan(&s.Uid, &s.DisplayName, &s.Created, &s.Creator.ID, &tag)
		if err != nil {
			return nil, err
		}
		s.Uid = strings.TrimSpace(s.Uid)
		ms, ok := soundsMap[s.Uid]
		if !ok {
			ms = &s
			soundsMap[s.Uid] = ms
		}
		if tag.Valid {
			ms.Tags = append(soundsMap[s.Uid].Tags, tag.String)
		}
	}

	sounds := make([]Sound, 0, len(soundsMap))
	for _, s := range soundsMap {
		sounds = append(sounds, *s)
	}

	return sounds, nil
}

func (t *Postgres) GetSound(uid string) (Sound, error) {
	rows, err := t.db.Query(`
	    SELECT "uid", "displayname", "created", "creatorid", "tag"
	    FROM sounds
	    LEFT JOIN sounds_tags
	    ON sounds."uid" = sounds_tags."sound"
		WHERE sounds."uid" = $1
	`, uid)
	if err != nil {
		return Sound{}, t.wrapErr(err)
	}

	var s Sound
	for rows.Next() {
		var tag sql.NullString
		err = rows.Scan(&s.Uid, &s.DisplayName, &s.Created, &s.Creator.ID, &tag)
		if err != nil {
			return Sound{}, err
		}
		s.Uid = strings.TrimSpace(s.Uid)
		if tag.Valid {
			s.Tags = append(s.Tags, tag.String)
		}
	}

	return s, nil
}

func (t *Postgres) GetGuildVolume(guildID string) (int, error) {
	return pg_getValue[int](t, "guilds", "volume", "id", guildID)
}

func (t *Postgres) SetGuildVolume(guildID string, volume int) error {
	return pg_setValue(t, "guilds", "volume", volume, "id", guildID)
}

func (t *Postgres) GetUserFastTrigger(userID string) (string, error) {
	return pg_getValue[string](t, "users", "fasttrigger", "id", userID)
}

func (t *Postgres) SetUserFastTrigger(userID, ident string) error {
	return pg_setValue(t, "users", "fasttrigger", ident, "id", userID)
}

func (t *Postgres) GetGuildFilters(guildID string) (GuildFilters, error) {
	rows, err := t.db.Query(`
		SELECT "exclude", "tag"
		FROM guild_filters
		WHERE "guildid" = $1
	`, guildID)
	if err != nil {
		return GuildFilters{}, t.wrapErr(err)
	}

	var gf GuildFilters
	for rows.Next() {
		var (
			exclude bool
			tag     string
		)
		err = rows.Scan(&exclude, &tag)
		if err != nil {
			return GuildFilters{}, err
		}
		if exclude {
			gf.Exclude = append(gf.Exclude, tag)
		} else {
			gf.Include = append(gf.Include, tag)
		}
	}

	return gf, nil
}

func (t *Postgres) SetGuildFilters(guildID string, f GuildFilters) error {
	before, err := t.GetGuildFilters(guildID)
	if err != nil {
		return err
	}

	excludeAdded, excludeRemoved := util.Diff(before.Exclude, f.Exclude)
	includeAdded, includeRemoved := util.Diff(before.Include, f.Include)

	err = t.tx(func(tx *sql.Tx) error {
		for _, tag := range excludeRemoved {
			_, err := tx.Exec(
				`DELETE FROM guild_filters WHERE "guildid" = $1 AND "tag" = $2 AND "exclude" = 'true'`,
				guildID, tag)
			if err != nil {
				return err
			}
		}
		for _, tag := range includeRemoved {
			_, err := tx.Exec(
				`DELETE FROM guild_filters WHERE "guildid" = $1 AND "tag" = $2 AND "exclude" = 'false'`,
				guildID, tag)
			if err != nil {
				return err
			}
		}
		for _, tag := range excludeAdded {
			_, err := tx.Exec(
				`INSERT INTO guild_filters ("guildid", "tag", "exclude") VALUES ($1, $2, 'true')`,
				guildID, tag)
			if err != nil {
				return err
			}
		}
		for _, tag := range includeAdded {
			_, err := tx.Exec(
				`INSERT INTO guild_filters ("guildid", "tag", "exclude") VALUES ($1, $2, 'false')`,
				guildID, tag)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

func (t *Postgres) PutPlaybackLog(e PlaybackLogEntry) error {
	_, err := t.db.Exec(`
		INSERT INTO playbacklog ("id", "sound", "guildid", "userid", "timestamp")
		VALUES ($1, $2, $3, $4, $5)
	`, e.Id, e.Ident, e.GuildID, e.UserID, e.Timestamp)
	return err
}

func (t *Postgres) GetPlaybackLog(guildID, ident, userID string, limit, offset int) ([]PlaybackLogEntry, error) {
	filter := "WHERE 'true'"
	var args []any
	args = append(args, limit, offset)

	if guildID != "" {
		args = append(args, guildID)
		filter += fmt.Sprintf(` AND "guildid" = $%d`, len(args))
	}

	if userID != "" {
		args = append(args, userID)
		filter += fmt.Sprintf(` AND "userid" = $%d`, len(args))
	}

	if ident != "" {
		args = append(args, ident)
		filter += fmt.Sprintf(` AND "sound" = $%d`, len(args))
	}

	rows, err := t.db.Query(fmt.Sprintf(`
		SELECT "id", "sound", "guildid", "userid", "timestamp"
		FROM playbacklog
		%s
		ORDER BY "timestamp" DESC
		LIMIT $1 OFFSET $2
	`, filter), args...)
	if err != nil {
		return nil, t.wrapErr(err)
	}

	var logs []PlaybackLogEntry
	for rows.Next() {
		var log PlaybackLogEntry
		err = rows.Scan(&log.Id, &log.Ident, &log.GuildID, &log.UserID, &log.Timestamp)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (t *Postgres) GetPlaybackLogSize() (int, error) {
	var n int
	err := t.db.QueryRow(`SELECT COUNT("id") FROM playbacklog`).Scan(&n)
	return n, t.wrapErr(err)
}

func (t *Postgres) GetPlaybackStats(guildID, userID string) ([]PlaybackStats, error) {
	filter := "WHERE 'true'"
	var args []any

	if guildID != "" {
		args = append(args, guildID)
		filter += fmt.Sprintf(` AND "guildid" = $%d`, len(args))
	}

	if userID != "" {
		args = append(args, userID)
		filter += fmt.Sprintf(` AND "userid" = $%d`, len(args))
	}

	rows, err := t.db.Query(fmt.Sprintf(`
		SELECT "sound", count("sound")
		FROM playbacklog
		%s
		GROUP BY "sound"
		ORDER BY "count" DESC
	`, filter), args...)
	if err != nil {
		return nil, t.wrapErr(err)
	}

	var logs []PlaybackStats
	for rows.Next() {
		var log PlaybackStats
		err = rows.Scan(&log.Ident, &log.Count)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (t *Postgres) GetAdmins() ([]string, error) {
	rows, err := t.db.Query(`
		SELECT "id" FROM users
		WHERE "admin" = 'true'
	`)
	if err != nil {
		return nil, t.wrapErr(err)
	}

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (t *Postgres) AddAdmin(userID string) error {
	return pg_setValue(t, "users", "admin", true, "id", userID)

}

func (t *Postgres) RemoveAdmin(userID string) error {
	return pg_setValue(t, "users", "admin", false, "id", userID)

}

func (t *Postgres) IsAdmin(userID string) (bool, error) {
	return pg_getValue[bool](t, "users", "admin", "id", userID)
}

func (t *Postgres) GetFavorites(userID string) ([]string, error) {
	rows, err := t.db.Query(`SELECT "sound" FROM user_favorites WHERE "userid" = $1`, userID)
	if err != nil {
		return nil, t.wrapErr(err)
	}

	var sounds []string
	for rows.Next() {
		var sound string
		err = rows.Scan(&sound)
		if err != nil {
			return nil, err
		}
		sounds = append(sounds, sound)
	}

	return sounds, nil
}

func (t *Postgres) AddFavorite(userID, ident string) error {
	_, err := t.db.Exec(`
		INSERT INTO user_favorites ("userid", "sound")
		VALUES ($1, $2)
	`, userID, ident)
	return err

}

func (t *Postgres) RemoveFavorite(userID, ident string) error {
	_, err := t.db.Exec(`DELETE FROM user_favorites WHERE "userid" = $1 AND "sound" = $2`,
		userID, ident)
	return t.wrapErr(err)

}

func (t *Postgres) GetApiKey(userID string) (string, error) {
	return pg_getValue[string](t, "users", "apikey", "id", userID)
}

func (t *Postgres) GetUserByApiKey(token string) (string, error) {
	return pg_getValue[string](t, "users", "id", "apikey", token)
}

func (t *Postgres) SetApiKey(userID, token string) error {
	return pg_setValue(t, "users", "apikey", token, "id", userID)
}

func (t *Postgres) RemoveApiKey(userID string) error {
	return pg_delete(t, "users", "id", userID)
}

func (t *Postgres) SetTwitchSettings(s TwitchSettings) error {
	filterInclude := strings.Join(s.Filters.Include, ",")
	filterExclude := strings.Join(s.Filters.Exclude, ",")
	blocklist := strings.Join(s.Blocklist, ",")

	res, err := t.db.Exec(`
		UPDATE twitchsettings
		SET "twitchusername" = $1,
		    "prefix" = $2,
			"ratelimitburst" = $3,
			"ratelimitreset" = $4,
			"filtersinclude" = $5,
			"filtersexclude" = $6,
			"blocklist" = $7
		WHERE "userid" = $8;
	`, s.TwitchUserName, s.Prefix, s.RateLimit.Burst, s.RateLimit.ResetSeconds,
		filterInclude, filterExclude, blocklist, s.UserID)
	if err != nil {
		return err
	}

	ar, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ar == 0 {
		_, err = t.db.Exec(`
			INSERT INTO twitchsettings (
				"userid", "twitchusername", "prefix", "ratelimitburst",
				"ratelimitreset", "filtersinclude", "filtersexclude", "blocklist"
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
		`, s.UserID, s.TwitchUserName, s.Prefix, s.RateLimit.Burst,
			s.RateLimit.ResetSeconds, filterInclude, filterExclude, blocklist)
		if err == nil {
			return err
		}
	}

	return err
}

func (t *Postgres) GetTwitchSettings(twitchname string) (TwitchSettings, error) {
	var (
		s                                       TwitchSettings
		filterInclude, filterExclude, blockList string
	)
	err := t.db.QueryRow(`
		SELECT "twitchusername", "prefix", "ratelimitburst", "ratelimitreset",
		       "filtersinclude", "filtersexclude", "blocklist"
		FROM twitchsettings
		WHERE "userid" = $1;
	`, twitchname).Scan(&s.TwitchUserName, &s.Prefix, &s.RateLimit.Burst,
		&s.RateLimit.ResetSeconds, &filterInclude, &filterExclude, &blockList)
	if err != nil {
		return TwitchSettings{}, t.wrapErr(err)
	}

	if filterInclude != "" {
		s.Filters.Include = strings.Split(filterInclude, ",")
	}
	if filterExclude != "" {
		s.Filters.Exclude = strings.Split(filterExclude, ",")
	}
	if blockList != "" {
		s.Blocklist = strings.Split(blockList, ",")
	}

	return s, nil
}

// --- Helpers ---

func (t *Postgres) tx(f func(*sql.Tx) error) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (t *Postgres) wrapErr(err error) error {
	if err != nil && err == sql.ErrNoRows {
		return dberrors.ErrNotFound
	}
	return err
}

func pg_getValue[TVal, TWv any](t *Postgres, table, vk, wk string, wv TWv) (TVal, error) {
	var v TVal
	err := t.db.QueryRow(
		fmt.Sprintf(`SELECT "%s" FROM %s WHERE "%s" = $1`, vk, table, wk), wv).Scan(&v)
	return v, t.wrapErr(err)
}

func pg_setValue[TVal, TWv any](t *Postgres, table, vk string, val TVal, wk string, wv TWv) error {
	res, err := t.db.Exec(
		fmt.Sprintf(`UPDATE %s SET "%s" = $1 WHERE "%s" = $2`, table, vk, wk), val, wv)
	if err != nil {
		return err
	}

	ar, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if ar == 0 {
		_, err = t.db.Exec(
			fmt.Sprintf(`INSERT INTO %s ("%s", "%s") VALUES ($1, $2)`, table, wk, vk), wv, val)
		if err == nil {
			return err
		}
	}

	return err
}

func pg_delete[TWv any](t *Postgres, table, wk string, wv TWv) error {
	_, err := t.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE "%s" = $1`, table, wk), wv)
	return t.wrapErr(err)
}
