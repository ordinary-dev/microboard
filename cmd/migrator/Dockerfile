FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o migrator cmd/migrator/*.go

FROM alpine:3.18.3

WORKDIR /app
COPY --from=builder /app/migrator /usr/local/bin/microboard-migrator

RUN mkdir db
COPY db/migrations ./db/migrations

ENV MB_LOGLEVEL="info"

ENTRYPOINT ["/usr/local/bin/microboard-migrator"]
