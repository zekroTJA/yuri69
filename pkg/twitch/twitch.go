package twitch

import (
	"fmt"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
	"github.com/sirupsen/logrus"
	"github.com/zekroTJA/ratelimit"
	"github.com/zekroTJA/timedmap"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/util"
)

var prefixes = []string{"!yuri", "!y"}

type TwitchConfig struct {
	Username        string
	OAuthToken      string
	JoinChannels    []string
	ImpersonateUser string
}

type Twitch struct {
	client *twitch.Client
	ct     *controller.Controller

	rateLimits *timedmap.TimedMap[string, *ratelimit.Limiter]

	impersonatedUser string
	publicAddress    string
}

func New(config TwitchConfig, ct *controller.Controller, publicAddress string) (Twitch, error) {
	var t Twitch
	t.client = twitch.NewClient(config.Username, config.OAuthToken)
	t.rateLimits = timedmap.New[string, *ratelimit.Limiter](5 * time.Minute)
	t.ct = ct
	t.impersonatedUser = config.ImpersonateUser
	t.publicAddress = publicAddress

	t.client.Join(config.JoinChannels...)

	t.client.OnSelfJoinMessage(func(message twitch.UserJoinMessage) {
		logrus.WithField("channel", message.Channel).Info("Joined twitch channel")
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
		return Twitch{}, err
	}

	return t, nil
}

func (t Twitch) onMessage(message twitch.PrivateMessage) {
	msg := strings.TrimSpace(message.Message)

	ok, prefix := util.StartsWithAny(msg, prefixes)
	if !ok {
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
		if !t.rateLimit(message.User.ID) {
			return
		}
		t.ct.PlayRandom(t.impersonatedUser, nil, []string{"nsft"}) // TODO: un-hardcode filters

	default:
		if !t.rateLimit(message.User.ID) {
			return
		}
		err := t.ct.Play(t.impersonatedUser, invoke)
		if err != nil {
			logrus.WithError(err).WithField("invoke", invoke).Error("Failed playing sound via twitch")
		}
	}
}

func (t Twitch) rateLimit(userID string) bool {
	rl := t.rateLimits.GetValue(userID)
	if rl == nil {
		rl = ratelimit.NewLimiter(30*time.Second, 3)
		t.rateLimits.Set(userID, rl, 3*30*time.Second)
	}
	return rl.Allow()
}
