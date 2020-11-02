FROM golang:1.15.3 AS build

WORKDIR /go/src/github.com/motoki317/traq-message-indexer

COPY ./go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app -ldflags "-s -w"

FROM alpine:latest

RUN apk add --update ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=build /go/src/github.com/motoki317/traq-message-indexer/app /app
COPY ./mysql/init /mysql/init

CMD ["/app"]
