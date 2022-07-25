package main

import (
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

var (
	upgrader = websocket.Upgrader{}
)

type channelHandler struct {
	// NATS connection
	nc *nats.Conn

	// NATS JetStream context
	stream nats.JetStreamContext

	// Consumers hub
	hub *ConsumersHub
}

func NewChannelHandler(nc *nats.Conn, hub *ConsumersHub) *channelHandler {
	// Use the JetStream context to produce and consumer messages
	// that have been persisted.
	js, err := nc.JetStream()
	if err != nil {
		return nil
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "CHITCHAT",
		Subjects: []string{"CHITCHAT.*.message"},
	})

	return &channelHandler{
		nc: nc,
		stream: js,
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

func Connect() (*nats.Conn, error) {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		return nil, err
	}

	return nc, nil
}