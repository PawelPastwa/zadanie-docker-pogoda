FROM golang:1.26.4-alpine AS builder

WORKDIR /app

RUN go mod init weatherapp

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:3.24

LABEL org.opencontainers.image.authors="Pawel Pastwa"
LABEL org.opencontainers.image.title="Zadanie 1 - Aplikacja Pogodowa"

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /app/app .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./app"]
