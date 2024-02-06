FROM golang:alpine AS builder

RUN apk add --no-cache zlib-dev sqlite-dev pkgconfig build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 go build -tags "fts5" -o knowhere

FROM alpine:latest

RUN apk add --no-cache openssh sqlite curl

WORKDIR /app
COPY --from=builder /app/knowhere .

EXPOSE 3000
ENTRYPOINT ["/app/knowhere", "server", "--port", "3000", "--db", "/var/osm/colorado.db" ]