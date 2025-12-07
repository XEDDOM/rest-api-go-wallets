FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/server .
COPY cmd/api/config.env config.env

EXPOSE 3000

CMD ["./server"]
