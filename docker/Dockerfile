FROM golang:alpine AS builder

LABEL stage=gobuilder
ENV CGO_ENABLED=0 GOOS=linux
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

COPY ../go.mod ../go.sum ./

RUN go mod download
COPY ../ ./

RUN go build -ldflags="-s -w" -o /app/core ./core/main.go
FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York

WORKDIR /app
COPY --from=builder /app/core /app/core

CMD ["/app/core"]