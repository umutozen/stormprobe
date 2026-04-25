FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o stormprobe ./cmd

FROM projectdiscovery/katana:latest AS katana
FROM projectdiscovery/httpx:latest AS httpx

FROM alpine:latest
RUN apk --no-cache add ca-certificates chromium
WORKDIR /app
COPY --from=builder /app/stormprobe .
COPY --from=katana /usr/local/bin/katana /usr/local/bin/katana
COPY --from=httpx /usr/local/bin/httpx /usr/local/bin/httpx
RUN mkdir -p /app/outputs
VOLUME ["/app/outputs"]

ENTRYPOINT ["./stormprobe"]
