FROM golang:1.14 as builder
WORKDIR /bot
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

FROM alpine:3.11
RUN apk add --no-cache ca-certificates
WORKDIR /bot
COPY --from=builder /bot/server server
ENTRYPOINT ["/bot/server"]
