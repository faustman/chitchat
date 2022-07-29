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
