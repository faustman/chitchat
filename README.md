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

## About The Project

### Built With

## Getting Started

### Prerequisites

### Installation

## Usage

## Contact

### Andrii Tytar

- Telegram [@junglecore](https://t.me/junglecore)
- Email: andrii.tytar@gmail.com
