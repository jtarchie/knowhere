FROM golang:alpine AS builder

RUN apk add --no-cache zlib-dev sqlite-dev pkgconfig build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN --mount=type=cache,target="/root/.cache/go-build" go build -ldflags="-w -s" -o zstdseek github.com/SaveTheRbtz/zstd-seekable-format-go/cmd/zstdseek
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_CFLAGS="-O3" CGO_ENABLED=1 go build -ldflags="-w -s" -tags "fts5" -o knowhere

FROM alpine:latest

RUN apk add --no-cache \
  bash \
  curl \
  openssh \
  rclone \
  sqlite

WORKDIR /app
COPY --from=builder /app/knowhere .
COPY --from=builder /app/zstdseek .
COPY --from=builder /app/bin/entrypoint.sh .

EXPOSE 3000
ENTRYPOINT [ "bash" ]
CMD ["/app/entrypoint.sh"]
