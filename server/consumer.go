package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
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

func NewConsumersHub() *ConsumersHub {
	return &ConsumersHub{
		consumers: make(map[*Consumer]bool),
		register: make(chan *Consumer),
		unregister: make(chan *Consumer),
	}
}

func (h *ConsumersHub) run() {
	for {
		select {
		case consumer := <- h.register:
			h.consumers[consumer] = true
		case consumer := <- h.unregister:
			if _, ok := h.consumers[consumer]; ok {
				delete(h.consumers, consumer)
			}
		}
	}
}

type Consumer struct {
	Channel string
	User *User
	ws *websocket.Conn
	stream nats.JetStreamContext
	hub *ConsumersHub
}

func (c Consumer) Register() {
	c.hub.register <- &c

	// TODO: notify new user
}

func (c Consumer) Unregister() {
	c.hub.unregister <- &c

	// TODO: notify remove user
}

// TODO: maybe better to name Listen ?
func (c Consumer) Listen() {
	fmt.Println("Listen")

	msgSubject := fmt.Sprintf("CHITCHAT.%s.message", c.Channel)

	c.stream.Subscribe(msgSubject,  func(msg *nats.Msg) {
		c.ws.WriteMessage(websocket.TextMessage, msg.Data)
	})

	for {
		fmt.Println("Read")
		// Read
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			// TODO: unregister
			break
		}

		c.stream.Publish(msgSubject, msg)
	}

	// TODO: for user precense
	// kv, err := c.stream.KeyValue(c.Channel)
	// if err != nil {
	// 	panic(err)
	// }

	// kv.WatchAll()
}