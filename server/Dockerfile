# This Dockerfile from official golang docker image example
# TODO: for production there is two step build layer should be implemented
# TODO: for better development experience, use server auto-reload on file change
# like https://github.com/cespare/reflex or similar

FROM golang:1.18-alpine

WORKDIR /go/src/server

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/server ./...

EXPOSE 4000

CMD ["server"]