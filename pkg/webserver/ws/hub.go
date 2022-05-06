package ws

import (
	"encoding/json"
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/zekrotja/yuri69/pkg/controller"
	"github.com/zekrotja/yuri69/pkg/generic"
	. "github.com/zekrotja/yuri69/pkg/models"
	"github.com/zekrotja/yuri69/pkg/util"
	"github.com/zekrotja/yuri69/pkg/webserver/auth"
)

type Hub struct {
	upgrader    websocket.Upgrader
	authHandler *auth.AuthHandler
	ct          *controller.Controller

	subs generic.SyncMap[string, *subscription]
}

func NewHub(authHandler *auth.AuthHandler, ct *controller.Controller) *Hub {
	return &Hub{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		authHandler: authHandler,
		ct:          ct,
	}
}

func (t *Hub) Upgrade(ctx *routing.Context) error {
	conn, err := t.upgrader.Upgrade(ctx.Response, ctx.Request, nil)
	if err != nil {
		logrus.WithError(err).Error("Upgrading web socket request failed")
		return err
	}

	sub := newSubscription(t, conn)
	logrus.WithField("id", sub.Id()).Debug("WS connection opened")
	sub.OnClose(func(external bool) {
		logrus.
			WithField("id", sub.Id()).
			WithField("external", external).
			Debug("WS connection closed")
		t.subs.Delete(sub.Id())
	})

	t.subs.Store(sub.Id(), sub)

	return nil
}

func (t *Hub) Broadcast(event Event[any]) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	t.subs.Range(func(_ string, sub *subscription) bool {
		sub.Publish(data)
		return true
	})

	return nil
}

func (t *Hub) BroadcastScoped(event Event[any], userIDs []string) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	t.subs.Range(func(_ string, sub *subscription) bool {
		if sub.UserId() != "" && util.Contains(userIDs, sub.UserId()) {
			sub.Publish(data)
		}
		return true
	})

	return nil
}
