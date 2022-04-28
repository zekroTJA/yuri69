package player

import (
	"io"
	"net/http"
	"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/yuri69/pkg/discord"
	"github.com/zekrotja/yuri69/pkg/generic"
	"github.com/zekrotja/yuri69/pkg/lavalink"
	"github.com/zekrotja/yuri69/pkg/storage"
)

type Manager struct {
	dc *discord.Discord
	st storage.IStorage
	ll *lavalink.Lavalink

	players generic.SyncMap[string, *Player]

	router *routing.Router
	server *http.Server

	loadedFiles *timedmap.TimedMap[string, io.ReadCloser]
}

func NewManager(dc *discord.Discord, st storage.IStorage, ll *lavalink.Lavalink) (*Manager, error) {
	var t Manager

	t.dc = dc
	t.st = st
	t.ll = ll

	t.router = routing.New()
	t.router.Get("/file/<id>", t.handleGetFile)

	t.loadedFiles = timedmap.New[string, io.ReadCloser](5 * time.Minute)

	return &t, nil
}

func (t *Manager) ListenAndServeBlocking() error {
	t.server = &http.Server{
		Addr:    "0.0.0.0:6969",
		Handler: t.router,
	}
	return t.server.ListenAndServe()
}

func (t *Manager) GetPlayer(guildID, channelID string) (*Player, error) {
	var err error

	p, ok := t.players.Load(guildID)
	if !ok {
		p, err = t.createPlayer(guildID, channelID)
		if err != nil {
			return nil, err
		}
		t.players.Store(guildID, p)
	}

	return p, nil
}

func (t *Manager) createPlayer(guildID, channelID string) (*Player, error) {
	err := t.dc.Session().ChannelVoiceJoinManual(guildID, channelID, false, true)
	if err != nil {
		return nil, err
	}

	p := &Player{
		guildID: guildID,
		mgr:     t,
	}

	return p, nil
}

func (t *Manager) loadFile(name string) (string, error) {
	rc, _, err := t.st.GetObject("sounds", name)
	if err != nil {
		return "", err
	}

	id := xid.New().String()
	t.loadedFiles.Set(id, rc, 5*time.Minute)
	return id, nil
}

func (t *Manager) handleGetFile(ctx *routing.Context) error {
	id := ctx.Param("id")

	r := t.loadedFiles.GetValue(id)
	if r == nil {
		return ctx.WriteWithStatus("", http.StatusBadRequest)
	}
	defer r.Close()

	_, err := io.Copy(ctx.Response, r)
	if err != nil {
		logrus.WithError(err).Error("Writing file to response body failed")
		return ctx.WriteWithStatus("", http.StatusInternalServerError)
	}

	ctx.Response.WriteHeader(http.StatusOK)
	return nil
}
