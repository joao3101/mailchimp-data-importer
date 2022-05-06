FROM golang:1.18 as builder

# Enable modules support
ENV GO111MODULE=on

RUN mkdir -p /build
WORKDIR /build

COPY go.mod go.sum ./

COPY . .

RUN go get ./...
RUN go mod download

ENV CGO_ENABLED=0
# Compile static binary
RUN go build -ldflags "-s -w" ./cmd/

##################################################
# Final image
FROM alpine

RUN mkdir -p /app
WORKDIR /app

CMD ["/app"]
