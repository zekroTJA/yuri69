package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/traefik/paerser/env"
	"github.com/traefik/paerser/file"
	"github.com/zekrotja/yuri69/pkg/util"
)

func Parse[T any](cfgFile string, envPrefix string, def ...T) (cfg T, err error) {
	cfg = util.Opt(def)

	if err = file.Decode(cfgFile, &cfg); err != nil && !os.IsNotExist(err) {
		return
	}

	godotenv.Load()
	if err = env.Decode(os.Environ(), envPrefix, &cfg); err != nil {
		return
	}

	return
}
