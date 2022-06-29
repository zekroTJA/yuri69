package twitch

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/yuri69/pkg/errs"
	"github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/rlhandler"
	"github.com/zekrotja/yuri69/pkg/util"
)

var prefixes = []string{"!yuri", "!y"}

type TwitchConfig struct {
	Username   string
	OAuthToken string
}

type Instance struct {
	Settings models.TwitchSettings

	rlh    rlhandler.RatelimitHandler
	userID string
}

type InternalEvent struct {
	Type    string
	Payload string
}

type PlayEvent struct {
	UserID  string
	Sound   string
	Filters models.GuildFilters
}

type Twitch struct {
	*util.EventBus[PlayEvent]

	client *twitch.Client

	instances *timedmap.TimedMap[string, *Instance]
	eventbus  *util.EventBus[InternalEvent]

	publicAddress string
}

func New(config TwitchConfig, publicAddress string) (*Twitch, error) {
	var t Twitch
	t.EventBus = util.NewEventBus[PlayEvent](100)
	t.client = twitch.NewClient(config.Username, config.OAuthToken)
	t.instances = timedmap.New[string, *Instance](1 * time.Hour)
	t.eventbus = util.NewEventBus[InternalEvent](100)
	t.publicAddress = publicAddress

	t.client.OnSelfJoinMessage(func(message twitch.UserJoinMessage) {
		logrus.WithField("channel", message.Channel).Info("Joined twitch channel")
		t.eventbus.Publish(InternalEvent{
			Type:    "join",
			Payload: message.Channel,
		})
	})

	t.client.OnSelfPartMessage(func(message twitch.UserPartMessage) {
		logrus.WithField("channel", message.Channel).Info("Left twitch channel")
		t.instances.Remove(message.Channel)
		t.eventbus.Publish(InternalEvent{
			Type:    "leave",
			Payload: message.Channel,
		})
	})

	t.client.OnPrivateMessage(t.onMessage)

	cErr := make(chan error)

	t.client.OnConnect(func() {
		logrus.Info("Twitch client connteded")
		cErr <- nil
	})

	go func() {
		cErr <- t.client.Connect()
	}()

	if err := <-cErr; err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Twitch) Update(s models.TwitchSettings) {
	instance := t.instances.GetValue(s.TwitchUserName)
	if instance == nil {
		return
	}

	instance.Settings = s
	instance.rlh.Update(s.RateLimit.Burst, time.Duration(s.RateLimit.ResetSeconds)*time.Second)
}

func (t *Twitch) Join(userid string, s models.TwitchSettings) error {
	instance := t.instances.GetValue(s.TwitchUserName)

	if instance != nil {
		if instance.userID != userid {
			return errors.New("already joined same channel by another user")
		}

		t.Update(s)
		return nil
	}

	instance = &Instance{
		userID:   userid,
		Settings: s,
		rlh:      rlhandler.New(s.RateLimit.Burst, time.Duration(s.RateLimit.ResetSeconds)*time.Second),
	}

	t.instances.Set(s.TwitchUserName, instance, 24*time.Hour, func(value *Instance) {
		t.client.Depart(s.TwitchUserName)
	})

	ch, unsub := t.eventbus.Subscribe()
	defer unsub()

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	t.client.Join(s.TwitchUserName)

	for {
		select {
		case <-timeout.C:
			return errors.New("timed out")
		case e := <-ch:
			if e.Type == "join" && e.Payload == s.TwitchUserName {
				return nil
			}
		}
	}
}

func (t *Twitch) Leave(twitchname string) error {
	instance := t.instances.GetValue(twitchname)
	if instance == nil {
		return errs.WrapUserError("no instance initialized")
	}

	ch, unsub := t.eventbus.Subscribe()
	defer unsub()

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	t.client.Depart(instance.Settings.TwitchUserName)

	for {
		select {
		case <-timeout.C:
			return errors.New("timed out")
		case e := <-ch:
			if e.Type == "leave" && e.Payload == instance.Settings.TwitchUserName {
				return nil
			}
		}
	}
}

func (t *Twitch) Joined(userid string) bool {
	return t.instances.Contains(userid)
}

func (t *Twitch) GetConnectedChannel(username string) (string, *Instance, error) {
	for channel, instance := range t.instances.Snapshot() {
		users, err := t.client.Userlist(channel)
		if err != nil {
			return "", instance, err
		}
		if util.Contains(users, username) {
			return channel, instance, nil
		}
	}

	return "", nil, errs.WrapUserError("not connected to any twitch chat")
}

func (t *Twitch) Play(username, ident string) (bool, ratelimit.Reservation, error) {
	_, instance, err := t.GetConnectedChannel(username)
	if err != nil {
		return false, ratelimit.Reservation{}, err
	}
	ok, res := t.play(instance, username, ident)
	return ok, res, nil
}

func (t *Twitch) play(instance *Instance, username string, ident string) (bool, ratelimit.Reservation) {
	ok, res := instance.rlh.Get(username).Reserve()
	if ok {
		t.Publish(PlayEvent{
			UserID:  instance.userID,
			Sound:   ident,
			Filters: instance.Settings.Filters,
		})
	}
	return ok, res
}

func (t *Twitch) onMessage(message twitch.PrivateMessage) {
	msg := strings.TrimSpace(message.Message)

	instance := t.instances.GetValue(message.Channel)
	if instance == nil {
		return
	}

	prefix := instance.Settings.Prefix

	if !strings.HasPrefix(message.Message, prefix) {
		return
	}

	if msg == prefix {
		t.client.Reply(message.Channel, message.ID,
			"You can play sounds live on stream when using on of the following commands: `!y rand` - Play a random sound | `!y <sound>` - Paly a specific sound | `!y list` - Get a list of available sounds.")
		return
	}

	if len(msg) <= len(prefix) {
		return
	}

	split := strings.Split(msg[len(prefix)+1:], " ")
	invoke := strings.ToLower(split[0])

	switch invoke {

	case "list", "sounds", "ls":
		t.client.Reply(message.Channel, message.ID, fmt.Sprintf(
			"Here you can find a list of available sounds (Yes, this will be improved some day 😅): %s/api/v1/public/twitch/sounds",
			t.publicAddress))

	case "r", "rand", "random":
		t.play(instance, message.User.Name, "")

	default:
		t.play(instance, message.User.Name, invoke)
	}
}
