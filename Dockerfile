FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/notifyhub cmd/notifyhub/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/notifyhub .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./notifyhub"]