FROM golang:alpine AS builder

RUN apk add --no-cache zlib-dev sqlite-dev pkgconfig build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=1 go build -tags "fts5" -o knowhere

FROM alpine:latest

RUN apk add --no-cache bash curl openssh sqlite

WORKDIR /app
COPY --from=builder /app/knowhere .
COPY --from=builder /app/bin/entrypoint.sh .

EXPOSE 3000
ENTRYPOINT [ "bash" ]
CMD ["/app/entrypoint.sh"]
