package main

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

const StreamName = "CHITCHAT"

// NewStream connect to Nats and create new JetStream stream
// There is a few models that we could use Streams for chat:
// - Stream per channel - could be usefull with massive channels, expensive
// - One stream, Subject per channel - lightwight solution, but could reach JetStream limits
// See https://docs.nats.io/using-nats/developer/develop_jetstream/model_deep_dive
func NewStream() (nats.JetStreamContext, error) {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		return nil, err
	}

	// Use the JetStream context to produce and consumer messages
	// that have been persisted.
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Adding a stream is an idempotent function
	// which means that if a stream does not exist,
	// it will be created, and if a stream already exists,
	// then the add operation will succeed only
	// if the existing stream matches exactly the attributes specified in the 'add' call.
	// Default Retention policy: LimitsPolicy
	js.AddStream(&nats.StreamConfig{
		Name:     StreamName,
		Subjects: []string{StreamName + ".*.*"},
		NoAck: true,
	})

	return js, nil
}

func GetConsumersDescription(js nats.JetStreamContext) []string {
	cs := js.ConsumersInfo(StreamName)

	var cd []string

	for c := range cs {
		cd = append(cd, c.Config.Description)
	}

	return cd
}

// GetPresenceBucket create if not exists and return KeyValue presence bucket.
func GetPresenceBucket(js nats.JetStreamContext, channel string) (nats.KeyValue, error) {
	return js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: channel + "-presence",
	})
}

// Generate subject based on channel and message
func subject(c, m string) string {
	return fmt.Sprintf("%s.%s.%s", StreamName, c, m)
}

// PresenceSubject for presence messages
func PresenceSubject(channel string) string {
	return subject(channel, "presence")
}

// MessageSubject for text messages
func MessageSubject(channel string) string {
	return subject(channel, "message")
}

// ChannelSubject for all channel messages
func ChannelSubject(channel string) string {
	return subject(channel, "*")
}