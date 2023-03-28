-- +goose Up

CREATE TABLE IF NOT EXISTS sounds (
  uid VARCHAR(30) NOT NULL,
  displayname TEXT NOT NULL DEFAULT '',
  created TIMESTAMP NOT NULL,
  creatorid TEXT NOT NULL,
  PRIMARY KEY (uid)
);

CREATE TABLE IF NOT EXISTS sounds_tags (
  id INT GENERATED ALWAYS AS IDENTITY,
  sound VARCHAR(30) NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_tags
    FOREIGN KEY (sound)
    REFERENCES sounds(uid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS guilds (
  id VARCHAR(32) NOT NULL,
  volume INT NOT NULL DEFAULT '0',
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS guild_filters (
  id INT GENERATED ALWAYS AS IDENTITY,
  guildid VARCHAR(32) NOT NULL,
  exclude BOOLEAN NOT NULL,
  tag TEXT NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(32) NOT NULL,
  fasttrigger TEXT NOT NULL DEFAULT '',
  admin BOOLEAN NOT NULL DEFAULT 'false',
  apikey TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS user_favorites (
  id INT GENERATED ALWAYS AS IDENTITY,
  userid VARCHAR(32) NOT NULL,
  sound VARCHAR(30) NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT fk_favorites
    FOREIGN KEY(sound)
    REFERENCES sounds(uid)
    ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS playbacklog (
  id VARCHAR(20) NOT NULL,
  sound VARCHAR(30) NOT NULL,
  guildid VARCHAR(32) NOT NULL,
  userid VARCHAR(32) NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS twitchsettings (
  userid VARCHAR(32) NOT NULL,
  twitchusername VARCHAR(30) NOT NULL DEFAULT '',
  prefix VARCHAR(20) NOT NULL DEFAULT '',
  ratelimitburst INT NOT NULL DEFAULT '0',
  ratelimitreset INT NOT NULL DEFAULT '0',
  filtersinclude TEXT NOT NULL DEFAULT '',
  filtersexclude TEXT NOT NULL DEFAULT '',
  blocklist TEXT NOT NULL DEFAULT '',
  PRIMARY KEY (userid)
);

-- +goose Down

DROP TABLE twitchsettings;
DROP TABLE playbacklog;
DROP TABLE user_favorites;
DROP TABLE users;
DROP TABLE guild_filters;
DROP TABLE guilds;
DROP TABLE sounds_tags;
DROP TABLE sounds;
