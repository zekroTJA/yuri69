package controllers

import (
	"encoding/json"
	"time"

	routing "github.com/zekrotja/ozzo-routing/v2"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/webserver/middleware"
)

type publicController struct {
	ct *controller.Controller
}

func NewPublicController(r *routing.RouteGroup, ct *controller.Controller) {
	t := publicController{ct: ct}
	r.Get("/twitch/sounds", middleware.Cache(10*time.Minute, false, true), t.handleGetTwitchSounds)
	return
}

func (t *publicController) handleGetTwitchSounds(ctx *routing.Context) error {
	sounds, err := t.ct.ListSounds(string(models.SortOrderCreated), nil, []string{"nsft"})
	if err != nil {
		return err
	}

	soundNames := make([]string, len(sounds))
	for i, sound := range sounds {
		soundNames[i] = sound.Uid
	}

	raw, err := json.MarshalIndent(soundNames, "", "  ")
	if err != nil {
		return err
	}

	ctx.Response.Header().Set("Content-Type", "application/json")
	_, err = ctx.Response.Write(raw)

	return err
}
