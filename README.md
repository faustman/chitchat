# CHITCHAT

Modern chat in GO

## Idea

This project intends to reveal Go's best parts in real-time communication and discover new opportunities in modern streaming protocols.

### Protocol + transport

Having a successful experience with a few chat projects using WebSockets (even Long Pulling) and Realtime Databases (firebase) technologies and realizing that each of them works great for its purpose, I decided to dig deeper and find out what modern technologies can offer me.

- [WebSocket](https://websockets.spec.whatwg.org/) - pretty old protocol (RFC 6455, 2011), Living Standard. Supported across amlost all web platforms, Bidirectional.
- HTTP/2 + [SSE](https://html.spec.whatwg.org/multipage/server-sent-events.html#server-sent-events) - newest technology. Semi-Bidirectional. Could be good for small chat. Limited by open connections
- gRPC-Web - it's a promising solution since gRPC has bidirectional streams, but my investigation leads me to the fact that it's too raw to use it in production. See [gRPC-Web Streaming Roadmap](https://github.com/grpc/grpc-web/blob/master/doc/streaming-roadmap.md)
- WebRTC (p2p) - kind a also interesting to use, but requires some orchestration. But potentialy could save some infrastructure.
- [WebTransport](https://w3c.github.io/webtransport/) - most modern API that will sits on top of HTTP/3 and potentialy could replace WebSockets in the future. In draft
- Socket.IO or Cetrifugo - looks nice, especialy for prototyping. Mostly brain free solutions, setup and forget. Not fit for project purpose.

Seems that **WebSocket** is the most reliable and stable protocol and for chat purposes fits better than everything else.

### TODO + ideas

Chat ideas:

- Using NATS JetStream for pub/sub and persisting comments
- Auth via JWT by auth service (jwt contains user, channel)
- service for web socket connection and publish message
  - Publish message could be done by WS msg
  - We do need api for getting all messages (?)
  - Think about how to deal with connection lags ? (Internet connection)
- By WS connection state we could maintain a Users statuses
  - We do need api for getting all users before connection

## About The Project

The project build on a NATS's built-in distributed persistence system called JetStream.
The idea behind that is that stream could be used for distribution and for persisting messaging. Which makes the system pretty lightweight and reliable.

After the implementation, I think that JetStream isn't the best choice for that, the fact that it is still in beta limits the functionalities and reflects on usability.

A more strict forward solution could just be using Redis PubSun for distribution messages and Redis KeyValue for persisting messages and users.

### Main flow

```js
// First u need to authorized
const body = new FormData();
body.append('name', 'John Snow');
body.append('email', 'john.snow@gmail.com');
body.append('channel', 'lobby');

const token = await fetch("/auth", {
                            method: "POST",
                            body
                          }).then((response) => {
                            return response.json();
                          }).then((body) => body['token']);
// 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoiOTdmOTYwYTQ4YjQxOWZmNTVhMmFiZmJhOWQ4ZGUyMjMiLCJuYW1lIjoiSm9obiBTbm93IiwiYXZhdGFyIjoiaHR0cHM6Ly93d3cuZ3JhdmF0YXIuY29tL2F2YXRhci81MzIzOTFkNWEyYmFkYTEyYzVmMTYzNWU5Y2JkNDljYz9zPTEyOCJ9LCJjaGFubmVsIjoibG9iYnkiLCJleHAiOjE2NTkzNjI0NDUsImlzcyI6ImNoaXRjaGF0In0.Mq9baOmORsG-UP-TMee5DM9ff1YfaqKMivrwO1xhDoc'

// Then u can get auth object
const auth = await fetch("/auth?token=" + token).then((response) => {
                            return response.json();
                          });
// {"user":{"id":"97f960a48b419ff55a2abfba9d8de223","name":"John Snow","avatar":"https://www.gravatar.com/avatar/532391d5a2bada12c5f1635e9cbd49cc?s=128"},"channel":"lobby","exp":1659362445,"iss":"chitchat"}

// Then let's get channel history and active users
const messages = await fetch("/messages?token=" + token).then((response) => {
                            return response.json();
                          });
// {"messages": [{"from_user": {"id": "<id>", "name": "Andrii", avatar: "<avatar-url>"}, "sent_at": "2022-07-29T11:39:33.185636614Z", "text": "Hello", "type": "message"}]}

const users = await fetch("/users?token=" + token).then((response) => {
                            return response.json();
                          });
// {"users": [{"id": "<id>", "name": "Andrii", avatar: "<avatar-url>"}]}

// For getting new messages we need to open a WebSocket
const socket = new WebSocket("ws://localhost:8080/channel?token=" + token);

// Connection opened
socket.addEventListener('open', function (event) {
    // all good
});

// Listen for messages
socket.addEventListener('message', function (event) {
    console.log('Message from server ', JSON.parse(event.data));

    /**
     * Text messages
     *
     * {
     *  "type": "message",
     *  "from_user": {
     *    "id": "f73dd4c5575e2b89770f2dead82b537c",
     *    "name": "Andrii",
     *       "avatar": "https://www.gravatar.com/avatar/d5b1e46917bbeccee72465a56878a3e9?s=128"
     *  },
     *  "sent_at": "2022-07-29T14:10:50.683901889Z",
     *  "text": "hey"
     * }
     */

    /**
     * Join messages. When someone new joining
     *
     * {
     *   "type": "join",
     *   "from_user": {
     *       "id": "f73dd4c5575e2b89770f2dead82b537c",
     *       "name": "Andrii",
     *       "avatar": "https://www.gravatar.com/avatar/d5b1e46917bbeccee72465a56878a3e9?s=128"
     *    },
     * }
     */

    /**
     * Leave messages. When someone leaving
     *
     * {
     *   "type": "leave",
     *   "from_user": {
     *       "id": "f73dd4c5575e2b89770f2dead82b537c",
     *       "name": "Andrii",
     *       "avatar": "https://www.gravatar.com/avatar/d5b1e46917bbeccee72465a56878a3e9?s=128"
     *    },
     * }
     */
});

// To post a message to the chat, just
socket.send("Hello!");
```

### Built With

- Core
  - NATS JetStream - for streaming/persisting messages, track user presence, etc.
  - Envoy - as edge proxy for WS and HTTP, load balancing server
  - Docker Compose - for orchestration
- Server
  - Echo - High performance, extensible, minimalist Go web framework.
  - gorilla/websocket - A fast, well-tested and widely used WebSocket implementation for Go.
  - golang-jwt/jwt - For JWT Auth flow.
- Client
  - create-react-app + typescript - scaffolding React app
  - react-use-websocket - React Hook for WebSocket communication
  - ChakraUI - for building UI

## Getting Started

## Installation

Just make sure that you have Docker Compose in place.

### Prerequisites

- Docker Compose

## Usage

```sh
git clone git@github.com:faustman/chitchat.git
cd chitchat

docker compose up
```

then open `http://localhost:8080/` in browser.

## Contact

### Andrii Tytar

- Telegram [@junglecore](https://t.me/junglecore)
- Email: andrii.tytar@gmail.com
