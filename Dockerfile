# Билд стадия
FROM golang:1.24.2-alpine AS builder

WORKDIR /booking

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /booking/mybooking cmd/server/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /booking/mybooking .

COPY --from=builder /booking/config ./config/

COPY --from=builder /booking/.env .

EXPOSE 8080

CMD ["./mybooking"]