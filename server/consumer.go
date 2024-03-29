package main

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

// Default configuration for WS connections.
// From https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// consumerHub suppose to organize cosumers in one place
// idea from: https://github.com/gorilla/websocket/blob/master/examples/chat/hub.go
type ConsumersHub struct {
	// Registered consumers.
	consumers map[*Consumer]bool

	// Register consumer from consumers.
	register chan *Consumer

	// Unregister consumer from consumers.
	unregister chan *Consumer
}

// NewConsumersHub create new hub.
func NewConsumersHub() *ConsumersHub {
	return &ConsumersHub{
		consumers:  make(map[*Consumer]bool),
		register:   make(chan *Consumer),
		unregister: make(chan *Consumer),
	}
}

// run listen to register/unregister signals and mutate consumers map.
func (h *ConsumersHub) run() {
	for {
		select {
		case consumer := <-h.register:
			h.consumers[consumer] = true
		case consumer := <-h.unregister:
			if _, ok := h.consumers[consumer]; ok {
				delete(h.consumers, consumer)
			}
		}
	}
}

func (h *ConsumersHub) Shutdown() {
	for consumer := range h.consumers {
		consumer.Shutdown()
	}
}

// Consumer represent User consumers.
type Consumer struct {
	// Channel name to listen.
	Channel string

	// Current user.
	User *User

	// WebSocket connection.
	ws *websocket.Conn

	// Channel stream.
	stream nats.JetStreamContext

	// Presence store.
	presence nats.KeyValue

	// ConsumersHub for self management.
	hub *ConsumersHub

	// Binded Logger for tracking internal process.
	Logger echo.Logger

	shutdown chan bool
}

// NewConsumer build new Consumer
// TODO: let's reduce number of agruments
func NewConsumer(channel string, user *User, ws *websocket.Conn, hub *ConsumersHub, stream nats.JetStreamContext, presence nats.KeyValue, logger echo.Logger) *Consumer {
	return &Consumer{
		Channel:  channel,
		User:     user,
		ws:       ws,
		hub:      hub,
		stream:   stream,
		presence: presence,
		Logger:   logger,
		shutdown: make(chan bool),
	}
}

// Register consumer in hub, managing join presence.
func (c Consumer) Register() {
	c.hub.register <- &c

	if c.Join() {
		data, err := json.Marshal(c.User)
		if err != nil {
			c.Logger.Errorf("Consumer JSON Marshall error: %v", err)
		}

		// Put user in presence store
		c.presence.Create(c.User.Id, data)
	}
}

// Unregister consumer from hub, managing leave presence.
func (c Consumer) Unregister() {
	c.hub.unregister <- &c

	if c.Leave() {
		// Remove user from presence store
		c.presence.Purge(c.User.Id)
	}
}

// Listen create new listener for incomming and ongoing channel messages for User consumer.
func (c Consumer) Listen() {
	defer c.Unregister()

	var subscription *nats.Subscription

	// Without nats.DeliverNew() the subscriber will get all messages in the stream
	// It allows us for free fill the chat history in UI, but I found that it hard to manage it
	// if it'll passible to pass last message time to WS reconnect state,
	// then we'll send to client only missed messages.
	//
	// Consumer description here used for track users subscription. See Unregister for details.
	subscription, err := c.stream.Subscribe(ChannelSubject(c.Channel), func(msg *nats.Msg) {
		c.ws.SetWriteDeadline(time.Now().Add(writeWait))
		c.ws.WriteMessage(websocket.TextMessage, msg.Data)
	}, nats.DeliverNew(), nats.Description(c.User.Id))

	if err != nil {
		c.Logger.Errorf("Consumer Listener error: %v", err)
		return
	}

	defer subscription.Unsubscribe()

	// Read pump
	go func() {
		c.ws.SetReadLimit(maxMessageSize)
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		for {
			// Read msg
			_, msg, err := c.ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					c.Logger.Errorf("Consumer read error: %v", err)
				}

				// no ws connection here anymore, exit
				break
			}

			c.PublishMsg(MessageSubject(c.Channel), NewChannelMessage(c.User, time.Now(), string(msg)))
		}
	}()

	// PingPong ticker
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	// Service stuff
	for {
		select {
		case <-c.shutdown:
			// Shutdown
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			c.ws.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

func (c Consumer) Shutdown() {
	close(c.shutdown)
}

// Join to channel if no consumers for User
func (c Consumer) Join() bool {
	consumers := GetConsumersDescription(c.stream)

	// if there still no more consumers for user
	// send leave msg to the channel
	for _, consumer := range consumers {
		if consumer == c.User.Id {
			return false
		}
	}

	c.PublishMsg(PresenceSubject(c.Channel), NewChannelJoinMessage(c.User, time.Now()))

	return true
}

// Leave from channel if no more consumers for User
func (c Consumer) Leave() bool {
	consumers := GetConsumersDescription(c.stream)

	// if there still no more consumers for user
	// send leave msg to the channel
	for _, consumer := range consumers {
		if consumer == c.User.Id {
			return false
		}
	}

	c.PublishMsg(PresenceSubject(c.Channel), NewChannelLeaveMessage(c.User, time.Now()))

	return true
}

// Publish channel message to the channel
func (c Consumer) PublishMsg(subject string, msg ChannelMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		c.Logger.Errorf("Consumer JSON Marshall error: %v", err)
	}

	c.stream.Publish(subject, data)
}
