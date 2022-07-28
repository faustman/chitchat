// This code is for testing NATS JS feature for getting all messages from the stream
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     "STORE",
		Subjects: []string{"STORE.*.message"},
		NoAck: true,
	})

	subj := "STORE.mystore2.message"

	fillStream := func(w http.ResponseWriter, req *http.Request) {

		for i := 0; i < 5; i++ {
			js.Publish(subj, []byte(fmt.Sprintf("Hello %d!", i)))
		}

		io.WriteString(w, "OK!\n")
	}

	http.HandleFunc("/fill", fillStream)

	test1 := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sub, err := js.SubscribeSync(subj)
		if err != nil {
			log.Fatal(err)
		}

		info, err := js.StreamInfo("STORE")
		if err != nil {
			log.Fatal(err)
		}
		lastSeq := info.State.LastSeq

		var messages []string

		for i := uint64(0); i < lastSeq; i++  {
			m, err := sub.NextMsg(1 * time.Second)

			fmt.Println(m)

			if err == nil {
				messages = append(messages, string(m.Data))
			} else {
				fmt.Println(err)
				break
			}
		}

		resp := map[string]any{"messages": messages, "pending": lastSeq}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		w.Write(jsonResp)
	}

	http.HandleFunc("/test1", test1)

	log.Fatal(http.ListenAndServe(":8091", nil))
}
