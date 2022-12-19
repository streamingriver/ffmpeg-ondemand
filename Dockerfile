FROM golang:1.17-alpine as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /ffmpeg-wrapper


FROM alpine:latest

RUN \
  apk add --update bash supervisor inotify-tools && \
  rm -rf /var/cache/apk/*

COPY --from=mwader/static-ffmpeg:4.4.1 /ffmpeg  /ffmpeg
COPY --from=builder /ffmpeg-wrapper /ffmpeg-wrapper

ENV APP_ROOT=/dev/shm
ENV APP_NAME="default"
ENV APP_BIND=":9999"

EXPOSE 9999

ENTRYPOINT ["/ffmpeg-wrapper"]
