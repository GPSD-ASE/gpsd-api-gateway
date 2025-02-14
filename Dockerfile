FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o gpsd-api-gateway ./internal

FROM debian:bookworm

WORKDIR /app
COPY --from=builder /app/gpsd-api-gateway .

EXPOSE 3000

CMD ["./gpsd-api-gateway"]