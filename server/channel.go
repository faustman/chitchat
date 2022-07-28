package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

var (
	upgrader = websocket.Upgrader{}
)

// ChannelMessage unsing for communication in channel
type ChannelMessage struct {
	Type string `json:"type"`
	FromUser *User `json:"from_user"`
	SentAt time.Time `json:"sent_at"`
	Text string `json:"text"`
}

func NewChannelMessage(user *User, sentAt time.Time, text string) ChannelMessage {
	return ChannelMessage{
		Type: "message",
		FromUser: user,
		SentAt: sentAt,
		Text: text,
	}
}

type channelHandler struct {
	// NATS JetStream context
	stream nats.JetStreamContext

	// Consumers hub
	hub *ConsumersHub
}

func NewChannelHandler(stream nats.JetStreamContext, hub *ConsumersHub) *channelHandler {
	return &channelHandler{
		stream: stream,
		hub: hub,
	}
}

// const pingPeriod = 10 * time.Second
// const writeWait = 10 * time.Second

// Listen to incoming websocket connection and register new consumer.
func (h channelHandler) Listen(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	token := c.Get("token").(*jwt.Token)
	auth := token.Claims.(*Auth)

	consumer := &Consumer{
		Channel: auth.Channel,
		User: auth.User,
		ws: ws,
		hub: h.hub,
		stream: h.stream,
	}

	consumer.Register()

	go consumer.Listen()

	return nil
}

func (h channelHandler) GetMessages(c echo.Context) error {
	token := c.Get("token").(*jwt.Token)
	auth := token.Claims.(*Auth)

	msgSubject := fmt.Sprintf("CHITCHAT.%s.message", auth.Channel)

	opts := []nats.SubOpt{}

	st := c.QueryParam("start_time")

	// fetch("/messages?token=" + token + "&start_time=" + Math.round(1658877495926 / 1000)).then((r) => r.json()).then(console.log).catch(console.error)
	if len(st) > 0 {
		unix, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			return nil
		}

		opts = append(opts, nats.StartTime(time.Unix(unix, 0)))
	}

	sub, err := h.stream.SubscribeSync(msgSubject, opts...)
	if err != nil {
		return err
	}

	info, err := h.stream.StreamInfo("CHITCHAT")
	if err != nil {
		return nil
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
				return nil
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