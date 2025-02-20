FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o gpsd-api-gateway ./internal

FROM debian:bookworm

WORKDIR /app

RUN apt-get update && apt-get install -y curl lsof net-tools ca-certificates libc6

COPY --from=builder /app/gpsd-api-gateway .

EXPOSE 3000

CMD ["./gpsd-api-gateway"]