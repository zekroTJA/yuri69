package ws

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	. "github.com/zekrotja/yuri69/pkg/models"
)

const authenticationDeadline = 15 * time.Second

type subscription struct {
	id  string
	hub *Hub

	onClose func(external bool)

	conn   *websocket.Conn
	userId string

	cOut   chan []byte
	cClose chan struct{}
}

func newSubscription(h *Hub, conn *websocket.Conn) *subscription {
	t := &subscription{
		id:     xid.New().String(),
		hub:    h,
		conn:   conn,
		cOut:   make(chan []byte),
		cClose: make(chan struct{}),
	}

	go t.reader()
	go t.writer()

	t.awaitAuthentication()

	return t
}

func (t *subscription) Id() string {
	return t.id
}

func (t *subscription) UserId() string {
	return t.userId
}

func (t *subscription) OnClose(f func(external bool)) {
	t.onClose = f
}

func (t *subscription) Close() error {
	t.onClose(false)
	t.cClose <- struct{}{}
	return nil
}

func (t *subscription) Publish(data []byte) {
	if !t.authenticated() {
		return
	}
	t.publish(data)
}

// --- Internal ---

func (t *subscription) publish(data []byte) {
	t.cOut <- data
}

func (t *subscription) publishEvent(e Event[any]) {
	data, _ := json.Marshal(e)
	t.publish(data)
}

func (t *subscription) awaitAuthentication() {
	t.publishEvent(Event[any]{
		Type: "authpromp",
		Payload: EventAuthPromptPayload{
			Deadline:  time.Now().Add(authenticationDeadline),
			TokenType: "accesstoken",
		},
	})
	time.AfterFunc(authenticationDeadline, func() {
		if !t.authenticated() {
			t.publishEvent(Event[any]{
				Type: "authpromptfailed",
			})
			t.Close()
		}
	})
}

func (t *subscription) writer() {
	for {
		select {

		case <-t.cClose:
			t.conn.WriteMessage(websocket.CloseMessage, nil)
			return

		case msg := <-t.cOut:
			err := t.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logrus.WithField("id", t.id).WithError(err).Error("Writing to web socket conn failed")
				t.Close()
				return
			}
		}
	}
}

func (t *subscription) reader() {
	for {
		_, message, err := t.conn.ReadMessage()
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err) {
				logrus.WithField("id", t.id).WithError(err).Error("Unexpected WS error")
			}
			t.onClose(true)
			break
		}
		go t.onMessage(message)
	}
}

func (t *subscription) onMessage(msg []byte) {
	typ, err := unmarshal[EventType](msg)
	if err != nil {
		t.publishEvent(WrapErrorEvent(err, http.StatusBadRequest))
		return
	}

	switch strings.ToLower(typ.Type) {
	case "auth":
		if t.authenticated() {
			t.publishEvent(WrapErrorEvent(errors.New("already authenticated"), http.StatusBadRequest))
			return
		}
		authPayload, err := unmarshal[Event[EventAuthRequest]](msg)
		if err != nil {
			t.publishEvent(Event[any]{
				Type: "authpromptfailed",
				Payload: StatusModel{
					Status:  http.StatusBadRequest,
					Message: err.Error(),
				},
			})
			t.Close()
			return
		}
		if authPayload.Payload.Token == "" {
			t.publishEvent(Event[any]{
				Type: "authpromptfailed",
				Payload: StatusModel{
					Status:  http.StatusBadRequest,
					Message: "token is empty",
				},
			})
			t.Close()
			return
		}
		claims, err := t.hub.authHandler.CheckAuthRaw(authPayload.Payload.Token)
		if err != nil {
			t.publishEvent(Event[any]{
				Type: "authpromptfailed",
				Payload: StatusModel{
					Status:  http.StatusBadRequest,
					Message: err.Error(),
				},
			})
			t.Close()
			return
		}
		t.userId = claims.UserID
		payload, err := t.hub.ct.GetCurrentState(t.userId)
		if err != nil {
			logrus.WithError(err).Error("Getting current state failed")
		}
		t.publishEvent(Event[any]{Type: "authok", Payload: payload})
	}
}

func (t *subscription) authenticated() bool {
	return t.userId != ""
}

func unmarshal[T any](data []byte) (T, error) {
	var v T
	err := json.Unmarshal(data, &v)
	return v, err
}
