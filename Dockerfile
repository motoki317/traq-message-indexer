FROM golang:1.14

WORKDIR /go/src/github.com/motoki317/traq-message-indexer
COPY . .

RUN go mod download
RUN go build -o app

CMD ["./app"]
