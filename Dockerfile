FROM golang:1.23-alpine AS builder

RUN apk add --no-cache build-base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p ./build && \
    CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o ./build/main ./src/core

FROM alpine:latest

WORKDIR /app/build

COPY --from=builder /app/build/main .
COPY --from=builder /app/.env ../ 
COPY --from=builder /app/src/config/configs/base.yaml ../src/config/configs/base.yaml 

EXPOSE 8080

CMD ["./main"]
