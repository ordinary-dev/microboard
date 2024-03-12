FROM golang:1.22-alpine AS builder

RUN apk add gcc
RUN apk add musl-dev
RUN apk add vips-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD . .

RUN go build -o microboard cmd/api/*.go

FROM alpine:3.18.3

RUN apk add vips
RUN apk add ffmpeg

WORKDIR /app
COPY --from=builder /app/microboard /usr/local/bin/microboard
COPY assets ./assets
COPY templates ./templates

RUN mkdir db
COPY db/migrations ./db/migrations

RUN mkdir /app/uploads

ENV MB_ISPRODUCTION="true"
ENV MB_UPLOADDIR="/app/uploads"
ENV MB_PREVIEWDIR="/app/previews"
ENV MB_LOGLEVEL="warning"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/microboard"]
