FROM --platform=linux/amd64 golang:alpine AS builder

RUN apk add --no-cache zlib-dev sqlite-dev pkgconfig build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 go build -o knowhere

FROM --platform=linux/amd64 alpine:latest

RUN apk add --no-cache openssh sqlite

WORKDIR /app
COPY --from=builder /app/knowhere .

EXPOSE 3000
ENTRYPOINT ["/app/knowhere", "server", "--port", "3000", "--db", "/var/osm/test.db" ]