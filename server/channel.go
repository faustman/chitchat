package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// For dev purpose only: if u want to skip origin checking
		// CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// ChannelMessage unsing for communication in channel.
type ChannelMessage struct {
	Type     string    `json:"type"`
	FromUser *User     `json:"from_user"`
	SentAt   time.Time `json:"sent_at"`
	Text     string    `json:"text,omitempty"`
}

// NewChannelMessage build new text ChannelMessage.
func NewChannelMessage(user *User, sentAt time.Time, text string) ChannelMessage {
	return ChannelMessage{
		Type:     "message",
		FromUser: user,
		SentAt:   sentAt,
		Text:     text,
	}
}

// NewChannelMessage build new join ChannelMessage.
func NewChannelJoinMessage(user *User, sentAt time.Time) ChannelMessage {
	return ChannelMessage{
		Type:     "join",
		FromUser: user,
		SentAt:   sentAt,
	}
}

// NewChannelMessage build new leave ChannelMessage.
func NewChannelLeaveMessage(user *User, sentAt time.Time) ChannelMessage {
	return ChannelMessage{
		Type:     "leave",
		FromUser: user,
		SentAt:   sentAt,
	}
}

// channelHandler handle channel stuff.
type channelHandler struct {
	// NATS JetStream context
	stream nats.JetStreamContext

	// Consumers hub
	hub *ConsumersHub
}

// NewChannelHandler build new channelHandler.
func NewChannelHandler(stream nats.JetStreamContext, hub *ConsumersHub) *channelHandler {
	return &channelHandler{
		stream: stream,
		hub:    hub,
	}
}

// Listen to incoming websocket connection and register new consumer.
func (h channelHandler) Listen(c echo.Context) error {
	// Upgrade to ws
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	auth := ExtactAuth(c)

	presence, err := GetPresenceBucket(h.stream, auth.Channel)
	if err != nil {
		return err
	}

	consumer := NewConsumer(auth.Channel, auth.User, ws, h.hub, h.stream, presence, c.Logger())

	consumer.Register()

	go consumer.Listen()

	return nil
}

// GetMessages extract text messages from stream and respond with JSON
func (h channelHandler) GetMessages(c echo.Context) error {
	auth := ExtactAuth(c)

	opts := []nats.SubOpt{nats.OrderedConsumer()}

	st := c.QueryParam("start_time")

	// fetch("/messages?token=" + token + "&start_time=" + Math.round(1658877495926 / 1000)).then((r) => r.json()).then(console.log).catch(console.error)
	if len(st) > 0 {
		unix, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			return err
		}

		opts = append(opts, nats.StartTime(time.Unix(unix, 0)))
	}

	sub, err := h.stream.SubscribeSync(MessageSubject(auth.Channel), opts...)
	if err != nil {
		return err
	}

	info, err := h.stream.StreamInfo(StreamName)
	if err != nil {
		return err
	}
	// We need to know how much to iterate until the end
	lastSeq := info.State.LastSeq

	// Geeting all messages from the stream
	var messages []ChannelMessage

	for i := uint64(0); i < lastSeq; i++ {
		m, err := sub.NextMsg(1 * time.Second)

		if err == nil {
			msg := ChannelMessage{}

			if err := json.Unmarshal(m.Data, &msg); err != nil {
				return err
			}

			messages = append(messages, msg)
		} else {
			break
		}
	}

	sub.Unsubscribe()

	return c.JSON(http.StatusCreated, echo.Map{
		"messages": messages,
	})
}

// GetUsers from channel presence store
func (h channelHandler) GetUsers(c echo.Context) error {
	auth := ExtactAuth(c)

	presence, err := GetPresenceBucket(h.stream, auth.Channel)
	if err != nil {
		return err
	}

	var users []User

	uids, _ := presence.Keys()

	for _, uid := range uids {
		entry, err := presence.Get(uid)
		if err != nil {
			return err
		}

		user := User{}
		if err := json.Unmarshal(entry.Value(), &user); err != nil {
			return err
		}

		users = append(users, user)
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"users": users,
	})
}
