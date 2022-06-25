package main

import (
	"context"
	"flag"

	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/config"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/database"
	"github.com/zekrotja/yuri69/pkg/debug"
	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/player"
	"github.com/zekrotja/yuri69/pkg/storage"
	"github.com/zekrotja/yuri69/pkg/twitch"
	"github.com/zekrotja/yuri69/pkg/util"
	"github.com/zekrotja/yuri69/pkg/webserver"
)

var (
	fConfigFile = flag.String("c", "config.yml", "The location of the config file.")
	fLogLevel   = flag.Int("l", int(logrus.InfoLevel), "The log level (0 - 6)")
	fDebug      = flag.Bool("debug", false, "Enable debug mode")
	fVerbose    = flag.Bool("verbose", false, "Show code files in logs")
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

	logrus.SetReportCaller(*fVerbose)

	// --- Load Config ---
	cfg, err := config.Parse(*fConfigFile, "YURI_", config.DefaultConfig)
	if err != nil {
		logrus.WithError(err).Fatal("Config parsing failed")
	}
	logrus.WithField("file", *fConfigFile).Info("Config loaded")
	logrus.Debugf("Config Content: %+v", cfg)

	// --- Setup Database Module
	db, err := database.WrapCache(database.New(cfg.Database))
	if err != nil {
		logrus.WithError(err).Fatal("Database initialization failed")
	}
	defer func() {
		logrus.Info("Shutting down database connection ...")
		db.Close()
	}()
	logrus.WithField("typ", cfg.Database.Type).Info("Database initialized")

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

	// --- Setup Player ---
	pl, err := player.NewPlayer(cfg.Player, dc, st)
	if err != nil {
		logrus.WithError(err).Fatal("Player creation failed")
	}
	go func() {
		err = pl.ListenAndServeBlocking()
		if err != nil {
			logrus.WithError(err).Fatal("Player startup failed")
		}
	}()
	logrus.Info("Player initialized")
	defer func() {
		logrus.Info("Shutting down player ...")
		pl.Close()
	}()

	// --- Setup Controller ---
	ct, err := controller.New(db, st, pl, dc, cfg.Discord.OwnerID)
	if err != nil {
		logrus.WithError(err).Fatal("Controller initialization failed")
	}
	defer ct.Close()
	logrus.Info("Controller initialized")

	// --- Twitch Client ---
	if cfg.Twitch != nil {
		_, err = twitch.New(*cfg.Twitch, ct, cfg.Webserver.PublicAddress)
		if err != nil {
			logrus.WithError(err).Fatal("Twitch client creation failed")
		}
		logrus.Info("Twitch client initialized")
	}

	// --- Setup Web Server ---
	ws, err := webserver.New(cfg.Webserver, ct)
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

	_ = ct
}
