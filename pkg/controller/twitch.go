package controller

import (
	"github.com/zekroTJA/ratelimit"
	. "github.com/zekrotja/yuri69/pkg/models"
)

func (t *Controller) TwitchState(username string) (TwitchAPIState, error) {
	var res TwitchAPIState

	channel, instance, err := t.tw.GetConnectedChannel(username)
	if err != nil {
		return res, err
	}

	res.Channel = channel
	res.RateLimit.Burst = instance.Settings.RateLimit.Burst
	res.RateLimit.Burst = instance.Settings.RateLimit.ResetSeconds

	return res, nil
}

func (t *Controller) TwitchListSounds(username string, order string) ([]Sound, error) {
	_, instance, err := t.tw.GetConnectedChannel(username)
	if err != nil {
		return nil, err
	}

	return t.ListSounds(order,
		instance.Settings.Filters.Include,
		instance.Settings.Filters.Exclude)
}

func (t *Controller) TwitchPlay(username string, ident string) (bool, ratelimit.Reservation, error) {
	return t.tw.Play(username, ident)
}
