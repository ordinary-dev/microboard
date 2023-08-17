FROM golang:1.21.0-alpine3.18 AS builder

RUN apk add gcc
RUN apk add musl-dev
RUN apk add vips

WORKDIR /app

COPY src/go.mod src/go.sum ./
RUN go mod download

ADD src .

RUN go build -o main

FROM alpine:3.18.3

RUN apk add vips
RUN apk add ffmpeg

WORKDIR /app
COPY --from=builder /app/main /usr/local/bin/microboard
COPY src/assets ./assets
COPY src/templates ./templates
COPY src/migrations ./migrations

RUN mkdir /app/uploads

ENV MB_ISPRODUCTION="true"
ENV MB_UPLOADDIR="/app/uploads"
ENV MB_PREVIEWDIR="/app/previews"
ENV MB_LOGLEVEL="warning"

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/microboard"]
