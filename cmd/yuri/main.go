package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/config"
	"github.com/zekrotja/yuri69/pkg/debug"
	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/lavalink"
	"github.com/zekrotja/yuri69/pkg/player"
	"github.com/zekrotja/yuri69/pkg/storage"
	"github.com/zekrotja/yuri69/pkg/util"
	"github.com/zekrotja/yuri69/pkg/webserver"
)

var (
	fConfigFile = flag.String("c", "config.yml", "The location of the config file.")
	fLogLevel   = flag.Int("l", int(logrus.InfoLevel), "The log level (0 - 6)")
	fDebug      = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
	// --- Parse Flags ---
	flag.Parse()

	// --- Setup Formatter ---
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006/01/02 15:04:05 MST",
	})
	logrus.SetLevel(logrus.Level(*fLogLevel))

	debug.SetEnabled(*fDebug)
	if debug.Enabled() {
		logrus.Warn(
			"DEBUG MODE IS ENABLED! Using debug mode in production is a severe security risk!")
	}

	// --- Load Config ---
	cfg, err := config.Parse(*fConfigFile, "YURI_", config.DefaultConfig)
	if err != nil {
		logrus.WithError(err).Fatal("Config parsing failed")
	}
	logrus.WithField("file", *fConfigFile).Info("Config loaded")
	logrus.Debugf("Config Content: %+v", cfg)

	// --- Setup Storage Module ---
	st, err := storage.New(cfg.Storage)
	if err != nil {
		logrus.WithError(err).Fatal("Storage initialization failed")
	}
	logrus.WithField("typ", cfg.Storage.Type).Info("Storage initialized")

	// --- Setup Discord Session ---
	dc, err := discord.New(cfg.Discord)
	if err != nil {
		logrus.WithError(err).Fatal("Discord initialization failed")
	}
	err = dc.Open()
	if err != nil {
		logrus.WithError(err).Fatal("Discord connection failed")
	}
	logrus.Info("Discord connection initialized")
	defer func() {
		logrus.Info("Shutting down Discord connection ...")
		dc.Close()
	}()

	// --- Setup Lavalink Connection ---
	ll, err := lavalink.New(cfg.Lavalink, dc)
	if err != nil {
		logrus.WithError(err).Fatal("Lavalink connection failed")
	}
	logrus.Info("Lavalink connection initialized")
	defer func() {
		logrus.Info("Shutting down Lavalink connection ...")
		ll.Close()
	}()

	// --- Setup Player Manager ---
	mgr, err := player.NewPlayer(cfg.Player, dc, st, ll)
	if err != nil {
		logrus.WithError(err).Fatal("Manager creation failed")
	}
	go func() {
		err = mgr.ListenAndServeBlocking()
		if err != nil {
			logrus.WithError(err).Fatal("Manager file provider startup failed")
		}
	}()
	logrus.Info("Player manager initialized")

	fmt.Println(mgr.Play("526196711962705925", "959799153754509312", "mothers"))
	time.Sleep(3 * time.Second)
	mgr.Destroy("526196711962705925")

	// --- Setup Web Server ---
	ws, err := webserver.New(cfg.Webserver)
	if err != nil {
		logrus.WithError(err).Fatal("Webserver initialization failed")
	}
	go func() {
		err = ws.ListenAndServeBlocking()
		if err != nil {
			logrus.WithError(err).Fatal("Webserver startup failed")
		}
	}()
	logrus.WithField("addr", cfg.Webserver.BindAddress).Info("Webserver started")

	// Block either until passed context is done
	// or an exit signal was received.
	util.Block(context.Background())

	_ = st
	_ = ll
}
