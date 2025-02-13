FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Set architecture to match container's OS
RUN GOOS=linux GOARCH=amd64 go build -o gpsd-api-gateway .

FROM debian:bookworm

WORKDIR /app
COPY --from=builder /app/gpsd-api-gateway .

EXPOSE 3000

CMD ["./gpsd-api-gateway"]